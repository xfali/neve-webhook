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
	"github.com/xfali/neve-webhook/auth"
	"github.com/xfali/neve-webhook/errors"
	"github.com/xfali/neve-webhook/events"
	"github.com/xfali/neve-webhook/notifier"
	"github.com/xfali/neve-webhook/recorder"
	"github.com/xfali/neve-webhook/serialize"
	"github.com/xfali/xlog"
	"time"
)

const (
	NotifyTimeout = 15 * time.Second
	NotifySize    = 1024
)

type Opt func(m *defaultManager)

type SignatureFunc func(secret string) (string, error)

type defaultManager struct {
	logger xlog.Logger

	recorder recorder.Recorder
	notifier notifier.Notifier
	eventSvc events.Service

	ctx    context.Context
	cancel context.CancelFunc

	signFunc      SignatureFunc
	notifyTimeout time.Duration
}

func NewManager(recorder recorder.Recorder, opts ...Opt) *defaultManager {
	ret := &defaultManager{
		logger:        xlog.GetLogger(),
		recorder:      recorder,
		eventSvc:      events.NewEventService(-1),
		notifier:      notifier.NewHttpNotifier(nil),
		signFunc:      defaultSignFunc,
		notifyTimeout: NotifyTimeout,
	}
	for _, opt := range opts {
		opt(ret)
	}
	return ret
}

func (m *defaultManager) BeanAfterSet() error {
	return m.Start()
}

func (m *defaultManager) BeanDestroy() error {
	return m.Close()
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
				m.logger.Warnln("Get Event failed: ", err)
			} else {
				err = m.doNotify(m.ctx, e)
				if err != nil {
					m.logger.Warnln("Notify Event failed: ", err)
				}
			}
		}
	}
}

func (m *defaultManager) Notify(ctx context.Context, event events.IEvent, d serialize.Deserializer) (<-chan interface{}, error) {
	return nil, m.eventSvc.Put(ctx, event)
}

func (m *defaultManager) doNotify(ctx context.Context, event events.IEvent) error {
	offset := int64(0)
	var errList errors.ErrList
	now := time.Now()
	for {
		datas, _, err := m.recorder.Query(ctx, recorder.QueryCondition{
			EventType: event.GetType(),
			Offset:    offset,
			PageSize:  NotifySize,
		})
		if err != nil {
			return err
		}
		if len(datas) == 0 {
			break
		}
		offset++
		for _, d := range datas {
			secret, err := m.signFunc(d.Secret)
			if err != nil {
				errList.Add(err)
				m.logger.Errorln(err)
				continue
			}
			nCtx, _ := context.WithTimeout(ctx, m.notifyTimeout)
			_, err = m.notifier.Send(nCtx, d.Url, d.ContentType, secret, event)
			if err != nil {
				errList.Add(err)
				m.logger.Errorln("Notifier send message failed: ", err)
				err = m.recorder.UpdateNotifyStatus(ctx, d.ID, now, false)
				if err != nil {
					m.logger.Errorln("Recorder UpdateNotifyStatus failed: ", err)
				}
			} else {
				err = m.recorder.UpdateNotifyStatus(ctx, d.ID, now, true)
				if err != nil {
					m.logger.Errorln("Recorder UpdateNotifyStatus failed: ", err)
				}
			}
		}
	}

	if !errList.Empty() {
		return errList
	}
	return nil
}

func defaultSignFunc(secret string) (string, error) {
	return auth.HmacSignature(auth.DefaultSignatureKey, secret)
}

type opts struct{}

var Opts opts

func (o opts) SetEventService(s events.Service) Opt {
	return func(m *defaultManager) {
		m.eventSvc = s
	}
}

func (o opts) SetNotifier(n notifier.Notifier) Opt {
	return func(m *defaultManager) {
		m.notifier = n
	}
}

func (o opts) SetSignatureFunc(f SignatureFunc) Opt {
	return func(m *defaultManager) {
		m.signFunc = f
	}
}

func (o opts) SetNotifyTimeout(t time.Duration) Opt {
	return func(m *defaultManager) {
		m.notifyTimeout = t
	}
}
