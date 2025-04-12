package db

// GetUser retrieves a user by Telegram ID
func (db *DB) GetUser(tgID int64) (*User, error) {
	var user User
	result := db.conn.Where("tg_id = ?", tgID).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by Telegram username
func (db *DB) GetUserByUsername(username string) (*User, error) {
	var user User
	result := db.conn.Where("tg_username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// SetUser creates or updates a user
func (db *DB) SetUser(user *User) error {
	// Use upsert functionality (create if not exists, update if exists)
	result := db.conn.Save(user)
	return result.Error
}
