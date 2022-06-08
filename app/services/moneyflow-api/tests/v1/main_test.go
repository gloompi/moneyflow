package tests

import (
	"fmt"
	"testing"

	"github.com/gloompi/ultimate-service/business/data/tests"
	"github.com/gloompi/ultimate-service/foundation/docker"
)

var c *docker.Container

func TestMain(m *testing.M) {
	var err error
	c, err = tests.StartDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tests.StopDB(c)

	m.Run()
}
