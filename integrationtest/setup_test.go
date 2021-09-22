// +build integration

package integrationtest

import (
	"os"
	"testing"

	"github.com/pmadhvi/telness-manager/handlers"
	"github.com/pmadhvi/telness-manager/mock"
	"github.com/pmadhvi/telness-manager/service"
	"github.com/sirupsen/logrus"
)

var server handlers.Server

func setup() {
	var (
		log              = logrus.New()
		subscriptionRepo = &mock.DbMock{}
		client           = &mock.ClientMock{}
		subsvc           = service.SubscriptionSvc{Log: log, SubscriptionRepo: subscriptionRepo, PtsClient: client}
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
