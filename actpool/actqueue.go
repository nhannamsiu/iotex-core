// Copyright (c) 2018 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package actpool

import (
	"container/heap"
	"math/big"
	"sort"

	"github.com/pkg/errors"

	"github.com/iotexproject/iotex-core/blockchain/action"
)

type noncePriorityQueue []uint64

func (h noncePriorityQueue) Len() int           { return len(h) }
func (h noncePriorityQueue) Less(i, j int) bool { return h[i] < h[j] }
func (h noncePriorityQueue) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *noncePriorityQueue) Push(x interface{}) {
	in, ok := x.(uint64)
	if !ok {
		return
	}
	*h = append(*h, in)
}

func (h *noncePriorityQueue) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// ActQueue is the interface of actQueue
type ActQueue interface {
	Overlaps(action.Action) bool
	Put(action.Action) error
	FilterNonce(uint64) []action.Action
	SetStartNonce(uint64)
	StartNonce() uint64
	UpdateQueue(uint64) []action.Action
	SetPendingNonce(uint64)
	PendingNonce() uint64
	SetPendingBalance(*big.Int)
	PendingBalance() *big.Int
	Len() int
	Empty() bool
	PendingActs() []action.Action
	AllActs() []action.Action
}

// actQueue is a queue of actions from an account
type actQueue struct {
	// Map that stores all the actions belonging to an account associated with nonces
	items map[uint64]action.Action
	// Priority Queue that stores all the nonces belonging to an account. Nonces are used as indices for action map
	index noncePriorityQueue
	// Current nonce tracking the first action in queue
	startNonce uint64
	// Current pending nonce tracking previous actions that can be committed to the next block for the account
	pendingNonce uint64
	// Current pending balance for the account
	pendingBalance *big.Int
}

// NewActQueue create a new action queue
func NewActQueue() ActQueue {
	return &actQueue{
		items:          make(map[uint64]action.Action),
		index:          noncePriorityQueue{},
		startNonce:     uint64(1), // Taking coinbase Action into account, startNonce should start with 1
		pendingNonce:   uint64(1), // Taking coinbase Action into account, pendingNonce should start with 1
		pendingBalance: big.NewInt(0),
	}
}

// Overlap returns whether the current queue contains the given nonce
func (q *actQueue) Overlaps(act action.Action) bool {
	return q.items[act.Nonce()] != nil
}

// Put inserts a new action into the map, also updating the queue's nonce index
func (q *actQueue) Put(act action.Action) error {
	nonce := act.Nonce()
	if q.items[nonce] != nil {
		return errors.Wrapf(ErrNonce, "duplicate nonce")
	}
	heap.Push(&q.index, nonce)
	q.items[nonce] = act
	return nil
}

// FilterNonce removes all actions from the map with a nonce lower than the given threshold
func (q *actQueue) FilterNonce(threshold uint64) []action.Action {
	var removed []action.Action
	// Pop off priority queue and delete corresponding entries from map until the threshold is reached
	for q.index.Len() > 0 && (q.index)[0] < threshold {
		nonce := heap.Pop(&q.index).(uint64)
		removed = append(removed, q.items[nonce])
		delete(q.items, nonce)
	}
	return removed
}

// UpdateQueue updates the pending nonce and balance of the queue
func (q *actQueue) UpdateQueue(nonce uint64) []action.Action {
	// First, starting from the current pending nonce, incrementally find the next pending nonce
	// while updating pending balance if actions are payable
	for ; q.items[nonce] != nil; nonce++ {
		if !q.enoughBalance(q.items[nonce], true) {
			break
		}
	}
	q.pendingNonce = nonce

	// Find the index of new pending nonce within the queue
	sort.Sort(q.index)
	i := 0
	for ; i < q.index.Len(); i++ {
		if q.index[i] >= nonce {
			break
		}
	}
	// Case I: An unpayable action has been found while updating pending nonce/balance
	// Remove all the subsequent actions in the queue starting from the index of new pending nonce
	if q.items[nonce] != nil {
		return q.removeActs(i)
	}

	// Case II: All actions are payable while updating pending nonce/balance
	// Check all the subsequent actions in the queue starting from the index of new pending nonce
	// Find the nonce index of the first unpayable action
	// Remove all the subsequent actions in the queue starting from that index
	for ; i < q.index.Len(); i++ {
		nonce = q.index[i]
		act := q.items[nonce]
		if !q.enoughBalance(act, false) {
			break
		}
	}
	return q.removeActs(i)
}

// SetStartNonce sets the new start nonce for the queue
func (q *actQueue) SetStartNonce(nonce uint64) {
	q.startNonce = nonce
}

// StartNonce returns the current start nonce of the queue
func (q *actQueue) StartNonce() uint64 {
	return q.startNonce
}

// SetPendingNonce sets pending nonce for the queue
func (q *actQueue) SetPendingNonce(nonce uint64) {
	q.pendingNonce = nonce
}

// PendingNonce returns the current pending nonce of the queue
func (q *actQueue) PendingNonce() uint64 {
	return q.pendingNonce
}

// SetPendingBalance sets pending balance for the queue
func (q *actQueue) SetPendingBalance(balance *big.Int) {
	q.pendingBalance = balance
}

// PendingBalance returns the current pending balance of the queue
func (q *actQueue) PendingBalance() *big.Int {
	return q.pendingBalance
}

// Len returns the length of the action map
func (q *actQueue) Len() int {
	return len(q.items)
}

// Empty returns whether the queue of actions is empty or not
func (q *actQueue) Empty() bool {
	return q.Len() == 0
}

// PendingActs creates a consecutive nonce-sorted slice of actions
func (q *actQueue) PendingActs() []action.Action {
	if q.Len() == 0 {
		return []action.Action{}
	}
	acts := make([]action.Action, 0, len(q.items))
	nonce := q.startNonce
	for ; q.items[nonce] != nil; nonce++ {
		acts = append(acts, q.items[nonce])
	}
	return acts
}

// AllActs returns all the actions currently in queue
func (q *actQueue) AllActs() []action.Action {
	acts := make([]action.Action, 0, len(q.items))
	if q.Len() == 0 {
		return acts
	}
	sort.Sort(q.index)
	for _, nonce := range q.index {
		acts = append(acts, q.items[nonce])
	}
	return acts
}

// removeActs removes all the actions starting at idx from queue
func (q *actQueue) removeActs(idx int) []action.Action {
	removedFromQueue := make([]action.Action, 0)
	for i := idx; i < q.index.Len(); i++ {
		removedFromQueue = append(removedFromQueue, q.items[q.index[i]])
		delete(q.items, q.index[i])
	}
	q.index = q.index[:idx]
	heap.Init(&q.index)
	return removedFromQueue
}

// enoughBalance helps check whether queue's pending balance is sufficient for the given action
func (q *actQueue) enoughBalance(act action.Action, updateBalance bool) bool {
	cost, _ := act.Cost()
	if q.pendingBalance.Cmp(cost) < 0 {
		return false
	}

	if updateBalance {
		q.pendingBalance.Sub(q.pendingBalance, cost)
	}

	return true
}
