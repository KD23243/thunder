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

// Package service provides application-related business logic and operations.
package service

import (
	"errors"
	"fmt"

	"github.com/asgardeo/thunder/internal/application/model"
	"github.com/asgardeo/thunder/internal/application/store"
	dbprovider "github.com/asgardeo/thunder/internal/system/database/provider"
	"github.com/asgardeo/thunder/internal/system/log"
	"github.com/asgardeo/thunder/internal/system/utils"

)

// ApplicationServiceInterface defines the interface for the application service.
type ApplicationServiceInterface interface {
	GetOAuthApplication(clientID string) (*model.OAuthApplication, error)
	CreateApplication(app *model.Application) (*model.Application, error)
	GetApplicationList() ([]model.Application, error)
	GetApplication(appID string) (*model.Application, error)
	UpdateApplication(appID string, app *model.Application) (*model.Application, error)
	DeleteApplication(appID string) error
}

// ApplicationService is the default implementation of the ApplicationServiceInterface.
type ApplicationService struct{}

// GetApplicationService creates a new instance of ApplicationService.
func GetApplicationService() ApplicationServiceInterface {
	return &ApplicationService{}
}

// GetOAuthApplication retrieves the OAuth application based on the client id.
func (as *ApplicationService) GetOAuthApplication(clientID string) (*model.OAuthApplication, error) {
	if clientID == "" {
		return nil, errors.New("client ID cannot be empty")
	}
	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, "ApplicationService"))

	dbClient, err := dbprovider.NewDBProvider().GetDBClient("identity")
	if err != nil {
		logger.Error("Failed to get database client", log.Error(err))
		return nil, err
	}
	defer func() {
		if closeErr := dbClient.Close(); closeErr != nil {
			logger.Error("Failed to close database client", log.Error(closeErr))
			err = fmt.Errorf("failed to close database client: %w", closeErr)
		}
	}()

	results, err := dbClient.Query(store.QueryGetApplicationByClientID, clientID)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, errors.New("OAuth application not found")
	}

	row := results[0]

	clientID, ok := row["consumer_key"].(string)
	if !ok {
		return nil, errors.New("failed to parse consumer_key as string")
	}

	clientSecret, ok := row["consumer_secret"].(string)
	if !ok {
		return nil, errors.New("failed to parse consumer_secret as string")
	}

	var redirectURIs []string
	if row["callback_uris"] != nil {
		if uris, ok := row["callback_uris"].(string); ok {
			redirectURIs = utils.ParseStringArray(uris)
		} else {
			return nil, errors.New("failed to parse callback_uris as string")
		}
	}

	var allowedGrantTypes []string
	if row["grant_types"] != nil {
		if grants, ok := row["grant_types"].(string); ok {
			allowedGrantTypes = utils.ParseStringArray(grants)
		} else {
			return nil, errors.New("failed to parse grant_types as string")
		}
	}

	return &model.OAuthApplication{
		ClientID:          clientID,
		ClientSecret:      clientSecret,
		RedirectURIs:      redirectURIs,
		AllowedGrantTypes: allowedGrantTypes,
	}, nil
}

// CreateApplication creates the application.
func (as *ApplicationService) CreateApplication(app *model.Application) (*model.Application, error) {
	if app == nil {
		return nil, errors.New("application is nil")
	}
	if app.Name == "" {
		return nil, errors.New("application name cannot be empty")
	}
	if app.ClientID == "" {
		return nil, errors.New("client ID cannot be empty")
	}
	if app.ClientSecret == "" {
		return nil, errors.New("client secret cannot be empty")
	}
	if len(app.CallbackURLs) == 0 {
		return nil, errors.New("at least one callback URL is required")
	}

	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, "ApplicationService"))
	app.ID = utils.GenerateUUID()

	// Create the application in the database.
	err := store.CreateApplication(*app)
	if err != nil {
		logger.Error("Failed to create application", log.Error(err))
		return nil, err
	}
	return app, nil
}

// GetApplicationList list the applications.
func (as *ApplicationService) GetApplicationList() ([]model.Application, error) {
	applications, err := store.GetApplicationList()
	if err != nil {
		return nil, err
	}

	return applications, nil
}

// GetApplication get the application for given app id.
func (as *ApplicationService) GetApplication(appID string) (*model.Application, error) {
	if appID == "" {
		return nil, errors.New("application ID is empty")
	}

	application, err := store.GetApplication(appID)
	if err != nil {
		return nil, err
	}

	return &application, nil
}

// UpdateApplication update the application for given app id.
func (as *ApplicationService) UpdateApplication(appID string, app *model.Application) (*model.Application, error) {
	if appID == "" {
		return nil, errors.New("application ID is empty")
	}
	if app == nil {
		return nil, errors.New("application is nil")
	}
	if app.Name == "" {
		return nil, errors.New("application name cannot be empty")
	}
	if app.ClientID == "" {
		return nil, errors.New("client ID cannot be empty")
	}
	if app.ClientSecret == "" {
		return nil, errors.New("client secret cannot be empty")
	}
	if len(app.CallbackURLs) == 0 {
		return nil, errors.New("at least one callback URL is required")
	}

	err := store.UpdateApplication(app)
	if err != nil {
		return nil, err
	}

	return app, nil
}

// DeleteApplication delete the application for given app id.
func (as *ApplicationService) DeleteApplication(appID string) error {
	if appID == "" {
		return errors.New("application ID is empty")
	}

	err := store.DeleteApplication(appID)
	if err != nil {
		return err
	}

	return nil
}
