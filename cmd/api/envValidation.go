package main

func environmentValidation() {
	if databaseName == "" {
		panic("DB_DATABASE must be set")
	}
	if password == "" {
		panic("DB_PASSWORD must be set")
	}
	if username == "" {
		panic("DB_USERNAME must be set")
	}
	if databasePort == "" {
		panic("DB_PORT must be set")
	}
	if host == "" {
		panic("DB_HOST must be set")
	}
	if port == "" {
		panic("PORT must be set")
	}
	if secretKey == "" || len(secretKey) < 8 {
		panic("SECRET must be set and of at least 8 characters")
	}
}
