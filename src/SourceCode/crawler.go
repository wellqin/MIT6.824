package SourceCode

import (
	"fmt"
	"sync"
)

//
// Several solutions to the crawler exercise from the Go tutorial
// https://tour.golang.org/concurrency/10
//

//
// Serial crawler
//

func Into() { // 首字母大写才可以被其他函数调用
	fmt.Printf("=== into ===\n")
}

/*
一个serial的代码例子，它用一个map来存储，爬取过的url就是true，没fetch过的就先设置url对应true并fetch它
但是很明显这不是parallel的，跑完发现结果是串行返回的

加go变成协程后，发现先爬取部分，然后后续依旧返回，但这不是真正的并发
*/
func Serial(url string, fetcher Fetcher, fetched map[string]bool) {
	if fetched[url] { //已经爬取的忽略返回
		return
	}
	fetched[url] = true
	urls, err := fetcher.Fetch(url)
	if err != nil {
		return
	}
	for _, u := range urls {
		go Serial(u, fetcher, fetched) // 前面加go，变为协程
	}
	//fmt.Printf("= = = into  Serial = = =\n")
	return
}

//
// Concurrent crawler with shared state and Mutex
//

type fetchState struct {
	mu      sync.Mutex
	fetched map[string]bool
}

/*
ConcurrentMutex Crawler并发带互斥锁的爬虫：带有共享变量
加上lock() unlock()这样就可以保证每个thread在读写map的时候都是互斥的
waitGroup的使用可以保证等待所有条件满足再结束（有点unix里面的system call的意味，类似wait和signal等等，当然c里面的线程也不一样比如pthread等等）

go detector可以帮助检查出race脏数据的问题 而go会raise an error当同时read write一个变量

Thread Pool 线程池
threads成百上千还ok，但是millions of threads就不大可能了
所以更好的方式是搞一个fixed size pool of workser，让worker一个一个fetch url而不是给每个Url创建一个线程
*/

func ConcurrentMutex(url string, fetcher Fetcher, f *fetchState) {
	f.mu.Lock()
	already := f.fetched[url] // 保持这二行原子性
	f.fetched[url] = true
	f.mu.Unlock()

	if already {
		return
	}

	urls, err := fetcher.Fetch(url)
	if err != nil {
		return
	}
	var done sync.WaitGroup
	for _, u := range urls {
		done.Add(1)
		u2 := u
		go func() { // 给每个Url创建一个线程
			defer done.Done()
			ConcurrentMutex(u2, fetcher, f)
		}()
		//go func(u string) {
		//	defer done.Done()
		//	ConcurrentMutex(u, fetcher, f)
		//}(u)
	}
	done.Wait()
	return
}

func MakeState() *fetchState {
	f := &fetchState{}
	f.fetched = make(map[string]bool)
	return f
}

//
// Concurrent crawler with channels
//

func worker(url string, ch chan []string, fetcher Fetcher) {
	urls, err := fetcher.Fetch(url)
	if err != nil {
		ch <- []string{}
	} else {
		ch <- urls
	}
}

func master(ch chan []string, fetcher Fetcher) {
	n := 1
	fetched := make(map[string]bool)
	for urls := range ch {
		for _, u := range urls {
			if fetched[u] == false {
				fetched[u] = true
				n += 1
				go worker(u, ch, fetcher)
			}
		}
		n -= 1
		if n == 0 {
			break
		}
	}
}

/*
ConcurrentChannel方式的爬虫
Using channels instead of shared memory
没有map，也没用shared memory，也没有locks
这里调用master，然后worker与master通过ch这个channel来通信
*/

func ConcurrentChannel(url string, fetcher Fetcher) {
	ch := make(chan []string) // make(chan Type) //等价于make(chan Type, 0) 当 capacity= 0 时，channel 是无缓冲阻塞读写的
	go func() {
		ch <- []string{url} // channel通过操作符<-来接收和发送数据，位于左边接受右边发送
	}()
	master(ch, fetcher)
}

//
// Fetcher 接口只有方法声明，没有实现，没有数据字段，由别的（自定义）类型实现
// 接口可以匿名嵌入其它接口，或嵌入到结构中
// 接口就是一组抽象方法的集合，它必须由其他非接口类型实现，而不能自我实现
//

type Fetcher interface { //接⼝命名习惯以 er 结尾
	// Fetch returns a slice of URLs found on the page.
	Fetch(url string) (urls []string, err error)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct { // 结构体
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) ([]string, error) {
	if res, ok := f[url]; ok {
		fmt.Printf("found:   %s\n", url)
		return res.urls, nil
	}
	fmt.Printf("missing: %s\n", url)
	return nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
//var fetcher = fakeFetcher{
//	"http://golang.org/": &fakeResult{
//		"The Go Programming Language",
//		[]string{
//			"http://golang.org/pkg/",
//			"http://golang.org/cmd/",
//		},
//	},
//	"http://golang.org/pkg/": &fakeResult{
//		"Packages",
//		[]string{
//			"http://golang.org/",
//			"http://golang.org/cmd/",
//			"http://golang.org/pkg/fmt/",
//			"http://golang.org/pkg/os/",
//		},
//	},
//	"http://golang.org/pkg/fmt/": &fakeResult{
//		"Package fmt",
//		[]string{
//			"http://golang.org/",
//			"http://golang.org/pkg/",
//		},
//	},
//	"http://golang.org/pkg/os/": &fakeResult{
//		"Package os",
//		[]string{
//			"http://golang.org/",
//			"http://golang.org/pkg/",
//		},
//	},
//}
