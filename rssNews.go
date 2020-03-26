package main

import (
	"errors"
	"fmt"
	"github.com/mmcdole/gofeed"
	"sync"
	"time"
)

// RssNews news from Rss Feeds
type RssNews struct {
	SourceTile string
	Source     string
	UpdateTime time.Time
	SourceWeight
	Feed *gofeed.Feed
}

// UpdateWithTimeout Update a news feed and exit if timeout
func (r *RssNews) UpdateWithTimeout(idx int, wg *sync.WaitGroup, timeout time.Duration) error {
	ch := make(chan error)
	defer wg.Done()
	go func() {
		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(r.Source)
		if err != nil {
			ch <- err
		}
		r.Feed = feed
		ch <- nil
	}()
	select {
	case err := <-ch:
		fmt.Println("Read Feed Index", idx)
		if err != nil {
			fmt.Println("Read feed Error:", err)
		}
		return err
	case <-time.After(timeout):
		fmt.Println("update timeout")
		return errors.New("update timeout")
	}
}
