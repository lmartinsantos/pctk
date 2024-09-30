package pctk

import (
	"errors"
	"time"
)

var (
	// PromiseBroken is an error that indicates that the promise is broken.
	PromiseBroken = errors.New("broken promise")
)

// Future is a value that will be available in the future.
type Future interface {
	// Wait waits for the future to be completed.
	Wait() any

	// IsCompleted returns true if the future is completed.
	IsCompleted() bool

	// AndThen returns a future that is chained to this future.
	AndThen(func(any) Future) Future
}

// Promise is an instant when some event will be produced.
type Promise struct {
	done   chan struct{}
	result any
}

// NewPromise creates a new future.
func NewPromise() *Promise {
	done := make(chan struct{})
	return &Promise{done: done}
}

// Complete completes the future.
func (f *Promise) Complete() {
	close(f.done)
}

// CompleteWithValue completes the future with a value.
func (f *Promise) CompleteWithValue(v any) {
	f.result = v
	close(f.done)
}

// Break breaks the promise. This will complete the future with a PromiseBroken error as value.
func (f *Promise) Break() {
	f.CompleteWithValue(PromiseBroken)
}

// CompleteAfter completes the future after the given duration.
func (f *Promise) CompleteAfter(v any, d time.Duration) {
	if d == 0 {
		f.CompleteWithValue(v)
		return
	}
	time.AfterFunc(d, func() {
		f.CompleteWithValue(v)
	})
}

// CompleteWhen completes the future when the other future is completed.
func (f *Promise) CompleteWhen(other Future) {
	go func() {
		v := other.Wait()
		f.CompleteWithValue(v)
	}()
}

// Wait waits for the future to be completed.
func (f *Promise) Wait() any {
	<-f.done
	return f.result
}

// IsCompleted returns true if the future is completed.
func (f *Promise) IsCompleted() bool {
	select {
	case <-f.done:
		return true
	default:
		return false
	}
}

// AndThen returns a future that is chained to this future.
func (f *Promise) AndThen(g func(any) Future) Future {
	done := NewPromise()
	go func() {
		v := f.Wait()
		if v == PromiseBroken {
			done.Break()
			return
		}
		done.CompleteWhen(g(v))
	}()
	return done
}

// WithDelay returns a future that will be completed the given duration after f is completed.
func WithDelay(f Future, d time.Duration) Future {
	if d == 0 {
		return f
	}
	done := NewPromise()
	go func() {
		v := f.Wait()
		time.AfterFunc(d, func() {
			done.CompleteWithValue(v)
		})
	}()
	return done
}
