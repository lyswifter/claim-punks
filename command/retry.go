package command

import (
	"log"
	"time"
)

// Retry Retry
func Retry(attempts int, sleep time.Duration, fn func(string) error, str string) error {
	if err := fn(str); err != nil {
		if s, ok := err.(Stop); ok {
			return s.error
		}

		if attempts--; attempts > 0 {
			log.Printf("retry func error: %s. attemps #%d after %s.", err.Error(), attempts, sleep)
			time.Sleep(sleep)
			return Retry(attempts, sleep, fn, str)
		}
		return err
	}
	return nil
}

// Stop Stop
type Stop struct {
	error
}

// NoRetryError NoRetryError
func NoRetryError(err error) Stop {
	return Stop{err}
}
