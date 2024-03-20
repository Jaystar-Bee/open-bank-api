package utils

import "time"

func NowTime() string {
	return time.Now().Format(time.RFC3339)
}
