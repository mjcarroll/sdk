// Copyright 2023 Intrinsic Innovation LLC

package customer

import (
	"fmt"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	pb "intrinsic/kubernetes/accounts/service/api/accesscontrol/v1/accesscontrolv1_go_grpc_proto"
)

func init() {
	addUserInit(customerCmd)
}

var (
	flagEmail        string
	flagOrganization string
	flagRoleCSV      string
)

func parseCSV(s string) []string {
	parts := strings.Split(s, ",")
	for idx := range parts {
		parts[idx] = strings.TrimSpace(parts[idx])
	}
	return slices.DeleteFunc(parts, func(p string) bool { return p == "" })
}

func addUserInit(root *cobra.Command) {
	addUser.Flags().StringVar(&flagEmail, "email", "", "The email address of the user to invite.")
	addUser.Flags().StringVar(&flagOrganization, "organization", "", "The organization to invite the user to.")
	addUser.Flags().StringVar(&flagRoleCSV, "roles", "", "Optional comma-separated list of roles to assign to the user when they accept the invitation.")
	addUser.MarkFlagRequired("email")
	addUser.MarkFlagRequired("organization")
	root.AddCommand(addUser)
	pendingUsers.Flags().StringVar(&flagOrganization, "organization", "", "The organization to list pending invitations for.")
	pendingUsers.MarkFlagRequired("organization")
	root.AddCommand(pendingUsers)
}

var addUserHelp = `
Invite a user to an organization by email address.

Use the --roles flag (comma-separated list) to assign roles to the user after they accept the invitation.

Example:

		inctl customer add-user --email=user@example.com --organization=myorg --roles=owner
`

var addUser = &cobra.Command{
	Use:   "add-user",
	Short: "Invite a user to an organization by email address.",
	Long:  addUserHelp,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := withOrgID(cmd.Context())
		cl, err := newAccessControlV1Client(ctx)
		if err != nil {
			return err
		}
		req := pb.CreateOrganizationInvitationRequest{
			Parent: addPrefix(flagOrganization, "organizations/"),
			Invitation: &pb.OrganizationInvitation{
				Organization: flagOrganization,
				Email:        flagEmail,
				Roles:        addPrefixes(parseCSV(flagRoleCSV), "roles/"),
			},
		}
		if flagDebugRequests {
			protoPrint(&req)
		}
		op, err := cl.CreateOrganizationInvitation(ctx, &req)
		if err != nil {
			return fmt.Errorf("failed to create organization: %w", err)
		}
		if flagDebugRequests {
			protoPrint(op)
		}
		return nil
	},
}

var pendingUsersHelp = `
List all pending invitations for an organization.

		inctl customer pending-users --organization=myorg
`

var pendingUsers = &cobra.Command{
	Use:   "pending-users",
	Short: "List all pending invitations for an organization.",
	Long:  pendingUsersHelp,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := withOrgID(cmd.Context())
		cl, err := newAccessControlV1Client(ctx)
		if err != nil {
			return err
		}
		req := pb.ListOrganizationInvitationsRequest{
			Parent: addPrefix(flagOrganization, "organizations/"),
		}
		if flagDebugRequests {
			protoPrint(&req)
		}
		op, err := cl.ListOrganizationInvitations(ctx, &req)
		if err != nil {
			return fmt.Errorf("failed to list organization invitations: %w", err)
		}
		if flagDebugRequests {
			protoPrint(op)
		}
		for _, invitation := range op.GetInvitations() {
			fmt.Printf("%s\n", invitation.GetEmail())
		}
		return nil
	},
}
