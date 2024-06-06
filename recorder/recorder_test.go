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

func TestRecorder1(t *testing.T) {
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

	v, _, err := r.Query(ctx, QueryCondition{
		Id: id,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v)
	if len(v) != 1 {
		t.Fatalf("Expect 1 but get %d\n", len(v))
	}

	_, err = r.Create(ctx, Data{
		Url: "test",
		TriggerEventTypes: []string{
			"push",
		},
	})
	if err == nil {
		t.Fatalf("Expect error but get nil")
	}

	v, _, err = r.Query(ctx, QueryCondition{
		Url: "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v)
	if len(v) != 1 {
		t.Fatalf("Expect 1 but get %d\n", len(v))
	}

	v, _, err = r.Query(ctx, QueryCondition{
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
	v, _, _ = r.Query(ctx, QueryCondition{
		Id: id,
	})
	t.Log(v)
	if len(v) != 1 {
		t.Fatalf("Expect 1 but get %d\n", len(v))
	}
	if v[0].Url != "world" {
		t.Fatalf("Expect world but get %s\n", v[0].Url)
	}

	v, _, err = r.Query(ctx, QueryCondition{
		Url: "world",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v)
	if len(v) != 1 {
		t.Fatalf("Expect 1 but get %d\n", len(v))
	}

	err = r.Delete(ctx, id)
	if err != nil {
		t.Fatal(err)
	}

	v, _, err = r.Query(ctx, QueryCondition{})
	if err != nil {
		t.Fatal(err)
	}

	if len(v) != 0 {
		t.Fatalf("Expect 0 but get %d\n", len(v))
	}
}

func TestRecorder2(t *testing.T) {
	r := NewMemRecorder()
	ctx := context.Background()
	id, err := r.Create(ctx, Data{
		Url: "test1",
		TriggerEventTypes: []string{
			"push",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id)

	id, err = r.Create(ctx, Data{
		Url: "test2",
		TriggerEventTypes: []string{
			"push",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id)

	id, err = r.Create(ctx, Data{
		Url: "test3",
		TriggerEventTypes: []string{
			"push",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id)

	v, _, err := r.Query(ctx, QueryCondition{
		EventType: "push",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v)
	if len(v) != 3 {
		t.Fatalf("Expect 3 but get %d\n", len(v))
	}

	if v[0].Url != "test1" {
		t.Fatalf("Expect test1 but get %s\n", v[0].Url)
	}

	if v[1].Url != "test2" {
		t.Fatalf("Expect test1 but get %s\n", v[1].Url)
	}

	if v[2].Url != "test3" {
		t.Fatalf("Expect test1 but get %s\n", v[2].Url)
	}

	v, _, err = r.Query(ctx, QueryCondition{
		EventType: "push",
		Offset:    0,
		PageSize:  2,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v)
	if len(v) != 2 {
		t.Fatalf("Expect 2 but get %d\n", len(v))
	}

	if v[0].Url != "test1" {
		t.Fatalf("Expect test3 but get %s\n", v[0].Url)
	}

	if v[1].Url != "test2" {
		t.Fatalf("Expect test1 but get %s\n", v[1].Url)
	}

	v, _, err = r.Query(ctx, QueryCondition{
		EventType: "push",
		Offset:    1,
		PageSize:  2,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v)
	if len(v) != 1 {
		t.Fatalf("Expect 1 but get %d\n", len(v))
	}

	if v[0].Url != "test3" {
		t.Fatalf("Expect test3 but get %s\n", v[0].Url)
	}
}
