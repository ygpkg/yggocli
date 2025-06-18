package {{.DaoLayerName}}

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/openrpacloud/{{.ProjectName}}/apps/{{.AppName}}/models"
	"github.com/openrpacloud/{{.ProjectName}}/apps/{{.AppName}}/models/{{.AppName}}type"
	"github.com/ygpkg/yg-go/apis/apiobj"
	"github.com/ygpkg/yg-go/logs"
	"gorm.io/gorm"
)

type {{.StructName}}Cond struct {
	models.BaseCond
	Filters      []apiobj.Filter
	ID 			 uint
}

type {{.StructName}}Dao struct {
	models.BaseModel
}

func New{{.StructName}}Dao() *{{.StructName}}Dao {
	return &{{.StructName}}Dao{}
}

func (dao *{{.StructName}}Dao) TableName() string {
	return {{.AppName}}type.TableName{{.StructName}}
}

func (dao *{{.StructName}}Dao) WithTx(db *gorm.DB) *{{.StructName}}Dao {
	return &{{.StructName}}Dao{
		BaseModel: models.BaseModel{DBClient: db},
	}
}

func (dao *{{.StructName}}Dao) Insert(ctx *gin.Context, entity *{{.AppName}}type.{{.StructName}}) error {
	db := dao.DB(ctx).Table(dao.TableName())
	if err := db.Create(entity).Error; err != nil {
		return fmt.Errorf("[{{.StructName}}Dao] Insert fail, entity:%s, err: %v", logs.Json(entity), err)
	}
	return nil
}

func (dao *{{.StructName}}Dao) BatchInsert(ctx *gin.Context, entityList {{.AppName}}type.{{.StructName}}List) error {
	if len(entityList) == 0 {
		return fmt.Errorf("[{{.StructName}}Dao] BatchInsert fail, entityList is empty")
	}

	db := dao.DB(ctx).Table(dao.TableName())
	if err := db.Create(entityList).Error; err != nil {
		return fmt.Errorf("[{{.StructName}}Dao] BatchInsert fail, entityList:%s, err: %v", logs.Json(entityList), err)
	}
	return nil
}

func (dao *{{.StructName}}Dao) UpdateByID(ctx *gin.Context, id uint, entity *{{.AppName}}type.{{.StructName}}) error {
	db := dao.DB(ctx).Table(dao.TableName())
	if err := db.Where("id = ?", id).Updates(entity).Error; err != nil {
		return fmt.Errorf("[{{.StructName}}Dao] UpdateByID fail, id:%d, entity:%s, err: %v", id, logs.Json(entity), err)
	}
	return nil
}

func (dao *{{.StructName}}Dao) UpdateMap(ctx *gin.Context, id uint, updateMap map[string]interface{}) error {
	db := dao.DB(ctx).Table(dao.TableName())
	if err := db.Where("id = ?", id).Updates(updateMap).Error; err != nil {
		return fmt.Errorf("[{{.StructName}}Dao] UpdateMap fail, id:%d, updateMap:%s, err: %v", id, logs.Json(entity), err)
	}
	return nil
}

func (dao *{{.StructName}}Dao) Delete(ctx *gin.Context, id uint) error {
	db := dao.DB(ctx).Table(dao.TableName())
	updatedField := map[string]interface{}{
		"deleted_time": time.Now(),
	}
	if err := db.Where("id = ?", id).Updates(updatedField).Error; err != nil {
		return fmt.Errorf("[{{.StructName}}Dao] Delete fail, id:%d, err: %v", id, err)
	}
	return nil
}

func (dao *{{.StructName}}Dao) GetById(ctx *gin.Context, id uint) (*{{.AppName}}type.{{.StructName}}, error) {
	var entity {{.AppName}}type.{{.StructName}}
	db := dao.DB(ctx).Table(dao.TableName())
	if err := db.Where("id = ?", id).Find(&entity).Error; err != nil {
		return fmt.Errorf("[{{.StructName}}Dao] GetById fail, id:%d, err: %v", id, err)
	}
	return &entity, nil
}

func (dao *{{.StructName}}Dao) GetByCond(ctx *gin.Context, cond *{{.StructName}}Cond) (*{{.AppName}}type.{{.StructName}}, error) {
	var entity {{.AppName}}type.{{.StructName}}
	db := dao.DB(ctx).Table(dao.TableName())

	dao.BuildCondition(db, cond)

	if err := db.Find(&entity).Error; err != nil {
		return fmt.Errorf("[{{.StructName}}Dao] GetByCond fail, cond:%s, err: %v", logs.Json(cond), err)
	}
	return &entity, nil
}

func (dao *{{.StructName}}Dao) GetListByCond(ctx *gin.Context, cond *{{.StructName}}Cond) ({{.AppName}}type.{{.StructName}}List, error) {
	var entityList {{.AppName}}type.{{.StructName}}List
	db := dao.DB(ctx).Table(dao.TableName())

	dao.BuildCondition(db, cond)

	if err := db.Find(&entityList).Error; err != nil {
		return fmt.Errorf("[{{.StructName}}Dao] GetListByCond fail, cond:%s, err: %v", logs.Json(cond), err)
	}
	return entityList, nil
}

func (dao *{{.StructName}}Dao) GetPageListByCond(ctx *gin.Context, cond *{{.StructName}}Cond) ({{.AppName}}type.{{.StructName}}List, int64, error) {
	db := dao.DB(ctx).Table(dao.TableName())

	dao.BuildCondition(db, cond)

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return fmt.Errorf("[{{.StructName}}Dao] GetPageListByCond count fail, cond:%s, err: %v", logs.Json(cond), err)
	}
	if cond.Limit > 0 {
		db.Limit(cond.Limit)
	}
	if cond.Offset > 0 {
		db.Offset(cond.Offset)
	}
	var entityList {{.AppName}}type.{{.StructName}}List
	if err := db.Find(&entityList).Error; err != nil {
		return fmt.Errorf("[{{.StructName}}Dao] GetPageListByCond find fail, cond:%s, err: %v", logs.Json(cond), err)
	}
	return entityList, count, nil
}

func (dao *{{.StructName}}Dao) CountByCond(ctx *gin.Context, cond *{{.StructName}}Cond) (int64, error) {
	db := dao.DB(ctx).Table(dao.TableName())

	dao.BuildCondition(db, cond)
	var count int64
	if err := db.Count(&count).Error; err != nil {
		return fmt.Errorf("[{{.StructName}}Dao] CountByCond fail, cond:%s, err: %v", logs.Json(cond), err)
	}
	return count, nil
}

func (dao *{{.StructName}}Dao) BuildCondition(db *gorm.DB, cond *{{.StructName}}Cond) {
	db = dao.BaseModel.BuildBaseCondition(db, dao.TableName(), cond.BaseCond)
	if cond.ID > 0 {
		query := fmt.Sprintf("%s.id = ?", dao.TableName())
		db.Where(query, cond.ID)
	}

	return
} 