package consensus

import (
	"math/rand"
	"sort"
	"sync"
)

type Node struct {
	ID    string
	Coins int // Количество монет у узла
}

var nodes []Node
var mu sync.Mutex

func RegisterNode(id string, coins int) {
	mu.Lock()
	defer mu.Unlock()
	nodes = append(nodes, Node{ID: id, Coins: coins})
}

// Выбираем узел с минимальным балансом, при равенстве – случайный среди них
func SelectNode() string {
	mu.Lock()
	defer mu.Unlock()

	if len(nodes) == 0 {
		return ""
	}

	// Сортируем узлы по количеству монет
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Coins < nodes[j].Coins
	})

	// Находим все узлы с минимальным балансом
	minCoins := nodes[0].Coins
	var candidates []Node

	for _, node := range nodes {
		if node.Coins == minCoins {
			candidates = append(candidates, node)
		}
	}

	// Если несколько, выбираем случайный из них
	selected := candidates[rand.Intn(len(candidates))]
	return selected.ID
}
