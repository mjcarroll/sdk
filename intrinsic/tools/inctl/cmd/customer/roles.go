// Copyright 2023 Intrinsic Innovation LLC

package customer

import (
	"fmt"

	"github.com/spf13/cobra"

	accaccesscontrolv1pb "intrinsic/kubernetes/accounts/service/api/accesscontrol/v1/accesscontrol_go_grpc_proto"
	"intrinsic/tools/inctl/util/cobrautil"
)

var rolesCmd = cobrautil.ParentOfNestedSubcommands("roles", "List available roles.")

func init() {
	customerCmd.AddCommand(rolesCmd)
	rolesInit(rolesCmd)
}

func rolesInit(root *cobra.Command) {
	root.AddCommand(listRolesCmd)
}

var listRolesCmdHelp = `
List available roles.
`

var listRolesCmd = &cobra.Command{
	Use:   "list",
	Short: "List available roles.",
	Long:  listRolesCmdHelp,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		cl, err := newAccessControlV1Client(ctx)
		if err != nil {
			return err
		}
		rs, err := cl.ListRoles(ctx, &accaccesscontrolv1pb.ListRolesRequest{})
		if err != nil {
			return err
		}
		for _, r := range rs.GetRoles() {
			fmt.Printf("%s - %s - %s\n", r.GetName(), r.GetDisplayName(), r.GetDescription())
		}
		return nil
	},
}
