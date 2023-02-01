package main

import (
	logging "github.com/vantoan19/Petifies/server/libs/logging-config"
)

var logger = logging.NewLogger("ABC")

func main() {
	logger.Alert("Hello")
}
