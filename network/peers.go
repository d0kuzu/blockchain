package network

import "sync"

type Peers struct {
	List  []string
	Mutex sync.Mutex
}

func NewPeers(initialPeers []string) []string {
	return initialPeers
}

func (p *Peers) Add(peer string) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	for _, existing := range p.List {
		if existing == peer {
			return
		}
	}
	p.List = append(p.List, peer)
}

func (p *Peers) Get() []string {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	return p.List
}
