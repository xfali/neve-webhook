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
	"context"
	"fmt"
	"github.com/xfali/neve-webhook/recorder"
	"github.com/xfali/restclient/v2"
	"github.com/xfali/restclient/v2/request"
)

type webHooksClient struct {
	endpoint string
	client   restclient.RestClient

	CreatePath string `fig:"neve.web.hooks.routes.create"`
	UpdatePath string `fig:"neve.web.hooks.routes.update"`
	QueryPath  string `fig:"neve.web.hooks.routes.query"`
	DetailPath string `fig:"neve.web.hooks.routes.detail"`
	DeletePath string `fig:"neve.web.hooks.routes.delete"`
}

func NewWebHookClient(endpoint string, client restclient.RestClient) *webHooksClient {
	ret := &webHooksClient{
		endpoint: endpoint,
		client:   client,
	}
	return ret
}

func (s *webHooksClient) Create(ctx context.Context, rec recorder.Data) (string, error) {
	ret := Result[string]{}
	url := s.endpoint
	if s.CreatePath != "" {
		url = s.CreatePath
	}
	err := s.client.Exchange(url,
		request.WithRequestContext(ctx),
		request.MethodPost(),
		request.WithRequestBody(rec),
		request.WithResult(&ret))
	return ret.Data, err
}

func (s *webHooksClient) Update(ctx context.Context, id string, rec recorder.Data) error {
	url := s.endpoint + "/" + id
	if s.UpdatePath != "" {
		url = s.UpdatePath
	}
	err := s.client.Exchange(url,
		request.WithRequestContext(ctx),
		request.MethodPut(),
		request.WithRequestBody(rec))
	return err
}

func (s *webHooksClient) Get(ctx context.Context, cond recorder.QueryCondition) ([]recorder.Data, error) {
	url := s.endpoint
	if s.QueryPath != "" {
		url = s.QueryPath
	}
	ret := Result[[]recorder.Data]{}
	err := s.client.Exchange(fmt.Sprintf("%s?id=%s&event_type=%s", url, cond.Id, cond.EventType),
		request.WithRequestContext(ctx),
		request.MethodGet(),
		request.WithResult(&ret))
	return ret.Data, err
}

func (s *webHooksClient) Detail(ctx context.Context, id string) (recorder.Data, error) {
	url := s.endpoint + "/" + id
	if s.DetailPath != "" {
		url = s.DetailPath
	}
	ret := Result[recorder.Data]{}
	err := s.client.Exchange(url,
		request.WithRequestContext(ctx),
		request.MethodGet(),
		request.WithResult(&ret))
	return ret.Data, err
}

func (s *webHooksClient) Delete(ctx context.Context, id string) error {
	url := s.endpoint + "/" + id
	if s.DeletePath != "" {
		url = s.DeletePath
	}
	err := s.client.Exchange(url,
		request.WithRequestContext(ctx),
		request.MethodDelete())
	return err
}

type Result[T any] struct {
	Code int64  `json:"code"`
	Msg  string `json:"message"`
	Data T      `json:"data,omitempty"`
}
