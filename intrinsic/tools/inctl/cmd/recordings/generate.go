// Copyright 2023 Intrinsic Innovation LLC

package recordings

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	pb "intrinsic/logging/proto/bag_packager_service_go_grpc_proto"
	"intrinsic/tools/inctl/util/orgutil"
)

var generateRecordingE = func(cmd *cobra.Command, _ []string) error {
	client, err := newBagPackagerClient(cmd.Context())
	if err != nil {
		return err
	}
	req := &pb.GenerateBagRequest{
		Query: &pb.GenerateBagRequest_BagId{
			BagId: flagBagID,
		},
		OrganizationId: cmdFlags.GetString(orgutil.KeyOrganization),
	}
	resp, err := client.GenerateBag(cmd.Context(), req)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return fmt.Errorf("recording with id \"%s\" does not exist", flagBagID)
		}
		return err
	}

	fmt.Println(fmt.Sprintf("Generated ROS bag for ID %s", resp.GetBag().GetBagMetadata().GetBagId()))

	return nil
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates a ROS bag for a given recording id",
	Long:  "Generates a ROS bag for a given recording id",
	Args:  cobra.NoArgs,
	RunE:  generateRecordingE,
}

func init() {
	recordingsCmd.AddCommand(generateCmd)
	flags := generateCmd.Flags()

	flags.StringVar(&flagBagID, "recording_id", "", "The recording id to generate ROS bag for.")
}
