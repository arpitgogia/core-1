package backend

var Workers chan chan *RepoData

type Dispatcher struct{}

func NewDispatcher() Dispatcher {
	disp := Dispatcher{}
	return disp
}

func (d *Dispatcher) Start(count int) {
	Workers = make(chan chan *RepoData, count)
	for i := 0; i < count; i++ {
		worker := NewWorker(i+1, Workers)
		worker.Start()
	}

	go func() {
		for {
			work := <-Workload
			go func() {
				workers := <-Workers
				workers <- work
			}()
		}
	}()
}
