// Copyright 2023 Intrinsic Innovation LLC

package environments

import (
	"testing"
)

func TestAccountsProjectFromProject(t *testing.T) {
	tests := []struct {
		projects []string
		want     string
	}{
		{
			want: AccountsProjectDev,
			projects: []string{
				"intrinsic-portal-dev",
				"intrinsic-accounts-dev",
				"intrinsic-assets-dev",
			},
		},
		{
			want: AccountsProjectStaging,
			projects: []string{
				"intrinsic-portal-staging",
				"intrinsic-accounts-staging",
				"intrinsic-assets-staging",
			},
		},
		{
			want: AccountsProjectProd,
			projects: []string{
				"intrinsic-portal-prod",
				"intrinsic-accounts-prod",
				"intrinsic-assets-prod",
			},
		},
	}

	for _, tc := range tests {
		for _, project := range tc.projects {
			got := AccountsProjectFromProject(project)
			if got != tc.want {
				t.Errorf("AccountsProjectFromProject(%q) = %q, want %q", project, got, tc.want)
			}
		}
	}
}
