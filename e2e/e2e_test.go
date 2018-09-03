package e2e

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/coreos/etcd/pkg/expect"
	"github.com/coreos/etcd/pkg/fileutil"
)

var (
	etcdExecPrev   string
	etcdExecLatest string

	etcdctlExecPrev   string
	etcdctlExecLatest string

	defaultConfig = etcdProcessClusterConfig{
		execPath:     "",
		keepDataDir:  false,
		clusterSize:  5,
		initialToken: "test-token",
	}
)

func init() {
	versions := []string{}
	vs := os.Getenv("ETCD_VERSIONS")
	for _, ver := range strings.Split(vs, " ") {
		versions = append(versions, ver)
	}
	if len(versions) != 2 {
		log.Fatalf("expected 2 etcd versions, got %+v", versions)
	}

	etcdExecPrev = fmt.Sprintf("../bin/etcd-%s", versions[0])
	etcdExecLatest = fmt.Sprintf("../bin/etcd-%s", versions[1])

	etcdctlExecPrev = fmt.Sprintf("../bin/etcdctl-%s", versions[0])
	etcdctlExecLatest = fmt.Sprintf("../bin/etcdctl-%s", versions[1])

	if !fileutil.Exist(etcdExecPrev) {
		log.Fatalf("could not find %q", etcdExecPrev)
	}
	if !fileutil.Exist(etcdExecLatest) {
		log.Fatalf("could not find %q", etcdExecLatest)
	}
	if !fileutil.Exist(etcdctlExecPrev) {
		log.Fatalf("could not find %q", etcdctlExecPrev)
	}
	if !fileutil.Exist(etcdctlExecLatest) {
		log.Fatalf("could not find %q", etcdctlExecLatest)
	}
}

func tokenFunc(txt string) bool {
	if !strings.HasPrefix(txt, "http://localhost:") {
		return false
	}
	token := path.Base(txt)
	return isAlphanumeric(token)
}

