package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/olongfen/contrib/session"
	"github.com/olongfen/user_base/src/models"
	"github.com/olongfen/user_base/src/pkg/app"
	"github.com/olongfen/user_base/src/pkg/codes"
	"github.com/olongfen/user_base/src/service"
	"github.com/olongfen/user_base/src/utils"
	"strconv"
)

// @tags 管理员
// @Summary 管理员登录
// @Produce json
// @Accept json
// @Param username body string true "用户名"
// @Param password body string true "密码"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @router /api/v1/admin/login [post]
// AdminLogin
func AdminLogin(c *gin.Context) {
	var (
		form  = &utils.LoginForm{}
		err   error
		code  = codes.CodeProcessingFailed
		token string
	)
	defer func() {
		if err != nil {
			app.NewGinResponse(c).Fail(code, err.Error()).Response()
		} else {
			app.NewGinResponse(c).Success(map[string]string{"token": token}).Response()
		}
	}()

	form.IP = c.ClientIP()
	if err = c.ShouldBind(form); err != nil {
		code = codes.CodeParamInvalid
		return
	}
	if token, err = service.UserLogin(form, true); err != nil {
		return
	}
}

// AdminLogout 登出
// @tags 管理员
// @Summary 管理员登出
// @Produce json
// @Accept json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @router /api/v1/admin/logout [post]
func AdminLogout(c *gin.Context) {
	var (
		err  error
		code = codes.CodeProcessingFailed
		s    *session.Session
	)
	defer func() {
		if err != nil {
			app.NewGinResponse(c).Fail(code, err.Error()).Response()
		} else {
			app.NewGinResponse(c).Success(nil).Response()
		}
	}()
	if s, code, err = GetSession(c); err != nil {
		return
	}
	if err = service.UserLogout(s.UID); err != nil {
		return
	}

}

// @tags 管理员
// @Title 获取用户列表
// @Summary 获取用户列表
// @Description 获取用户列表
// @Accept json
// @Produce json
// @Param {} body utils.FormUserList true "查询数据"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Path /api/v1/admin/listUser [get]
func ListUser(c *gin.Context) {
	var (
		err  error
		code = codes.CodeProcessingFailed
		form = new(utils.FormUserList)
		data []*models.UserBase
		s    *session.Session
	)
	defer func() {
		if err != nil {
			app.NewGinResponse(c).Fail(code, err.Error()).Response()
		}
	}()

	if s, code, err = GetSession(c); err != nil {
		return
	}
	if err = c.ShouldBindQuery(form); err != nil {
		code = codes.CodeParamInvalid
		return
	}
	if data, err = service.GetUserList(s.UID, form); err != nil {
		return
	}
	app.NewGinResponse(c).Success(data).Response()
}

// @tags 管理员
// @Title 获取用户总数
// @Summary 获取用户总数
// @Description 获取用户总数
// @Accept json
// @Produce json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Path   /api/v1/admin/userTotal  []
func UserTotal(c *gin.Context) {
	var (
		err  error
		code = codes.CodeProcessingFailed
		data int64
		s    *session.Session
	)
	defer func() {
		if err != nil {
			app.NewGinResponse(c).Fail(code, err.Error()).Response()
		}
	}()

	if s, code, err = GetSession(c); err != nil {
		return
	}

	if data, err = service.GetUserCount(s.UID); err != nil {
		return
	}
	app.NewGinResponse(c).Success(data).Response()
}

// @tags 管理员
// @Title
// @Summary
// @Description
// @Accept json
// @Produce json
// @Param uid path string true "用户uid"
// @Param {} body utils.FormEditUser true "修改用户信息"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Path /api/v1/admin/editUser [post]
func EditUser(c *gin.Context) {
	var (
		err  error
		form = new(utils.FormEditUser)
		s    *session.Session

		code = codes.CodeProcessingFailed
	)
	defer func() {
		if err != nil {
			app.NewGinResponse(c).Fail(code, err.Error()).Response()
		}
	}()
	if s, code, err = GetSessionAndBindingForm(form, c); err != nil {
		return
	}
	if _, err = service.EditUserByRole(s.UID, form); err != nil {
		return
	}
}

// @tags 管理员
// @Summary 删除用户
// @Description 删除用户
// @Accept json
// @Produce json
// @Param uid query string true "用户UID"
// @Success 200  {object}
// @Failure 500  {object}
// @router /api/v1/admin/deleteUser [delete]
func DeleteUser(c *gin.Context) {
	var (
		err error
		uid string
		s   *session.Session

		code = codes.CodeProcessingFailed
	)
	uid = c.Query("uid")
	defer func() {
		if err != nil {
			app.NewGinResponse(c).Fail(code, err.Error()).Response()
		} else {
			app.NewGinResponse(c).Success(nil).Response()
		}
	}()
	if s, code, err = GetSession(c); err != nil {
		return
	}
	if err = service.DeleteUser(s.UID, uid); err != nil {
		return
	}
}

// @tags 管理员
// @Title 添加角色接口权限
// @Summary 添加角色接口权限
// @Description 添加角色接口权限
// @Accept json
// @Produce json
// @Param {} body utils.FormRoleAPIPerm true "添加api权限表单"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Path /api/v1/admin/addRoleAPIPerm [post]
func AddRoleAPIPerm(c *gin.Context) {
	var (
		f    = &utils.FormRoleAPIPerm{}
		err  error
		code = codes.CodeProcessingFailed
		ret  []int64
		s    *session.Session
	)
	defer func() {
		if err != nil {
			app.NewGinResponse(c).Fail(code, err.Error()).Response()
		}
	}()
	if s, code, err = GetSessionAndBindingForm(f, c); err != nil {
		return
	}
	if s.UID == f.Uid {
		err = utils.ErrActionNotAllow
		return
	}
	if ret, err = service.AddRuleAPI(f); err != nil {
		return
	}
	app.NewGinResponse(c).Success(ret).Response()
}

