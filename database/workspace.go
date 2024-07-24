package database

import (
	"database/sql"

	"github.com/viveknathani/nattukaka/types"
)

// SQL statements as constants
const (
	statementInsertWorkspace                           = `insert into workspaces (public_id, name) values ($1, $2) returning id`
	statementInsertWorkspaceUser                       = `insert into workspace_users (public_id, workspace_id, user_id, role) values ($1, $2, $3, $4) returning id`
	statementDeleteWorkspaceUsers                      = `delete from workspace_users where workspace_id = $1`
	statementDeleteWorkspace                           = `delete from workspaces where public_id = $1`
	statementSelectUsersByWorkspace                    = `select u.id, u.name, u.email, u.public_id from users u join workspace_users wu on u.id = wu.user_id where wu.workspace_id = $1`
	statementSelectWorkspaceByPublicID                 = `select id, public_id, name from workspaces where public_id = $1`
	statementSelectWorkspaceUserByUserIDAndWorkspaceID = `select id, public_id, workspace_id, user_id, role from workspace_users where user_id = $1 and workspace_id = (select id from workspaces where public_id = $2)`
	statementSelectWorkspacesByUserID                  = `select wu.id, wu.public_id, wu.workspace_id, wu.user_id, wu.role from workspace_users wu join workspaces w on wu.workspace_id = w.id where wu.user_id = $1`
)

// InsertWorkspace inserts a new workspace into the database
func (db *Database) InsertWorkspace(w *types.Workspace) (int, error) {
	var insertedID = -1

	err := db.query(statementInsertWorkspace, func(rows *sql.Rows) error {
		if rows.Next() {
			err := rows.Scan(&insertedID)
			if err != nil {
				return err
			}
		}
		return nil
	}, w.PublicID, w.Name)
	if err != nil {
		return -1, err
	}

	return insertedID, err
}

// InsertWorkspaceUser inserts a new user into a workspace
func (db *Database) InsertWorkspaceUser(wu *types.WorkspaceUser) error {
	err := db.execWithTransaction(statementInsertWorkspaceUser, wu.PublicID, wu.WorkspaceID, wu.UserID, wu.Role)
	return err
}

// DeleteWorkspaceUsers deletes all the users of a workspace
func (db *Database) DeleteWorkspaceUsers(workspaceID int) error {
	err := db.execWithTransaction(statementDeleteWorkspaceUsers, workspaceID)
	return err
}

// DeleteWorkspace deletes a workspace from the database
func (db *Database) DeleteWorkspace(publicID string) error {
	err := db.execWithTransaction(statementDeleteWorkspace, publicID)
	return err
}

// GetUsersByWorkspace gets all users in a workspace
func (db *Database) GetUsersByWorkspace(workspaceID int) ([]*types.User, error) {
	var users []*types.User
	err := db.query(statementSelectUsersByWorkspace, func(rows *sql.Rows) error {
		for rows.Next() {
			var u types.User
			err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.PublicID)
			if err != nil {
				return err
			}
			users = append(users, &u)
		}
		return nil
	}, workspaceID)

	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetWorkspaceByPublicID gets a workspace by its public ID
func (db *Database) GetWorkspaceByPublicID(publicID string) (*types.Workspace, error) {
	var w types.Workspace
	exists := false
	err := db.query(statementSelectWorkspaceByPublicID, func(rows *sql.Rows) error {
		if rows.Next() {
			err := rows.Scan(&w.ID, &w.PublicID, &w.Name)
			if err != nil {
				return err
			}
			exists = true
		}
		return nil
	}, publicID)

	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, nil
	}
	return &w, nil
}

// GetWorkspaceUserByUserIDAndWorkspaceID gets a workspace user by user ID and workspace ID
func (db *Database) GetWorkspaceUserByUserIDAndWorkspaceID(userID int, workspacePublicID string) (*types.WorkspaceUser, error) {
	var wu types.WorkspaceUser
	err := db.query(statementSelectWorkspaceUserByUserIDAndWorkspaceID, func(rows *sql.Rows) error {
		if rows.Next() {
			err := rows.Scan(&wu.ID, &wu.PublicID, &wu.WorkspaceID, &wu.UserID, &wu.Role)
			if err != nil {
				return err
			}
		}
		return nil
	}, userID, workspacePublicID)
	if err != nil {
		return nil, err
	}
	return &wu, nil
}

// GetWorkspacesByUserID gets workspaces a user is part of
func (db *Database) GetWorkspacesByUserID(userID int) ([]*types.WorkspaceUser, error) {
	var workspaces []*types.WorkspaceUser
	err := db.query(statementSelectWorkspacesByUserID, func(rows *sql.Rows) error {
		for rows.Next() {
			var wu types.WorkspaceUser
			err := rows.Scan(&wu.ID, &wu.PublicID, &wu.WorkspaceID, &wu.UserID, &wu.Role)
			if err != nil {
				return err
			}
			workspaces = append(workspaces, &wu)
		}
		return nil
	}, userID)
	if err != nil {
		return nil, err
	}
	return workspaces, nil
}
