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

package errors

import (
	"strings"
	"sync"
)

type ErrorList interface {
	Add(e ...error)
	Empty() bool
	Error() string
}

type ErrList []error

func (l *ErrList) Add(e ...error) {
	*l = append(*l, e...)
}

func (l ErrList) Empty() bool {
	return len(l) == 0
}

func (l ErrList) Error() string {
	buf := strings.Builder{}
	buf.WriteString("Have errors: ")
	for _, e := range l {
		buf.WriteString(e.Error())
		buf.WriteString(", ")
	}
	ret := buf.String()
	return ret[:len(ret)-2]
}

type LockedErrList struct {
	errs   []error
	locker sync.RWMutex
}

func (l *LockedErrList) Add(e ...error) {
	l.locker.Lock()
	defer l.locker.Unlock()
	l.errs = append(l.errs, e...)
}

func (l *LockedErrList) Empty() bool {
	l.locker.RLocker()
	defer l.locker.RUnlock()
	return len(l.errs) == 0
}

func (l *LockedErrList) Error() string {
	l.locker.RLocker()
	defer l.locker.RUnlock()
	buf := strings.Builder{}
	buf.WriteString("Have errors: ")
	for _, e := range l.errs {
		buf.WriteString(e.Error())
		buf.WriteString(", ")
	}
	ret := buf.String()
	return ret[:len(ret)-2]
}
