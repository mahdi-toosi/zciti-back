package schema

func Models() []any {
	return []any{
		User{},
		Post{},
		NotificationTemplate{},
		Notification{},
	}
}
