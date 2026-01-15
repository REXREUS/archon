package gemini

import (
	"context"
	"sync"

	"golang.org/x/time/rate"
)

type BatchProcessor struct {
	client      *Client
	rateLimiter *rate.Limiter
	workerCount int
}

func NewBatchProcessor(client *Client, rpm int) *BatchProcessor {
	// rpm: Requests Per Minute
	limit := rate.Limit(float64(rpm) / 60.0)

	return &BatchProcessor{
		client:      client,
		rateLimiter: rate.NewLimiter(limit, 1),
		workerCount: 5,
	}
}

func (bp *BatchProcessor) ProcessFiles(files []string, processFunc func(string) error) []error {
	jobs := make(chan string, len(files))
	results := make(chan error, len(files))

	var wg sync.WaitGroup
	for i := 0; i < bp.workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range jobs {
				err := bp.rateLimiter.Wait(context.Background())
				if err == nil {
					results <- processFunc(file)
				} else {
					results <- err
				}
			}
		}()
	}

	for _, f := range files {
		jobs <- f
	}
	close(jobs)
	wg.Wait()

	var errors []error
	close(results)
	for err := range results {
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}
