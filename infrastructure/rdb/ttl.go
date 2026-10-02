package rdb

import "time"

func ScaleTTL(ttl time.Duration, divisor int) time.Duration {
	if ttl <= 0 || divisor <= 1 {
		return ttl
	}

	scaled := ttl / time.Duration(divisor)

	switch {
	case scaled >= time.Hour:
		return scaled.Round(time.Minute)
	case scaled >= time.Minute:
		return scaled.Round(time.Second)
	default:
		return scaled.Round(100 * time.Millisecond)
	}
}
