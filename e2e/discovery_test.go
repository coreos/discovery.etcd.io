package e2e

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/coreos/etcd/pkg/fileutil"
)

var discoveryExec = "../bin/discovery"

func init() {
	if !fileutil.Exist(discoveryExec) {
		log.Fatalf("could not find %q", discoveryExec)
	}
}

func TestDiscovery_size_1(t *testing.T) { testDiscovery(t, 1) }
func TestDiscovery_size_3(t *testing.T) { testDiscovery(t, 3) }
func TestDiscovery_size_5(t *testing.T) { testDiscovery(t, 5) }
func TestDiscovery_size_7(t *testing.T) { testDiscovery(t, 7) }
func testDiscovery(t *testing.T, size int) {
	cfg := defaultConfig
	cfg.execPath = etcdExecLatest
	cfg.clusterSize = size
	epc, err := cfg.NewEtcdProcessCluster()
	if err != nil {
		t.Fatal(err)
	}
	if err = epc.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err = epc.Stop(10 * time.Second); err != nil {
			t.Fatal(err)
		}
	}()

	// set up discovery server for each node
	procs := make([]*discoveryProcess, size)
	errc := make(chan error)
	for i := 0; i < size; i++ {
		dcfg := discoveryProcessConfig{
			execPath:      discoveryExec,
			etcdEp:        epc.procs[i].cfg.curl.String(),
			discoveryHost: fmt.Sprintf("https://test-%d.etcd.io", i),
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

	// create a token from each discovery endpoint
	// simulate 'curl -L https://discovery.etcd.io/new?size=5'
	for i := 0; i < size; i++ {
		dhost := procs[i].cfg.discoveryHost
		tokenFunc := func(txt string) bool {
			if !strings.HasPrefix(txt, dhost+"/") {
				return false
			}
			token := strings.Replace(txt, dhost+"/", "", 1)
			return isAlphanumeric(token)
		}
		req := cURLReq{
			timeout:  5 * time.Second,
			endpoint: fmt.Sprintf("http://localhost:%d/new?size=5", procs[i].cfg.webPort),
			method:   http.MethodGet,
			expFunc:  tokenFunc,
		}
		token, err := req.Send()
		if err != nil {
			t.Fatal(err)
		}
		if !tokenFunc(token) {
			t.Fatalf("unexpected token %q", token)
		}
	}
}

var isAlphanumeric = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString
