package requests

// import "mime/multipart

var userErrorMsg = map[string]string{
	"Nickname.required": "nickname is required",
	"Nickname.min":      "nickname must be at least 1 characters long",
	"Nickname.max":      "nickname must be at most 50 characters long",
	"Avatar.required":   "avatar file is required",
}

type UpdateNicknameRequest struct {
	Nickname string `json:"nickname" form:"nickname" validate:"required,min=1,max=50"`
}

func (r *UpdateNicknameRequest) Validate() error {
	return FormatError(Validate.Struct(r), userErrorMsg)
}

type UpdateAvatarRequest struct {
	Avatar string `json:"avatar" form:"avatar" validate:"required"`
}

func (r *UpdateAvatarRequest) Validate() error {
	return FormatError(Validate.Struct(r), userErrorMsg)
}
