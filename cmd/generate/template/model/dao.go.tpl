package dao{{.PackageName}}

import (
	"fmt"
	"time"
	
    "{{.ProjectAppRelativePath}}/internal/code"
    {{- if isDefaultModelLayer .ModelLayerName}}
    "{{.ProjectAppRelativePath}}/internal/model"
    {{- else}}
    "{{.ProjectAppRelativePath}}/internal/model/{{.ModelLayerName}}"
    {{- end}}

    "github.com/gin-gonic/gin"
    "github.com/morehao/golib/gutils"
    "gorm.io/gorm"
)

type {{.StructName}}Cond struct {
	ID             uint64
	IDs            []uint64
	IsDelete       bool
	Page           int
	PageSize       int
	CreatedAtStart int64
	CreatedAtEnd   int64
	OrderField     string
}

type {{.StructName}}Dao struct {
	{{.ModelLayerName}}.Base
}

func New{{.StructName}}Dao() *{{.StructName}}Dao {
	return &{{.StructName}}Dao{}
}

func (dao *{{.StructName}}Dao) TableName() string {
	return dao.TableNamedao{{.StructName}}
}

func (dao *{{.StructName}}Dao) WithTx(db *gorm.DB) *{{.StructName}}Dao {
	return &{{.StructName}}Dao{
		Base: model.Base{Tx: db},
	}
}

func (dao *{{.StructName}}Dao) Insert(c *gin.Context, entity *{{.ModelLayerName}}.{{.StructName}}Entity) error {
	db := dao.Db(c).Model(&{{.ModelLayerName}}.{{.StructName}}Entity{})
	db = db.Table(dao.TableName())
	if err := db.Create(entity).Error; err != nil {
		return code.ErrorDbInsert.Wrapf(err, "[{{.StructName}}Dao] Insert fail, entity:%s", gutils.ToJsonString(entity))
	}
	return nil
}

func (dao *{{.StructName}}Dao) BatchInsert(c *gin.Context, entityList {{.ModelLayerName}}.{{.StructName}}EntityList) error {
	if len(entityList) == 0 {
		return code.ErrorDbInsert.Wrapf(nil, "[{{.StructName}}Dao] BatchInsert fail, entityList is empty")
	}

	db := dao.Db(c).Model(&{{.ModelLayerName}}.{{.StructName}}Entity{})
	db = db.Table(dao.TableName())
	if err := db.Create(entityList).Error; err != nil {
		return code.ErrorDbInsert.Wrapf(err, "[{{.StructName}}Dao] BatchInsert fail, entityList:%s", gutils.ToJsonString(entityList))
	}
	return nil
}

func (dao *{{.StructName}}Dao) Update(c *gin.Context, entity *{{.ModelLayerName}}.{{.StructName}}Entity) error {
	db := dao.Db(c).Model(&{{.ModelLayerName}}.{{.StructName}}Entity{})
	db = db.Table(dao.TableName())
	if err := db.Where("id = ?", entity.ID).Updates(entity).Error; err != nil {
		return code.ErrorDbUpdate.Wrapf(err, "[{{.StructName}}Dao] Update fail, entity:%s", gutils.ToJsonString(entity))
	}
	return nil
}

func (dao *{{.StructName}}Dao) UpdateMap(c *gin.Context, id uint64, updateMap map[string]interface{}) error {
	db := dao.Db(c).Model(&{{.ModelLayerName}}.{{.StructName}}Entity{})
	db = db.Table(dao.TableName())
	if err := db.Where("id = ?", id).Updates(updateMap).Error; err != nil {
		return code.ErrorDbUpdate.Wrapf(err, "[{{.StructName}}Dao] UpdateMap fail, id:%d, updateMap:%s", id, gutils.ToJsonString(updateMap))
	}
	return nil
}

func (dao *{{.StructName}}Dao) Delete(c *gin.Context, id, deletedBy uint64) error {
	db := dao.Db(c).Model(&{{.ModelLayerName}}.{{.StructName}}Entity{})
	db = db.Table(dao.TableName())
	updatedField := map[string]interface{}{
		"deleted_time": time.Now(),
		"deleted_by":   deletedBy,
	}
	if err := db.Where("id = ?", id).Updates(updatedField).Error; err != nil {
		return code.ErrorDbUpdate.Wrapf(err, "[{{.StructName}}Dao] Delete fail, id:%d, deletedBy:%d", id, deletedBy)
	}
	return nil
}

