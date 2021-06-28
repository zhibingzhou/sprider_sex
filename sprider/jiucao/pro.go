package jiucao

type Jiu_Request struct {
	Url      string
	Type     int
	PareFunc func(string, int) Jiu_PareResult
}

type Jiu_PareResult struct {
	Requests []Jiu_Request
}

func (r Jiu_Request) Do(port int) interface{} {
	pareResult := r.PareFunc(r.Url, port)
	return pareResult
}
