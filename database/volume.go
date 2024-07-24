package database

import (
	"database/sql"

	"github.com/viveknathani/nattukaka/types"
)

// SQL statements as constants
const (
	statementInsertVolume        = `insert into volumes (public_id) values ($1) returning id`
	statementSelectVolumeByID    = `select id, public_id from volumes where public_id = $1`
	statementDeleteVolume        = `delete from volumes where public_id = $1`
	statementDeleteServiceVolume = `delete from service_volumes where volume_id = $1`
	statementInsertServiceVolume = `insert into service_volumes (service_id, volume_id) values ($1, $2)`
	statementDetachServiceVolume = `delete from service_volumes where service_id = $1 and volume_id = $2`
)

// InsertVolume inserts a new volume into the database
func (db *Database) InsertVolume(publicID string) (int, error) {
	var insertedID = -1

	err := db.query(statementInsertVolume, func(rows *sql.Rows) error {
		if rows.Next() {
			err := rows.Scan(&insertedID)
			if err != nil {
				return err
			}
		}
		return nil
	}, publicID)
	if err != nil {
		return -1, err
	}

	return insertedID, err
}

// GetVolumeByID retrieves a volume by its public ID
func (db *Database) GetVolumeByID(publicID string) (*types.Volume, error) {
	var volume types.Volume
	err := db.query(statementSelectVolumeByID, func(rows *sql.Rows) error {
		if rows.Next() {
			err := rows.Scan(&volume.ID, &volume.PublicID)
			if err != nil {
				return err
			}
		}
		return nil
	}, publicID)
	if err != nil {
		return nil, err
	}
	return &volume, nil
}

// DeleteVolume deletes a volume and its associated service volumes
func (db *Database) DeleteVolume(publicID string) error {
	volume, err := db.GetVolumeByID(publicID)
	if err != nil {
		return err
	}

	// Delete associated service volumes
	_, err = db.pool.Exec(statementDeleteServiceVolume, volume.ID)
	if err != nil {
		return err
	}

	// Delete volume
	_, err = db.pool.Exec(statementDeleteVolume, publicID)
	return err
}

// AttachVolumeToService associates a volume with a service
func (db *Database) AttachVolumeToService(serviceID int, volumeID int) error {
	_, err := db.pool.Exec(statementInsertServiceVolume, serviceID, volumeID)
	return err
}

// DetachVolumeFromService disassociates a volume from a service
func (db *Database) DetachVolumeFromService(serviceID int, volumeID int) error {
	_, err := db.pool.Exec(statementDetachServiceVolume, serviceID, volumeID)
	return err
}
