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

package v1

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/goodrain/rainbond/util"

	v1 "k8s.io/api/apps/v1"
)

//SetUpgradePatch create and set upgrade pathch for deployment and statefulset
func (a *AppService) SetUpgradePatch(new *AppService) error {
	if a.statefulset != nil && new.statefulset != nil {
		statefulsetPatch, err := getStatefulsetModifiedConfiguration(a.statefulset, new.statefulset)
		if err != nil {
			return err
		}
		if len(statefulsetPatch) == 0 {
			return fmt.Errorf("no upgrade")
		}
		a.UpgradePatch["statefulset"] = statefulsetPatch
	}
	if a.deployment != nil && new.deployment != nil {
		deploymentPatch, err := getDeploymentModifiedConfiguration(a.deployment, new.deployment)
		if err != nil {
			return err
		}
		if len(deploymentPatch) == 0 {
			return fmt.Errorf("no upgrade")
		}
		a.UpgradePatch["deployment"] = deploymentPatch
	}
	return nil
}

//EncodeNode encode node
type EncodeNode struct {
	body  []byte
	value []byte
	Field map[string]EncodeNode
}

//UnmarshalJSON custom yaml decoder
func (e *EncodeNode) UnmarshalJSON(code []byte) error {
	e.body = code
	if len(code) < 1 {
		return nil
	}
	if code[0] != '{' {
		e.value = code
		return nil
	}
	var fields = make(map[string]EncodeNode)
	if err := json.Unmarshal(code, &fields); err != nil {
		return err
	}
	e.Field = fields
	return nil
}

//MarshalJSON custom marshal json
func (e *EncodeNode) MarshalJSON() ([]byte, error) {
	if e.value != nil {
		return e.value, nil
	}
	if e.Field != nil {
		var buffer = bytes.NewBufferString("{")
		count := 0
		length := len(e.Field)
		for k, v := range e.Field {
			buffer.WriteString(fmt.Sprintf("\"%s\":", k))
			value, err := v.MarshalJSON()
			if err != nil {
				return nil, err
			}
			buffer.Write(value)
			count++
			if count < length {
				buffer.WriteString(",")
			}
		}
		buffer.WriteByte('}')
		return buffer.Bytes(), nil
	}
	return nil, fmt.Errorf("marshal error")
}

//Contrast Compare value
func (e *EncodeNode) Contrast(endpoint *EncodeNode) bool {
	return util.BytesSliceEqual(e.value, endpoint.value)
}

//GetChange get change fields
func (e *EncodeNode) GetChange(endpoint *EncodeNode) *EncodeNode {
	if util.BytesSliceEqual(e.body, endpoint.body) {
		return nil
	}
	return getChange(*e, *endpoint)
}

func getChange(old, new EncodeNode) *EncodeNode {
	var result EncodeNode
	if util.BytesSliceEqual(old.body, new.body) {
		return nil
	}
	if old.Field == nil && new.Field == nil {
		if !util.BytesSliceEqual(old.value, new.value) {
			result.value = new.value
			return &result
		}
	}
	for k, v := range new.Field {
		if result.Field == nil {
			result.Field = make(map[string]EncodeNode)
		}
		if value := getChange(old.Field[k], v); value != nil {
			result.Field[k] = *value
		}
	}
	return &result
}

func getStatefulsetModifiedConfiguration(old, new *v1.StatefulSet) ([]byte, error) {
	old.Status = new.Status
	oldNeed := getAllowFields(old)
	newNeed := getAllowFields(new)
	return getchange(oldNeed, newNeed)
}

// updates to statefulset spec for fields other than 'replicas', 'template', and 'updateStrategy' are forbidden.
func getAllowFields(s *v1.StatefulSet) *v1.StatefulSet {
	return &v1.StatefulSet{
		Spec: v1.StatefulSetSpec{
			Replicas:       s.Spec.Replicas,
			Template:       s.Spec.Template,
			UpdateStrategy: s.Spec.UpdateStrategy,
		},
	}
}

func getDeploymentModifiedConfiguration(old, new *v1.Deployment) ([]byte, error) {
	old.Status = new.Status
	return getchange(old, new)
}

func getchange(old, new interface{}) ([]byte, error) {
	oldbuffer := bytes.NewBuffer(nil)
	newbuffer := bytes.NewBuffer(nil)
	err := json.NewEncoder(oldbuffer).Encode(old)
	if err != nil {
		return nil, fmt.Errorf("encode old body error %s", err.Error())
	}
	err = json.NewEncoder(newbuffer).Encode(new)
	if err != nil {
		return nil, fmt.Errorf("encode new body error %s", err.Error())
	}
	var en EncodeNode
	if err := json.NewDecoder(oldbuffer).Decode(&en); err != nil {
		return nil, err
	}
	var ennew EncodeNode
	if err := json.NewDecoder(newbuffer).Decode(&ennew); err != nil {
		return nil, err
	}
	change := en.GetChange(&ennew)
	changebody, err := json.Marshal(change)
	if err != nil {
		return nil, err
	}
	return changebody, nil
}
