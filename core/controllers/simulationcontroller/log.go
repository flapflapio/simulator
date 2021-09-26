package simulationcontroller

import (
	"log"
)

var logger = log.Default()

func println(v ...interface{}) {
	logControllerPrefix()
	logger.Println(v...)
}

func printf(format string, v ...interface{}) {
	logControllerPrefix()
	logger.Printf(format, v...)
}

func logControllerPrefix() {
	log.Print("SIMULATION CONTROLLER: ")
}
