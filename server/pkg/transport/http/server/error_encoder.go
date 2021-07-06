package server

func ServerErrorEncoder(errorEncoder ErrorEncoder) Option {
	return func(serv *Server) { serv.errorEncoder = errorEncoder }
}
