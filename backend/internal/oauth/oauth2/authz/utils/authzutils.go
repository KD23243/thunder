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

// Package utils provides utility functions for OAuth2 authorization operations.
package utils

import (
	"errors"
	"time"

	"github.com/asgardeo/thunder/internal/oauth/oauth2/authz/constants"
	"github.com/asgardeo/thunder/internal/oauth/oauth2/authz/model"

	"github.com/asgardeo/thunder/internal/system/utils"
)

// GetAuthorizationCode generates an authorization code based on the provided OAuth message.
func GetAuthorizationCode(oAuthMessage *model.OAuthMessage) (model.AuthorizationCode, error) {
	sessionData := oAuthMessage.SessionData
	clientID := sessionData.OAuthParameters.ClientID
	if clientID == "" {
		clientID = oAuthMessage.RequestQueryParams["client_id"]
	}
	redirectURI := sessionData.OAuthParameters.RedirectURI
	if redirectURI == "" {
		redirectURI = oAuthMessage.RequestQueryParams["redirect_uri"]
	}

	if clientID == "" || redirectURI == "" {
		return model.AuthorizationCode{}, errors.New("client_id or redirect_uri is missing")
	}

	authUserID := sessionData.AuthenticatedUser.UserID
	if authUserID == "" {
		return model.AuthorizationCode{}, errors.New("authenticated user not found")
	}

	authTime := sessionData.AuthTime
	if authTime.IsZero() {
		return model.AuthorizationCode{}, errors.New("authentication time is not set")
	}

	scope := sessionData.OAuthParameters.Scopes
	if scope == "" {
		scope = oAuthMessage.RequestQueryParams["scope"]
	}

	// TODO: Add expiry time logic.
	expiryTime := authTime.Add(10 * time.Minute)

	   return model.AuthorizationCode{
		CodeID:           utils.GenerateUUID(),
		Code:             utils.GenerateUUID(),
		ClientID:         clientID,
		RedirectURI:      redirectURI,
		AuthorizedUserID: authUserID,
		TimeCreated:      authTime,
		ExpiryTime:       expiryTime,
		Scopes:           scope,
		State:            constants.AuthCodeStateActive,
	}, nil
}
