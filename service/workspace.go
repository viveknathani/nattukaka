package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/viveknathani/nattukaka/types"
	"github.com/viveknathani/nattukaka/utils"
)

// CreateWorkspace creates a new workspace and makes the user an admin
func (srv *Service) CreateWorkspace(userEmail, workspaceName string) (int, string) {
	if workspaceName == "" {
		return fiber.StatusBadRequest, "invalid workspace name"
	}

	// Get user info
	user, err := srv.Db.GetUserByEmail(userEmail)
	if err != nil {
		srv.Logger.Error(err.Error())
		return fiber.StatusInternalServerError, "something went wrong"
	}
	if user == nil {
		return fiber.StatusBadRequest, "user does not exist"
	}

	// Generate public ID for workspace
	publicIDForWorkspace, err := utils.GeneratePublicId("workspace")
	if err != nil {
		srv.Logger.Error(err.Error())
		return fiber.StatusInternalServerError, "something went wrong"
	}

	// Insert workspace
	insertedID, err := srv.Db.InsertWorkspace(&types.Workspace{
		PublicID: publicIDForWorkspace,
		Name:     workspaceName,
	})
	if err != nil {
		srv.Logger.Error(err.Error())
		return fiber.StatusInternalServerError, "something went wrong"
	}

	// Generate public ID for workspace
	publicIDForWorkspaceUser, err := utils.GeneratePublicId("workspace_user")
	if err != nil {
		srv.Logger.Error(err.Error())
		return fiber.StatusInternalServerError, "something went wrong"
	}

	// Insert workspace user with role ADMIN
	err = srv.Db.InsertWorkspaceUser(&types.WorkspaceUser{
		PublicID:    publicIDForWorkspaceUser,
		WorkspaceID: int(insertedID),
		UserID:      user.ID,
		Role:        utils.UserRoleAdmin,
	})
	if err != nil {
		srv.Logger.Error(err.Error())
		return fiber.StatusInternalServerError, "something went wrong"
	}

	return fiber.StatusCreated, "workspace created"
}

// GetWorkspaceUsers returns the list of users in a workspace
func (srv *Service) GetWorkspaceUsers(workspacePublicID string) (int, string, []*types.User) {
	workspace, err := srv.Db.GetWorkspaceByPublicID(workspacePublicID)
	if err != nil {
		srv.Logger.Error(err.Error())
		return fiber.StatusInternalServerError, "something went wrong", nil
	}
	if workspace == nil {
		return fiber.StatusBadRequest, "workspace does not exist", nil
	}

	users, err := srv.Db.GetUsersByWorkspace(workspace.ID)
	if err != nil {
		srv.Logger.Error(err.Error())
		return fiber.StatusInternalServerError, "something went wrong", nil
	}

	return fiber.StatusOK, "workspace users fetched", users
}

// DeleteWorkspace deletes a workspace if the user is an admin
func (srv *Service) DeleteWorkspace(userEmail, workspacePublicID string) (int, string) {
	user, err := srv.Db.GetUserByEmail(userEmail)
	if err != nil {
		srv.Logger.Error(err.Error())
		return fiber.StatusInternalServerError, "something went wrong"
	}
	if user == nil {
		return fiber.StatusBadRequest, "user does not exist"
	}

	workspaceUser, err := srv.Db.GetWorkspaceUserByUserIDAndWorkspaceID(user.ID, workspacePublicID)
	if err != nil {
		srv.Logger.Error(err.Error())
		return fiber.StatusInternalServerError, "something went wrong"
	}
	if workspaceUser == nil || workspaceUser.Role != "ADMIN" {
		return fiber.StatusForbidden, "user is not admin"
	}

	err = srv.Db.DeleteWorkspaceUsers(workspaceUser.WorkspaceID)
	if err != nil {
		srv.Logger.Error(err.Error())
		return fiber.StatusInternalServerError, "something went wrong"
	}

	err = srv.Db.DeleteWorkspace(workspacePublicID)
	if err != nil {
		srv.Logger.Error(err.Error())
		return fiber.StatusInternalServerError, "something went wrong"
	}

	return fiber.StatusOK, "workspace deleted"
}

// GetMyWorkspaces returns the workspaces the user is a part of
func (srv *Service) GetMyWorkspaces(userEmail string) (int, string, []*types.WorkspaceUser) {
	user, err := srv.Db.GetUserByEmail(userEmail)
	if err != nil {
		srv.Logger.Error(err.Error())
		return fiber.StatusInternalServerError, "something went wrong", nil
	}
	if user == nil {
		return fiber.StatusBadRequest, "user does not exist", nil
	}

	workspaces, err := srv.Db.GetWorkspacesByUserID(user.ID)
	if err != nil {
		srv.Logger.Error(err.Error())
		return fiber.StatusInternalServerError, "something went wrong", nil
	}

	return fiber.StatusOK, "workspaces fetched", workspaces
}
