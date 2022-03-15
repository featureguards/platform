package error_codes

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

var mapping = map[codes.Code]int32{
	codes.InvalidArgument:    http.StatusBadRequest,
	codes.Unauthenticated:    http.StatusUnauthorized,
	codes.PermissionDenied:   http.StatusForbidden,
	codes.NotFound:           http.StatusNotFound,
	codes.Unavailable:        http.StatusServiceUnavailable,
	codes.DeadlineExceeded:   http.StatusGatewayTimeout,
	codes.FailedPrecondition: http.StatusPreconditionFailed,
	codes.OK:                 http.StatusOK,
	codes.ResourceExhausted:  http.StatusTooManyRequests,
}

var reverseMapping = make(map[int32]codes.Code)

func GrpcCode(httpCode int32) codes.Code {
	code, ok := reverseMapping[httpCode]
	if ok {
		return code
	}
	return codes.Internal
}

func HttpCode(code codes.Code) int32 {
	status, ok := mapping[code]
	if ok {
		return status
	}
	return http.StatusInternalServerError
}

func init() {
	for k, v := range mapping {
		reverseMapping[v] = k
	}
}
