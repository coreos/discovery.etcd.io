package e2e

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/coreos/etcd/pkg/expect"
)

type etcdProcessClusterConfig struct {
	execPath       string
	keepDataDir    bool
	clusterSize    int
	initialToken   string
	discoveryToken string
}

type etcdProcessConfig struct {
	execPath string
	args     []string

	dataDirPath string
	keepDataDir bool

	name string
	curl url.URL
	purl url.URL

	initialToken   string
	initialCluster string
	discoveryToken string
}

type etcdProcess struct {
	cfg   *etcdProcessConfig
	proc  *expect.ExpectProcess
	donec chan struct{} // closed when Interact() terminates
}

type etcdProcessCluster struct {
	cfg   *etcdProcessClusterConfig
	procs []*etcdProcess
}

var etcdBasePort int32 = 20000

// NewEtcdProcessCluster generates a new 'etcdProcessCluster'.
func (cfg *etcdProcessClusterConfig) NewEtcdProcessCluster() (*etcdProcessCluster, error) {
	cport := int(atomic.LoadInt32(&etcdBasePort))
	atomic.AddInt32(&etcdBasePort, int32(3*cfg.clusterSize))

	epc := &etcdProcessCluster{
		cfg:   cfg,
		procs: make([]*etcdProcess, cfg.clusterSize),
	}

	ics := make([]string, cfg.clusterSize)
	for i := 0; i < cfg.clusterSize; i++ {
		name := fmt.Sprintf("test%d%d.etcd", cport, i)
		curl := url.URL{Scheme: "http", Host: fmt.Sprintf("localhost:%d", cport+2*i)}
		purl := url.URL{Scheme: "http", Host: fmt.Sprintf("localhost:%d", cport+2*i+1)}
		dataDirPath, err := ioutil.TempDir(os.TempDir(), name)
		if err != nil {
			return nil, err
		}
		ics[i] = fmt.Sprintf("%s=%s", name, purl.String())

		epc.procs[i] = &etcdProcess{
			cfg: &etcdProcessConfig{
				execPath: cfg.execPath,
				args: []string{
					"--name", name,
					"--listen-client-urls", curl.String(),
					"--advertise-client-urls", curl.String(),
					"--listen-peer-urls", purl.String(),
					"--initial-advertise-peer-urls", purl.String(),
					"--initial-cluster-token", cfg.initialToken,
					"--data-dir", dataDirPath,
				},

				dataDirPath: dataDirPath,
				keepDataDir: cfg.keepDataDir,

				name: name,
				curl: curl,
				purl: purl,

				initialToken:   cfg.initialToken,
				initialCluster: "",
				discoveryToken: cfg.discoveryToken,
			},
		}
	}

	for i := 0; i < cfg.clusterSize; i++ {
		if epc.procs[i].cfg.discoveryToken == "" {
			cs := strings.Join(ics, ",")
			epc.procs[i].cfg.initialCluster = cs
			epc.procs[i].cfg.args = append(epc.procs[i].cfg.args, "--initial-cluster", cs)
		} else {
			epc.procs[i].cfg.args = append(epc.procs[i].cfg.args, "--discovery", epc.procs[i].cfg.discoveryToken)
		}
	}
	return epc, nil
}

func (ep *etcdProcess) Stop() error {
	if ep == nil {
		return nil
	}
	if err := ep.proc.Stop(); err != nil {
		return err
	}
	<-ep.donec
	return nil
}

func (epc *etcdProcessCluster) ClientEndpoints() []string {
	eps := make([]string, epc.cfg.clusterSize)
	for i, ep := range epc.procs {
		eps[i] = ep.cfg.curl.String()
	}
	return eps
}

func (epc *etcdProcessCluster) GRPCEndpoints() []string {
	eps := make([]string, epc.cfg.clusterSize)
	for i, ep := range epc.procs {
		eps[i] = ep.cfg.curl.Host
	}
	return eps
}

func (epc *etcdProcessCluster) Stop(d time.Duration) (err error) {
	errc := make(chan error, 1)
	go func() { errc <- epc.stop() }()
	select {
	case err := <-errc:
		return err
	case <-time.After(d):
		return fmt.Errorf("took longer than %v to Stop cluster", d)
	}
}

func (epc *etcdProcessCluster) stop() (err error) {
	for _, p := range epc.procs {
		if p == nil {
			continue
		}
		if curErr := p.Stop(); curErr != nil {
			if err != nil {
				err = fmt.Errorf("%v; %v", err, curErr)
			} else {
				err = curErr
			}
		}
		if !p.cfg.keepDataDir {
			os.RemoveAll(p.cfg.dataDirPath)
		}
	}
	return err
}

