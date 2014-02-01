// Lockstring implements a string protected by a RWLock
package lockstring

import (
	"sync"
)

type LockString struct {
	sync.RWMutex
	str string
}

func (r *LockString) String() string {
	r.RLock()
	str := r.str
	r.RUnlock()
	return str
}

func (r *LockString) Set(str string) {
	r.Lock()
	r.str = str
	r.Unlock()
}


