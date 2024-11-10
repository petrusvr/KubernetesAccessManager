/*
 * KubernetesAccessManager - API
 *
 * This is a backend API server documentation for KubernetesAccessManager  Some useful links: - [Jira](https://samuelus.atlassian.net/jira/software/projects/ZPI/boards/4) - [Confluence](https://samuelus.atlassian.net/wiki/spaces/ZPI/overview)
 *
 * API version: 0.0.4
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package models

// Error response
type ModelError struct {
	// Error code
	Code int32 `json:"code,omitempty"`
	// Error message
	Message string `json:"message,omitempty"`
}