func (dao *{{.StructName}}Dao) GetById(c *gin.Context, id uint64) (*{{.ModelLayerName}}.{{.StructName}}Entity, error) {
	var entity {{.StructName}}Entity
	db := dao.Db(c).Model(&{{.StructName}}Entity{})
	db = db.Table(dao.TableName())
	if err := db.Where("id = ?", id).Find(&entity).Error; err != nil {
		return nil, code.ErrorDbFind.Wrapf(err, "[{{.StructName}}Dao] GetById fail, id:%d", id)
	}
	return &entity, nil
}

func (dao *{{.StructName}}Dao) GetByCond(c *gin.Context, cond *{{.StructName}}Cond) (*{{.ModelLayerName}}.{{.StructName}}Entity, error) {
	var entity {{.ModelLayerName}}.{{.StructName}}Entity
	db := dao.Db(c).Model(&{{.ModelLayerName}}.{{.StructName}}Entity{})
	db = db.Table(dao.TableName())

	dao.BuildCondition(db, cond)

	if err := db.Find(&entity).Error; err != nil {
		return nil, code.ErrorDbFind.Wrapf(err, "[{{.StructName}}Dao] GetById fail, cond:%s", gutils.ToJsonString(cond))
	}
	return &entity, nil
}

func (dao *{{.StructName}}Dao) GetListByCond(c *gin.Context, cond *{{.StructName}}Cond) ({{.ModelLayerName}}.{{.StructName}}EntityList, error) {
	var entityList {{.ModelLayerName}}.{{.StructName}}EntityList
	db := dao.Db(c).Model(&{{.StructName}}Entity{})
	db = db.Table(dao.TableName())

	dao.BuildCondition(db, cond)

	if err := db.Find(&entityList).Error; err != nil {
		return nil, code.ErrorDbFind.Wrapf(err, "[{{.StructName}}Dao] GetListByCond fail, cond:%s", gutils.ToJsonString(cond))
	}
	return entityList, nil
}

func (dao *{{.StructName}}Dao) GetPageListByCond(c *gin.Context, cond *{{.StructName}}Cond) ({{.ModelLayerName}}.{{.StructName}}EntityList, int64, error) {
	db := dao.Db(c).Model(&{{.ModelLayerName}}.{{.StructName}}Entity{})
	db = db.Table(dao.TableName())

	dao.BuildCondition(db, cond)

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, code.ErrorDbFind.Wrapf(err, "[{{.StructName}}Dao] GetPageListByCond count fail, cond:%s", gutils.ToJsonString(cond))
	}
	if cond.PageSize > 0 && cond.Page > 0 {
		db.Offset((cond.Page - 1) * cond.PageSize).Limit(cond.PageSize)
	}
	var list {{.StructName}}EntityList
	if err := db.Find(&list).Error; err != nil {
		return nil, 0, code.ErrorDbFind.Wrapf(err, "[{{.StructName}}Dao] GetPageListByCond find fail, cond:%s", gutils.ToJsonString(cond))
	}
	return list, count, nil
}

func (dao *{{.StructName}}Dao) CountByCond(c *gin.Context, cond *{{.StructName}}Cond) (int64, error) {
	db := dao.Db(c).Model(&{{.StructName}}Entity{})
	db = db.Table(dao.TableName())

	dao.BuildCondition(db, cond)
	var count int64
	if err := db.Count(&count).Error; err != nil {
		return 0, code.ErrorDbFind.Wrapf(err, "[{{.StructName}}Dao] CountByCond fail, cond:%s", gutils.ToJsonString(cond))
	}
	return count, nil
}

func (dao *{{.StructName}}Dao) BuildCondition(db *gorm.DB, cond *{{.StructName}}Cond) {
	if cond.ID > 0 {
		query := fmt.Sprintf("%s.id = ?", dao.TableName())
		db.Where(query, cond.ID)
	}
	if len(cond.IDs) > 0 {
		query := fmt.Sprintf("%s.id in (?)", dao.TableName())
		db.Where(query, cond.IDs)
	}
	if cond.CreatedAtStart > 0 {
		query := fmt.Sprintf("%s.created_at >= ?", dao.TableName())
		db.Where(query, time.Unix(cond.CreatedAtStart, 0))
	}
	if cond.CreatedAtEnd > 0 {
		query := fmt.Sprintf("%s.created_at <= ?", dao.TableName())
		db.Where(query, time.Unix(cond.CreatedAtEnd, 0))
	}
	if cond.IsDelete {
		db.Unscoped()
	}

	if cond.OrderField != "" {
		db.Order(cond.OrderField)
	}

	return
} 