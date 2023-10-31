package schema

func Models() []any {
	// order matters
	return []any{
		User{},
		Post{},
		Business{},
		NotificationTemplate{},
		Notification{},
		Business{},
	}
}
