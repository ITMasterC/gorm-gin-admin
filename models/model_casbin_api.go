package models

import "github.com/olongfen/user_base/utils"

type RuleAPI struct {
	ID     uint   `json:"id" gorm:"column:id"`
	PType  string `json:"pType" gorm:"column:p_type"`
	Uid    string `json:"uid" gorm:"column:v0"`
	Path   string `json:"path" gorm:"column:v1"`
	Method string `json:"method" gorm:"column:v2"`
}

func RuleAPITableName() string {
	return "casbin_rule"
}

func (c *RuleAPI) InsertRuleAPI() (err error) {
	if err = DB.Table(RuleAPITableName()).Create(c).Error; err != nil {
		logModel.Errorln("[InsertRuleAPI] err: ", err)
		err = utils.ErrGetDataFailed
		return err
	}
	return
}

func (c *RuleAPI) UpdateRuleAPI(path string, method string, m map[string]interface{}) (err error) {
	if err = DB.Table(RuleAPITableName()).Updates(m).Where("v1 = ? and v2 = ?", path, method).Error; err != nil {
		logModel.Errorln("[UpdateRuleAPI] err: ", err)
		err = utils.ErrUpdateDataFailed
		return err
	}
	return
}

func (c *RuleAPI) DeleteRuleAPI(id int64) (err error) {
	if err = DB.Table(RuleAPITableName()).Delete(c, "id = ?", id).Error; err != nil {
		logModel.Errorln("[DeleteRuleAPI] err: ", err)
		err = utils.ErrDeleteDataFailed
		return err
	}
	return
}

func GetRuleAPIListByUID(uid string) (ret []*RuleAPI, err error) {
	if err = DB.Table(RuleAPITableName()).Where("uid = ?", uid).Find(&ret).Order("id asc").Error; err != nil {
		logModel.Errorln("[GetRuleAPIListByUID] err: ", err)
		err = utils.ErrGetDataFailed
		return nil, err
	}
	return
}

//

func (c *RuleAPI) DeleteRuleAPIData(path string, method string) (err error) {
	if err = DB.Table(RuleAPITableName()).Delete(c, "v1 = ? and v2 = ?", path, method).Error; err != nil {
		logModel.Errorln("[DeleteRuleAPIData] err: ", err)
		err = utils.ErrDeleteDataFailed
		return err
	}
	return
}