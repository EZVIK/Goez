package models

type User struct {
	ID 				int 		`gorm:"primary_key" json:"id"`
	Username 		string 		`json:"username"`
	Password 		string 		`json:"password"`
	NickName		string		`json:"nickname"`
	Avatar 			string		`json:"avatar"`
	Auth 			int			`json:"auth"`
	Articles		[]Article	`gorm:"-"`
	Records			[]Record	`gorm:"-"`
}

func CheckAuth(username, password string) (user User, err error) {

	err = db.Where("username = ? and password = ? ", username, password).First(&user).Error

	if err != nil {
	    return user, err
	}

	return
}

func Register(username, password, nickname string) (int, error) {
	user := User{Username: username, Password: password, NickName: nickname}

	err := db.Create(&user).Error

	if err != nil {
	    return 0, err
	}

	return user.ID, nil
}

func UpdateUser(id int, username, password, nickname string) (error)  {

	u := User{ID: id}
	if err := db.Model(&u).Updates(User{Username: username, Password: password, NickName: nickname}).Error; err != nil {
		return err
	}

	return nil
}

// 获取用户浏览记录
func (u User) GetUserRecords(keyword string) (re []Record, err error) {

	if err := db.Model(&re).Where("user_id = ?", u.ID).Error; err != nil {
		return re, err
	} else {
		return re, err
	}

}