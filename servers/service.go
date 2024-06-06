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
	"context"
	"fmt"
	"github.com/xfali/neve-webhook/recorder"
	"github.com/xfali/neve-webhook/service"
	"github.com/xfali/xlog"
)

type webHookServiceImpl struct {
	logger   xlog.Logger
	Recorder recorder.Recorder `inject:""`
}

func NewWebHookService() *webHookServiceImpl {
	return &webHookServiceImpl{
		logger: xlog.GetLogger(),
	}
}

func (s *webHookServiceImpl) Create(ctx context.Context, rec recorder.Input) (string, error) {
	return s.Recorder.Create(ctx, rec)
}

func (s *webHookServiceImpl) Update(ctx context.Context, id string, rec recorder.Input) error {
	return s.Recorder.Update(ctx, id, rec)
}

func (s *webHookServiceImpl) Get(ctx context.Context, cond recorder.QueryCondition) (service.ListData, error) {
	v, total, err := s.Recorder.Query(ctx, cond)
	return service.ListData{
		Webhooks: v,
		Total:    total,
	}, err
}

func (s *webHookServiceImpl) Detail(ctx context.Context, id string) (recorder.Data, error) {
	v, _, err := s.Recorder.Query(ctx, recorder.QueryCondition{Id: id})
	if err != nil {
		return recorder.Data{}, err
	}
	if len(v) == 0 {
		return recorder.Data{}, fmt.Errorf("ID %s not found ", id)
	}
	return v[0], nil
}

func (s *webHookServiceImpl) Delete(ctx context.Context, id string) error {
	return s.Recorder.Delete(ctx, id)
}
