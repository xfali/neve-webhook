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

package events

import (
	"context"
	"errors"
)

const (
	EventChanSize = 4096
)

var DisconnectedErr = errors.New("Service Disconnected ")

type defaultEventService struct {
	stopChan  chan struct{}
	eventChan chan IEvent
}

func NewEventService(bufSize int) *defaultEventService {
	if bufSize < 0 {
		bufSize = EventChanSize
	}
	ret := &defaultEventService{
		eventChan: make(chan IEvent, bufSize),
	}

	return ret
}

func (s *defaultEventService) Connect() error {
	s.stopChan = make(chan struct{})
	return nil
}

func (s *defaultEventService) Disconnect() error {
	select {
	case <-s.stopChan:
		return nil
	default:
		close(s.stopChan)
	}
	return nil
}

func (s *defaultEventService) Get(ctx context.Context) (IEvent, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-s.stopChan:
		return nil, DisconnectedErr
	case e := <-s.eventChan:
		return e, nil
	}
}

func (s *defaultEventService) Put(ctx context.Context, event IEvent) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-s.stopChan:
		return DisconnectedErr
	case s.eventChan <- event:
		return nil
	}
}
