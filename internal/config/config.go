package config

import "time"

type Config struct {
	Service
}

type Service struct {
	AccessTTL  time.Duration
	RefreshTTL time.Duration
	Secret     string
	Cost       int
}
