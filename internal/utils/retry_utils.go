package utils

import (
	"encoding/json"
	"math"
	"strings"
	"time"
)

type RetryConfig struct {
	MaxRetries    int
	InitialDelay  time.Duration
	MaxDelay      time.Duration
	BackoffFactor float64
}

func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:    5,
		InitialDelay:  500 * time.Millisecond,
		MaxDelay:      30 * time.Second,
		BackoffFactor: 2.0,
	}
}

func IsRetriableError(err error) bool {
	if err == nil {
		return false
	}

	if openAPIErr, ok := err.(interface{ Body() []byte }); ok {
		body := openAPIErr.Body()
		var respBody map[string]interface{}
		if jsonErr := json.Unmarshal(body, &respBody); jsonErr == nil {
			if payload, ok := respBody["payload"].(string); ok {
				return strings.Contains(strings.ToLower(payload), "system is currently being modified")
			}
		}
	}

	if httpErr, ok := err.(interface{ Code() int }); ok {
		code := httpErr.Code()
		return code >= 500 || code == 429 || code == 408
	}

	return false
}

func CalculateBackoff(attempt int, config RetryConfig) time.Duration {
	delay := config.InitialDelay * time.Duration(math.Pow(config.BackoffFactor, float64(attempt)))
	if delay > config.MaxDelay {
		delay = config.MaxDelay
	}
	return delay
}
