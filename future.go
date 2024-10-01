package pctk

import (
	"errors"
	"fmt"
	"time"
)

var (
	// PromiseBroken is an error that indicates that the promise is broken.
	PromiseBroken = errors.New("broken promise")
)

// Future is a value that will be available in the future.
type Future interface {
	// Wait waits for the future to be completed.
	Wait() (any, error)

	// IsCompleted returns true if the future is completed.
	IsCompleted() bool
}

// Continue continues a future with another future.
func Continue(f Future, g func(any) Future) Future {
	prom := NewPromise()
	go func() {
		v, err := f.Wait()
		if err != nil {
			prom.CompleteWithError(err)
			return
		}
		v, err = g(v).Wait()
		prom.CompleteWith(v, err)
	}()
	return prom
}

// Recover recovers from an error in a future. If the given future fails, the given function will be
// called with the error. If the function returns a future, it will be waited for and its value or
// its error will be returned.
func Recover(f Future, g func(error) Future) Future {
	prom := NewPromise()
	go func() {
		v, err := f.Wait()
		if err != nil {
			v, err = g(err).Wait()
		}
		prom.CompleteWith(v, err)
	}()
	return prom
}

// WaitAs waits for the future and returns the value as the given type.
func WaitAs[T any](f Future) (val T, err error) {
	var v any
	if v, err = f.Wait(); err != nil {
		return
	}
	var ok bool
	if val, ok = v.(T); !ok {
		err = fmt.Errorf("cannot convert %T to %T", v, val)
	}
	return
}

// Promise is an instant when some event will be produced.
type Promise struct {
	done   chan struct{}
	result any
	err    error
}

// NewPromise creates a new future.
func NewPromise() *Promise {
	done := make(chan struct{})
	return &Promise{done: done}
}

// Wait implements the Future interface.
func (f *Promise) Wait() (any, error) {
	<-f.done
	return f.result, f.err
}

// IsCompleted implements the Future interface.
func (f *Promise) IsCompleted() bool {
	select {
	case <-f.done:
		return true
	default:
		return false
	}
}

// Complete completes the future. This sets no value. The Wait function will return a zero value and
// no error.
func (f *Promise) Complete() {
	close(f.done)
}

// CompleteWith completes the future with the given value and error.
func (f *Promise) CompleteWith(v any, err error) {
	f.result = v
	f.err = err
	close(f.done)
}

// CompleteWithValue completes the future with a value.
func (f *Promise) CompleteWithValue(v any) {
	f.result = v
	close(f.done)
}

// CompleteWithError completes the future with an error.
func (f *Promise) CompleteWithError(err error) {
	f.err = err
	close(f.done)
}

// CompleteWithErrorf completes the future with an error formatted with the given format and args.
func (f *Promise) CompleteWithErrorf(format string, args ...any) {
	f.CompleteWithError(fmt.Errorf(format, args...))
}

// Bind binds the future to another future. The future will be completed with the value of the given
// future when it is completed.
func (f *Promise) Bind(other Future) {
	go func() {
		v, err := other.Wait()
		f.CompleteWith(v, err)
	}()
}

// Break breaks the promise. This will complete the future with a PromiseBroken error as value.
func (f *Promise) Break() {
	f.CompleteWithError(PromiseBroken)
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
