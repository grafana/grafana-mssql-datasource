package sqleng

import (
	"database/sql"
	"errors"
	"sync"
	"time"
)

const userPoolIdleTTL = 5 * time.Minute

var errDBClosed = errors.New("database connection is closed")

type userPool struct {
	db       *sql.DB
	lastUsed time.Time
}

type userPoolCache struct {
	mu       sync.Mutex
	pools    map[string]*userPool
	disposed bool
	stop     chan struct{}
	sweeper  sync.WaitGroup
}

func (c *userPoolCache) load(key string) (*sql.DB, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	pool, ok := c.pools[key]
	if !ok {
		return nil, false
	}
	pool.lastUsed = time.Now()
	return pool.db, true
}

func (c *userPoolCache) loadOrStore(key string, db *sql.DB) (*sql.DB, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.disposed {
		_ = db.Close()
		return nil, errDBClosed
	}
	if pool, ok := c.pools[key]; ok {
		_ = db.Close()
		pool.lastUsed = time.Now()
		return pool.db, nil
	}
	if c.pools == nil {
		c.pools = map[string]*userPool{}
		stop := make(chan struct{})
		c.stop = stop
		c.sweeper.Go(func() { c.sweepUntil(stop) })
	}
	c.pools[key] = &userPool{db: db, lastUsed: time.Now()}
	return db, nil
}

func (c *userPoolCache) sweepUntil(stop <-chan struct{}) {
	ticker := time.NewTicker(userPoolIdleTTL / 2)
	defer ticker.Stop()
	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			c.closeIdle()
		}
	}
}

func (c *userPoolCache) closeIdle() {
	c.mu.Lock()
	var idle []*sql.DB
	for key, pool := range c.pools {
		if time.Since(pool.lastUsed) >= userPoolIdleTTL {
			idle = append(idle, pool.db)
			delete(c.pools, key)
		}
	}
	c.mu.Unlock()

	for _, db := range idle {
		_ = db.Close()
	}
}

func (c *userPoolCache) dispose() {
	c.mu.Lock()
	c.disposed = true
	pools := c.pools
	c.pools = nil
	stop := c.stop
	c.stop = nil
	c.mu.Unlock()

	if stop != nil {
		close(stop)
	}
	c.sweeper.Wait()
	for _, pool := range pools {
		_ = pool.db.Close()
	}
}
