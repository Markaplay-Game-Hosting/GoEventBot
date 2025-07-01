package models

// UserRegisterRequest
// @Description User Registration Body
type UserRegisterRequest struct {
	// username of the user
	Name string `json:"name"`
	// email of the user
	Email string `json:"email"`
	// password
	Password string `json:"password"`
} // @name User.Register.Request

// UserActivationRequest
// @Description Activate token request
type UserActivationRequest struct {
	// Token
	TokenPlaintext string `json:"token"`
} // @name User.Activation.Request

// UserUpdatePasswordRequest
// @Description User Request to update password
type UserUpdatePasswordRequest struct {
	// User password
	Password string `json:"password"`
	// Authorization token
	TokenPlaintext string `json:"token"`
} // @name User.UpdatePassword.Request
