package job

// A Handler responds to a Request.
type Handler interface {
	Handle(*Request)
}

// A Chainer can accept a Handler as its next
// Handler to create a chain of Handlers.
type Chainer interface {
	SetNext(Handler)
}

// ChainHandler — a Handler that can be chained
// to a "next" handler to create a chain of Handlers.
type ChainHandler interface {
	Handler
	Chainer
}
