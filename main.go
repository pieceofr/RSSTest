package main

import ()
import "time"

func main() {
	pool := NewRssNewsPool()
	pool.SubscribeRss(hopkinsmedicine, Trusted)
	pool.SubscribeRss(taiwanCDC, Trusted)

	for {
		time.Sleep(120 * time.Second)
	}
}
