package schema

import "gorm.io/gorm"

func Models() []any {
	// order matters
	return []any{
		User{},
		Post{},
		Business{},
		NotificationTemplate{},
		Notification{},
	}
}

func DropExtraCommands(db *gorm.DB) {
	db.Exec("DROP TABLE business_users")
}
