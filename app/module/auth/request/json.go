package request

type Login struct {
	Mobile   uint64 `example:"09180338595" validate:"required,number"`
	Password string `example:"123456" validate:"required,min=6,max=255"`
}

type Register struct {
	FirstName string `example:"mahdi" validate:"required,min=2,max=255"`
	LastName  string `example:"lastname" validate:"required,min=2,max=255"`
	Mobile    uint64 `example:"9150338494" validate:"required,number"`
	Password  string `example:"123456" validate:"required,min=6,max=255"`
}
