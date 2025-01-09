// Goroutine middlewares
package inner

import (
	"context"
	"errors"
)

type InternalMiddlewareFn func(ctx context.Context) (interface{}, error)

type InternalMiddleware func(next InternalMiddlewareFn) InternalMiddlewareFn

var (
	ErrUnableToCast = errors.New("interface cast error")
)

func InternalMiddlewareChain(mws ...InternalMiddleware) InternalMiddleware {
	return func(next InternalMiddlewareFn) InternalMiddlewareFn {
		fn := next
		for mw := len(mws) - 1; mw >= 0; mw-- {
			fn = mws[mw](fn)
		}

		return fn
	}
}

func CastResult[T any](result interface{}, err error) (T, error) {
	var typedErr error
	typedResult, is := result.(T)
	if !is {
		typedErr = ErrUnableToCast
	}
	if err == nil {
		err = typedErr
	}
	return typedResult, err
}
