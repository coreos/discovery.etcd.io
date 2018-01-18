package handlers

import (
	"net/url"
	"sync"
)

// State is the discovery server configuration
// state shared between handlers.
type State struct {
	mu            sync.RWMutex
	etcdHost      string
	etcdCURL      *url.URL
	currentLeader string
	discHost      string
	useV3         bool // 'true' to use v3 API
}

func (st *State) endpoint() (ep string) {
	st.mu.RLock()
	ep = st.etcdHost
	st.mu.RUnlock()
	return ep
}

func (st *State) getCurrentLeader() (leader string) {
	st.mu.RLock()
	leader = st.currentLeader
	st.mu.RUnlock()
	return leader
}

func (st *State) setCurrentLeader(leader string) {
	st.mu.Lock()
	st.currentLeader = leader
	st.mu.Unlock()
}

func (st *State) isV3() (useV3 bool) {
	st.mu.RLock()
	useV3 = st.useV3
	st.mu.RUnlock()
	return useV3
}
