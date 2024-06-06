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
	"time"
)

const (
	HookStateNormal    = "normal"
	HookStateAbnormal  = "abnormal"
	HookStateForbidden = "forbidden"
)

type Data struct {
	ID                string    `json:"id" xml:"id" yaml:"id"`
	Url               string    `json:"url" xml:"url" yaml:"url"`
	ContentType       string    `json:"content_type" xml:"content_type" yaml:"content_type"`
	Secret            string    `json:"secret" xml:"secret" yaml:"secret"`
	TriggerEventTypes []string  `json:"event_type" xml:"event_type" yaml:"event_type"`
	State             string    `json:"state" xml:"state" yaml:"state"`
	FailureCount      int64     `json:"failure_count" xml:"failure_count" yaml:"failure_count"`
	SuccessCount      int64     `json:"success_count" xml:"success_count" yaml:"success_count"`
	LastFailureTime   time.Time `json:"last_failure_time" xml:"last_failure_time" yaml:"last_failure_time"`
	LastSuccessTime   time.Time `json:"last_success_time" xml:"last_success_time" yaml:"last_success_time"`
}

type Input struct {
	Url               string   `json:"url" xml:"url" yaml:"url"`
	ContentType       string   `json:"content_type" xml:"content_type" yaml:"content_type"`
	Secret            string   `json:"secret" xml:"secret" yaml:"secret"`
	TriggerEventTypes []string `json:"event_type" xml:"event_type" yaml:"event_type"`
	State             string   `json:"state" xml:"state" yaml:"state"`
}

func (i *Input) ToData() Data {
	return Data{
		Url:               i.Url,
		ContentType:       i.ContentType,
		Secret:            i.Secret,
		TriggerEventTypes: i.TriggerEventTypes,
		State:             i.State,
	}
}

type QueryCondition struct {
	Id        string
	EventType string
	Url       string
	State     string

	// Current page, start with 0
	Offset int64
	// Page size, default 20
	PageSize int64
}

type Recorder interface {
	Query(ctx context.Context, condition QueryCondition) ([]Data, int64, error)

	Create(ctx context.Context, data Input) (string, error)

	Update(ctx context.Context, id string, data Input) error

	UpdateNotifyStatus(ctx context.Context, id string, updateTime time.Time, success bool) error

	Delete(ctx context.Context, id string) error
}
