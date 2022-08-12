// package shared provides functions for handling the key-value pairs
// in the context variables being used
package shared

import "context"

type requestID string
type userID string

func WithRequestID(ctx context.Context, reqId string) context.Context {
	return context.WithValue(ctx, requestID("requestId"), reqId)
}

func WithUserID(ctx context.Context, userId string) context.Context {
	return context.WithValue(ctx, userID("userId"), userId)
}

func ExtractRequestID(ctx context.Context) string {

	result := ""
	data := ctx.Value(requestID("requestId"))
	if data != nil {
		result = data.(string)
	}

	return result
}

func ExtractUserID(ctx context.Context) string {

	result := ""
	data := ctx.Value(userID("userId"))
	if data != nil {
		result = data.(string)
	}

	return result
}
