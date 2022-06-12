package query

import "context"

type Bus interface {
	Dispatch(context.Context, Query) (Result, error)
	Register(Type, Handler)
}

type Type string

type Query interface {
	Type() Type
}

type Result interface {
}

type Handler interface {
	Handle(context.Context, Query) (Result, error)
}
