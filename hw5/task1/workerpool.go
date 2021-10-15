package main

import (
    "fmt"
    "time"
)

func worker(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        fmt.Printf("\nworker #%v started the job #%v", id, j)
        time.Sleep(time.Second / 2)
        fmt.Printf("\nworker #%v finished the job #%v", id, j)
        results<-j*2
    }
}

func main() {
	jobsCount := 5
	workersCount := 3
    jobs := make(chan int, jobsCount)
    results := make(chan int, jobsCount)

	for i := 1; i <= workersCount; i++ {
        go worker(i, jobs, results)
    }
	for i := 1; i <= jobsCount; i++ {
        jobs<-i
    }
	close(jobs)
	for i := 1; i <= jobsCount; i++ {
        fmt.Println("\ngot results:", <-results)
    }
}
