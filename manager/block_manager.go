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
	"github.com/xfali/neve-webhook/events"
	"github.com/xfali/neve-webhook/notifier"
	"github.com/xfali/neve-webhook/recorder"
	"github.com/xfali/neve-webhook/serialize"
	"github.com/xfali/xlog"
	"time"
)

const (
	ResponseChanBufferSize = 256
)

type BlockOpt func(m *blockManager)

type blockManager struct {
	logger xlog.Logger

	recorder recorder.Recorder
	notifier notifier.Notifier

	ctx    context.Context
	cancel context.CancelFunc

	signFunc      SignatureFunc
	notifyTimeout time.Duration
}

func NewSyncManager(recorder recorder.Recorder, opts ...BlockOpt) *blockManager {
	ret := &blockManager{
		logger:        xlog.GetLogger(),
		recorder:      recorder,
		notifier:      notifier.NewHttpNotifier(nil),
		signFunc:      defaultSignFunc,
		notifyTimeout: NotifyTimeout,
	}
	for _, opt := range opts {
		opt(ret)
	}
	return ret
}

func (m *blockManager) BeanAfterSet() error {
	return m.Start()
}

func (m *blockManager) BeanDestroy() error {
	return m.Close()
}

func (m *blockManager) Start() error {
	m.ctx, m.cancel = context.WithCancel(context.Background())
	return nil
}

func (m *blockManager) Close() error {
	m.cancel()

	return nil
}

func (m *blockManager) Notify(ctx context.Context, event events.IEvent, ds serialize.Deserializer) (<-chan interface{}, error) {
	offset := int64(0)
	respChan := make(chan interface{}, ResponseChanBufferSize)
	errList := &errors.LockedErrList{}
	for {
		datas, _, err := m.recorder.Query(ctx, recorder.QueryCondition{
			EventType: event.GetType(),
			Offset:    offset,
			PageSize:  ResponseChanBufferSize,
		})
		if err != nil {
			return nil, err
		}
		if len(datas) == 0 {
			break
		}
		offset++
		for _, d := range datas {
			go m.notify(ctx, d, event, ds, respChan, errList)
		}
	}

	return respChan, nil
}

func (m *blockManager) notify(ctx context.Context, d recorder.Data, event events.IEvent, ds serialize.Deserializer, respChan chan interface{}, errs errors.ErrorList) {
	var resp interface{}
	holder := &resp
	defer func(o *interface{}) {
		select {
		case <-ctx.Done():
			err := ctx.Err()
			if err != nil {
				m.logger.Warnln(err)
			}
			return
		case respChan <- *o:
		}
	}(holder)
	now := time.Now()
	secret, err := m.signFunc(d.Secret)
	if err != nil {
		m.logger.Errorln(err)
		return
	}
	nCtx, _ := context.WithTimeout(ctx, m.notifyTimeout)
	data, err := m.notifier.Send(nCtx, d.Url, d.ContentType, secret, event)
	if err != nil {
		errs.Add(err)
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
		if ds != nil {
			v, err := ds.Deserialize(data)
			if err != nil {
				errs.Add(err)
				m.logger.Errorln("Deserialize hook response failed: ", err)
			}
			holder = &v
		}
	}
}

type blockOpts struct{}

var BlockOpts blockOpts

func (o blockOpts) SetNotifier(n notifier.Notifier) BlockOpt {
	return func(m *blockManager) {
		m.notifier = n
	}
}

func (o blockOpts) SetSignatureFunc(f SignatureFunc) BlockOpt {
	return func(m *blockManager) {
		m.signFunc = f
	}
}

func (o blockOpts) SetNotifyTimeout(t time.Duration) BlockOpt {
	return func(m *blockManager) {
		m.notifyTimeout = t
	}
}
