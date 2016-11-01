package blocks

type Request struct {
	InitTime    int
	ServiceTime int
}

func (r *Request) GetServiceTime() int {
	return r.ServiceTime
}
