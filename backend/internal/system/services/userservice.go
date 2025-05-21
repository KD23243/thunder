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

	"github.com/asgardeo/thunder/internal/system/server"
	"github.com/asgardeo/thunder/internal/user/handler"
)

// UserService is the service for user management operations.
type UserService struct {
	userHandler *handler.UserHandler
}

// NewUserService creates a new instance of UserService.
func NewUserService(mux *http.ServeMux) *UserService {
	instance := &UserService{
		userHandler: &handler.UserHandler{},
	}
	instance.RegisterRoutes(mux)

	return instance
}

// RegisterRoutes registers the routes for user management operations.
//
//nolint:dupl // Ignoring false positive duplicate code
func (s *UserService) RegisterRoutes(mux *http.ServeMux) {
	opts1 := server.RequestWrapOptions{
		Cors: &server.Cors{
			AllowedMethods:   "GET, POST",
			AllowedHeaders:   "Content-Type, Authorization",
			AllowCredentials: true,
		},
	}
	server.WrapHandleFunction(mux, "POST /users", &opts1, s.userHandler.HandleUserPostRequest)
	server.WrapHandleFunction(mux, "GET /users", &opts1, s.userHandler.HandleUserListRequest)
	server.WrapHandleFunction(mux, "OPTIONS /users", &opts1, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	opts2 := server.RequestWrapOptions{
		Cors: &server.Cors{
			AllowedMethods:   "GET, PUT, DELETE",
			AllowedHeaders:   "Content-Type, Authorization",
			AllowCredentials: true,
		},
	}
	server.WrapHandleFunction(mux, "GET /users/", &opts2, s.userHandler.HandleUserGetRequest)
	server.WrapHandleFunction(mux, "PUT /users/", &opts2, s.userHandler.HandleUserPutRequest)
	server.WrapHandleFunction(mux, "DELETE /users/", &opts2, s.userHandler.HandleUserDeleteRequest)
	server.WrapHandleFunction(mux, "OPTIONS /users/", &opts2, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
}
