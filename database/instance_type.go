package database

import (
	"database/sql"

	"github.com/viveknathani/nattukaka/types"
)

// SQL statement to fetch all instance types
const statementSelectAllInstanceTypes = `select id, public_id, name, cpu, memory, disk from instance_types`

// GetAllInstanceTypes fetches all instance types from the database
func (db *Database) GetAllInstanceTypes() ([]types.InstanceType, error) {
	var instanceTypes []types.InstanceType

	err := db.query(statementSelectAllInstanceTypes, func(rows *sql.Rows) error {
		for rows.Next() {
			var instanceType types.InstanceType
			err := rows.Scan(&instanceType.ID, &instanceType.PublicID, &instanceType.Name, &instanceType.CPU, &instanceType.Memory, &instanceType.Disk)
			if err != nil {
				return err
			}
			instanceTypes = append(instanceTypes, instanceType)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return instanceTypes, nil
}
