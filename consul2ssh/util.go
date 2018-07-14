package consul2ssh

import (
	"os"
)

func GetEnvKey(key, defaultVal string) string {
	if value, isSet := os.LookupEnv(key); isSet {
		return value
	}
	return defaultVal
}
