// SPDX-License-Identifier: Apache-2.0
// Copyright (C) 2024 The Falco Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

//go:build linux

package push

import (
	"fmt"
	"runtime"

	"github.com/falcosecurity/plugin-sdk-go/pkg/loader"
	"github.com/pterm/pterm"

	"github.com/falcosecurity/falcoctl/pkg/oci"
	"github.com/falcosecurity/falcoctl/pkg/options"
)

func pluginConfigLayer(logger *pterm.Logger, filePath, platform string, artifactOptions *options.Artifact) (*oci.ArtifactConfig, error) {
	config := &oci.ArtifactConfig{
		Name:    artifactOptions.Name,
		Version: artifactOptions.Version,
	}

	// Parse the requirements.
	// Check if the user has provided any.
	if len(artifactOptions.Requirements) != 0 {
		logger.Info("Requirements provided by user", logger.Args("plugin", filePath))
		if err := config.ParseRequirements(artifactOptions.Requirements...); err != nil {
			return nil, err
		}
	} else {
		logger.Info("Parsing requirements from: ", logger.Args("plugin", filePath))
		sysPlatform := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
		// If no user provided requirements then try to parse them from the plugin.
		if platform != sysPlatform {
			logger.Info("Skipping, incompatible platform", logger.Args("plugin", platform, "current system", sysPlatform))
			return nil, nil
		}
		req, err := pluginRequirement(filePath)
		if err != nil {
			return nil, err
		}
		config.SetRequirement(req.Name, req.Version)
	}

	return config, nil
}

// pluginRequirement given a plugin as a shared library it loads it and gets the api version
// required by the plugin.
func pluginRequirement(filePath string) (*oci.ArtifactRequirement, error) {
	plugin, err := loader.NewPlugin(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open plugin %q: %w", filePath, err)
	}

	return &oci.ArtifactRequirement{
		Name:    pluginRequirementKey,
		Version: plugin.Info().RequiredAPIVersion,
	}, nil
}
