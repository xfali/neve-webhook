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

package manager

import (
	"context"
	"github.com/xfali/neve-webhook/errors"
	events2 "github.com/xfali/neve-webhook/events"
	notifier2 "github.com/xfali/neve-webhook/notifier"
	"github.com/xfali/neve-webhook/recorder"
	"github.com/xfali/xlog"
)

type Opt func(m *defaultManager)

type defaultManager struct {
	logger xlog.Logger

	recorder recorder.Recorder
	notifier notifier2.Notifier
	eventSvc events2.Service

	ctx    context.Context
	cancel context.CancelFunc
}

func NewManager(recorder recorder.Recorder, opts ...Opt) *defaultManager {
	ret := &defaultManager{
		logger:   xlog.GetLogger(),
		recorder: recorder,
		eventSvc: events2.NewEventService(-1),
		notifier: notifier2.NewHttpNotifier(nil),
	}
	for _, opt := range opts {
		opt(ret)
	}
	return ret
}

func (m *defaultManager) Start() error {
	err := m.eventSvc.Connect()
	if err != nil {
		return err
	}
	m.ctx, m.cancel = context.WithCancel(context.Background())
	go m.loop()
	return nil
}

func (m *defaultManager) Close() error {
	m.cancel()

	return m.eventSvc.Disconnect()
}

func (m *defaultManager) loop() {
	for {
		select {
		case <-m.ctx.Done():
			return
		default:
			e, err := m.eventSvc.Get(m.ctx)
			if err != nil {
				m.logger.Errorln("Get Event failed: ", err)
			} else {
				err = m.doNotify(m.ctx, e)
				if err != nil {
					m.logger.Errorln("Notify Event failed: ", err)
				}
			}
		}
	}
}

func (m *defaultManager) Notify(ctx context.Context, event *events2.Event) error {
	return m.eventSvc.Put(ctx, event)
}

func (m *defaultManager) doNotify(ctx context.Context, event *events2.Event) error {
	datas, err := m.recorder.Query(ctx, recorder.QueryCondition{
		EventType: event.Type,
	})
	if err != nil {
		return err
	}

	var errList errors.ErrorList

	for _, d := range datas {
		err = m.notifier.Send(ctx, d.Url, d.ContentType, d.Secret, event.Type, event.PayLoad)
		if err != nil {
			_ = errList.Add(err)
		}
	}

	if !errList.Empty() {
		return errList
	}
	return nil
}

type opts struct{}

var Opts opts

func (o opts) SetEventService(s events2.Service) Opt {
	return func(m *defaultManager) {
		m.eventSvc = s
	}
}

func (o opts) SetNotifier(n notifier2.Notifier) Opt {
	return func(m *defaultManager) {
		m.notifier = n
	}
}
