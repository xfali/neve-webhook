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
	"time"
)

type Service struct {
	Manager  manager.Manager `inject:""`
	stopChan chan struct{}
}

func NewService() *Service {
	return &Service{
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
				_ = o.Manager.Notify(context.Background(), &events.Event{
					Type:    "push",
					PayLoad: result.Ok("This is a test"),
				})
			}
		}
	}()
	return nil
}

func (o *Service) BeanDestroy() error {
	close(o.stopChan)
	return nil
}
