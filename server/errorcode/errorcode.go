package errorcode

import "google.golang.org/grpc/status"

var (
	ErrUnknown = status.New(codes.Unkown)
)
