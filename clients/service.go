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
}

func NewWebHookClient(endpoint string, client restclient.RestClient) *webHooksClient {
	return &webHooksClient{
		endpoint: endpoint,
		client:   client,
	}
}

func (s *webHooksClient) Create(ctx context.Context, rec recorder.Data) (string, error) {
	ret := Result[string]{}
	err := s.client.Exchange(s.endpoint,
		request.WithRequestContext(ctx),
		request.MethodPost(),
		request.WithRequestBody(rec),
		request.WithResult(&ret))
	return ret.Data, err
}

func (s *webHooksClient) Update(ctx context.Context, id string, rec recorder.Data) error {
	err := s.client.Exchange(s.endpoint+"/"+id,
		request.WithRequestContext(ctx),
		request.MethodPut(),
		request.WithRequestBody(rec))
	return err
}

func (s *webHooksClient) Get(ctx context.Context, cond recorder.QueryCondition) ([]recorder.Data, error) {
	ret := Result[[]recorder.Data]{}
	err := s.client.Exchange(fmt.Sprintf("%s?id=%s&event_type=%s", s.endpoint, cond.Id, cond.EventType),
		request.WithRequestContext(ctx),
		request.MethodGet(),
		request.WithResult(&ret))
	return ret.Data, err
}

func (s *webHooksClient) Detail(ctx context.Context, id string) (recorder.Data, error) {
	ret := Result[recorder.Data]{}
	err := s.client.Exchange(s.endpoint+"/"+id,
		request.WithRequestContext(ctx),
		request.MethodGet(),
		request.WithResult(&ret))
	return ret.Data, err
}

func (s *webHooksClient) Delete(ctx context.Context, id string) error {
	err := s.client.Exchange(s.endpoint+"/"+id,
		request.WithRequestContext(ctx),
		request.MethodDelete())
	return err
}

type Result[T any] struct {
	Code int64  `json:"code"`
	Msg  string `json:"message"`
	Data T      `json:"data,omitempty"`
}
