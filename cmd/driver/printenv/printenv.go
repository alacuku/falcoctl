// SPDX-License-Identifier: Apache-2.0
// Copyright (C) 2023 The Falco Authors
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

package driverprintenv

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/net/context"

	driverdistro "github.com/falcosecurity/falcoctl/pkg/driver/distro"
	driverkernel "github.com/falcosecurity/falcoctl/pkg/driver/kernel"
	"github.com/falcosecurity/falcoctl/pkg/options"
)

type driverPrintenvOptions struct {
	*options.Common
	*options.Driver
}

// NewDriverPrintenvCmd print info about driver falcoctl config as env vars.
func NewDriverPrintenvCmd(ctx context.Context, opt *options.Common, driver *options.Driver) *cobra.Command {
	o := driverPrintenvOptions{
		Common: opt,
		Driver: driver,
	}

	cmd := &cobra.Command{
		Use:                   "printenv [flags]",
		DisableFlagsInUseLine: true,
		Short:                 "[Preview] Print env vars",
		Long: `[Preview] Print variables used by driver as env vars.
** This command is in preview and under development. **`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.RunDriverPrintenv(ctx)
		},
	}
	return cmd
}

func (o *driverPrintenvOptions) RunDriverPrintenv(_ context.Context) error {
	o.Printer.DefaultText.Printf("DRIVER=%q\n", o.Driver.Type.String())
	o.Printer.DefaultText.Printf("DRIVERS_REPO=%q\n", strings.Join(o.Driver.Repos, ", "))
	o.Printer.DefaultText.Printf("DRIVER_VERSION=%q\n", o.Driver.Version)
	o.Printer.DefaultText.Printf("DRIVER_NAME=%q\n", o.Driver.Name)
	o.Printer.DefaultText.Printf("HOST_ROOT=%q\n", o.Driver.HostRoot)

	kr, err := driverkernel.FetchInfo("", "")
	if err != nil {
		return err
	}

	d, err := driverdistro.Discover(kr, o.Driver.HostRoot)
	if err != nil {
		if !errors.Is(err, driverdistro.ErrUnsupported) {
			return err
		}
	}
	o.Printer.DefaultText.Printf("TARGET_ID=%q\n", d.String())

	o.Printer.DefaultText.Printf("ARCH=%q\n", kr.Architecture.ToNonDeb())
	o.Printer.DefaultText.Printf("KERNEL_RELEASE=%q\n", kr.String())
	o.Printer.DefaultText.Printf("KERNEL_VERSION=%q\n", kr.KernelVersion)

	fixedKr := d.FixupKernel(kr)
	o.Printer.DefaultText.Printf("FIXED_KERNEL_RELEASE=%q\n", fixedKr.String())
	o.Printer.DefaultText.Printf("FIXED_KERNEL_VERSION=%q\n", fixedKr.KernelVersion)

	return nil
}
