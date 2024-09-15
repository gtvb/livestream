package repository

import (
	"log"

	"github.com/gtvb/livestream/utils"
)

var (
	container *utils.TestContainer
)

func setupDatabase() *utils.TestContainer {
	container, err := utils.NewTestContainer("ls-db-test")
	if err != nil {
		log.Panicf("Error: could not start container, reason -> %s\n", err)
	}

	err = container.SetupDatabaseWrapper()
	if err != nil {
		log.Panicf("Error: could not start database wrapper, reason -> %s\n", err)
	}

	return container
}
