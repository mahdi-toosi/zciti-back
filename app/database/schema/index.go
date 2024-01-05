package schema

import "gorm.io/gorm"

func MainDBModels() []any {
	// order matters
	return []any{
		User{},
		Post{},
		Business{},
		NotificationTemplate{},
		Notification{},
	}
}

func ChatDBModels() []any {
	// order matters
	return []any{
		MessageRoom{},
		Message{},
	}
}

func MainDBDropExtraCommands(db *gorm.DB) {
	db.Exec("DROP TABLE business_users")
}

func ChatDBDropExtraCommands(db *gorm.DB) {
}
