package requests

import (
	"errors"
	"regexp"
)

var authErrorMsg = map[string]string{
	"Username.required": "Username is required",
	"Username.min":      "Username must be at least 2 characters",
	"Username.max":      "Username must be at most 32 characters",
	"Username.regexp":   "Username must contain only letters, numbers, and underscores",
	"Email.required":    "Email is required",
	"Email.email":       "Invalid email format",
	"Password.required": "Password is required",
	"Password.min":      "Password must be at least 6 characters",
	"FlowId.required":   "FlowId is required",
	"Nickname.min":      "Nickname must be at least 2 characters",
	"Nickname.max":      "Nickname must be at most 32 characters",
	"Nickname.regexp":   "Nickname must contain only letters, numbers, and underscores",
	"Code.regexp":       "Code format is invalid",
	"Code.len":          "Code must be 6 characters",
	"Code.required":     "Code is required",
	"Device.max":        "Device must be at most 100 characters",
}

// register
type EmailRegisterLinkRequest struct {
	Username string `json:"username" form:"username" validate:"required,min=2,max=32"`
	Email    string `json:"email" form:"email" validate:"required,email"`
}

type EmailRegisterLinkVerifyRequest struct {
	FlowId   string `json:"flowId" form:"flowId" validate:"required"`
	Password string `json:"password" form:"password" validate:"required,min=6"`
	Nickname string `json:"nickname" form:"nickname"`
}

type EmailRegisterCodeRequest struct {
	Username string `json:"username" form:"username" validate:"required,min=2,max=32"`
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=6"`
}

type EmailRegisterCodeVerifyRequest struct {
	Email string `json:"email" form:"email" validate:"required,email"`
	Code  string `json:"code" form:"code" validate:"required,len=6"`
}

// login
type EmailLoginWithPasswordRequest struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=6"`
	Device   string `json:"device" form:"device" validate:"max=100"`
}

type UsernameLoginWithPasswordRequest struct {
	Username string `json:"username" form:"username" validate:"required"`
	Password string `json:"password" form:"password" validate:"required,min=6"`
	Device   string `json:"device" form:"device" validate:"max=100"`
}

// reset password
type EmailPasswordResetLinkRequest struct {
	Email string `json:"email" form:"email" validate:"required,email"`
}

type EmailPasswordResetCodeRequest struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=6"`
}

type EmailPasswordResetLinkVerifyRequest struct {
	FlowId   string `json:"flowId" form:"flowId" validate:"required"`
	Password string `json:"password" form:"password" validate:"required,min=6"`
}

type EmailPasswordResetCodeVerifyRequest struct {
	Email string `json:"email" form:"email" validate:"required,email"`
	Code  string `json:"code" form:"code" validate:"required,len=6"`
}

// token
type RefreshTokenRequest struct {
	AccessToken  string `json:"accessToken" form:"accessToken" validate:"required"`
	RefreshToken string `json:"refreshToken" form:"refreshToken" validate:"required"`
	Device       string `json:"device" form:"device" validate:"max=100"`
}

// logout
type LogoutRequest struct {
	Device string `json:"device" form:"device" validate:"required,max=100"`
}

// register
func (r *EmailRegisterLinkRequest) Validate() error {
	err := FormatError(Validate.Struct(r), authErrorMsg)
	if err != nil {
		return err
	}

	// email format check
	if !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(r.Email) {
		return errors.New(authErrorMsg["Email.email"])
	}

	// username format check
	if !regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`).MatchString(r.Username) {
		return errors.New(authErrorMsg["Username.regexp"])
	}

	return nil
}

func (r *EmailRegisterCodeRequest) Validate() error {
	return FormatError(Validate.Struct(r), authErrorMsg)
}

func (r *EmailRegisterLinkVerifyRequest) Validate() error {
	err := FormatError(Validate.Struct(r), authErrorMsg)
	if err != nil {
		return err
	}

	if r.Nickname != "" {
		// only letters, numbers, and underscores
		if !regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`).MatchString(r.Nickname) {
			return errors.New(authErrorMsg["Nickname.regexp"])
		}
	}

	// regexp check, only numbers
	// if !regexp.MustCompile(`^[0-9]+$`).MatchString(r.Code) {
	// 	return errors.New(authErrorMsg["Code.regexp"])
	// }

	return nil
}

func (r *EmailRegisterCodeVerifyRequest) Validate() error {
	return FormatError(Validate.Struct(r), authErrorMsg)
}

// login
func (r *EmailLoginWithPasswordRequest) Validate() error {
	err := FormatError(Validate.Struct(r), authErrorMsg)
	if err != nil {
		return err
	}

	if r.Device == "" {
		r.Device = "unknown"
	}

	return nil
}

func (r *UsernameLoginWithPasswordRequest) Validate() error {
	return FormatError(Validate.Struct(r), authErrorMsg)
}

// reset password
func (r *EmailPasswordResetLinkRequest) Validate() error {
	return FormatError(Validate.Struct(r), authErrorMsg)
}

func (r *EmailPasswordResetCodeRequest) Validate() error {
	return FormatError(Validate.Struct(r), authErrorMsg)
}

func (r *EmailPasswordResetLinkVerifyRequest) Validate() error {
	return FormatError(Validate.Struct(r), authErrorMsg)
}

func (r *EmailPasswordResetCodeVerifyRequest) Validate() error {
	return FormatError(Validate.Struct(r), authErrorMsg)
}

// token
func (r *RefreshTokenRequest) Validate() error {
	return FormatError(Validate.Struct(r), authErrorMsg)
}

// logout
func (r *LogoutRequest) Validate() error {
	return FormatError(Validate.Struct(r), authErrorMsg)
}
