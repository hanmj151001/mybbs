package mysql

import (
	"bluebell_backend/models"
	"bluebell_backend/pkg/gofunc"
	"bluebell_backend/pkg/snowflake"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
)

// 把每一步数据库操作封装成函数
// 待logic层根据业务需求调用

const secret = "huchao.vip"

/**
 * @Author mengjie.han
 * @Description //TODO 自己已经写过了这里需要删除一下
 * @Date 21:50 2022/2/10
 **/
func encryptPassword(data []byte) (result string) {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum(data))
}

/**
 * @Author mengjie.han
 * @Description //TODO 检查指定用户名的用户是否存在
 * @Date 21:50 2022/2/10
 **/
func CheckUserExist(username string) (error error) {
	sqlstr := `select count(user_id) from user where username = ?`
	var count int
	if err := db.Get(&count, sqlstr, username); err != nil {
		return err
	}
	if count > 0 {
		return errors.New("用户已存在")
	}
	return
}
func CheckEmailExist(Email string) (count int, err error) {
	sqlStr := `select count(user_id) from user where Email = ?`
	if err := db.Get(&count, sqlStr, Email); err != nil {
		return
	}
	if count > 0 {
		return count, errors.New("邮箱已经被占用")
	}
	return
}

/**
 * @Author hucaho
 * @Description //TODO 注册业务-向数据库中插入一条新的用户
 * @Date 21:51 2022/2/10
 **/
func InsertUser(user models.User) (error error) {
	// 对密码进行加密
	user.Password = encryptPassword([]byte(user.Password))
	// 执行SQL语句入库
	sqlstr := `insert into user(user_id,username,password) values(?,?,?)`
	_, err := db.Exec(sqlstr, user.UserID, user.UserName, user.Password)
	return err
}

func InsertUser2(user models.User2) (err error) {
	// 执行SQL语句入库
	sqlStr := `insert into user(user_id,email,password,nickname) values(?,?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Email, user.Password, user.NickName)
	return err
}

func Register(user *models.User) (err error) {
	sqlStr := "select count(user_id) from user where username = ?"
	var count int64
	err = db.Get(&count, sqlStr, user.UserName)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if count > 0 {
		// 用户已存在
		return ErrorUserExit
	}
	// 生成user_id
	userID, err := snowflake.GetID()
	if err != nil {
		return ErrorGenIDFailed
	}
	// 生成加密密码
	password := encryptPassword([]byte(user.Password))
	// 把用户插入数据库
	sqlStr = "insert into user(user_id, username, password) values (?,?,?)"
	_, err = db.Exec(sqlStr, userID, user.UserName, password)
	return
}

/**
 * @Author huchao
 * @Description //TODO 登录业务
 * @Date 21:52 2022/2/10
 **/
func Login(user *models.User2) (err error) {
	originPassword := user.Password // 记录一下原始密码(用户登录的密码)
	sqlStr := "select user_id, email, password from user where email = ?"
	err = db.Get(user, sqlStr, user.Email)
	if err != nil && err != sql.ErrNoRows {
		// 查询数据库出错
		return
	}
	if err == sql.ErrNoRows {
		// 用户不存在
		return ErrorUserNotExit
	}
	// 生成加密密码与查询到的密码比较
	password := gofunc.EncodePassword(originPassword)
	if user.Password != password {
		return ErrorPasswordWrong
	}
	return
}

/**
 * @Author huchao
 * @Description //TODO 根据ID查询作者信息
 * @Date 22:05 2022/2/12
 **/
func GetUserByID(id uint64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id, username from user where user_id = ?`
	err = db.Get(user, sqlStr, id)
	return
}

//TODO 这个逻辑不知道写的对不对，记得检查

func GetUserByEmail(email string) (user *models.User2, err error) {
	user = &models.User2{}
	sqlStr := `select * from user where user_id = ?`
	err = db.Get(user, sqlStr, email)
	if err != nil && err != sql.ErrNoRows {
		// 查询数据库出错
		return
	}
	if err == sql.ErrNoRows {
		// 用户不存在
		return user, ErrorUserNotExit
	}
	return
}

func UpdatePassword(password string, email string) (err error) {
	sqlStr := `UPDATE user SET password = ? WHERE email = ?`
	_, err = db.Exec(sqlStr, password, email)
	if err != nil {
		errors.New("db exec failed")
	}
	return
}
