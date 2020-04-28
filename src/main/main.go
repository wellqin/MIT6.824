package main

// 不在go module下，就利用go path，import src下面的目录/文件，比如“SourceCode/crawler”
// 在go module下 你源码中 import …/ 这样的引入形式不支持了， 应该改成 import 模块名/ 。 这样就ok了
// go module前删除go path，并且开启module环境变量+设置代理变量，同时在pycharm中也进行设置，最后才完成
// import mod文件中的module，并加上文件在src下的目录
import (
	"fmt"
	crawler "github.com/wellqin/MIT6.824/src/SourceCode"
)

//
// Fetcher
//

type Fetcher interface {
	// Fetch returns a slice of URLs found on the page.
	Fetch(url string) (urls []string, err error)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
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
var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}

//
// main
//

func main() {
	fmt.Printf("=== Serial===\n")
	crawler.Serial("http://golang.org/", fetcher, make(map[string]bool))
	//crawler.Into()

	fmt.Printf("=== ConcurrentMutex ===\n")
	crawler.ConcurrentMutex("http://golang.org/", fetcher, crawler.MakeState())

	fmt.Printf("=== ConcurrentChannel ===\n")
	crawler.ConcurrentChannel("http://golang.org/", fetcher)
}
