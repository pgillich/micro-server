// Goroutine middlewares
package inner

import (
	"context"
	"errors"
)

type InternalMiddlewareFn func(ctx context.Context) (interface{}, error)

type InternalMiddlewareFnTyped[T any] func(ctx context.Context) (T, error)

type InternalMiddleware func(next InternalMiddlewareFn) InternalMiddlewareFn

type InternalMiddlewareTyped[T any] func(next InternalMiddlewareFn) InternalMiddlewareFnTyped[T]

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

func InternalMiddlewareChainTyped[T any](mws ...InternalMiddleware) InternalMiddlewareTyped[T] {
	return func(next InternalMiddlewareFn) InternalMiddlewareFnTyped[T] {
		fn := next
		for mw := len(mws) - 1; mw >= 0; mw-- {
			fn = mws[mw](fn)
		}

		return func(ctx context.Context) (T, error) {
			result, err := fn(ctx)

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
	}
}
