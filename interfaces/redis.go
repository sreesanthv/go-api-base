package interfaces

import "time"

type Redis interface {
	Set(key string, value string, expiry time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
}
