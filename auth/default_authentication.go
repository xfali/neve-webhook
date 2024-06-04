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

package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
)

const (
	DefaultSignatureKey = "neve-webhook-key-H*o2A@_M"
)

type SimpleAuthentication string

func NewSimpleAuthentication(secret string) *SimpleAuthentication {
	ret := SimpleAuthentication(secret)
	return &ret
}

func (a *SimpleAuthentication) Signature(ctx context.Context) (string, error) {
	return HmacSignature(DefaultSignatureKey, string(*a))
}

func HmacSignature(key, data string) (string, error) {
	h := hmac.New(sha1.New, []byte(key))
	_, err := h.Write([]byte(data))
	if err != nil {
		return "", err
	}

	signature := h.Sum(nil)

	return hex.EncodeToString(signature), nil
}
