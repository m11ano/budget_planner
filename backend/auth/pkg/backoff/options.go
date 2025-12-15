package backoff

import "time"

type controllerOptions struct {
	ClearInterval time.Duration
}

func defaultControllerOptions() *controllerOptions {
	return &controllerOptions{
		ClearInterval: 1 * time.Minute,
	}
}

func WithControllerClearInterval(d time.Duration) ControllerOption {
	return func(o *controllerOptions) {
		o.ClearInterval = d
	}
}

type sessionOptions struct {
	Ttl             time.Duration
	InitialInterval time.Duration
	Multiplier      float64
	MaxInterval     time.Duration
}

func defaultSessionOptions() *sessionOptions {
	return &sessionOptions{
		Ttl:             1 * time.Hour,
		InitialInterval: 1 * time.Second,
		Multiplier:      2,
		MaxInterval:     60 * time.Second,
	}
}

func WithTtl(d time.Duration) SessionOption {
	return func(o *sessionOptions) {
		o.Ttl = d
	}
}

func WithInitialInterval(d time.Duration) SessionOption {
	return func(o *sessionOptions) {
		o.InitialInterval = d
	}
}

func WithMultiplier(m float64) SessionOption {
	return func(o *sessionOptions) {
		o.Multiplier = m
	}
}

func WithMaxInterval(d time.Duration) SessionOption {
	return func(o *sessionOptions) {
		o.MaxInterval = d
	}
}
