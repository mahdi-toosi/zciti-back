package request

type Login struct {
	Mobile   uint64 `example:"9380338494" validate:"required,number"`
	Password string `example:"123456" validate:"required,min=6,max=100"`
}

type Register struct {
	Login
	FirstName string `example:"mahdi" validate:"required,min=2,max=255"`
	LastName  string `example:"lastname" validate:"required,min=2,max=255"`
}

type SendOtp struct {
	Mobile uint64 `example:"9380338494" validate:"required,number"`
}

type ResetPass struct {
	Login
	Otp string `example:"1234567" validate:"required,min=5,max=10"`
}
