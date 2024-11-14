package my

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"gorm.io/datatypes"
)

// result表

type Result struct {
	//gorm.Model
	Id         int            `gorm:"primary_key" json:"id"`
	Vul        bool           `gorm:"column:vul" json:"vul"`
	Detail     datatypes.JSON `gorm:"column:detail" json:"detail"`
	PluginId   string         `gorm:"plugin_id" json:"plugin_id"`
	PluginName string         `gorm:"plugin_name" json:"plugin_name"`
	TaskId     int            `gorm:"task_id" json:"task_id"`
}

func (Result) TableName() string {
	return "my_results"
}

func AddResult(result Result) bool {
	global.GVA_DB.Model(&Result{}).Create(&result)
	return true
}

type ResultSearchField struct {
	Search    string
	TaskField int
	VulField  int
}

func GetResultTotal(field *ResultSearchField) (total int64) {
	db := global.GVA_DB.Model(&Result{})
	if field.TaskField != -1 {
		db = db.Where("task_id = ?", field.TaskField)
	}
	if field.VulField != -1 {
		db = db.Where("vul = ?", field.VulField)
	}
	if field.Search != "" {
		db = db.Where(
			global.GVA_DB.Where("detail like ?", "%"+field.Search+"%")).
			Or("plugin_name like ?", "%"+field.Search+"%")
	}
	db.Count(&total)
	return
}

func GetResult(page int, pageSize int, field *ResultSearchField) (results []Result) {
	db := global.GVA_DB.Model(&Result{})
	if field.TaskField != -1 {
		db = db.Where("task_id = ?", field.TaskField)
	}
	if field.VulField != -1 {
		db = db.Where("vul = ?", field.VulField)
	}
	if field.Search != "" {
		db = db.Where(
			global.GVA_DB.Where("detail like ?", "%"+field.Search+"%"))
	}
	//	分页
	if page > 0 && pageSize > 0 {
		db = db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&results)
	}
	return
}

func DeleteResult(id int) bool {
	global.GVA_DB.Model(&Result{}).Where("id = ?", id).Delete(&Result{})
	return true
}

func ExistResultByID(id int) bool {
	var result Result
	global.GVA_DB.Model(&Result{}).Where("id = ?", id).First(&result)
	if result.Id > 0 {
		return true
	}
	return false
}
