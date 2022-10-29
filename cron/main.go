/*
Thanks to Stephen Afamo for the helper functions:
https://stephenafamo.com/blog/posts/better-scheduling-in-go
https://github.com/stephenafamo/kronika/blob/master/kronika.go

*/

package cron

import (
	"context"
	"time"
)

// Every sends the time to the returned channel at the specified intervals
func Every(ctx context.Context, start time.Time, interval time.Duration) <-chan time.Time {
	// Create the channel which we will return
	stream := make(chan time.Time, 3)

	// Calculating the first start time in the future
	// Need to check if the time is zero (e.g. if time.Time{} was used)
	if !start.IsZero() {
		diff := time.Until(start)
		if diff < 0 {
			total := diff - interval
			times := total / interval * -1

			start = start.Add(times * interval)
		}
	}

	// Run this in a goroutine, or our function will block until the first event
	go func() {
		// Run the first event after it gets to the start time
		timer := time.NewTimer(time.Until(start))
		defer timer.Stop() // Make sure to stop the timer when we're done

		// Listen on both the timer and the context done channel.
		// Useful if the context is closed before the first timer
		select {
		case t := <-timer.C:
			stream <- t
		case <-ctx.Done():
			close(stream)
			return
		}

		// Open a new ticker
		ticker := time.NewTicker(interval)
		defer ticker.Stop() // Make sure to stop the ticker when we're done

		// Listen on both the ticker and the context done channel to know when to stop
		for {
			select {
			case t2 := <-ticker.C:
				stream <- t2
			case <-ctx.Done():
				close(stream)
				return
			}
		}
	}()

	return stream
}
