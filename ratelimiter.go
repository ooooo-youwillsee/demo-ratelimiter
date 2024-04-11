package main

import (
	"time"
)

type RateLimiter interface {
	acquire() (bool, error)
}

type FixedWindowRateLimiter struct {
	// 固定窗口大小, 单位ms
	windowInterval time.Duration
	// 限制
	limit int
	// 窗口开始时间
	prevTime time.Time
	// 当前限制
	curLimit int
}

func NewFixedRateLimiter(windowSize time.Duration, limit int) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		windowInterval: windowSize,
		limit:          limit,
		curLimit:       0,
		prevTime:       time.Now(),
	}
}

func (s *FixedWindowRateLimiter) acquire() (bool, error) {
	// 不在一个时间窗口，重置
	if time.Until(s.prevTime) > s.windowInterval {
		s.curLimit = 0
	}
	s.curLimit++
	s.prevTime = time.Now()
	return s.curLimit < s.limit, nil
}

type TokenBucketRateLimiter struct {
	// 桶大小
	bucketSize int
	// 速率，单位ms
	rate int
	// 剩余令牌数
	remainTokens int
	// 时间
	prevTime time.Time
}

func NewTokenBucketRateLimiter(bucketSize int, rate int) *TokenBucketRateLimiter {
	return &TokenBucketRateLimiter{
		bucketSize:   bucketSize,
		rate:         rate,
		remainTokens: bucketSize,
		prevTime:     time.Now(),
	}
}

func (t *TokenBucketRateLimiter) acquire() (bool, error) {
	// 计算新的令牌数
	newTokens := int(time.Until(t.prevTime).Milliseconds()) * t.rate
	t.remainTokens += newTokens
	if t.remainTokens >= t.bucketSize {
		t.remainTokens = t.bucketSize
	}
	t.remainTokens--
	t.prevTime = time.Now()
	return t.remainTokens > 0, nil
}

type LeakyBucketRateLimiter struct {
	// 速率，单位ms
	rate int
	// 剩余令牌数
	remainTokens int
	// 时间
	prevTime time.Time
}

func NewLeakyBucketRateLimiter(rate int) *LeakyBucketRateLimiter {
	return &LeakyBucketRateLimiter{
		rate:         rate,
		remainTokens: rate,
		prevTime:     time.Now(),
	}
}

func (l *LeakyBucketRateLimiter) acquire() (bool, error) {
	// 不是同一毫秒，重置令牌数
	if time.Now().UnixMilli() != l.prevTime.UnixMilli() {
		l.remainTokens = l.rate
	}
	l.remainTokens--
	l.prevTime = time.Now()
	return l.remainTokens > 0, nil
}
