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

	"github.com/asgardeo/thunder/internal/authn"
	"github.com/asgardeo/thunder/internal/system/server"
)

// AuthenticationService defines the service for handling authentication requests.
type AuthenticationService struct {
	authHandler *authn.AuthenticationHandler
}

// NewAuthenticationService creates a new instance of AuthenticationService.
func NewAuthenticationService(mux *http.ServeMux) *AuthenticationService {
	instance := &AuthenticationService{
		authHandler: authn.NewAuthenticationHandler(),
	}
	instance.RegisterRoutes(mux)
	return instance
}

// RegisterRoutes registers the routes for the AuthenticationService.
func (s *AuthenticationService) RegisterRoutes(mux *http.ServeMux) {
	opts := server.RequestWrapOptions{
		Cors: &server.Cors{
			AllowedMethods:   "POST, GET",
			AllowedHeaders:   "Content-Type, Authorization",
			AllowCredentials: true,
		},
	}
	server.WrapHandleFunction(mux, "POST /flow/authn", &opts, s.authHandler.HandleAuthenticationRequest)
	server.WrapHandleFunction(mux, "GET /flow/authn", &opts, s.authHandler.HandleAuthenticationRequest)
	server.WrapHandleFunction(mux, "OPTIONS /flow/authn", &opts, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
}
