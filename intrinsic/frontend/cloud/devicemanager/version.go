// Copyright 2023 Intrinsic Innovation LLC

// Package version provides functions to translate between inversion and UI versions
package version

import (
	"strings"
)

// TranslateOSInversionToUI translates an inversion style OS version to an UI style OS version
//
// This is the inverse of TranslateOSUIToInversion
func TranslateOSInversionToUI(inversion string) string {
	return strings.TrimPrefix(inversion, "xfa.")
}

// TranslateOSUIToInversion translates an UI style OS version to an inversion style OS version
//
// This is the inverse of TranslateOSInversionToUI
func TranslateOSUIToInversion(ui string) string {
	return "xfa." + ui
}

// TranslateBaseInversionToUI translates an inversion style base version to an UI style base version
//
// This is the inverse of TranslateBaseUIToInversion
func TranslateBaseInversionToUI(inversion string) string {
	return strings.ToUpper(inversion)
}

// TranslateBaseUIToInversion translates an UI style base version to an inversion style base version
//
// This is the inverse of TranslateBaseInversionToUI
func TranslateBaseUIToInversion(ui string) string {
	return strings.ToLower(ui)
}
