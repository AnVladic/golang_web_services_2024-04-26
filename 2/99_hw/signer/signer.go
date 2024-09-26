package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

func SingleHash(in, out chan interface{}) {
	var wgMain sync.WaitGroup
	var mu = &sync.Mutex{}
	for inResult := range in {
		if val, ok := inResult.(int); ok {
			wgMain.Add(1)

			go func(val int) {
				defer wgMain.Done()

				valText := strconv.Itoa(val)
				var ch = make(chan string, 2)

				go func(data string, ch chan string) {
					ch <- DataSignerCrc32(data)
				}(valText, ch)
				go func(data string, ch chan string) {
					mu.Lock()
					md5 := DataSignerMd5(data)
					mu.Unlock()
					ch <- DataSignerCrc32(md5)
				}(valText, ch)
				out <- <-ch + "~" + <-ch
			}(val)
		}
	}
	wgMain.Wait()
}

func MultiHash(in, out chan interface{}) {
	var wgMain sync.WaitGroup
	for inResult := range in {
		if val, ok := inResult.(string); ok {
			wgMain.Add(1)

			go func(value string) {
				defer wgMain.Done()

				const n = 6
				var results [n]string
				var wg = &sync.WaitGroup{}
				wg.Add(n)
				for th := 0; th < n; th++ {
					go func(i int) {
						defer wg.Done()
						results[i] = DataSignerCrc32(strconv.Itoa(i) + value)
					}(th)
				}
				wg.Wait()
				out <- strings.Join(results[:], "")
			}(val)
		}
	}
	wgMain.Wait()
}

func CombineResults(in, out chan interface{}) {
	var results []string
	for inResult := range in {
		if val, ok := inResult.(string); ok {
			results = append(results, val)
		}
	}
	sort.Strings(results)
	out <- strings.Join(results, "_")
}

func ExecutePipeline(jobs ...job) {
	var in = make(chan interface{})
	var out = make(chan interface{})
	var wg = &sync.WaitGroup{}
	wg.Add(len(jobs))
	for _, jobFunc := range jobs {
		out = make(chan interface{})
		go func(in, out chan interface{}, jobFunc job) {
			defer wg.Done()
			defer close(out)
			jobFunc(in, out)
		}(in, out, jobFunc)
		in = out
	}
	wg.Wait()
}