func TestE2E_size_1_no_discovery(t *testing.T) { testE2E(t, 1, false) }
func TestE2E_size_3_no_discovery(t *testing.T) { testE2E(t, 3, false) }
func TestE2E_size_5_no_discovery(t *testing.T) { testE2E(t, 5, false) }
func TestE2E_size_7_no_discovery(t *testing.T) { testE2E(t, 7, false) }
func TestE2E_size_1_discovery(t *testing.T)    { testE2E(t, 1, true) }
func TestE2E_size_3_discovery(t *testing.T)    { testE2E(t, 3, true) }
func TestE2E_size_5_discovery(t *testing.T)    { testE2E(t, 5, true) }
func TestE2E_size_7_discovery(t *testing.T)    { testE2E(t, 7, true) }
func testE2E(t *testing.T, size int, useDiscovery bool) {
	cfg := defaultConfig
	cfg.execPath = etcdExecPrev
	cfg.clusterSize = size

	etcdClus1, err := cfg.NewEtcdProcessCluster()
	if err != nil {
		t.Fatal(err)
	}
	if err = etcdClus1.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err = etcdClus1.Stop(10 * time.Second); err != nil {
			t.Fatal(err)
		}
	}()

	etcdServer := etcdClus1
	if useDiscovery {
		// set up discovery server for each node
		procs := make([]*discoveryProcess, size)
		errc := make(chan error)
		for i := 0; i < size; i++ {
			dcfg := discoveryProcessConfig{
				execPath: discoveryExec,
				etcdEp:   etcdClus1.procs[i].cfg.curl.String(),
			}
			procs[i] = dcfg.NewDiscoveryProcess()
			go func(dp *discoveryProcess) {
				errc <- dp.Start()
			}(procs[i])
		}
		for i := 0; i < size; i++ {
			if err = <-errc; err != nil {
				t.Fatal(err)
			}
		}
		req := cURLReq{
			timeout:  5 * time.Second,
			endpoint: fmt.Sprintf("http://localhost:%d/new?size=%d", procs[0].cfg.webPort, size),
			method:   http.MethodGet,
			expFunc:  tokenFunc,
		}
		var token string
		token, err = req.Send()
		if err != nil {
			t.Fatal(err)
		}
		if !tokenFunc(token) {
			t.Fatalf("unexpected token %q", token)
		}

		// run etcd on top of discovery
		etcdCfg := defaultConfig
		etcdCfg.execPath = etcdExecPrev
		etcdCfg.clusterSize = size
		etcdCfg.discoveryToken = token
		etcdClus2, err := etcdCfg.NewEtcdProcessCluster()
		if err != nil {
			t.Fatal(err)
		}
		if err = etcdClus2.Start(); err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err = etcdClus2.Stop(10 * time.Second); err != nil {
				t.Fatal(err)
			}
		}()

		// validate token writes on etcd starts
		ctlv2Discovery := etcdctlCtl{
			api:         2,
			execPath:    etcdctlExecLatest,
			endpoints:   etcdClus1.ClientEndpoints(),
			dialTimeout: 7 * time.Second,
		}
		if err = ctlv2Discovery.Get("l", KV{
			Key: path.Join("_etcd", "registry", path.Base(token), "_config", "size"),
			Val: strconv.Itoa(size),
		}); err != nil {
			t.Fatal(err)
		}

		// validate discovery server register
		var resp *http.Response
		resp, err = http.Get(token)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		var cresp client.Response
		if err = json.NewDecoder(resp.Body).Decode(&cresp); err != nil {
			t.Fatal(err)
		}
		if len(cresp.Node.Nodes) != size {
			t.Fatalf("len(cresp.Node.Nodes) expected %d, got %d", size, len(cresp.Node.Nodes))
		}

		etcdServer = etcdClus2
	}

	ctlv2 := etcdctlCtl{
		api:         2,
		execPath:    etcdctlExecLatest,
		endpoints:   etcdServer.ClientEndpoints(),
		dialTimeout: 7 * time.Second,
	}
	ctlv3 := ctlv2
	ctlv3.api = 3
	kvs := []KV{}
	for i := 0; i < 5; i++ {
		kvs = append(kvs, KV{Key: fmt.Sprintf("foo%d", i), Val: fmt.Sprintf("var%d", i)})
	}
	for _, kv := range kvs {
		if err = ctlv2.Put(kv.Key, kv.Val); err != nil {
			t.Fatal(err)
		}
		if err = ctlv3.Put(kv.Key, kv.Val); err != nil {
			t.Fatal(err)
		}
	}
	time.Sleep(3 * time.Second) // wait for log replication
	for _, kv := range kvs {
		if err = ctlv2.Get("l", kv); err != nil {
			t.Fatal(err)
		}
		if err = ctlv3.Get("l", kv); err != nil {
			t.Fatal(err)
		}
	}

	// test etcd version upgrades
	for i := 0; i < etcdServer.cfg.clusterSize; i++ {
		if err = etcdServer.procs[i].Stop(); err != nil {
			t.Fatal(err)
		}

		// restart with latest etcd (upgrade)
		etcdServer.procs[i].cfg.execPath = etcdExecLatest
		child, err := expect.NewExpect(etcdServer.procs[i].cfg.execPath, etcdServer.procs[i].cfg.args...)
		if err != nil {
			t.Fatal(err)
		}
		etcdServer.procs[i].proc = child
		etcdServer.procs[i].donec = make(chan struct{})
		if err = etcdServer.procs[i].waitReady(); err != nil {
			t.Fatal(err)
		}

		// check data after upgrade
		ctlv2 := etcdctlCtl{
			api:         2,
			execPath:    etcdctlExecLatest,
			endpoints:   []string{etcdServer.procs[i].cfg.curl.String()},
			dialTimeout: 7 * time.Second,
		}
		ctlv3 := ctlv2
		ctlv3.api = 3
		for _, kv := range kvs {
			if err = ctlv2.Get("s", kv); err != nil {
				t.Fatal(err)
			}
			if err = ctlv3.Get("s", kv); err != nil {
				t.Fatal(err)
			}
		}
	}
}
