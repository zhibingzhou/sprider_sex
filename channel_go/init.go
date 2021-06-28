package channel_go

import "fmt"

type Job interface {
	Do(int) interface{}
}

type ManagerChannel struct {
	Count    int
	Request  chan Job
	Requests chan chan Job
	Result   chan interface{}
	EndJob   chan int
}

var Manager *ManagerChannel

func NewManagerChannel(count int) *ManagerChannel {
	return &ManagerChannel{
		Count:    count,
		Request:  make(chan Job),
		Requests: make(chan chan Job),
		Result:   make(chan interface{}),
		EndJob:   make(chan int),
	}
}

func (x *ManagerChannel) Run() {
	for i := 0; i < x.Count; i++ {
		NewWorker().Run(x.Requests, x.Result, i)
	}
	var ArryReq []Job
	var ArryReqs []chan Job

	go func() {

		for {
			var requestch chan Job
			var request Job

			if len(ArryReq) > 0 && len(ArryReqs) > 0 {
				requestch = ArryReqs[0]
				request = ArryReq[0]
				fmt.Println("通道", len(ArryReqs))
				fmt.Println("结果", len(ArryReq))
			}

			select {
			case Reqs := <-x.Requests:
				ArryReqs = append(ArryReqs, Reqs)

			case Req := <-x.Request:
				ArryReq = append(ArryReq, Req)

			case requestch <- request:
				ArryReq = ArryReq[1:]
				ArryReqs = ArryReqs[1:]

				if len(ArryReq) == 1 {
					fmt.Println("任务结束")
					x.EndJob <- 1
				}

			}

		}

	}()
}




