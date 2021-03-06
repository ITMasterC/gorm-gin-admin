package models

import (
	"github.com/olongfen/gorm-gin-admin/src/pkg/query"
	"github.com/olongfen/gorm-gin-admin/src/utils"
	"gorm.io/gorm"
)

const (
	UserStatusRegister = iota // 0 注册状态
	UserStatusLogin           // 1 登录状态
	UserStatusLogout          // 2 登出状态
	UserStatusLock            // 3 封锁状态
)

const (
	UserRoleNormal     = iota // 普通用户
	UserRoleAdmin             // 普通管理员
	UserRoleSuperAdmin        // 超级管理员
)

// UserBase 用户信息
type UserBase struct {
	Model
	Uid      string `json:"uid" gorm:"type:varchar(36); unique_index"`
	Username string `json:"username" gorm:"type:varchar(16); unique_index"`
	Phone    string `json:"phone" gorm:"type:varchar(11); unique_index"`
	LoginPwd string `json:"-" gorm:"type:varchar(128)"`
	PayPwd   string `json:"-" gorm:"type:varchar(128)"`
	Email    string `json:"email" gorm:"type:varchar(32)"`
	Nickname string `json:"nickname" gorm:"type:varchar(12)"`
	HeadIcon string `json:"headIcon"`
	Sign     string `json:"sign" gorm:"type:varchar(256)"`
	Status   int    `json:"status"`
	//
	Role int `json:"role"  gorm:"default:0; index"` // 0 普通用户，1 普通管理员， 2 超级管理员或者开发人员

	// 外键
	// UserCard UserCard `json:"userCard" gorm:"foreignkey:ID"`
}

func NewUserBase() *UserBase {
	return new(UserBase)
}

// Insert 插入一条数据
func (u *UserBase) Insert(options ...*gorm.DB) (err error) {
	if err = getDB(options...).Create(u).Error; err != nil {
		logModel.Errorln(err)
		err = utils.ErrInsertDataFailed
		return
	}
	return
}

func (u *UserBase) Update(uid string, data interface{}, options ...*gorm.DB) (err error) {
	if err = getDB(options...).Model(u).Where("uid = ?", uid).Updates(data).Error; err != nil {
		logModel.Errorln(err)
		err = utils.ErrInsertDataFailed
		return
	}
	return
}

// UpdateUser 更新数据
func (u *UserBase) UpdateUser(uid string, options ...*gorm.DB) (err error) {
	if err = getDB(options...).Model(u).Where("uid = ?", uid).Updates(u).Error; err != nil {
		logModel.Errorln(err)
		err = utils.ErrUpdateDataFailed
		return
	}
	return
}

// UpdateOne 更新一个字段
func (u *UserBase) UpdateOne(uid string, column string, val interface{}, options ...*gorm.DB) (err error) {
	if err = getDB(options...).Model(u).Where("uid = ?", uid).Update(column, val).Error; err != nil {
		logModel.Errorln(err)
		err = utils.ErrUpdateDataFailed
		return
	}
	return
}

// GetById 通过id获取数据
func (u *UserBase) GetById(id int64) (err error) {
	if err = DB.First(u, "id = ?", id).Error; err != nil {
		logModel.Errorln(err)
		err = utils.ErrGetDataFailed
		return
	}
	return
}

// Delete 删除数据
func (u *UserBase) Delete(uid string, options ...*gorm.DB) (err error) {
	if err = getDB(options...).Model(u).Where("uid = ?", uid).Delete(u).Error; err != nil {
		logModel.Errorln(err)
		err = utils.ErrDeleteDataFailed
		return
	}
	return
}

// GetByUId 通过id获取数据
func (u *UserBase) GetByUId(uid string) (err error) {
	if err = DB.First(u, "uid = ?", uid).Error; err != nil {
		logModel.Errorln(err)
		err = utils.ErrGetDataFailed
		return
	}
	return
}

// GetByUsername 通过username获取用户信息
func (u *UserBase) GetByUsername(username string) (err error) {
	if err = DB.Find(u, "username = ?", username).Error; err != nil {
		logModel.Errorln(err)
		err = utils.ErrGetDataFailed
		return
	}
	return
}

// GetByPhone 通过phone获取信息
func (u *UserBase) GetByPhone(phone string) (err error) {
	if err = DB.First(u, "phone = ?", phone).Error; err != nil {
		logModel.Errorln(err)
		err = utils.ErrGetDataFailed
		return
	}
	return
}

// GetUserList 获取用户列表
func GetUserList(q *query.Query) (ret []*UserBase, err error) {
	if err = DB.Model(&UserBase{}).Where(q.Cond, q.Values...).Offset(q.PageNum).Limit(q.PageSize).Order("id asc").Find(&ret).Error; err != nil {
		logModel.Errorln(err)
		err = utils.ErrGetDataFailed
		return
	}
	return
}

// GetUserTotal 获取总数
func GetUserTotal(cond interface{}) (ret int64, err error) {
	var count int64
	if err = DB.Model(&UserBase{}).Where(cond).Count(&count).Error; err != nil {
		logModel.Errorln(err)
		err = utils.ErrGetDataFailed
		return 0, err
	}

	return count, nil
}
