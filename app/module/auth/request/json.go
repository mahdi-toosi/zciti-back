package request

type LoginRequest struct {
	Mobile   uint64 `json:"mobile" example:"09180338595" validate:"regex:09(1\[0-9\]|3\[1-9\]|2\[1-9\]\)-?\[0-9\]{3}-?\[0-9\]{4}"` //nolint:govet
	Password string `json:"password" example:"123456" validate:"required,min=6,max=255"`
}
