// Copyright 2023 Intrinsic Innovation LLC

// Package assetscmd contains the root command for the assets command.
package assetscmd

import (
	"github.com/spf13/cobra"
	"intrinsic/assets/inctl/listreleased"
	"intrinsic/assets/inctl/listreleasedversions"
	"intrinsic/tools/inctl/cmd/root"
)

var assetsCmd = &cobra.Command{
	Use:   root.AssetsCmdName,
	Short: "Manages assets",
	Long:  "Manages assets",
}

func init() {
	assetsCmd.AddCommand(listreleased.GetCommand())
	assetsCmd.AddCommand(listreleasedversions.GetCommand())

	root.RootCmd.AddCommand(assetsCmd)
}
