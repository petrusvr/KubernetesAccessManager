/*
 * KubernetesUserManager - API
 *
 * This is a backend API server documentation for KubernetesUserManager  Some useful links: - [Jira](https://samuelus.atlassian.net/jira/software/projects/ZPI/boards/4) - [Confluence](https://samuelus.atlassian.net/wiki/spaces/ZPI/overview)
 *
 * API version: 0.0.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package models

// Authentication tokens
type AuthTokens struct {
	// JWT access token
	AccessToken string `json:"accessToken,omitempty"`
	// Refresh token
	RefreshToken string `json:"refreshToken,omitempty"`
	// Access token expiration time in seconds
	ExpiresIn int32 `json:"expiresIn,omitempty"`
}
