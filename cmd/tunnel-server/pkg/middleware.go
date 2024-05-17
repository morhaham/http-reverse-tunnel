package pkg

import (
	"context"
	"net"
)

type Handler func(ctx context.Context, message Message, conn net.Conn) Handler

type Middleware struct {
	handlers []Handler
}

func (mw *Middleware) Add(h Handler) {
	mw.handlers = append(mw.handlers, h)
}
