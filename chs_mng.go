package newbing

import "sync"

type StrChMng struct {
	cacheChs map[string]chan<- string
	cacheMu  *sync.RWMutex
}

func NewStrChMng() (mng *StrChMng) {
	mng = &StrChMng{
		cacheChs: make(map[string]chan<- string, 64),
		cacheMu:  &sync.RWMutex{},
	}
	return
}

func (mng *StrChMng) AddCh(uid string, ch chan<- string) {
	mng.cacheMu.Lock()
	defer mng.cacheMu.Unlock()
	mng.cacheChs[uid] = ch
}

func (mng *StrChMng) GetCh(uid string) (ch chan<- string, has bool) {
	mng.cacheMu.RLock()
	defer mng.cacheMu.RUnlock()
	ch, has = mng.cacheChs[uid]
	return
}

func (mng *StrChMng) DelCh(uid string) {
	mng.cacheMu.Lock()
	defer mng.cacheMu.Unlock()
	delete(mng.cacheChs, uid)
}
