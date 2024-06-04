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

import "sync/atomic"

type IdGenerator interface {
	Next() int64
}

func NewIdGenerator() *defaultIdGenerator {
	return &defaultIdGenerator{
		id: 0,
	}
}

type defaultIdGenerator struct {
	id int64
}

func (g *defaultIdGenerator) Next() int64 {
	return atomic.AddInt64(&g.id, 1)
}
