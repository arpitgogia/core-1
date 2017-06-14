package backend

import "time"

type BackendServer struct {
	Database MemSQL
	Repos    map[int]*ArchRepo
}

func (bs *BackendServer) OpenSQL() {
	bs.Database.Open()
}

func (bs *BackendServer) CloseSQL() {
	bs.Database.Close()
}

func (bs *BackendServer) Timer() {
	ticker := time.NewTicker(time.Millisecond * 500)
	go func() {
		for range ticker.C {

			data, err := bs.Database.Read()
			if err != nil {
				panic(err) // TODO: Implement proper error handling.
			}

			disp := NewDispatcher()
			disp.Start(10)

			// Collector(data)
			// TODO: Implement the rest of the logic here.
		}
	}()
}
