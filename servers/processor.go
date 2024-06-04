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
	"github.com/xfali/fig"
	"github.com/xfali/neve-core/bean"
	"github.com/xfali/neve-webhook/manager"
	"github.com/xfali/neve-webhook/recorder"
)

type ProcessorOpt func(*neveGinProcessor)

type RecorderCreator func() recorder.Recorder

type ManagerCreator func(r recorder.Recorder) manager.Manager

type neveGinProcessor struct {
	recorderCreator RecorderCreator
	managerCreator  ManagerCreator
}

func NewWebhooksServerProcessor(opts ...ProcessorOpt) *neveGinProcessor {
	ret := &neveGinProcessor{
		recorderCreator: func() recorder.Recorder {
			return recorder.NewMemRecorder()
		},
		managerCreator: func(r recorder.Recorder) manager.Manager {
			return manager.NewManager(r)
		},
	}
	for _, opt := range opts {
		opt(ret)
	}
	return ret
}

func (p *neveGinProcessor) Init(conf fig.Properties, container bean.Container) error {
	recorder := p.recorderCreator()
	if err := container.Register(recorder); err != nil {
		return err
	}
	manager := p.managerCreator(recorder)
	if err := container.Register(manager); err != nil {
		return err
	}
	if err := container.Register(NewWebHookService()); err != nil {
		return err
	}
	if err := container.Register(NewWebHookHandler()); err != nil {
		return err
	}
	return nil
}

func (p *neveGinProcessor) Classify(o interface{}) (bool, error) {
	return false, nil
}

func (p *neveGinProcessor) Process() error {
	return nil
}

func (p *neveGinProcessor) BeanDestroy() error {
	return nil
}

type processorOpts struct {
}

var ProcessorOpts processorOpts

func (o processorOpts) RecorderCreator(recorderCreator RecorderCreator) ProcessorOpt {
	return func(processor *neveGinProcessor) {
		processor.recorderCreator = recorderCreator
	}
}

func (o processorOpts) ManagerCreator(creator ManagerCreator) ProcessorOpt {
	return func(processor *neveGinProcessor) {
		processor.managerCreator = creator
	}
}