// @tags 管理员
// @Title 删除角色接口权限
// @Summary 删除角色接口权限
// @Description 删除角色接口权限
// @Accept json
// @Produce json
// @Param {} body utils.FormRoleAPIPerm true "删除api权限表单"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Path /api/v1/admin/removeRoleAPIPerm [delete]
func RemoveRolePermAPI(c *gin.Context) {
	var (
		f    = &utils.FormRoleAPIPerm{}
		err  error
		code = codes.CodeProcessingFailed
		ret  []int64
		s    *session.Session
	)
	defer func() {
		if err != nil {
			app.NewGinResponse(c).Fail(code, err.Error()).Response()
		}
	}()
	if s, code, err = GetSessionAndBindingForm(f, c); err != nil {
		return
	}
	if s.UID == f.Uid {
		err = utils.ErrActionNotAllow
		return
	}
	if ret, err = service.RemoveRuleAPI(f); err != nil {
		return
	}
	app.NewGinResponse(c).Success(ret).Response()
}

// @tags 管理员
// @Title 获取用户权限
// @Summary 获取用户权限
// @Description 获取用户权限
// @Accept json
// @Produce json
// @Param uid query string false "用户uid,不输入默认返回自己uid"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Path /api/v1/getRoleApiList [get]
func GetRoleApiList(c *gin.Context) {
	var (
		err  error
		code = codes.CodeProcessingFailed
		s    *session.Session
		uid  string
		data []struct {
			Path   string
			Method string
		}
	)
	defer func() {
		if err != nil {
			app.NewGinResponse(c).Fail(code, err.Error()).Response()
		}
	}()
	if s, code, err = GetSession(c); err != nil {
		return
	}
	uid = c.Query("uid")
	if len(uid) == 0 {
		uid = s.UID
	}
	if data, err = service.GetRuleApiList(uid); err != nil {
		return
	}
	app.NewGinResponse(c).Success(data).Response()

}

// @tags 管理员
// @Title 获取全部api
// @Summary 获取全部api
// @Description 获取全部api
// @Accept json
// @Produce json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Path  /api/v1/admin/getAllApiGroup [get]
func GetAllAPIGroup(c *gin.Context) {
	var (
		err  error
		code = codes.CodeProcessingFailed
		ret  []*models.APIGroup
	)
	defer func() {
		if err != nil {
			app.NewGinResponse(c).Fail(code, err.Error()).Response()
		}
	}()
	if _, code, err = GetSession(c); err != nil {
		return
	}
	if ret, err = service.GetAPIGroupList(); err != nil {
		return
	}
	app.NewGinResponse(c).Success(ret).Response()
}

// @tags 管理员
// @Title 创建api
// @Summary 创建api
// @Description
// @Accept json
// @Produce json
// @Param {array} body utils.FormAPIGroupAdd true "api数组"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Path /api/v1/admin/addApiGroup [post]
func AddApiGroup(c *gin.Context) {
	var (
		err  error
		code = codes.CodeProcessingFailed
		f    []*utils.FormAPIGroupAdd
		ret  []*models.APIGroup
	)
	defer func() {
		if err != nil {
			app.NewGinResponse(c).Fail(code, err.Error()).Response()
		}
	}()
	if _, code, err = GetSessionAndBindingForm(f, c); err != nil {
		return
	}
	if ret, err = service.AddAPIGroup(f); err != nil {

		return
	}
	app.NewGinResponse(c).Success(ret).Response()
}

// @tags 管理员
// @Title 删除api
// @Summary  删除api
// @Description 删除api
// @Accept json
// @Produce json
// @Param id query int true "id"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Path /api/v1/admin/removeApiGroup [delete]
func RemoveApiGroup(c *gin.Context) {
	var (
		err  error
		code = codes.CodeProcessingFailed
		id   string
	)
	defer func() {
		if err != nil {
			app.NewGinResponse(c).Fail(code, err.Error()).Response()
		}
	}()
	id = c.Query("id")
	_id, err_ := strconv.ParseUint(id, 10, 64)
	if err_ != nil {
		code = codes.CodeParamInvalid
		err = err_
		return
	}
	if err = service.DelAPIGroup(int64(_id)); err != nil {
		return
	}
	app.NewGinResponse(c).Success(nil).Response()
}

// @tags 管理员
// @Title 修改api
// @Summary 修改api
// @Description 修改api
// @Accept json
// @Produce json
// @Param {} body utils.FormAPIGroupEdit true "表单"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Path  /api/v1/admin/editApiGroup [put]
func EditApiGroup(c *gin.Context) {
	var (
		f    = &utils.FormAPIGroupEdit{}
		err  error
		code = codes.CodeProcessingFailed
		ret  *models.APIGroup
	)
	defer func() {
		if err != nil {
			app.NewGinResponse(c).Fail(code, err.Error()).Response()
		}
	}()
	if _, code, err = GetSessionAndBindingForm(f, c); err != nil {
		return
	}

	if ret, err = service.EditAPIGroup(f); err != nil {

		return
	}
	app.NewGinResponse(c).Success(ret).Response()

}
