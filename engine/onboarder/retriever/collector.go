package retriever

import "github.com/google/go-github/github"

var Workload = make(chan github.Issue, 100)

func Collector(issues []github.Issue) {
	for _, i := range issues {
		Workload <- i
	}
}

// NOTE: This particular function will likely need some logic regarding what
// is passed into it - this will then determine which particular channel the
// objects are then passed into.
