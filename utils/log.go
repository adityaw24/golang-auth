package utils

import "log"

func LogError(layer string, message string, err error) {
	log.Printf("[ERROR] [%s] : %s => %v\n", layer, message, err)
}
