package gdisktype

import (
	boshlog "github.com/cloudfoundry/bosh-agent/logger"

	"google.golang.org/api/compute/v1"
)

const googleDiskTypeServiceLogTag = "GoogleDiskTypeService"

type GoogleDiskTypeService struct {
	project        string
	computeService *compute.Service
	logger         boshlog.Logger
}

func NewGoogleDiskTypeService(
	project string,
	computeService *compute.Service,
	logger boshlog.Logger,
) GoogleDiskTypeService {
	return GoogleDiskTypeService{
		project:        project,
		computeService: computeService,
		logger:         logger,
	}
}
