package nevewebhook

import (
	"github.com/xfali/neve-core/processor"
	"github.com/xfali/neve-webhook/clients"
	"github.com/xfali/neve-webhook/servers"
	"github.com/xfali/neve-webhook/service"
	"github.com/xfali/restclient/v2"
)

func NewWebhooksServerProcessor(opts ...servers.ProcessorOpt) processor.Processor {
	return servers.NewWebhooksServerProcessor(opts...)
}

func NewWebhooksClient(endpoint string, client restclient.RestClient) service.WebHookService {
	return clients.NewWebHookClient(endpoint, client)
}
