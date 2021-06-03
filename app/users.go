package main

type User struct {
	Name      string
	Password  string
}

func (user *User) checkUserExist() bool {
	var id int64
	row := db.QueryRow("SELECT id FROM users WHERE name = ? and password = md5(?)", user.Name, user.Password)
	err := row.Scan(&id)

	if err != nil {
		return false
	}

	return true
}