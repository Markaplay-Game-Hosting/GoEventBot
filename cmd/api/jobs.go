package main

import "net/http"

// getJob
// @Summary      Get Job
// @Description  Get job information
// @Tags         Jobs
// @Produce      json
// @Success      200 {object} data.Job
// @Router       /jobs/{id} [get]
func (app *application) getJob(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}
	jobs, err := app.models.Jobs.Get(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"jobs": jobs}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// listJobs
// @Summary      List Jobs
// @Description  List all jobs
// @Tags         Jobs
// @Produce      json
// @Success      200 {object} models.ListJobsResponse
// @Router       /jobs [get]
func (app *application) listJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := app.models.Jobs.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"jobs": jobs}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// skipJob
// @Summary      Skip Job
// @Description  Skip a job that has been loaded
// @Tags         Jobs
// @Produce      json
// @Success      204
// @Router       /jobs/{id}/skip [patch]
func (app *application) skipJob(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}
	err = app.models.Jobs.SkipJob(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	app.eventWorker.CancelJob(id)

	if err := app.writeJSON(w, http.StatusNoContent, nil, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
