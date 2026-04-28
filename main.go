package main

import (
	"fmt"
	"sandbox/internal/mahjong"
	"sort"
	"strings"
	"sync"
	"time"
)

const exCount = 100_000_000
const procCount = 8

func main() {

	var wg sync.WaitGroup
	sem := make(chan struct{}, procCount)
	t := time.Now()
	var mu sync.Mutex

	results := make(map[string]int)

	for i := range exCount {
		sem <- struct{}{}
		wg.Go(func() {
			defer func(i int) { <-sem }(i)
			game := mahjong.NewGame(4)
			melds := mahjong.FindMelds(game.Players[0].Hand)
			var res []string
			for _, meld := range melds {
				if meld.KindString() != "" {
					res = append(res, meld.KindString())
				}
			}
			mu.Lock()
			results[strings.Join(res, ", ")]++
			mu.Unlock()
		})
	}

	wg.Wait()

	keys := make([]string, 0, len(results))
	for k := range results {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("Result: %s, count: %d\n", k, results[k])
	}

	fmt.Printf("Total executions: %d \n", exCount)
	fmt.Printf("Duration: %s \n", time.Since(t))
}
