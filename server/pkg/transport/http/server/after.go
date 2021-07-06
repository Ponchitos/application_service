package server

func ServerAfter(after ...ResponseFunc) Option {
	return func(serv *Server) { serv.after = append(serv.after, after...) }
}
