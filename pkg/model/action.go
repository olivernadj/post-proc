/*
 * Post process
 *
 * This is a sample server Processing and saving incoming requests
 *
 * API version: 1.0.0
 * Contact: mr.oliver.nadj@gmail.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package model

type NewAction struct {

	Action string `json:"action,omitempty"`

	// Status of the item
	State string `json:"state,omitempty"`
}