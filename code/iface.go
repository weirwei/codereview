package code

import "context"

type ICode interface {
	GetCode(ctx context.Context) ([]CodePatch, error)
}
