package filter

import (
	"github.com/khitrov-aleksandr/proxyguard/filter/base"
)

type PhoneRateLimiter struct {
	rl *base.RateLimiter
}

func NewPhonesRateLimiter(rl *base.RateLimiter) *PhoneRateLimiter {
	return &PhoneRateLimiter{
		rl: rl,
	}
}

func (prl *PhoneRateLimiter) Allow(key string, count int64, timeout int) bool {
	return prl.rl.Allow(prl.getKey(key), count, timeout)
}

func (prl *PhoneRateLimiter) getKey(key string) string {
	return "rate_limiter:phone:" + key
}
