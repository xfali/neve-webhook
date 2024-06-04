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
)

type Data struct {
	Url               string   `json:"url" xml:"url" yaml:"url"`
	ContentType       string   `json:"content_type" xml:"content_type" yaml:"content_type"`
	Secret            string   `json:"secret" xml:"secret" yaml:"secret"`
	TriggerEventTypes []string `json:"event_type" xml:"event_type" yaml:"event_type"`
}

type QueryCondition struct {
	Id        string
	EventType string
	Url       string
}

type Recorder interface {
	Query(ctx context.Context, condition QueryCondition) ([]Data, error)

	Create(ctx context.Context, data Data) (string, error)

	Update(ctx context.Context, id string, data Data) error

	Delete(ctx context.Context, id string) error
}
