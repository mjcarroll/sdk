// Copyright 2023 Intrinsic Innovation LLC

// Package listreleased defines the list_released command that lists assets in catalog.
package listreleased

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
	acgrpcpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_grpc_proto"
	acpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_grpc_proto"
	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"
	"intrinsic/assets/idutils"
	"intrinsic/assets/listutils"
	atpb "intrinsic/assets/proto/asset_type_go_proto"
	viewpb "intrinsic/assets/proto/view_go_proto"
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/printer"
)

const pageSize int64 = 50

// GetCommand returns a command to list released assets.
func GetCommand() *cobra.Command {
	flags := cmdutils.NewCmdFlags()
	cmd := &cobra.Command{
		Use:   "list_released",
		Short: "List assets from the catalog.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			conn, err := clientutils.DialCatalogFromInctl(cmd, flags)
			if err != nil {
				return fmt.Errorf("cannot create client connection: %w", err)
			}
			defer conn.Close()
			client := acgrpcpb.NewAssetCatalogClient(conn)
			prtr, err := printer.NewPrinter(root.FlagOutput)
			if err != nil {
				return err
			}

			at, err := flags.GetFlagAssetType()
			if err != nil {
				return err
			}

			assets, err := listutils.ListAllAssets(
				cmd.Context(),
				client,
				pageSize,
				viewpb.AssetViewType_ASSET_VIEW_TYPE_BASIC,
				&acpb.ListAssetsRequest_AssetFilter{
					AssetTypes:  []atpb.AssetType{at},
					OnlyDefault: proto.Bool(true),
				},
			)
			if err != nil {
				return err
			}
			idVersions := make([]string, len(assets))
			for i, asset := range assets {
				idVersion, err := idutils.IDVersionFromProto(asset.GetMetadata().GetIdVersion())
				if err != nil {
					return err
				}
				idVersions[i] = idVersion
			}
			sort.Strings(idVersions)
			prtr.Print(strings.Join(idVersions, "\n"))

			return nil
		},
	}
	flags.SetCommand(cmd)
	flags.AddFlagAssetType()

	return cmd
}
