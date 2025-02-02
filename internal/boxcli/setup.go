// Copyright 2022 Jetpack Technologies Inc and contributors. All rights reserved.
// Use of this source code is governed by the license in the LICENSE file.

package boxcli

import (
	"os"

	"github.com/fatih/color"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"go.jetpack.io/devbox/internal/nix"
	"go.jetpack.io/devbox/internal/ux"
)

const nixDaemonFlag = "daemon"

func setupCmd() *cobra.Command {
	setupCommand := &cobra.Command{
		Use:    "setup",
		Short:  "Setup devbox dependencies",
		Hidden: true,
	}

	installNixCommand := &cobra.Command{
		Use:   "nix",
		Short: "Install Nix",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInstallNixCmd(cmd)
		},
	}

	installNixCommand.Flags().Bool(nixDaemonFlag, false, "Install Nix in multi-user mode.")
	setupCommand.AddCommand(installNixCommand)
	return setupCommand
}

func runInstallNixCmd(cmd *cobra.Command) error {
	if nix.BinaryInstalled() {
		color.New(color.FgYellow).Fprint(
			cmd.ErrOrStderr(),
			"Nix is already installed. If this is incorrect please remove the "+
				"nix-shell binary from your path.\n",
		)
		return nil
	}
	return nix.Install(cmd.ErrOrStderr(), nixDaemonFlagVal(cmd))
}

func ensureNixInstalled(cmd *cobra.Command, _args []string) error {
	return nix.EnsureNixInstalled(cmd.ErrOrStderr(), nixDaemonFlagVal(cmd))
}

func nixDaemonFlagVal(cmd *cobra.Command) *bool {
	if !cmd.Flags().Changed(nixDaemonFlag) {
		if os.Geteuid() == 0 {
			ux.Fwarning(
				cmd.ErrOrStderr(),
				"Running as root. Installing Nix in multi-user mode.\n",
			)
			return lo.ToPtr(true)
		}
		return nil
	}

	val, err := cmd.Flags().GetBool(nixDaemonFlag)
	if err != nil {
		return nil
	}
	return &val
}
