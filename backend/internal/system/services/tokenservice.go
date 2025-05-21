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

package services

import (
	"net/http"

	"github.com/asgardeo/thunder/internal/oauth/oauth2/token"
	"github.com/asgardeo/thunder/internal/system/server"
)

// TokenService defines the service for handling OAuth2 token requests.
type TokenService struct {
	tokenHandler *token.TokenHandler
}

// NewTokenService creates a new instance of TokenService.
func NewTokenService(mux *http.ServeMux) *TokenService {
	instance := &TokenService{
		tokenHandler: &token.TokenHandler{},
	}
	instance.RegisterRoutes(mux)

	return instance
}

// RegisterRoutes registers the routes for the TokenService.
func (s *TokenService) RegisterRoutes(mux *http.ServeMux) {
	opts := server.RequestWrapOptions{
		Cors: &server.Cors{
			AllowedMethods:   "POST",
			AllowedHeaders:   "Content-Type, Authorization",
			AllowCredentials: true,
		},
	}
	server.WrapHandleFunction(mux, "POST /oauth2/token", &opts, s.tokenHandler.HandleTokenRequest)
}
