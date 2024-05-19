package closer

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/relby/diva.back/internal/logger"
)

var globalCloser = New()

func Add(callbacks ...func() error) {
	globalCloser.Add(callbacks...)
}

func Wait() {
	globalCloser.Wait()
}

func CloseAll() {
	globalCloser.CloseAll()
}

type Closer struct {
	mu        sync.Mutex
	once      sync.Once
	done      chan struct{}
	callbacks []func() error
}

func New() *Closer {
	closer := &Closer{
		done: make(chan struct{}),
	}
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		signal.Stop(ch)
		closer.CloseAll()
	}()

	return closer
}

func (closer *Closer) Add(callbacks ...func() error) {
	closer.mu.Lock()
	defer closer.mu.Unlock()

	closer.callbacks = append(closer.callbacks, callbacks...)
}

func (closer *Closer) Wait() {
	<-closer.done
}

func (closer *Closer) CloseAll() {
	closer.once.Do(func() {
		defer close(closer.done)

		closer.mu.Lock()
		defer closer.mu.Unlock()

		var wg sync.WaitGroup
		errs := make(chan error, len(closer.callbacks))
		for _, callback := range closer.callbacks {
			wg.Add(1)
			go func(callback func() error) {
				defer wg.Done()
				err := callback()
				if err != nil {
					errs <- err
				}
			}(callback)
		}
		go func() {
			wg.Wait()
			close(errs)
		}()

		for err := range errs {
			logger.Err().Printf("could not close resource: %v", err)
		}
	})
}
