package schema

import "gorm.io/gorm"

func MainDBModels() []any {
	// order matters
	return []any{
		User{},
		Business{},
		Asset{},
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
	db.Exec("DROP TABLE IF EXISTS business_users")
}

func ChatDBDropExtraCommands(db *gorm.DB) {
}
