package ingestor

var Workers chan chan interface{}

type Dispatcher struct {
}

func (d *Dispatcher) Start(count int) {
	Workers = make(chan chan interface{}, count)
	for i := 0; i < count; i++ {
		worker := NewWorker(i+1, Workers)
		worker.Start()
	}

	go func() {
		for {
			work := <-Workload
			go func() {
				worker := <-Workers
				worker <- work
			}()
		}
	}()
}
