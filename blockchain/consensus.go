package blockchain

import (
	"math/big"
	"math/rand"
	"sync"
	"time"
)

type Stakeholder struct {
	Address string
	Stake   *big.Int
}

type Consensus struct {
	Stakeholders []Stakeholder
	Mutex        sync.Mutex
}

func (c *Consensus) AddStakeholder(address string, stake *big.Int) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.Stakeholders = append(c.Stakeholders, Stakeholder{Address: address, Stake: stake})
}

func (c *Consensus) SelectLeader() string {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if len(c.Stakeholders) == 0 {
		return ""
	}

	totalStake := big.NewInt(0)
	for _, s := range c.Stakeholders {
		totalStake.Add(totalStake, s.Stake)
	}

	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	threshold := big.NewInt(r.Int63n(totalStake.Int64()))

	accumulator := big.NewInt(0)
	for _, s := range c.Stakeholders {
		accumulator.Add(accumulator, s.Stake)
		if accumulator.Cmp(threshold) > 0 {
			return s.Address
		}
	}

	return ""
}
