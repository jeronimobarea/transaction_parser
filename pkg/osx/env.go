package osx

import "os"

func GetEnvFallback(key, fallback string) string {
	if res, ok := os.LookupEnv(key); ok {
		return res
	}
	return fallback
}
