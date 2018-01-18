package e2e

import (
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/coreos/etcd/pkg/expect"
)

type discoveryProcessConfig struct {
	execPath      string
	etcdEp        string
	discoveryHost string
	webPort       int
}

type discoveryProcess struct {
	cfg   *discoveryProcessConfig
	proc  *expect.ExpectProcess
	donec chan struct{} // closed when Interact() terminates
}

var discoveryBasePort int32 = 30000

// NewDiscoveryProcess creates a new 'discoveryProcess'.
func (cfg *discoveryProcessConfig) NewDiscoveryProcess() *discoveryProcess {
	port := int(atomic.LoadInt32(&discoveryBasePort))
	atomic.AddInt32(&discoveryBasePort, 2)

	copied := *cfg
	copied.webPort = port

	return &discoveryProcess{cfg: &copied}
}

func (dp *discoveryProcess) Stop(d time.Duration) (err error) {
	errc := make(chan error, 1)
	go func() { errc <- dp.proc.Stop() }()
	select {
	case err := <-errc:
		return err
	case <-time.After(d):
		return fmt.Errorf("took longer than %v to Stop cluster", d)
	}
}

func (dp *discoveryProcess) Start() error {
	args := []string{
		dp.cfg.execPath,
		"--etcd", dp.cfg.etcdEp,
		"--host", dp.cfg.discoveryHost,
		"--addr", fmt.Sprintf(":%d", dp.cfg.webPort),
	}
	child, err := expect.NewExpect(args[0], args[1:]...)
	if err != nil {
		return err
	}
	dp.proc = child
	dp.donec = make(chan struct{})

	readyC := make(chan error)
	go func() {
		readyC <- dp.waitReady()
	}()
	select {
	case err = <-readyC:
		if err != nil {
			return err
		}
	case <-time.After(10 * time.Second):
		return fmt.Errorf("timed out waiting for discover server")
	}
	return nil
}

func (dp *discoveryProcess) waitReady() error {
	defer close(dp.donec)
	return waitReadyExpectProcDiscovery(dp.proc)
}

func waitReadyExpectProcDiscovery(exproc *expect.ExpectProcess) error {
	readyStrs := []string{"discovery server started", "discovery serving on"}
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

type cURLReq struct {
	timeout  time.Duration
	endpoint string
	method   string
	data     string
	expFunc  func(txt string) bool
}

func (req cURLReq) Send() (string, error) {
	cmdArgs := []string{"curl"}
	if req.timeout > 0 {
		cmdArgs = append(cmdArgs, "--max-time", fmt.Sprint(int(req.timeout.Seconds())))
	}
	cmdArgs = append(cmdArgs, "-L", req.endpoint)
	switch req.method {
	case http.MethodPost, http.MethodPut:
		dt := req.data
		if !strings.HasPrefix(dt, "{") { // for non-JSON value
			dt = "value=" + dt
		}
		cmdArgs = append(cmdArgs, "-X", req.method, "-d", dt)
	}

	proc, err := expect.NewExpect(cmdArgs[0], cmdArgs[1:]...)
	if err != nil {
		return "", err
	}

	var line string
	for {
		line, err = proc.ExpectFunc(req.expFunc)
		if err != nil {
			proc.Close()
			return "", err
		}
		break
	}
	return line, proc.Close()

}
