package shutdowner

import (
	"context"
	"sync"
	"time"
)

type AfterFunc func(name string, err error)

type Shutdowner struct {
	wg        sync.WaitGroup
	afterFunc AfterFunc
	procs     []Process
}

func NewShutdowner(afterFunc AfterFunc, procs ...Process) *Shutdowner {
	return &Shutdowner{
		wg:        sync.WaitGroup{},
		afterFunc: afterFunc,
		procs:     procs,
	}
}

func (sd *Shutdowner) Wait(ctx context.Context, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for _, process := range sd.procs {
		sd.wg.Add(1)

		go sd.terminate(ctx, process)
	}

	sd.wg.Wait()
}

func (sd *Shutdowner) terminate(ctx context.Context, process Process) {
	defer sd.wg.Done()

	complete := make(chan error)
	go func() {
		complete <- process.cancelFunc(ctx)
	}()

	select {
	case err := <-complete:
		sd.afterFunc(process.name, err)
	case <-ctx.Done():
		sd.afterFunc(process.name, ctx.Err())
	}
}
