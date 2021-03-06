package service

import (
	"github.com/olongfen/gorm-gin-admin/src/models"
	"github.com/olongfen/gorm-gin-admin/src/utils"
)

// AddAPIGroup
func AddAPIGroup(uid string,f []*utils.FormAPIGroupAdd) (ret []*models.APIGroup, err error) {
	var (

		datas []*models.APIGroup
	)
	role := new(models.UserBase)
	if err = role.GetByUId(uid);err!=nil{
		return
	}
	// role must super admin
	if role.Role!=models.UserRoleSuperAdmin{
		err =utils.ErrActionNotAllow
		return
	}
	for _, v := range f {
		var (
			data = new(models.APIGroup)
		)
		data.ApiGroup = v.ApiGroup
		data.Description = v.Description
		data.Path = v.Path
		data.Method = v.Method

		datas = append(datas, data)
	}
	if err = models.BatchInsertAPIGroup(datas); err != nil {
		return nil, err
	}
	return models.GetAPIGroupList()
}

func EditAPIGroup(uid string,f *utils.FormAPIGroupEdit) (ret *models.APIGroup, err error) {
	var (
		data       = new(models.APIGroup)
		dataCasbin = new(models.RuleAPI)
		m          = map[string]interface{}{}
	)
	role := new(models.UserBase)
	if err = role.GetByUId(uid);err!=nil{
		return
	}
	// role must super admin
	if role.Role!=models.UserRoleSuperAdmin{
		err =utils.ErrActionNotAllow
		return
	}
	if err = data.Get(f.Id); err != nil {
		return nil, err
	}
	path := data.Path
	method := data.Method
	if len(f.ApiGroup) != 0 {
		data.ApiGroup = f.ApiGroup
	}
	if len(f.Description) != 0 {
		data.Description = f.Description
	}
	if len(f.Path) != 0 {
		m["v1"] = f.Path
		data.Path = f.Path
	}
	if len(f.Method) != 0 {
		data.Method = f.Method
		m["v2"] = f.Method
	}
	if len(m) > 0 {
		if err = dataCasbin.Update(path, method, m); err != nil {
			logServe.Errorln(err)
			err = nil
		}
	}
	if err = data.Update(f.Id, data); err != nil {
		return nil, err
	}

	ret = data
	return
}

func DelAPIGroup(uid string,id int64) (err error) {
	var (
		data       = new(models.APIGroup)
		dataCasbin = new(models.RuleAPI)
	)
	role := new(models.UserBase)
	if err = role.GetByUId(uid);err!=nil{
		return
	}
	// role must super admin
	if role.Role!=models.UserRoleSuperAdmin{
		err =utils.ErrActionNotAllow
		return
	}
	if err = data.Get(id); err != nil {
		return
	}

	if err = dataCasbin.DeleteByPathAndMethod(data.Path, data.Method); err != nil {
		return err
	}

	return data.Delete(id)
}

// GetAPIGroupList
func GetAPIGroupList() (ret []*models.APIGroup, err error) {
	return models.GetAPIGroupList()
}
