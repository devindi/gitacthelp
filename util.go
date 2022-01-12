package main

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func requireEnvironmentVariable(key string) string {
	value, present := os.LookupEnv(key)
	if !present {
		log.Fatal("missed variable ", key)
	}
	return value
}

func getEnvironmentVariable(key string) string {
	value, present := os.LookupEnv(key)
	if !present {
		log.Warnf("skipped variable %s", key)
	}
	return value
}
