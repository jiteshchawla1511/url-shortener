package models

import "time"

type Response struct {
	URL            string        `json:"url"`
	CustomShort    string        `json:"custom_short"`
	Expiry         time.Duration `json:"expiry"`
	RateRemaining  int           `json:"rate_reamining"`
	RateLimitReset time.Duration `json:"rate_limit_reset"`
}