func (epc *etcdProcessCluster) Start() error {
	for i := 0; i < epc.cfg.clusterSize; i++ {
		cfg := epc.procs[i].cfg
		if !cfg.keepDataDir {
			os.RemoveAll(cfg.dataDirPath)
		}
	}

	for i := 0; i < epc.cfg.clusterSize; i++ {
		child, err := expect.NewExpect(epc.procs[i].cfg.execPath, epc.procs[i].cfg.args...)
		if err != nil {
			return err
		}
		epc.procs[i].proc = child
		epc.procs[i].donec = make(chan struct{})
	}

	readyC := make(chan error, epc.cfg.clusterSize)
	for i := range epc.procs {
		go func(n int) {
			readyC <- epc.procs[n].waitReady()
		}(i)
	}
	for range epc.procs {
		if err := <-readyC; err != nil {
			epc.Stop(10 * time.Second)
			return err
		}
	}
	return nil
}

func (ep *etcdProcess) waitReady() error {
	defer close(ep.donec)
	return waitReadyExpectProc(ep.proc)
}

func waitReadyExpectProc(exproc *expect.ExpectProcess) error {
	readyStrs := []string{"enabled capabilities for version", "published"}
	c := 0
	matchSet := func(l string) bool {
		for _, s := range readyStrs {
			if strings.Contains(l, s) {
				c++
				break
			}
		}
		return c == len(readyStrs)
	}
	_, err := exproc.ExpectFunc(matchSet)
	return err
}

type etcdctlCtl struct {
	api         int
	execPath    string
	endpoints   []string
	dialTimeout time.Duration
}

func (ctl etcdctlCtl) PrefixArgs() []string {
	cmdArgs := []string{ctl.execPath}
	switch ctl.api {
	case 2:
		cmdArgs = append(cmdArgs, fmt.Sprintf("--total-timeout=%s", ctl.dialTimeout))
	case 3:
		cmdArgs = append(cmdArgs, fmt.Sprintf("--dial-timeout=%s", ctl.dialTimeout))
	default:
		panic(fmt.Errorf("unknown API version %d", ctl.api))
	}
	cmdArgs = append(cmdArgs, fmt.Sprintf("--endpoints=%s", strings.Join(ctl.endpoints, ",")))
	return cmdArgs
}

func (ctl etcdctlCtl) Put(k, v string) error {
	cmdArgs := ctl.PrefixArgs()
	var exp string
	switch ctl.api {
	case 2:
		cmdArgs = append(cmdArgs, "set", k, fmt.Sprintf("%q", v))
		exp = v
	case 3:
		os.Setenv("ETCDCTL_API", "3")
		defer os.Unsetenv("ETCDCTL_API")
		cmdArgs = append(cmdArgs, "put", k, fmt.Sprintf("%q", v))
		exp = "OK"
	}
	return spawnWithExpects(cmdArgs, exp)
}

func (ctl etcdctlCtl) Delete(key, val string, num int) error {
	cmdArgs := ctl.PrefixArgs()
	var exp string
	switch ctl.api {
	case 2:
		cmdArgs = append(cmdArgs, "rm", key)
		exp = fmt.Sprintf("PrevNode.Value: %s", val)
	case 3:
		os.Setenv("ETCDCTL_API", "3")
		defer os.Unsetenv("ETCDCTL_API")
		cmdArgs = append(cmdArgs, "del", key)
		exp = fmt.Sprintf("%d", num)
	}
	return spawnWithExpects(cmdArgs, exp)
}

// KV is key-value pair.
type KV struct {
	Key, Val string
}

func (ctl etcdctlCtl) Get(consistency string, kv KV) error {
	cmdArgs := ctl.PrefixArgs()
	var lines []string
	switch ctl.api {
	case 2:
		cmdArgs = append(cmdArgs, "get", kv.Key)
		if consistency == "l" {
			cmdArgs = append(cmdArgs, "--quorum")
		}
		lines = append(lines, kv.Val)
	case 3:
		os.Setenv("ETCDCTL_API", "3")
		defer os.Unsetenv("ETCDCTL_API")
		cmdArgs = append(cmdArgs, "get", kv.Key)
		cmdArgs = append(cmdArgs, fmt.Sprintf("--consistency=%s", consistency))
		lines = append(lines, kv.Key, kv.Val)
	}
	return spawnWithExpects(cmdArgs, lines...)
}

const noOutputLineCount = 2 // cov-enabled binaries emit PASS and coverage count lines

func spawnWithExpects(args []string, xs ...string) error {
	proc, err := expect.NewExpect(args[0], args[1:]...)
	if err != nil {
		return err
	}
	// process until either stdout or stderr contains
	// the expected string
	var (
		lines    []string
		lineFunc = func(txt string) bool { return true }
	)
	for _, txt := range xs {
		for {
			l, lerr := proc.ExpectFunc(lineFunc)
			if lerr != nil {
				proc.Close()
				return fmt.Errorf("%v (expected %q, got %q)", lerr, txt, lines)
			}
			lines = append(lines, l)
			if strings.Contains(l, txt) {
				break
			}
		}
	}
	perr := proc.Close()
	if len(xs) == 0 && proc.LineCount() != noOutputLineCount { // expect no output
		return fmt.Errorf("unexpected output (got lines %q, line count %d)", lines, proc.LineCount())
	}
	return perr
}
