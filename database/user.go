package database

// CreateUser 保存用户到数据库
func CreateUser(user *User) error {
	result := DB.Create(user) // 使用 GORM 将数据持久化
	return result.Error       // 返回错误（如果有）
}

// DeleteUserByID 根据 ID 删除用户
func DeleteUserByID(id int) error {
	result := DB.Delete(&User{}, id)

	if result.RowsAffected == 0 {
		return result.Error
	}

	return nil
}

func GetAllUsers() ([]User, error) {
	var users []User
	result := DB.Find(&users)
	return users, result.Error
}

func UpdateUser(id int, updatedUser *User) error {
	var user User
	result := DB.First(&user, id)
	if result.Error != nil {
		return result.Error
	}

	user.Name = updatedUser.Name
	user.Email = updatedUser.Email

	saveResult := DB.Save(&user)
	return saveResult.Error
}

func GetUserByID(id int) (*User, error) {
	var user User
	result := DB.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
