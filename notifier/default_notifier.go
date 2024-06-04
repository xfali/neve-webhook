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

package notifier

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/xfali/xlog"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

func defaultTransportDialContext(dialer *net.Dialer) func(context.Context, string, string) (net.Conn, error) {
	return dialer.DialContext
}

var DefaultTransport http.RoundTripper = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: defaultTransportDialContext(&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}),
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
	TLSClientConfig: &tls.Config{
		InsecureSkipVerify: true,
	},
}

type httpNotifier struct {
	logger xlog.Logger
	client *http.Client
}

func NewHttpNotifier(client *http.Client) *httpNotifier {
	ret := &httpNotifier{
		logger: xlog.GetLogger(),
		client: &http.Client{
			Timeout:   30 * time.Second,
			Transport: DefaultTransport,
		},
	}
	if client != nil {
		ret.client = client
	}
	return ret
}

func (n *httpNotifier) Send(ctx context.Context, url string, contentType string, secretSign string, eventType string, payLoad interface{}) error {
	var data []byte
	var err error
	if contentType == "" {
		contentType = "application/json"
	}
	if strings.Index(contentType, "application/json") == 0 {
		if payLoad != nil {
			data, err = json.Marshal(payLoad)
			if err != nil {
				return err
			}
		}
	} else if strings.Index(contentType, "application/xml") == 0 {
		if payLoad != nil {
			data, err = xml.Marshal(payLoad)
			if err != nil {
				return err
			}
		}
	}

	var r io.Reader
	if len(data) > 0 {
		r = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, r)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("X-Neve-WebHook-Event", eventType)
	req.Header.Set("X-Neve-WebHook-Signature", secretSign)
	resp, err := n.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		err = fmt.Errorf("Notify to URL %s failed, http status: %d ", url, resp.StatusCode)
		d, _ := ioutil.ReadAll(resp.Body)
		respStr := ""
		if len(d) > 0 {
			respStr = string(d)
		}
		n.logger.Errorf("Notify error: %v, response data: %s \n", err, respStr)
		return err
	}
	return nil
}
