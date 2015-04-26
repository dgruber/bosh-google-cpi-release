package ginstance

import (
	"reflect"
	"sort"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/api"
	"google.golang.org/api/compute/v1"
)

func (i GoogleInstanceService) ConfigureNetworks(id string, instanceNetworks GoogleInstanceNetworks) error {
	instance, found, err := i.Find(id, "")
	if err != nil {
		return err
	}
	if !found {
		return bosherr.Errorf("Google Instance '%s' not found", id)
	}

	// TODO: Configure VIP network

	if err := i.configureTargetPool(instance, instanceNetworks); err != nil {
		return err
	}

	return nil
}

func (i GoogleInstanceService) UpdateNetworks(id string, instanceNetworks GoogleInstanceNetworks) error {
	instance, found, err := i.Find(id, "")
	if err != nil {
		return err
	}
	if !found {
		return bosherr.Errorf("Google Instance '%s' not found", id)
	}

	if err = i.updateNetwork(instance, instanceNetworks); err != nil {
		return err
	}

	if err = i.updateIpForwarding(instance, instanceNetworks); err != nil {
		return err
	}

	if err = i.updateEphemeralExternalIp(instance, instanceNetworks); err != nil {
		return err
	}

	if err = i.updateTags(instance, instanceNetworks); err != nil {
		return err
	}

	i.ConfigureNetworks(id, instanceNetworks)

	return nil
}

func (i GoogleInstanceService) configureTargetPool(instance *compute.Instance, instanceNetworks GoogleInstanceNetworks) error {
	targetPoolName := instanceNetworks.TargetPool()
	if targetPoolName == "" {
		return nil
	}

	targetPool, found, err := instanceNetworks.targetPoolService.Find(targetPoolName, "")
	if err != nil {
		return err
	}
	if !found {
		return bosherr.WrapErrorf(err, "Google Target Pool '%s' does not exists", targetPoolName)
	}

	for _, tpInstance := range targetPool.Instances {
		if tpInstance == instance.SelfLink {
			// Instance already attached to the target pool
			return nil
		}
	}

	i.logger.Debug(googleInstanceServiceLogTag, "Attaching Google Instance '%s' to Google Target Pool '%s'", instance.Name, targetPoolName)
	err = instanceNetworks.targetPoolService.AddInstance(targetPoolName, instance.Name)
	if err != nil {
		return err
	}

	return nil
}

func (i GoogleInstanceService) updateNetwork(instance *compute.Instance, instanceNetworks GoogleInstanceNetworks) error {
	// If the network has changed we need to recreate the VM
	dynamicNetwork := instanceNetworks.DynamicNetwork()
	if instance.NetworkInterfaces[0].Network != dynamicNetwork.NetworkName {
		return api.NotSupportedError{}
	}

	return nil
}

func (i GoogleInstanceService) updateIpForwarding(instance *compute.Instance, instanceNetworks GoogleInstanceNetworks) error {
	// If IP Forwarding has changed we need to recreate the VM
	if instance.CanIpForward != instanceNetworks.CanIpForward() {
		return api.NotSupportedError{}
	}

	return nil
}

func (i GoogleInstanceService) updateEphemeralExternalIp(instance *compute.Instance, instanceNetworks GoogleInstanceNetworks) error {
	var instanceExternalIp string

	if len(instance.NetworkInterfaces[0].AccessConfigs) > 0 {
		instanceExternalIp = instance.NetworkInterfaces[0].AccessConfigs[0].NatIP
	}

	if instanceNetworks.EphemeralExternalIP() {
		if instanceExternalIp == "" {
			networkInterface := instance.NetworkInterfaces[0].Name
			accessConfig := &compute.AccessConfig{
				Name: "External NAT",
				Type: "ONE_TO_ONE_NAT",
			}
			err := i.AddAccessConfig(instance.Name, instance.Zone, networkInterface, accessConfig)
			if err != nil {
				return err
			}
		}
	} else {
		if instanceExternalIp != "" {
			// TODO: Only if network has no vip
			networkInterface := instance.NetworkInterfaces[0].Name
			accessConfig := instance.NetworkInterfaces[0].AccessConfigs[0].Name
			err := i.DeleteAccessConfig(instance.Name, instance.Zone, networkInterface, accessConfig)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (i GoogleInstanceService) updateTags(instance *compute.Instance, instanceNetworks GoogleInstanceNetworks) error {
	// Parset network tags
	networkTags, err := instanceNetworks.Tags()
	if err != nil {
		return err
	}

	// Check if tags have changed
	sort.Strings(networkTags.Items)
	sort.Strings(instance.Tags.Items)
	if reflect.DeepEqual(networkTags.Items, instance.Tags.Items) {
		return nil
	}

	// Override the instance tags preserving the original fingerprint
	instanceTags := &compute.Tags{
		Fingerprint: instance.Tags.Fingerprint,
		Items:       networkTags.Items,
	}

	// Update the instance tags
	err = i.SetTags(instance.Name, instance.Zone, instanceTags)
	if err != nil {
		return err
	}

	return nil
}
