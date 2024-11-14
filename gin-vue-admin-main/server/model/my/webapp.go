package my

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type Webapp struct {
	Id       int    `gorm:"primary_key" json:"id"`
	Name     string `gorm:"column:name" json:"name" binding:"required"`
	Provider string `gorm:"column:provider" json:"provider"`
	Remarks  string `gorm:"column:remarks" json:"remarks"`
}

func (Webapp) TableName() string {
	return "my_webapps"
}

type WebappSearchField struct {
	Search string
}

func GetWebappsTotal(field *WebappSearchField) (total int64) {
	db := global.GVA_DB.Model(&Webapp{})
	if field.Search != "" {
		db = db.Where(
			global.GVA_DB.Where("name like ?", "%"+field.Search+"%").
				Or("provider like ?", "%"+field.Search+"%").
				Or("remarks like ?", "%"+field.Search+"%"))
	}
	db.Count(&total)
	return
}

func GetWebapps(page int, pageSize int, field *WebappSearchField) (apps []Webapp) {
	db := global.GVA_DB.Model(&Webapp{})
	if field.Search != "" {
		db = db.Where(
			global.GVA_DB.Where("name like ?", "%"+field.Search+"%").
				Or("provider like ?", "%"+field.Search+"%").
				Or("remarks like ?", "%"+field.Search+"%"))
	}
	//	分页
	if page > 0 && pageSize > 0 {
		db = db.Offset((page - 1) * pageSize).Order("name asc").Limit(pageSize).Find(&apps)
	}
	return
}

func GetWebapp(id int) (app Webapp) {
	global.GVA_DB.Model(&Webapp{}).Where("id = ?", id).First(&app)
	return
}

func AddWebapp(app Webapp) bool {
	global.GVA_DB.Create(&app)
	return true
}

func EditWebapp(id int, app Webapp) bool {
	global.GVA_DB.Model(&Webapp{}).Where("id = ?", id).Updates(app)
	return true
}

func DeleteWebapp(id int) bool {
	global.GVA_DB.Model(&Webapp{}).Where("id = ?", id).Delete(&Webapp{})
	return true
}

func ExistWebappByName(name string) bool {
	var app Webapp
	global.GVA_DB.Model(&Webapp{}).Where("name = ?", name).First(&app)
	if app.Id > 0 {
		return true
	}
	return false
}

func ExistWebappById(id int) bool {
	var app Webapp
	global.GVA_DB.Model(&Webapp{}).Where("id = ?", id).First(&app)
	if app.Id > 0 {
		return true
	}
	return false
}
