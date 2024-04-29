package schema

import "gorm.io/gorm"

func MainDBModels() []any {
	// order matters
	return []any{
		User{},
		Business{},
		Asset{},
		Post{},
		Product{},
		Order{},
		Taxonomy{},
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
	db.Exec("DROP TABLE IF EXISTS post_taxonomy")
}

func ChatDBDropExtraCommands(db *gorm.DB) {
}
