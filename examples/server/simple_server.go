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

package main

import (
	"github.com/xfali/neve-core"
	"github.com/xfali/neve-core/processor"
	"github.com/xfali/neve-utils/neverror"
	neveweb "github.com/xfali/neve-web"
	"github.com/xfali/neve-webhook/servers"
	"github.com/xfali/xlog"
	"os"
)

func main() {
	xlog.Infoln(os.Getwd())
	app := neve.NewFileConfigApplication("examples/server/config.yaml")
	neverror.PanicError(app.RegisterBean(processor.NewValueProcessor()))
	neverror.PanicError(app.RegisterBean(neveweb.NewGinProcessor()))
	neverror.PanicError(app.RegisterBean(servers.NewWebhooksServerProcessor()))
	neverror.PanicError(app.RegisterBean(NewService()))
	neverror.PanicError(app.Run())
}
