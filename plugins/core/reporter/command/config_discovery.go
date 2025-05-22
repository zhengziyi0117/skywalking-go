// Licensed to Apache Software Foundation (ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation (ASF) licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package command

import (
	"github.com/apache/skywalking-go/plugins/core/reporter"
	common "skywalking.apache.org/repo/goapi/collect/common/v3"
)

const ConfigurationDiscoveryCommandName = "ConfigurationDiscoveryCommand"

func NewConfigDiscoveryService() *ConfigDiscoveryService {
	return &ConfigDiscoveryService{}
}

type ConfigDiscoveryService struct {
	UUID     string
	watchers map[string]reporter.AgentConfigChangeWatcher
}

func (s *ConfigDiscoveryService) BindWatchers(watchers []reporter.AgentConfigChangeWatcher) {
	// bind watchers
	s.watchers = make(map[string]reporter.AgentConfigChangeWatcher)
	for _, watcher := range watchers {
		s.watchers[watcher.Key()] = watcher
	}
}

func (s *ConfigDiscoveryService) HandleCommand(command *common.Command) {
	var uuid string
	var newConfigs = make(map[string]*common.KeyStringValuePair)
	for _, pair := range command.GetArgs() {
		if pair.Key == "SerialNumber" {
		} else if pair.Key == "UUID" {
			uuid = pair.Value
		} else {
			newConfigs[pair.Key] = pair
		}
	}

	// check same uuid
	if s.UUID == uuid {
		return
	}

	// notify to all watchers
	for key, watcher := range s.watchers {
		pair := newConfigs[key]
		if pair == nil || pair.Value == "" {
			watcher.Notify(reporter.DELETED, "")
		} else if pair.Value != watcher.Value() {
			watcher.Notify(reporter.MODIFY, pair.Value)
		}
	}

	// update uuid
	s.UUID = uuid
}
