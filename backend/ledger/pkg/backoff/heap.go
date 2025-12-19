package backoff

import (
	"container/heap"
	"sync"
	"time"
)

type heapItem struct {
	session   *Session
	expiredAt time.Time
	index     int
}

type heapController struct {
	mu    sync.Mutex
	items []*heapItem
	idx   map[string]*heapItem
}

func newHeapController() *heapController {
	h := &heapController{idx: make(map[string]*heapItem)}
	heap.Init(h)
	return h
}

func (h *heapController) Len() int {
	return len(h.items)
}

func (h *heapController) Less(i, j int) bool {
	return h.items[i].expiredAt.Before(h.items[j].expiredAt)
}

func (h *heapController) Swap(i, j int) {
	h.items[i], h.items[j] = h.items[j], h.items[i]
	h.items[i].index = i
	h.items[j].index = j
}

func (h *heapController) Push(x interface{}) {
	// nolint
	item := x.(*heapItem)
	h.items = append(h.items, item)
	item.index = len(h.items) - 1
	h.idx[item.session.Key()] = item
}

func (h *heapController) Pop() interface{} {
	n := len(h.items)
	item := h.items[n-1]
	h.items = h.items[:n-1]
	delete(h.idx, item.session.Key())
	return item
}

func (h *heapController) add(s *Session) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.Push(&heapItem{
		session:   s,
		expiredAt: s.expiredAt,
	})
}

func (h *heapController) update(s *Session) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if item, ok := h.idx[s.key]; ok {
		item.expiredAt = s.expiredAt
		heap.Fix(h, item.index)
	}
}
