// Copyright 2023 Intrinsic Innovation LLC

package customer

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
	grpccredentials "google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	"intrinsic/config/environments"
	accaccesscontrolv1grpcpb "intrinsic/kubernetes/accounts/service/api/accesscontrol/v1/accesscontrolv1_go_grpc_proto"
	accresourcemanagerv1grpcpb "intrinsic/kubernetes/accounts/service/api/resourcemanager/v1/resourcemanagerv1_go_grpc_proto"
	"intrinsic/tools/inctl/auth/auth"
	"intrinsic/tools/inctl/util/orgutil"
)

func authFromVipr() (string, string) {
	authOrg := vipr.GetString(orgutil.KeyOrganization)
	authProject := vipr.GetString(orgutil.KeyProject)
	authEnv := vipr.GetString(orgutil.KeyEnvironment)
	org := authOrg
	if authProject != "" {
		org = authOrg + "@" + authProject
	}
	return authEnv, org
}

// Aliases for convenience
var newResourceManagerV1Client = func(ctx context.Context) (resourceManagerV1Client, error) {
	env, org := authFromVipr()
	return newSecureAccountsResourceManagerAPIKeyClient(ctx, env, org)
}

var newAccessControlV1Client = func(ctx context.Context) (accessControlV1Client, error) {
	env, org := authFromVipr()
	return newSecureAccountsAccessControlAPIKeyClient(ctx, env, org)
}

// Aliases for convenience
type resourceManagerV1Client = accresourcemanagerv1grpcpb.ResourceManagerServiceClient
type accessControlV1Client = accaccesscontrolv1grpcpb.AccessControlServiceClient

func newSecureAccountsAccessControlAPIKeyClient(ctx context.Context, env, org string) (accaccesscontrolv1grpcpb.AccessControlServiceClient, error) {
	conn, err := newConnAuthStore(ctx, environments.AccountsDomain(env), org)
	if err != nil {
		return nil, err
	}
	return accaccesscontrolv1grpcpb.NewAccessControlServiceClient(conn), nil
}

// newSecureAccountsTokensServiceAPIKeyClient creates a new secure ResourceManagerClient using API keys.
// Suitable for calling the ressourcemanager via HTTPS from any environment.
func newSecureAccountsResourceManagerAPIKeyClient(ctx context.Context, env, org string) (accresourcemanagerv1grpcpb.ResourceManagerServiceClient, error) {
	conn, err := newConnAuthStore(ctx, environments.AccountsDomain(env), org)
	if err != nil {
		return nil, err
	}
	return accresourcemanagerv1grpcpb.NewResourceManagerServiceClient(conn), nil
}

// Can be overwridden/injected in tests.
var authStore = auth.NewStore()

func newConnAuthStore(ctx context.Context, addr, org string) (*grpc.ClientConn, error) {
	// determine address
	orgInfo, err := authStore.ReadOrgInfo(org)
	if err != nil {
		return nil, fmt.Errorf("failed to read org info for %q: %v", org, err)
	}
	// fetch API key
	project := orgInfo.Project
	cfg, err := auth.NewStore().GetConfiguration(project)
	if err != nil {
		return nil, fmt.Errorf("failed to get project configuration for project %q: %v", project, err)
	}
	creds, err := cfg.GetDefaultCredentials()
	if err != nil {
		return nil, fmt.Errorf("failed to get default credentials for project %q: %v", project, err)
	}
	return newConn(ctx, addr, grpc.WithPerRPCCredentials(creds))
}

func newConn(ctx context.Context, addr string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	// create connection
	var grpcOpts = []grpc.DialOption{
		grpc.WithStatsHandler(new(ocgrpc.ClientHandler)),
		grpc.WithTransportCredentials(grpccredentials.NewTLS(&tls.Config{})),
	}
	grpcOpts = append(grpcOpts, opts...)
	conn, err := grpc.NewClient(addr+":443", grpcOpts...)
	if err != nil {
		return nil, errors.Wrapf(err, "grpc.Dial(%q)", addr)
	}
	return conn, nil
}

const (
	// cookieHeaderName is the name of the header / metadata field used for cookies
	cookieHeaderName = "Cookie"
	// orgIDCookie is the cookie key for the organization identifier.
	orgIDCookie = "org-id"
)

// withOrgID adds the org ID to the outgoing RCP context.
func withOrgID(ctx context.Context) context.Context {
	o := vipr.GetString(orgutil.KeyOrganization)
	md := ToMDString(&http.Cookie{Name: orgIDCookie, Value: o})
	return metadata.AppendToOutgoingContext(ctx, md...)
}

// ToMDString converts a list of http.Cookie objects to a string that can be used as a metadata
// value.
func ToMDString(cs ...*http.Cookie) []string {
	cookiesKV := []string{}
	for _, c := range cs {
		cookiesKV = append(cookiesKV, (&http.Cookie{Name: c.Name, Value: c.Value}).String())
	}
	return []string{cookieHeaderName, strings.Join(cookiesKV, "; ")}
}
