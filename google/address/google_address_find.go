package gaddress

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/google/util"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

func (a GoogleAddressService) Find(id string, region string) (*compute.Address, bool, error) {
	if region == "" {
		a.logger.Debug(googleAddressServiceLogTag, "Finding Google Address '%s'", id)
		filter := fmt.Sprintf("name eq .*%s", id)
		addresses, err := a.computeService.Addresses.AggregatedList(a.project).Filter(filter).Do()
		if err != nil {
			return &compute.Address{}, false, bosherr.WrapErrorf(err, "Failed to find Google Address '%s'", id)
		}

		for _, addressItems := range addresses.Items {
			for _, address := range addressItems.Addresses {
				// Return the first address (it can only be 1 address with the same name across all regions)
				return address, true, nil
			}
		}

		return &compute.Address{}, false, nil
	}

	a.logger.Debug(googleAddressServiceLogTag, "Finding Google Address '%s' in region '%s'", id, region)
	address, err := a.computeService.Addresses.Get(a.project, gutil.ResourceSplitter(region), id).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return &compute.Address{}, false, nil
		}

		return &compute.Address{}, false, bosherr.WrapErrorf(err, "Failed to find Google Address '%s' in region '%s'", id, region)
	}

	return address, true, nil
}
