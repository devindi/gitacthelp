package gitacthelp

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func RequireEnvironmentVariable(key string) string {
	value, present := os.LookupEnv(key)
	if !present {
		log.Fatal("missed variable ", key)
	}
	return value
}

func GetEnvironmentVariable(key string) string {
	value, present := os.LookupEnv(key)
	if !present {
		log.Warnf("skipped variable %s", key)
	}
	return value
}
