package service

import (
	"github.com/viveknathani/nattukaka/types"
)

// GetAllInstanceTypes fetches all instance types
func (srv *Service) GetAllInstanceTypes() ([]types.InstanceType, error) {
	instanceTypes, err := srv.Db.GetAllInstanceTypes()
	if err != nil {
		srv.Logger.Error(err.Error())
		return nil, err
	}
	return instanceTypes, nil
}
