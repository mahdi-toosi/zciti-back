package schema

import "gorm.io/gorm"

func MainDBModels() []any {
	// order matters
	return []any{
		User{},
		Business{},
		Post{},
		Comment{},
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
