package request

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils/paginator"
	"strings"
	"time"
)

type User struct {
	ID          uint64
	FirstName   string `example:"mahdi" validate:"required,min=2,max=255"`
	LastName    string `example:"lastname" validate:"required,min=2,max=255"`
	Mobile      uint64 `example:"9380338494" validate:"required,number"`
	Password    string
	Permissions schema.UserPermissions `example:"{1:['operator']}"`
}

type UpdateUserAccount struct {
	ID          uint64 `example:"1" validate:"required,min=1"`
	FirstName   string `example:"mahdi" validate:"required,min=2,max=255"`
	LastName    string `example:"lastname" validate:"required,min=2,max=255"`
	Mobile      uint64 `example:"9380338494" validate:"required,number"`
	CityID      uint64 `example:"1" validate:"number"`
	WorkspaceID uint64 `example:"1" validate:"number"`
	DormitoryID uint64 `example:"1" validate:"number"`
}

type BusinessUsers struct {
	BusinessID  uint64
	Username    string
	FullName    string
	CityID      uint64
	WorkspaceID uint64
	DormitoryID uint64
	IsSuspended string
	CountUsing  uint64
	Export      string
	UserIDs     []uint64
	StartTime   *time.Time
	EndTime     *time.Time
	Role        string
	Pagination  *paginator.Pagination
}

type BusinessUsersStoreRole struct {
	Roles               []schema.UserRole                  `example:"[user]" validate:"required"`
	UserID              uint64                             `example:"1" validate:"required,number,min=1"`
	BusinessID          uint64                             `example:"1" validate:"required,number,min=1"`
	PostsToObserve      []uint64                           `json:",omitempty" example:"[1,2,3]"`
	TaxonomiesToObserve schema.UserMetaTaxonomiesToObserve `json:",omitempty" example:"{1: { checked: true; partialChecked: false }}"`
}

type BusinessUsersToggleSuspense struct {
	IsSuspended    bool   `example:"true"`
	SuspenseReason string `example:"reason"`
	UserID         uint64 `example:"1" validate:"required,number,min=1"`
	BusinessID     uint64 `example:"1" validate:"required,number,min=1"`
}
type OrderStatus struct {
	Token      string `example:"token"`
	ResNum     string `example:"resnum"` // order id
	RefNum     string `example:"refnum"` // این پارامتر کدی است تا 50 حرف یا عدد که برای هر تراکنش ایجاد می شود.
	State      string `example:"state"`  // OK or FAILED
	Status     string `example:"OK" validate:"required"`
	TraceNo    string `example:"traceno"`   //  این پارامتر شماره پیگیری تولید شده توسط سپ می باشد.
	SecurePan  string `example:"securepan"` //  این پارامتر شماره پیگیری تولید شده توسط سپ می باشد.
	Amount     int    `example:"10000"`
	Rrn        string `example:"rrn"`        //  شماره مرجع تراکنش
	MID        string `example:"mid"`        // شماره ترمینال پذیرنده
	TerminalID string `example:"terminalid"` // شماره ترمینال پذیرنده
}

type Users struct {
	Pagination *paginator.Pagination
	Keyword    string
}

func (req *User) ToDomain() *schema.User {
	return &schema.User{
		ID:          req.ID,
		Mobile:      req.Mobile,
		Permissions: req.Permissions,
		LastName:    strings.TrimSpace(req.LastName),
		FirstName:   strings.TrimSpace(req.FirstName),
	}
}

func (req *UpdateUserAccount) ToDomain() *schema.User {
	p := &schema.User{
		ID:        req.ID,
		Mobile:    req.Mobile,
		LastName:  strings.TrimSpace(req.LastName),
		FirstName: strings.TrimSpace(req.FirstName),
	}

	if req.CityID != 0 {
		p.CityID = &req.CityID
	}
	if req.WorkspaceID != 0 {
		p.WorkspaceID = &req.WorkspaceID
	}
	if req.DormitoryID != 0 {
		p.DormitoryID = &req.DormitoryID
	}

	return p
}
