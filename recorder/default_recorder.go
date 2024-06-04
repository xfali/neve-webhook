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

package recorder

import (
	"context"
	"fmt"
	"strconv"
	"sync"
)

type Opt func(r *simpleRecorder)

type memRecorder struct {
	locker      sync.RWMutex
	idGenerator IdGenerator
	eventMap    map[string]map[string]struct{}
	idMap       map[string]*Data
}

type ContextFilter interface {
	Filter(ctx context.Context) (Recorder, error)
}

type singleFilter struct {
	r Recorder
}

func NewSingleMemFilter() *singleFilter {
	return &singleFilter{
		r: NewMemRecorder(),
	}
}

func NewSimpleRecorder() *simpleRecorder {
	ret := &simpleRecorder{
		filter: NewSingleMemFilter(),
	}
	return ret
}

func NewMemRecorder() *memRecorder {
	ret := &memRecorder{
		eventMap:    map[string]map[string]struct{}{},
		idMap:       map[string]*Data{},
		idGenerator: NewIdGenerator(),
	}
	return ret
}

func (r *singleFilter) Filter(ctx context.Context) (Recorder, error) {
	return r.r, nil
}

func (r *memRecorder) Create(ctx context.Context, data Data) (string, error) {
	r.locker.Lock()
	defer r.locker.Unlock()

	id := r.idGenerator.Next()
	idStr := strconv.FormatInt(id, 10)
	if v, ok := r.idMap[idStr]; ok {
		for _, e := range v.TriggerEventTypes {
			delete(r.eventMap[e], data.Url)
		}
		v.Url = data.Url
		v.Secret = data.Secret
		v.ContentType = data.ContentType
		v.TriggerEventTypes = data.TriggerEventTypes
	} else {
		r.idMap[idStr] = &data
	}
	for _, e := range data.TriggerEventTypes {
		if m, ok := r.eventMap[e]; ok {
			m[idStr] = struct{}{}
		} else {
			r.eventMap[e] = map[string]struct{}{
				idStr: {},
			}
		}
	}
	return idStr, nil
}

func (r *memRecorder) Update(ctx context.Context, id string, data Data) error {
	r.locker.Lock()
	defer r.locker.Unlock()

	idStr := id
	if v, ok := r.idMap[idStr]; ok {
		for _, e := range v.TriggerEventTypes {
			delete(r.eventMap[e], data.Url)
		}
		if data.Url != "" {
			v.Url = data.Url
		}
		if data.Secret != "" {
			v.Secret = data.Secret
		}
		if data.ContentType != "" {
			v.ContentType = data.ContentType
		}
		v.TriggerEventTypes = data.TriggerEventTypes
	} else {
		return fmt.Errorf("ID %s not found ", id)
	}
	for _, e := range data.TriggerEventTypes {
		if m, ok := r.eventMap[e]; ok {
			m[idStr] = struct{}{}
		} else {
			r.eventMap[e] = map[string]struct{}{
				idStr: {},
			}
		}
	}
	return nil
}

func (r *memRecorder) Delete(ctx context.Context, id string) error {
	r.locker.Lock()
	defer r.locker.Unlock()

	if v, ok := r.idMap[id]; ok {
		for _, e := range v.TriggerEventTypes {
			delete(r.eventMap[e], v.Url)
		}
		delete(r.idMap, id)
	}
	return nil
}

func (r *memRecorder) Query(ctx context.Context, condition QueryCondition) ([]Data, error) {
	r.locker.RLock()
	defer r.locker.RUnlock()

	if condition.Id != "" {
		if v, ok := r.idMap[condition.Id]; ok {
			return []Data{*v}, nil
		} else {
			return nil, fmt.Errorf("ID %s not found ", condition.Id)
		}
	}

	if condition.EventType != "" {
		return r.queryByEventType(ctx, condition.EventType)
	}

	ret := make([]Data, 0, len(r.idMap))
	for _, v := range r.idMap {
		ret = append(ret, *v)
	}
	return ret, nil
}

func (r *memRecorder) queryByEventType(ctx context.Context, eventType string) ([]Data, error) {
	r.locker.RLock()
	defer r.locker.RUnlock()

	maps := r.eventMap[eventType]

	if len(maps) > 0 {
		ret := make([]Data, 0, len(maps))
		for k := range maps {
			ret = append(ret, *r.idMap[k])
		}
		return ret, nil
	}
	return nil, nil
}

type simpleRecorder struct {
	filter ContextFilter
}

func (r *simpleRecorder) Create(ctx context.Context, data Data) (string, error) {
	rr, err := r.filter.Filter(ctx)
	if err != nil {
		return "", err
	}
	return rr.Create(ctx, data)
}

func (r *simpleRecorder) Query(ctx context.Context, condition QueryCondition) ([]Data, error) {
	rr, err := r.filter.Filter(ctx)
	if err != nil {
		return nil, err
	}
	return rr.Query(ctx, condition)
}

func (r *simpleRecorder) Update(ctx context.Context, id string, data Data) error {
	rr, err := r.filter.Filter(ctx)
	if err != nil {
		return err
	}
	return rr.Update(ctx, id, data)
}

func (r *simpleRecorder) Delete(ctx context.Context, id string) error {
	rr, err := r.filter.Filter(ctx)
	if err != nil {
		return err
	}
	return rr.Delete(ctx, id)
}

type opts struct{}

var Opts opts

func (o opts) SetFilter(f ContextFilter) Opt {
	return func(r *simpleRecorder) {
		r.filter = f
	}
}
