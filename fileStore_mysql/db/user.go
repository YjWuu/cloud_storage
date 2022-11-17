package db

import (
	"cloud_storage/fileStore_mysql/db/mysql"
	"fmt"
)

func UserSignup(username string, passwd string) bool {
	stmt, err := mysql.DBConn().Prepare("insert ignore into tbl_user(`user_name`, `user_pwd`) values (?, ?)")
	if err != nil {
		fmt.Println("Failed into insert, err:" + err.Error())
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(username, passwd)
	if err != nil {
		fmt.Println("Failed to insert, err:" + err.Error())
		return false
	}
	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	return false
}

// UserSignin 判断密码是否一致
func UserSignin(username string, encpwd string) bool {
	stmt, err := mysql.DBConn().Prepare("select * from tbl_user where user_name = ? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Println(err.Error())
	} else if rows == nil {
		fmt.Println("username not found:" + username)
		return false
	}
	pRows := mysql.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd {
		return true
	}
	return false
}

// UpdateToken 刷新用户登录的token
func UpdateToken(username, token string) bool {
	stmt, err := mysql.DBConn().Prepare("replace into tbl_user_token(`user_name`, `user_token`) values (?, ?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

type User struct {
	Username     string
	Email        string
	Phone        string
	SignupAt     string
	LastActiveAt string
	Status       int
}

func GetUserInfo(username string) (User, error) {
	user := User{}
	stmt, err := mysql.DBConn().Prepare("select user_name, signup_at from tbl_user where user_name = ? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}
	defer stmt.Close()
	//执行查询的操作
	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		return user, err
	}
	return user, nil
}
