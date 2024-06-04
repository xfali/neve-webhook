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

package servers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xfali/neve-web/gineve/midware/loghttp"
	"github.com/xfali/neve-web/result"
	"github.com/xfali/neve-webhook/recorder"
	"github.com/xfali/neve-webhook/service"
	"github.com/xfali/xlog"
	"net/http"
)

type ResponseFunc func(ctx *gin.Context, o interface{}) (abort bool)

type webHookHandler struct {
	logger xlog.Logger

	HLog loghttp.HttpLogger `inject:""`

	Service service.WebHookService `inject:""`

	respFunc ResponseFunc
}

func NewWebHookHandler() *webHookHandler {
	return &webHookHandler{
		logger:   xlog.GetLogger(),
		respFunc: defaultResponse,
	}
}

func (o *webHookHandler) HttpRoutes(engine gin.IRouter) {
	engine.POST("/webhooks", o.create)
	engine.PUT("/webhooks/:id", o.update)
	engine.GET("/webhooks", o.get)
	engine.GET("/webhooks/:id", o.detail)
	engine.DELETE("/webhooks/:id", o.delete)
}

func (o *webHookHandler) create(ctx *gin.Context) {
	d := recorder.Data{}
	err := ctx.Bind(&d)
	if err != nil {
		if o.respFunc(ctx, err) {
			return
		}
	}
	id, err := o.Service.Create(ctx, d)
	if err != nil {
		if o.respFunc(ctx, err) {
			return
		}
	}

	_ = o.respFunc(ctx, id)
}

func (o *webHookHandler) update(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		if o.respFunc(ctx, fmt.Errorf("Path param id invalid ")) {
			return
		}
	}
	d := recorder.Data{}
	err := ctx.Bind(&d)
	if err != nil {
		if o.respFunc(ctx, err) {
			return
		}
	}
	err = o.Service.Update(ctx, id, d)
	if err != nil {
		if o.respFunc(ctx, err) {
			return
		}
	}

	_ = o.respFunc(ctx, nil)
}

func (o *webHookHandler) get(ctx *gin.Context) {
	id := ctx.Query("id")
	eventType := ctx.Query("event_type")
	v, err := o.Service.Get(ctx, recorder.QueryCondition{
		Id:        id,
		EventType: eventType,
	})
	if err != nil {
		if o.respFunc(ctx, err) {
			return
		}
	}

	_ = o.respFunc(ctx, v)
}

func (o *webHookHandler) detail(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		if o.respFunc(ctx, fmt.Errorf("Path param id invalid ")) {
			return
		}
	}
	v, err := o.Service.Detail(ctx, id)
	if err != nil {
		if o.respFunc(ctx, err) {
			return
		}
	}
	_ = o.respFunc(ctx, v)
}

func (o *webHookHandler) delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		if o.respFunc(ctx, fmt.Errorf("Path param id invalid ")) {
			return
		}
	}
	err := o.Service.Delete(ctx, id)
	if err != nil {
		if o.respFunc(ctx, err) {
			return
		}
	}

	_ = o.respFunc(ctx, nil)
}

func defaultResponse(ctx *gin.Context, o interface{}) bool {
	if o == nil {
		ctx.Status(http.StatusOK)
		return false
	}
	if e, ok := o.(error); ok {
		_ = ctx.AbortWithError(http.StatusBadRequest, e)
		return true
	} else {
		ctx.JSON(http.StatusOK, result.Ok(o))
		return false
	}
}
