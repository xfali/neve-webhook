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

import "strings"

type ErrorList []error

func (l *ErrorList) Add(e ...error) *ErrorList {
	*l = append(*l, e...)
	return l
}

func (l ErrorList) Empty() bool {
	return len(l) == 0
}

func (l ErrorList) Error() string {
	buf := strings.Builder{}
	buf.WriteString("Have errors: ")
	for _, e := range l {
		buf.WriteString(e.Error())
		buf.WriteString(", ")
	}
	ret := buf.String()
	return ret[:len(ret)-2]
}