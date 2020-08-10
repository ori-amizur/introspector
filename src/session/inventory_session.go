package session

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/openshift/assisted-service/client"
	"github.com/openshift/assisted-service/pkg/auth"
	"github.com/openshift/assisted-service/pkg/requestid"
	"github.com/sirupsen/logrus"
)

func createUrl(inventoryUrl string) string {
	return fmt.Sprintf("%s/%s", inventoryUrl, client.DefaultBasePath)
}

type InventorySession struct {
	ctx    context.Context
	logger logrus.FieldLogger
	client *client.AssistedInstall
}

func (i *InventorySession) Context() context.Context {
	return i.ctx
}

func (i *InventorySession) Logger() logrus.FieldLogger {
	return i.logger
}

func (i *InventorySession) Client() *client.AssistedInstall {
	return i.client
}

func createBmInventoryClient(inventoryUrl string, pullSecretToken string) *client.AssistedInstall {
	clientConfig := client.Config{}
	clientConfig.URL, _ = url.Parse(createUrl(inventoryUrl))
	clientConfig.Transport = requestid.Transport(http.DefaultTransport)
	clientConfig.AuthInfo = auth.AgentAuthHeaderWriter(pullSecretToken)
	bmInventory := client.New(clientConfig)
	return bmInventory
}

func New(inventoryUrl string, pullSecretToken string) *InventorySession {
	id := requestid.NewID()
	ret := InventorySession{
		ctx:    requestid.ToContext(context.Background(), id),
		logger: requestid.RequestIDLogger(logrus.StandardLogger(), id),
		client: createBmInventoryClient(inventoryUrl, pullSecretToken),
	}
	return &ret
}
