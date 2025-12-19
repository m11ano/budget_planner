package backoff

import (
	"container/heap"
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type ControllerOption func(*controllerOptions)

type SessionOption func(*sessionOptions)

const DefaultGroupID = "default"

type Controller struct {
	ctx        context.Context
	cancel     context.CancelFunc
	heapCtrl   *heapController
	sessions   map[string]*Session
	groupCfg   map[string]*sessionOptions
	mu         sync.RWMutex
	opts       *controllerOptions
	evictGoing int32
}

func NewController(ctrlOpts ...ControllerOption) *Controller {
	ctx, cancel := context.WithCancel(context.Background())

	cOpts := defaultControllerOptions()
	for _, opt := range ctrlOpts {
		opt(cOpts)
	}

	ctrl := &Controller{
		ctx:        ctx,
		cancel:     cancel,
		heapCtrl:   newHeapController(),
		sessions:   make(map[string]*Session),
		groupCfg:   make(map[string]*sessionOptions),
		opts:       cOpts,
		evictGoing: 1,
	}

	ctrl.groupCfg[DefaultGroupID] = defaultSessionOptions()

	go ctrl.evictLoop()
	return ctrl
}

func (c *Controller) SetConfigForGroup(groupID string, opts ...SessionOption) {
	cfg := defaultSessionOptions()
	for _, opt := range opts {
		opt(cfg)
	}

	c.mu.Lock()
	c.groupCfg[groupID] = cfg
	c.mu.Unlock()
}

func (c *Controller) GetIfExists(sessionKey string, groupID ...string) (*Session, bool) {
	useGroup := DefaultGroupID
	if len(groupID) > 0 && groupID[0] != "" {
		useGroup = groupID[0]
	}
	key := sessionKey + "_" + useGroup

	c.mu.RLock()
	sess, sessOK := c.sessions[key]
	_, cfgOK := c.groupCfg[useGroup]
	c.mu.RUnlock()

	if sessOK && cfgOK && !sess.isExpired() {
		return sess, true
	}

	return nil, false
}

func (c *Controller) Get(sessionKey string, groupIDs ...string) *Session {
	useGroup := DefaultGroupID
	if len(groupIDs) > 0 && groupIDs[0] != "" {
		useGroup = groupIDs[0]
	}
	key := sessionKey + "_" + useGroup

	c.mu.RLock()
	sess, sessOK := c.sessions[key]
	cfg, cfgOK := c.groupCfg[useGroup]
	c.mu.RUnlock()

	if sessOK && cfgOK && !sess.isExpired() {
		return sess
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if !cfgOK {
		cfg = defaultSessionOptions()
		c.groupCfg[useGroup] = cfg
	}

	sess = newSession(c.heapCtrl, key, cfg)
	c.sessions[key] = sess

	return sess
}

func (c *Controller) removeSession(s *Session) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.sessions, s.Key())
	atomic.StoreInt32(&s.expired, 1)
}

func (c *Controller) evictLoop() {
	defer atomic.StoreInt32(&c.evictGoing, 0)

	ticker := time.NewTicker(c.opts.ClearInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			var toRemove []*Session

			c.heapCtrl.mu.Lock()
			for c.heapCtrl.Len() > 0 {
				item := c.heapCtrl.items[0]
				if !item.session.isExpired() {
					break
				}
				heap.Pop(c.heapCtrl)
				toRemove = append(toRemove, item.session)
			}
			c.heapCtrl.mu.Unlock()

			for _, s := range toRemove {
				c.removeSession(s)
			}
		}
	}
}

func (c *Controller) Stop(ctx context.Context) {
	c.cancel()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if atomic.LoadInt32(&c.evictGoing) == 0 {
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}
