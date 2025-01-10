package util

import "github.com/khitrov-aleksandr/proxyguard/repository"

type RateLimiter struct {
	r repository.Repository
}

func NewRateLimiter(r repository.Repository) *RateLimiter {
	return &RateLimiter{r}
}

func (rl *RateLimiter) Check(key string, count int64, timeout int) bool {
	if rl.r.Incr(key) > count {
		return true
	}

	rl.r.Expr(key, timeout)

	return false
}
