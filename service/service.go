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

package service

import (
	"context"
	"github.com/xfali/neve-webhook/recorder"
)

type WebHookService interface {
	Create(ctx context.Context, rec recorder.Input) (string, error)

	Update(ctx context.Context, id string, rec recorder.Input) error

	Get(ctx context.Context, cond recorder.QueryCondition) (ListData, error)

	Detail(ctx context.Context, id string) (recorder.Data, error)

	Delete(ctx context.Context, id string) error
}
