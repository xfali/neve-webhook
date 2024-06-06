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
	"github.com/xfali/neve-web/result"
	"github.com/xfali/neve-webhook/events"
	"github.com/xfali/neve-webhook/manager"
	"github.com/xfali/neve-webhook/serialize"
	"github.com/xfali/xlog"
	"time"
)

type Service struct {
	logger   xlog.Logger
	Manager  manager.Manager `inject:""`
	stopChan chan struct{}
}

func NewService() *Service {
	return &Service{
		logger:   xlog.GetLogger(),
		stopChan: make(chan struct{}),
	}
}

func (o *Service) BeanAfterSet() error {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-o.stopChan:
				return
			case <-ticker.C:
				ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
				resps, err := o.Manager.Notify(ctx, &events.Event{
					Type:    "push",
					PayLoad: result.Ok("This is a test"),
				}, serialize.DeserializeFunc(func(bytes []byte) (interface{}, error) {
					return string(bytes), nil
				}))
				if resps != nil {
					end := false
					for {
						select {
						case <-ctx.Done():
							o.logger.Errorln(ctx.Err())
							end = true
						case v := <-resps:
							o.logger.Infoln("Response: ", v)
						}
						if end {
							break
						}
					}
				}
				o.logger.Infoln("Notify error: ", err)
			}
		}
	}()
	return nil
}

func (o *Service) BeanDestroy() error {
	close(o.stopChan)
	return nil
}
