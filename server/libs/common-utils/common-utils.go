package common

import (
	"os"
)

func IsDevEnv() bool {
	return os.Getenv("SERVER_MODE") == "development"
}
