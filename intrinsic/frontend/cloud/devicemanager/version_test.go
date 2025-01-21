// Copyright 2023 Intrinsic Innovation LLC

package version

import "testing"

const (
	inversionOS   = "xfa.20241221.RC01"
	uiOS          = "20241221.RC01"
	inversionBase = "20250121.rc00"
	uiBase        = "20250121.RC00"
)

var (
	osI2P   = TranslateOSInversionToUI
	osP2I   = TranslateOSUIToInversion
	baseI2P = TranslateBaseInversionToUI
	baseP2I = TranslateBaseUIToInversion
)

func TestInverse(t *testing.T) {
	if got, want := osP2I(osI2P(inversionOS)), inversionOS; got != want {
		t.Errorf("TranslateOSUIToInversion is not the inverse of TranslateOSInversionToUI. got: %q want: %q", got, want)
	}
	if got, want := baseP2I(baseI2P(inversionBase)), inversionBase; got != want {
		t.Errorf("TranslateBaseUIToInversion is not the inverse of TranslateBaseInversionToUI. got: %q want: %q", got, want)
	}
}

func TestOS(t *testing.T) {
	if got, want := osP2I(uiOS), inversionOS; got != want {
		t.Errorf("TranslateOSUIToInversion(%q)=%q, want %q", uiOS, got, want)
	}
	if got, want := osI2P(inversionOS), uiOS; got != want {
		t.Errorf("TranslateOSInversionToUI(%q)=%q, want %q", inversionOS, got, want)
	}
}

func TestBase(t *testing.T) {
	if got, want := baseP2I(uiBase), inversionBase; got != want {
		t.Errorf("TranslateBaseUIToInversion(%q)=%q, want %q", uiBase, got, want)
	}
	if got, want := baseI2P(inversionBase), uiBase; got != want {
		t.Errorf("TranslateBaseInversionToUI(%q)=%q, want %q", inversionBase, got, want)
	}
}
