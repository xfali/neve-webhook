/*
 * Copyright (C) 2024, Xiongfa Li.
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package clients

import (
	"github.com/gin-gonic/gin"
	"github.com/xfali/neve-web/gineve/midware/loghttp"
	"github.com/xfali/neve-webhook/notifier"
	"github.com/xfali/xlog"
	"net/http"
)

type SignatureVerifier interface {
	VerifySignature(signature string) (httpStatus int, err error)
}

type EventProcessor interface {
	ProcessWebhookEvent(eventType string, payload []byte) error
}

type webHookHandler struct {
	logger xlog.Logger

	HLog loghttp.HttpLogger `inject:""`

	SignatureVerifier SignatureVerifier `inject:""`

	EventProcessor EventProcessor `inject:""`

	EventsPath string `fig:"neve.web.hooks.routes.events"`
}

func NewWebHookHandler() *webHookHandler {
	return &webHookHandler{
		logger: xlog.GetLogger(),
	}
}

func (o *webHookHandler) HttpRoutes(engine gin.IRouter) {
	if o.EventsPath == "" {
		o.EventsPath = "/events"
	}
	engine.POST(o.EventsPath, o.HLog.LogHttp(), o.events)
}

func (o *webHookHandler) events(ctx *gin.Context) {
	payload, err := ctx.GetRawData()
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	t := ctx.GetHeader(notifier.EventTypeHeader)
	s := ctx.GetHeader(notifier.EventSignatureHeader)
	o.logger.Debugf("Event Type: %s, Event Sign: %s, Event Payload: %s\n", t, s, string(payload))
	code, err := o.SignatureVerifier.VerifySignature(s)
	if err != nil {
		_ = ctx.AbortWithError(code, err)
		return
	}

	err = o.EventProcessor.ProcessWebhookEvent(t, payload)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
}
