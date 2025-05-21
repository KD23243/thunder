/*
 * Copyright (c) 2025, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

// Package service provides the implementation for user management operations.
package service

import (
	"errors"

	"github.com/asgardeo/thunder/internal/system/utils"

	"github.com/asgardeo/thunder/internal/system/log"
	"github.com/asgardeo/thunder/internal/user/model"
	"github.com/asgardeo/thunder/internal/user/store"
)

// UserServiceInterface defines the interface for the user service.
type UserServiceInterface interface {
	CreateUser(user *model.User) (*model.User, error)
	GetUserList() ([]model.User, error)
	GetUser(userID string) (*model.User, error)
	UpdateUser(userID string, user *model.User) (*model.User, error)
	DeleteUser(userID string) error
}

// UserService is the default implementation of the UserServiceInterface.
type UserService struct{}

// GetUserService creates a new instance of UserService.
func GetUserService() UserServiceInterface {
	return &UserService{}
}

// CreateUser creates the user.
func (as *UserService) CreateUser(user *model.User) (*model.User, error) {
	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, "UserService"))

	user.ID = utils.GenerateUUID()

	// Create the user in the database.
	err := store.CreateUser(*user)
	if err != nil {
		logger.Error("Failed to create user", log.Error(err))
		return nil, err
	}
	return user, nil
}

// GetUserList list the users.
func (as *UserService) GetUserList() ([]model.User, error) {
	users, err := store.GetUserList()
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetUser get the user for given user id.
func (as *UserService) GetUser(userID string) (*model.User, error) {
	if userID == "" {
		return nil, errors.New("user ID is empty")
	}

	user, err := store.GetUser(userID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUser update the user for given user id.
func (as *UserService) UpdateUser(userID string, user *model.User) (*model.User, error) {
	if userID == "" {
		return nil, errors.New("user ID is empty")
	}

	err := store.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser delete the user for given user id.
func (as *UserService) DeleteUser(userID string) error {
	if userID == "" {
		return errors.New("user ID is empty")
	}

	err := store.DeleteUser(userID)
	if err != nil {
		return err
	}

	return nil
}
