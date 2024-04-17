package grpcerror

import "artk.dev/core/apperror"

func CreateFakeDependency(err error) int {
	return int(apperror.KindOf(err))
}
