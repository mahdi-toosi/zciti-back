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
		OrderItem{},
		Reservation{},
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
	db.Exec("DROP TABLE IF EXISTS posts_taxonomies")
	db.Exec("DROP TABLE IF EXISTS products_taxonomies")
}

func ChatDBDropExtraCommands(db *gorm.DB) {
}
