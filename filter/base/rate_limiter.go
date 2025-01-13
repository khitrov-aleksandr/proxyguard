package base

import (
	"github.com/khitrov-aleksandr/proxyguard/repository"
)

type RateLimiter struct {
	r repository.Repository
}

func NewRateLimiter(r repository.Repository) *RateLimiter {
	return &RateLimiter{r}
}

func (rl *RateLimiter) Allow(key string, count int64, timeout int) bool {
	res := rl.r.Incr(key) <= count
	rl.r.Expr(key, timeout)

	return res
}
