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

package conversion

import (
	"github.com/goodrain/rainbond/db"
	v1 "github.com/goodrain/rainbond/worker/appm/types/v1"
)

func init() {
	//first conv service source
	RegistConversion(ServiceSource)
	//step2 conv service base
	RegistConversion(TenantServiceBase)
	//step3 conv service pod base info
	RegistConversion(TenantServiceVersion)
	//step4 conv service plugin
	RegistConversion(TenantServicePlugin)
	//step5 conv service inner and outer regist
	RegistConversion(TenantServiceRegist)
}

//Conversion conversion function
//Any application attribute implementation is similarly injected
type Conversion func(*v1.AppService, db.Manager) error

//conversionList conversion function list
var conversionList []Conversion

//RegistConversion regist conversion function list
func RegistConversion(fun Conversion) {
	conversionList = append(conversionList, fun)
}

//InitAppService init a app service
func InitAppService(dbmanager db.Manager, serviceID string) (*v1.AppService, error) {
	appService := &v1.AppService{
		AppServiceBase: v1.AppServiceBase{
			ServiceID: serviceID,
		},
		UpgradePatch: make(map[string][]byte, 2),
	}
	for _, c := range conversionList {
		if err := c(appService, dbmanager); err != nil {
			return nil, err
		}
	}
	return appService, nil
}

//InitCacheAppService init cache app service.
//if store manager receive a kube model belong with service and not find in store,will create
func InitCacheAppService(dbmanager db.Manager, serviceID, version, createrID string) (*v1.AppService, error) {
	appService := &v1.AppService{
		AppServiceBase: v1.AppServiceBase{
			ServiceID:    serviceID,
			CreaterID:    createrID,
			ExtensionSet: make(map[string]string),
		},
		UpgradePatch: make(map[string][]byte, 2),
	}
	if err := TenantServiceBase(appService, dbmanager); err != nil {
		return nil, err
	}
	return appService, nil
}
