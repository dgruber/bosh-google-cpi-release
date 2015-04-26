package ginstance

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/google/util"
	"google.golang.org/api/compute/v1"
)

func (i GoogleInstanceService) SetMetadata(id string, vmMetadata GoogleInstanceMetadata) error {
	// Find the instance
	instance, found, err := i.Find(id, "")
	if err != nil {
		return err
	}
	if !found {
		return bosherr.Errorf("Google Instance '%s' not found", id)
	}

	// We need to reuse the original instance metadata fingerprint and items
	metadata := instance.Metadata
	metadataMap := make(map[string]string)

	// Grab the original metadata items
	for _, item := range metadata.Items {
		metadataMap[item.Key] = item.Value
	}

	// Add or override the new metadata items
	for key, value := range vmMetadata {
		metadataMap[key] = value.(string)
	}

	// Set the new metadata items
	var metadataItems []*compute.MetadataItems
	for key, value := range metadataMap {
		metadataItems = append(metadataItems, &compute.MetadataItems{Key: key, Value: value})
	}
	metadata.Items = metadataItems

	i.logger.Debug(googleInstanceServiceLogTag, "Setting metadata for Google Instance '%s'", id)
	operation, err := i.computeService.Instances.SetMetadata(i.project, gutil.ResourceSplitter(instance.Zone), id, metadata).Do()
	if err != nil {
		return bosherr.WrapErrorf(err, "Failed to set metadata for Google Instance '%s'", id)
	}

	if _, err = i.operationService.Waiter(operation, instance.Zone, ""); err != nil {
		return bosherr.WrapErrorf(err, "Failed to set metadata for Google Instance '%s'", id)
	}

	return nil
}
