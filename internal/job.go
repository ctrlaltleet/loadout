package internal

import (
	"fmt"
	"sync"
)

type Job struct {
	Name    string
	Version string
	URL     string
	Hash    string
	OutPath string
}

type Result struct {
	Name string
	Err  error
}

func Worker(wg *sync.WaitGroup, jobs <-chan Job, results chan<- Result) {
	defer wg.Done()
	for job := range jobs {
		fmt.Printf("[*] fetching package %s (%s)\n", job.Name, formatVersion(job.Version))
		err := FetchAsset(job.URL, job.Hash, job.OutPath, job.Name)
		results <- Result{Name: job.Name, Err: err}
	}
}

func formatVersion(v string) string {
	if v == "" {
		return "version unset"
	}
	return v
}