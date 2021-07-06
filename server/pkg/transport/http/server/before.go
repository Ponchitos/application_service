package server

func ServerBefore(before ...RequestFunc) Option {
	return func(serv *Server) { serv.before = append(serv.before, before...) }
}
