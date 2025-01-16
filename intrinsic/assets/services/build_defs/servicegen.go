// Copyright 2023 Intrinsic Innovation LLC

// Package servicegen implements creation of the service type bundle.
package servicegen

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	anypb "google.golang.org/protobuf/types/known/anypb"
	"intrinsic/assets/bundleio"
	"intrinsic/assets/idutils"
	smpb "intrinsic/assets/services/proto/service_manifest_go_proto"
	"intrinsic/util/proto/protoio"
	"intrinsic/util/proto/registryutil"
	"intrinsic/util/proto/sourcecodeinfoview"
)

// ServiceData holds the data needed to create a service bundle.
type ServiceData struct {
	// Optional path to default config proto.
	DefaultConfig string
	// Paths to binary file descriptor set protos to be used to resolve the configuration and behavior tree messages.
	FileDescriptorSets []string
	// Paths to tar archives for images.
	ImageTars []string
	// The deserialized ServiceManifest.
	Manifest *smpb.ServiceManifest
	// Bundle tar path.
	OutputBundle string
}

func validateManifest(m *smpb.ServiceManifest) error {
	if err := idutils.ValidateIDProto(m.GetMetadata().GetId()); err != nil {
		return fmt.Errorf("invalid name or package: %v", err)
	}
	if m.GetMetadata().GetVendor().GetDisplayName() == "" {
		return fmt.Errorf("vendor.display_name must be specified")
	}
	if m.GetServiceDef() != nil && m.GetServiceDef().GetSimSpec() == nil {
		return fmt.Errorf("a sim_spec must be specified if a service_def is provided;  see go/intrinsic-specifying-sim for more information")
	}
	return nil
}

func setDifference(slice1, slice2 []string) []string {
	var difference []string
	for _, val := range slice1 {
		if !slices.Contains(slice2, val) {
			difference = append(difference, val)
		}
	}
	return difference
}

// validateImageTars validates the provided images from the BUILD rule match the correct
// images specified in the manifest.
func validateImageTars(manifest *smpb.ServiceManifest, imgTarsList []string) error {
	var imagesInManifest []string
	if name := manifest.GetServiceDef().GetSimSpec().GetImage().GetArchiveFilename(); name != "" {
		imagesInManifest = append(imagesInManifest, name)
	}
	if name := manifest.GetServiceDef().GetRealSpec().GetImage().GetArchiveFilename(); name != "" {
		imagesInManifest = append(imagesInManifest, name)
	}
	basenameImageTarsList := []string{}
	for _, val := range imgTarsList {
		basenameImageTarsList = append(basenameImageTarsList, filepath.Base(val))
	}
	if diff := setDifference(basenameImageTarsList, imagesInManifest); len(diff) != 0 {
		return fmt.Errorf("images listed in the BUILD rule are not provided in the manifest: %v", diff)
	}
	if diff := setDifference(imagesInManifest, basenameImageTarsList); len(diff) != 0 {
		return fmt.Errorf("images listed in the manifest are not provided in the BUILD rule: %v", diff)
	}
	return nil
}

func pruneSourceCodeInfo(defaultConfig *anypb.Any, fds *dpb.FileDescriptorSet) error {
	if fds == nil {
		return nil
	}

	var fullNames []string
	if defaultConfig != nil {
		typeURLParts := strings.Split(defaultConfig.GetTypeUrl(), "/")
		if len(typeURLParts) < 1 {
			return fmt.Errorf("cannot extract default proto name from type URL: %v", defaultConfig.GetTypeUrl())
		}
		fullNames = append(fullNames, typeURLParts[len(typeURLParts)-1])
	}

	// Note that a nil default config will cause all source code info fields to be
	// stripped out.
	sourcecodeinfoview.PruneSourceCodeInfo(fullNames, fds)
	return nil
}

// CreateService bundles the data needed for software services.
func CreateService(d *ServiceData) error {
	if err := validateManifest(d.Manifest); err != nil {
		return fmt.Errorf("invalid manifest: %v", err)
	}

	set, err := registryutil.LoadFileDescriptorSets(d.FileDescriptorSets)
	if err != nil {
		return fmt.Errorf("unable to build FileDescriptorSet: %v", err)
	}

	types, err := registryutil.NewTypesFromFileDescriptorSet(set)
	if err != nil {
		return fmt.Errorf("failed to populate the registry: %v", err)
	}

	var defaultConfig *anypb.Any
	if d.DefaultConfig != "" {
		defaultConfig = &anypb.Any{}
		if err := protoio.ReadTextProto(d.DefaultConfig, defaultConfig, protoio.WithResolver(types)); err != nil {
			return fmt.Errorf("failed to read default config proto: %v", err)
		}
	}

	if err := validateImageTars(d.Manifest, d.ImageTars); err != nil {
		return fmt.Errorf("unable to retrieve image tars: %v", err)
	}

	if err := pruneSourceCodeInfo(defaultConfig, set); err != nil {
		return fmt.Errorf("unable to process source code info: %v", err)
	}
	if err := bundleio.WriteService(d.OutputBundle, bundleio.WriteServiceOpts{
		Manifest:    d.Manifest,
		Descriptors: set,
		Config:      defaultConfig,
		ImageTars:   d.ImageTars,
	}); err != nil {
		return fmt.Errorf("unable to write service bundle: %v", err)
	}

	return nil
}
