package utils

import "os"

func IsDebug() bool {
	_, ok := os.LookupEnv("FANGS_DEV")
	return ok
}
