package backend

type Worker struct {
	ID    int
	Work  chan *RepoData
	Queue chan chan *RepoData
	Repos map[int]*ArchRepo
	Quit  chan bool
}

func NewWorker(id int, queue chan chan *RepoData) Worker {
	return Worker{
		ID:    id,
		Work:  make(chan *RepoData),
		Queue: queue,
		Repos: make(map[int]*ArchRepo),
		Quit:  make(chan bool),
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			w.Queue <- w.Work
			select {
			case repodata := <-w.Work:
				if len(repodata.Open) != 0 {
					w.Repos[repodata.RepoID].Hive.Blender.Models[0].Conflator.SetIssueRequests(repodata.Open)
				}
				if len(repodata.Closed) != 0 {
					w.Repos[repodata.RepoID].Hive.Blender.Models[0].Conflator.SetIssueRequests(repodata.Closed)
				}
				if len(repodata.Pulls) != 0 {
					w.Repos[repodata.RepoID].Hive.Blender.Models[0].Conflator.SetPullRequests(repodata.Pulls)
				}
				w.Repos[repodata.RepoID].Hive.Blender.Models[0].Conflator.Conflate()
				// TODO: Call Learn/Predict.
                // TODO: Add in Assignment call.
			case <-w.Quit:
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	go func() {
		w.Quit <- true
	}()
}
