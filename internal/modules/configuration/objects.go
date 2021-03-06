/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package configuration

import "github.com/kubenext/lissio/internal/describer"

var (
	pluginDescriber = &PluginListDescriber{}

	rootDescriber = describer.NewSection(
		"/",
		"Configuration",
		pluginDescriber,
	)
)
