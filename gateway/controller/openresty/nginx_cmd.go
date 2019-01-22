// RAINBOND, Application Management Platform
// Copyright (C) 2014-2017 Goodrain Co., Ltd.

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version. For any non-GPL usage of Rainbond,
// one or multiple Commercial Licenses authorized by Goodrain Co., Ltd.
// must be obtained first.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package openresty

import (
	"os"
	"os/exec"
	"path"

	"github.com/goodrain/rainbond/gateway/controller/openresty/template"
)

var (
	nginxBinary      = "nginx"
	defaultNginxConf = "/run/nginx/conf/nginx.conf"
)

func init() {
	nginxBinary = path.Join(os.Getenv("OPENRESTY_HOME"), "/nginx/sbin/nginx")
	ngx := os.Getenv("NGINX_BINARY")
	if ngx != "" {
		nginxBinary = ngx
	}
	customConfig := os.Getenv("NGINX_CUSTOM_CONFIG")
	if customConfig != "" {
		template.CustomConfigPath = customConfig
	}
}
func nginxExecCommand(args ...string) *exec.Cmd {
	var cmdArgs []string
	cmdArgs = append(cmdArgs, "-c", defaultNginxConf)
	cmdArgs = append(cmdArgs, args...)
	return exec.Command(nginxBinary, cmdArgs...)
}
