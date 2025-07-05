package models

import "github.com/Markaplay-Game-Hosting/GoEventBot/internal/data"

// ListJobsResponse
// @Description response with a list of jobs
type ListJobsResponse struct {
	// list of jobs
	Jobs []data.Job `json:"jobs"`
} // @name Jobs.List.Response
