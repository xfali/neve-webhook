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

package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/xfali/neve-webhook/clients"
	"github.com/xfali/neve-webhook/notifier"
	"github.com/xfali/neve-webhook/recorder"
	"github.com/xfali/neve-webhook/service"
	"github.com/xfali/restclient/v2"
	"github.com/xfali/xlog"
	"net/http"
)

type TestHandler struct {
	logger xlog.Logger
	cli    service.WebHookService
	id     string
}

func NewTestHandler() *TestHandler {
	cli := clients.NewWebHookClient("http://localhost:8080/webhooks", restclient.New())
	return &TestHandler{
		logger: xlog.GetLogger(),
		cli:    cli,
	}
}

func (o *TestHandler) BeanAfterSet() error {
	id, err := o.cli.Create(context.Background(), recorder.Data{
		Url:    "http://localhost:8081/events",
		Secret: "just-test",
		TriggerEventTypes: []string{
			"push",
		},
	})
	if err != nil {
		return err
	}
	o.id = id
	return nil
}

func (o *TestHandler) BeanDestroy() error {
	if o.id != "" {
		return o.cli.Delete(context.Background(), o.id)
	}
	return nil
}

func (o *TestHandler) HttpRoutes(engin gin.IRouter) {
	engin.POST("/events", o.events)
}

func (o *TestHandler) events(ctx *gin.Context) {
	v := clients.Result[string]{}
	err := ctx.Bind(&v)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	t := ctx.GetHeader(notifier.EventTypeHeader)
	s := ctx.GetHeader(notifier.EventSignatureHeader)
	o.logger.Infof("EventType: %s Event Sign: %s\n", t, s)
	o.logger.Infoln(v)
}