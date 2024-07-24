package service

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/viveknathani/nattukaka/types"
	"github.com/viveknathani/nattukaka/utils"
)

// CreateVolume creates a new volume
func (srv *Service) CreateVolume() (int, string, *types.Volume) {
	publicID, err := utils.GeneratePublicId("volume")
	if err != nil {
		return fiber.StatusInternalServerError, err.Error(), nil
	}
	id, err := srv.Db.InsertVolume(publicID)
	return fiber.StatusCreated, "", &types.Volume{
		ID:       id,
		PublicID: publicID,
	}
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

	fmt.Println(service, volume)
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

// DeleteVolume deletes a volume and its associated service volumes
func (srv *Service) DeleteVolume(publicID string) error {
	return srv.Db.DeleteVolume(publicID)
}
