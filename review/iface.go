package review

import "context"

type IReviewer interface {
	// llm response handler
	SetHandler(handler func(context.Context, string) error)
	// execute review
	Exec() error
}
