// +build integration

package integrationtest

import (
	"os"
	"testing"

	"github.com/pmadhvi/iban-validator/handlers"
	"github.com/sirupsen/logrus"
)

func setup() {
	var log = logrus.New()
	log.SetOutput(os.Stdout)

	server := handlers.Server{
		Log:  log,
		Port: "7000",
	}

	go func() {
		err := server.Start()
		log.Error("Test Server could not be started:", err)
	}()
}

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}
