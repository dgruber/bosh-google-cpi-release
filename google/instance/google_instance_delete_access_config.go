package ginstance

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/google/util"
)

func (i GoogleInstanceService) DeleteAccessConfig(id string, zone string, networkInterface string, accessConfig string) error {
	i.logger.Debug(googleInstanceServiceLogTag, "Deleting access config for Google Instance '%s'", id)
	operation, err := i.computeService.Instances.DeleteAccessConfig(i.project, gutil.ResourceSplitter(zone), id, accessConfig, networkInterface).Do()
	if err != nil {
		return bosherr.WrapErrorf(err, "Failed to delete access config for Google Instance '%s'", id)
	}

	if _, err = i.operationService.Waiter(operation, zone, ""); err != nil {
		return bosherr.WrapErrorf(err, "Failed to delete access config for Google Instance '%s'", id)
	}

	return nil
}
