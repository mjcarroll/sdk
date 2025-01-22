// Copyright 2023 Intrinsic Innovation LLC

// Package assetscmd contains the root command for the assets command.
package assetscmd

import (
	"github.com/spf13/cobra"
	"intrinsic/tools/inctl/cmd/root"
)

var assetsCmd = &cobra.Command{
	Use:   root.AssetsCmdName,
	Short: "Manages assets",
	Long:  "Manages assets",
}

func init() {

	root.RootCmd.AddCommand(assetsCmd)
}
