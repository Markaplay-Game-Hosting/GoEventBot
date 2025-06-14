package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Permissions PermissionModel
	Tokens      TokenModel
	Users       UserModel
	Events      EventModel
	Webhooks    WebhookModel
	Jobs        JobModel
	OAuth       OAuthModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Permissions: PermissionModel{DB: db},
		Tokens:      TokenModel{DB: db},
		Users:       UserModel{DB: db},
		Events:      EventModel{DB: db},
		Webhooks:    WebhookModel{DB: db},
		Jobs:        JobModel{DB: db},
		OAuth:       OAuthModel{},
	}
}
