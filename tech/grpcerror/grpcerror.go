package grpcerror

import "artk.dev/core/apperror"

func createFakeDependency(err error) int {
	return int(apperror.KindOf(err))
}
