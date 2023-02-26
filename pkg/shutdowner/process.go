package shutdowner

import "context"

type CancelFunc func(ctx context.Context) error

type Process struct {
	name       string
	cancelFunc CancelFunc
}

func WithProcess(name string, cancelFunc CancelFunc) Process {
	return Process{
		name:       name,
		cancelFunc: cancelFunc,
	}
}
