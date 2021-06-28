package channel_go

type Worker struct {
	Work chan Job
}

func NewWorker() Worker {
	return Worker{Work: make(chan Job)}
}

//i 为传入的线程名称，chorme 需要不同的端口，用此来标记不同的端口
func (w Worker) Run(manager chan chan Job, re chan interface{}, i int) {
	go func() {
		for {
			manager <- w.Work
			select {
			case work := <-w.Work:
				result := work.Do(i)
				re <- result
			}
		}
	}()
}

