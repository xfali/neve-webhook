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
	"testing"
)

func TestRecorder(t *testing.T) {
	r := NewMemRecorder()
	ctx := context.Background()
	id, err := r.Create(ctx, Data{
		Url: "test",
		TriggerEventTypes: []string{
			"push",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id)

	v, err := r.Query(ctx, QueryCondition{
		Id: id,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v)
	if len(v) != 1 {
		t.Fatalf("Expect 1 but get %d\n", len(v))
	}

	v, err = r.Query(ctx, QueryCondition{
		EventType: "push",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v)
	if len(v) != 1 {
		t.Fatalf("Expect 1 but get %d\n", len(v))
	}

	err = r.Update(ctx, id, Data{
		Url: "world",
	})
	if err != nil {
		t.Fatal(err)
	}
	v, _ = r.Query(ctx, QueryCondition{
		Id: id,
	})
	t.Log(v)
	if len(v) != 1 {
		t.Fatalf("Expect 1 but get %d\n", len(v))
	}
	if v[0].Url != "world" {
		t.Fatalf("Expect world but get %s\n", v[0].Url)
	}

	err = r.Delete(ctx, id)
	if err != nil {
		t.Fatal(err)
	}

	v, err = r.Query(ctx, QueryCondition{})
	if err != nil {
		t.Fatal(err)
	}

	if len(v) != 0 {
		t.Fatalf("Expect 0 but get %d\n", len(v))
	}
}
