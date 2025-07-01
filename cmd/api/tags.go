package main

import (
	"database/sql"
	"errors"
	"github.com/Markaplay-Game-Hosting/GoEventBot/cmd/api/models"
	"github.com/Markaplay-Game-Hosting/GoEventBot/internal/data"
	"github.com/Markaplay-Game-Hosting/GoEventBot/internal/validator"
	"net/http"
)

// createTagHandler
// @Summary      Create Tag
// @Description  Create a new tag
// @Tags         Tags
// @param  tag body models.CreateTag true "required body for creating a new tag"
// @Produce      json
// @Success      201 {object} data.Tag "tag information"
// @Router       /tags [post]
// @Security ApiKeyAuth
func (app *application) createTagHandler(w http.ResponseWriter, r *http.Request) {
	var input models.CreateTag

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	tag := &data.Tag{
		Name:        input.Name,
		Description: input.Description,
	}

	v := validator.New()
	data.ValidateTag(v, tag)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Tags.Insert(tag)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"tag": tag}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// getTagHandler
// @Summary      Get Tag
// @Description  Get tag
// @Tags         Tags
// @param  id path string true "id of the tag"
// @Produce      json
// @Success      200 {object} data.Tag "tag information"
// @Router       /tags/{id} [get]
func (app *application) getTagHandler(w http.ResponseWriter, r *http.Request) {
	tagID, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	tag, err := app.models.Tags.GetByID(tagID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if tag == nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"tag": tag}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// updateTagHandler
// @Summary      Update Tag
// @Description  Update tag
// @Tags         Tags
// @param  id path string true "id of the tag"
// @param tag body  models.UpdateTagRequest true "body to update the tag"
// @Produce      json
// @Success      200 {object} data.Tag "tag information"
// @Router       /tags/{id} [put]
// @Security ApiKeyAuth
func (app *application) updateTagHandler(w http.ResponseWriter, r *http.Request) {
	tagID, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var input models.UpdateTagRequest

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	currentTag, err := app.models.Tags.GetByID(tagID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if currentTag == nil {
		app.notFoundResponse(w, r)
		return
	}

	if input.Description != "" {
		currentTag.Description = input.Description
	}
	if input.Name != "" {
		currentTag.Name = input.Name
	}

	v := validator.New()
	data.ValidateTag(v, currentTag)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Tags.Update(currentTag)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"tag": currentTag}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// deleteTagHandler
// @Summary      Delete Tag
// @Description  Delete tag
// @Tags         Tags
// @param  id path string true "id of the tag"
// @Produce      json
// @Success      204
// @Router       /tags/{id} [delete]
// @Security ApiKeyAuth
func (app *application) deleteTagHandler(w http.ResponseWriter, r *http.Request) {
	tagID, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.models.Tags.Delete(tagID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return

	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "Tag deleted successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// listTagsHandler
// @Summary      List Tags
// @Description  List all tags
// @Tags         Tags
// @Produce      json
// @Success      200 {array} data.Tag "tag list"
// @Router       /tags [get]
func (app *application) listTagsHandler(w http.ResponseWriter, r *http.Request) {
	tags, err := app.models.Tags.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if len(tags) == 0 {
		err := app.writeJSON(w, http.StatusNoContent, envelope{"tags": []data.Tag{}}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"tags": tags}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
