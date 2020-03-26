package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"sync"
	"time"
)

const feedUpdateTimeout = 30 * time.Second

// SourceWeight  the weight  of a  source
type SourceWeight int

const (
	// General the level of a source is general no good / no bad
	General = iota
	Trusted
	Signed
)

// RssNewsPool news from anySource
type RssNewsPool struct {
	Pool           []*RssNews
	UpdateInterval time.Duration
	sync.RWMutex
}

// NewRssNewsPool is to new a RssNewsPool
func NewRssNewsPool() *RssNewsPool {
	pool := RssNewsPool{UpdateInterval: 15 * time.Second}
	go pool.updateNewsPool()
	return &pool
}

// SubscribeRss subscribe to a new Rss feed according to its url
func (p *RssNewsPool) SubscribeRss(url string, weight SourceWeight) error {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return err
	}
	p.RLock()
	defer p.RUnlock()
	p.Pool = append(p.Pool, &RssNews{SourceTile: feed.Title, Source: url, SourceWeight: weight, UpdateTime: time.Now(), Feed: feed})
	return nil
}

// updateNewsPool the routine to update all news
func (p *RssNewsPool) updateNewsPool() {
	delay := time.After(5 * time.Second)
	for {
		select {
		case <-delay:
			fmt.Println("TimeUp:", time.Now().String())
			var wg sync.WaitGroup
			for idx, feed := range p.Pool {
				wg.Add(1)
				go feed.UpdateWithTimeout(idx, &wg, feedUpdateTimeout)
				fmt.Println("idx:", idx)
			}
			wg.Wait()
			delay = time.After(p.UpdateInterval) // periodical process
		}
	}
}
