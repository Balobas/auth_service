package shutdown

import (
	"context"
	"log"
	"sync"
)

var cl closer

type closer struct {
	mx    sync.Mutex
	once  sync.Once
	funcs []CloserFunc
}

type CloserFunc func(ctx context.Context) error

func Add(f ...CloserFunc) {
	cl.mx.Lock()
	cl.funcs = append(cl.funcs, f...)
	cl.mx.Unlock()
}

func CloseAll(ctx context.Context) {
	cl.mx.Lock()
	defer cl.mx.Unlock()

	cl.once.Do(func() {
		errs := make(chan error, len(cl.funcs))

		for _, f := range cl.funcs {
			go func(f CloserFunc) {
				errs <- f(ctx)
			}(f)
		}

		done := make(chan struct{}, 2)
		defer close(done)

		go func(errs chan error, done chan struct{}) {
			for i := 0; i < cap(errs); i++ {
				if err := <-errs; err != nil {
					log.Printf("failed to close: %v", err)
				}
			}
			done <- struct{}{}
		}(errs, done)

		select {
		case <-done:
			log.Printf("shutdown successfuly finished")
		case <-ctx.Done():
			log.Printf("failed to finish shutdown: %v", ctx.Err())
			<-done
			return
		}
	})
}

func WrapClose(f func(ctx context.Context)) CloserFunc {
	return func(ctx context.Context) error {
		f(ctx)
		return nil
	}
}
