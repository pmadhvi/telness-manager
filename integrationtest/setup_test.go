// +build integration

package integrationtest

import (
	"github.com/pmadhvi/telness-manager/handlers"
	"github.com/pmadhvi/telness-manager/mock"
	"github.com/pmadhvi/telness-manager/service"
	"github.com/sirupsen/logrus"
	"os"
	"testing"
)

var server handlers.Server

func setup() {
	var (
		log              = logrus.New()
		subscriptionRepo = &mock.DbMock{}
		subsvc           = service.SubscriptionSvc{SubscriptionRepo: subscriptionRepo, Log: log}
	)
	log.SetOutput(os.Stdout)
	server = handlers.Server{Log: log, Port: "7000", SubscriptionService: subsvc}
	go func() {
		err := server.Start()
		log.Error("Test Server could not be started:", err)
	}()
}

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}
