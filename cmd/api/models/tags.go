package models

// CreateTag
// @Description Create Tag Model
type CreateTag struct {
	// Name of the tag
	Name string `json:"name" example:"Warriors"`
	// Short description of the tag
	Description string `json:"description" example:"regroup events for warriors"`
} // @name Tag.Create.Request

// UpdateTagRequest
// @Description Request body to update the tag
type UpdateTagRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
} // @name Tag.Update.Request
