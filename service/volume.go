package service

import (
	"github.com/viveknathani/nattukaka/types"
)

// CreateVolume creates a new volume
func (srv *Service) CreateVolume(publicID string) error {
	return srv.Db.InsertVolume(&types.Volume{PublicID: publicID})
}

// AttachVolumeToService attaches a volume to a service
func (srv *Service) AttachVolumeToService(servicePublicID, volumePublicID string) error {
	service, err := srv.Db.GetServiceByID(servicePublicID)
	if err != nil {
		return err
	}

	volume, err := srv.Db.GetVolumeByID(volumePublicID)
	if err != nil {
		return err
	}

	return srv.Db.AttachVolumeToService(service.ID, volume.ID)
}

// DetachVolumeFromService detaches a volume from a service
func (srv *Service) DetachVolumeFromService(servicePublicID, volumePublicID string) error {
	service, err := srv.Db.GetServiceByID(servicePublicID)
	if err != nil {
		return err
	}

	volume, err := srv.Db.GetVolumeByID(volumePublicID)
	if err != nil {
		return err
	}

	return srv.Db.DetachVolumeFromService(service.ID, volume.ID)
}

// GetVolumeByID retrieves a volume by its public ID
func (srv *Service) GetVolumeByID(publicID string) (*types.Volume, error) {
	return srv.Db.GetVolumeByID(publicID)
}

// GetVolumesByWorkspace retrieves all volumes associated with a workspace
func (srv *Service) GetVolumesByWorkspace(workspaceID int, page int) ([]types.Volume, error) {
	return srv.Db.GetVolumesByWorkspace(workspaceID, page)
}

// DeleteVolume deletes a volume and its associated service volumes
func (srv *Service) DeleteVolume(publicID string) error {
	return srv.Db.DeleteVolume(publicID)
}
