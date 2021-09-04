package runnable

type Runnable interface {
	Init() error
	Run(errC chan<- error)
	Shutdown()
}
