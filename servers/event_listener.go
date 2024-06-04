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
	"github.com/xfali/neve-core/appcontext"
	"github.com/xfali/neve-webhook/events"
	"github.com/xfali/neve-webhook/manager"
	"github.com/xfali/xlog"
)

type WebhookEvent interface {
	appcontext.ApplicationEvent
	events.IEvent
}

type eventListener struct {
	logger  xlog.Logger
	Manager manager.Manager `inject:""`
}

func NewEventListener() *eventListener {
	ret := &eventListener{
		logger: xlog.GetLogger(),
	}
	return ret
}

func (l *eventListener) RegisterConsumer(registry appcontext.ApplicationEventConsumerRegistry) error {
	return registry.RegisterApplicationEventConsumer(l.handlerEvent)
}

// 当事件为*customerEvent类型时自动匹配并调用该方法
// customerEvent需实现ApplicationEvent接口
func (l *eventListener) handlerEvent(event WebhookEvent) {
	err := l.Manager.Notify(context.Background(), event)
	if err != nil {
		l.logger.Errorln(err)
	}
}
