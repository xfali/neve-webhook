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
	"github.com/xfali/goutils/container/xmap"
	"strconv"
	"sync"
	"time"
)

type Opt func(r *simpleRecorder)

type memRecorder struct {
	locker      sync.RWMutex
	idGenerator IdGenerator
	eventMap    map[string]*xmap.LinkedMap
	urlMap      map[string]string
	idMap       *xmap.LinkedMap
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
		eventMap:    map[string]*xmap.LinkedMap{},
		idMap:       xmap.NewLinkedMap(),
		urlMap:      map[string]string{},
		idGenerator: NewIdGenerator(),
	}
	return ret
}

func (r *singleFilter) Filter(ctx context.Context) (Recorder, error) {
	return r.r, nil
}

func (r *memRecorder) Create(ctx context.Context, input Input) (string, error) {
	r.locker.Lock()
	defer r.locker.Unlock()

	if input.Url == "" {
		return "", fmt.Errorf("Url cannot be empty ")
	}

	if url := r.urlMap[input.Url]; url != "" {
		return "", fmt.Errorf("Url have been exists ")
	}

	id := r.idGenerator.Next()
	idStr := strconv.FormatInt(id, 10)

	data := input.ToData()
	data.ID = idStr
	data.State = HookStateNormal
	r.idMap.Put(idStr, &data)
	r.urlMap[data.Url] = idStr
	for _, e := range data.TriggerEventTypes {
		if m, ok := r.eventMap[e]; ok {
			m.Put(idStr, struct {
			}{})
		} else {
			lm := xmap.NewLinkedMap()
			lm.Put(idStr, struct {
			}{})
			r.eventMap[e] = lm
		}
	}
	return idStr, nil
}

func (r *memRecorder) Update(ctx context.Context, idStr string, data Input) error {
	r.locker.Lock()
	defer r.locker.Unlock()

	if x, ok := r.idMap.Get(idStr); ok {
		v := x.(*Data)
		for _, e := range v.TriggerEventTypes {
			r.eventMap[e].Delete(data.Url)
		}
		if data.Url != "" {
			if data.Url != v.Url {
				delete(r.urlMap, v.Url)
				r.urlMap[data.Url] = idStr
			}
			v.Url = data.Url
		}
		if data.Secret != "" {
			v.Secret = data.Secret
		}
		if data.ContentType != "" {
			v.ContentType = data.ContentType
		}
		if data.State != "" {
			v.State = data.State
		}
		v.TriggerEventTypes = data.TriggerEventTypes
	} else {
		return fmt.Errorf("ID %s not found ", idStr)
	}
	for _, e := range data.TriggerEventTypes {
		if m, ok := r.eventMap[e]; ok {
			m.Put(idStr, struct {
			}{})
		} else {
			lm := xmap.NewLinkedMap()
			lm.Put(idStr, struct {
			}{})
			r.eventMap[e] = lm
		}
	}
	return nil
}

func (r *memRecorder) UpdateNotifyStatus(ctx context.Context, id string, updateTime time.Time, success bool) error {
	r.locker.Lock()
	defer r.locker.Unlock()

	idStr := id
	if x, ok := r.idMap.Get(idStr); ok {
		v := x.(*Data)
		if success {
			v.SuccessCount += 1
			v.LastSuccessTime = updateTime
		} else {
			v.FailureCount += 1
			v.LastFailureTime = updateTime
		}
	} else {
		return fmt.Errorf("ID %s not found ", id)
	}
	return nil
}

func (r *memRecorder) Delete(ctx context.Context, id string) error {
	r.locker.Lock()
	defer r.locker.Unlock()

	if x, ok := r.idMap.Get(id); ok {
		v := x.(*Data)
		for _, e := range v.TriggerEventTypes {
			r.eventMap[e].Delete(v.Url)
		}
		r.idMap.Delete(id)
		delete(r.urlMap, v.Url)
	}
	return nil
}

func (r *memRecorder) Query(ctx context.Context, condition QueryCondition) ([]Data, int64, error) {
	r.locker.RLock()
	defer r.locker.RUnlock()

	if condition.PageSize == 0 {
		condition.PageSize = 20
	}
	total := int64(r.idMap.Size())
	if condition.Id != "" {
		if v, ok := r.idMap.Get(condition.Id); ok {
			return []Data{*v.(*Data)}, total, nil
		} else {
			return nil, total, fmt.Errorf("ID %s not found ", condition.Id)
		}
	}

	if condition.Url != "" {
		ret, err := r.queryByUrl(ctx, condition.Url)
		return ret, total, err
	}

	if condition.EventType != "" {
		ret, err := r.queryByEventType(ctx, condition.EventType, condition.State, condition.Offset, condition.PageSize)
		return ret, total, err
	}

	ret := make([]Data, 0, condition.PageSize)
	skip := condition.Offset * condition.PageSize
	current := int64(0)
	r.idMap.Foreach(func(key interface{}, value interface{}) bool {
		if skip > current {
			current++
			return true
		}
		if current-skip < condition.PageSize {
			hd := value.(*Data)
			if condition.State != "" {
				if hd.State == condition.State {
					ret = append(ret, *hd)
				}
			} else {
				ret = append(ret, *hd)
			}
			current++
			return true
		} else {
			return false
		}
	})
	return ret, total, nil
}

func (r *memRecorder) queryByUrl(ctx context.Context, url string) ([]Data, error) {
	id := r.urlMap[url]
	if id == "" {
		return nil, fmt.Errorf("Url %s not found ", url)
	}

	v, have := r.idMap.Get(id)
	if !have {
		return nil, fmt.Errorf("Url %s not found ", url)
	}
	return []Data{*v.(*Data)}, nil
}

func (r *memRecorder) queryByEventType(ctx context.Context, eventType, state string, offset, pageSize int64) ([]Data, error) {
	maps := r.eventMap[eventType]

	if maps != nil && maps.Size() > 0 {
		skip := offset * pageSize
		current := int64(0)

		ret := make([]Data, 0, 16)
		maps.Foreach(func(key interface{}, value interface{}) bool {
			if skip > current {
				current++
				return true
			}
			if current-skip < pageSize {
				if v, have := r.idMap.Get(key); have {
					hd := v.(*Data)
					if state != "" {
						if hd.State == state {
							ret = append(ret, *hd)
						}
					} else {
						ret = append(ret, *hd)
					}
				}
				current++
				return true
			} else {
				return false
			}
		})

		return ret, nil
	}
	return nil, nil
}

type simpleRecorder struct {
	filter ContextFilter
}

func (r *simpleRecorder) Create(ctx context.Context, data Input) (string, error) {
	rr, err := r.filter.Filter(ctx)
	if err != nil {
		return "", err
	}
	return rr.Create(ctx, data)
}

func (r *simpleRecorder) Query(ctx context.Context, condition QueryCondition) ([]Data, int64, error) {
	rr, err := r.filter.Filter(ctx)
	if err != nil {
		return nil, 0, err
	}
	return rr.Query(ctx, condition)
}

func (r *simpleRecorder) Update(ctx context.Context, id string, data Input) error {
	rr, err := r.filter.Filter(ctx)
	if err != nil {
		return err
	}
	return rr.Update(ctx, id, data)
}

func (r *simpleRecorder) UpdateNotifyStatus(ctx context.Context, id string, updateTime time.Time, success bool) error {
	rr, err := r.filter.Filter(ctx)
	if err != nil {
		return err
	}
	return rr.UpdateNotifyStatus(ctx, id, updateTime, success)
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
