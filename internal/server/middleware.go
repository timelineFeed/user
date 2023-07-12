package server

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware/selector"
)

var whiteApi = []string{"/user/register"}

func newWhiteListMatcher() selector.MatchFunc {

	whiteList := make(map[string]struct{})
	for _, v := range whiteApi {
		whiteList[v] = struct{}{}
	}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}
