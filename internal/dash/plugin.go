/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package dash

import (
	"github.com/pkg/errors"

	"github.com/kubenext/lissio/internal/module"
	"github.com/kubenext/lissio/pkg/action"
	"github.com/kubenext/lissio/pkg/plugin"
	"github.com/kubenext/lissio/pkg/plugin/api"
)

func initPlugin(moduleManager module.ManagerInterface, actionManager *action.Manager, service api.Service) (*plugin.Manager, error) {
	apiService, err := api.New(service)
	if err != nil {
		return nil, errors.Wrap(err, "create dashboard api")
	}

	m := plugin.NewManager(apiService, moduleManager, actionManager)

	pluginList, err := plugin.AvailablePlugins(plugin.DefaultConfig)
	if err != nil {
		return nil, errors.Wrap(err, "finding available plugins")
	}

	for _, pluginPath := range pluginList {
		if err := m.Load(pluginPath); err != nil {
			return nil, errors.Wrapf(err, "initialize plugin %q", pluginPath)
		}

	}

	return m, nil
}
