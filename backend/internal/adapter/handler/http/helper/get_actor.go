package httpx

import (
	"context"
	"errors"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

type ctxKey int

const actorKey ctxKey = iota

func ActorFromContext(ctx context.Context) (user.Actor, error) {
	v := ctx.Value(actorKey)
	if v == nil {
		return user.Actor{}, errors.New("missing actor in context")
	}
	a, ok := v.(user.Actor)
	if !ok {
		return user.Actor{}, errors.New("invalid actor type in context")
	}
	return a, nil
}