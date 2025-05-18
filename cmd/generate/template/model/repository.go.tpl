package dao{{.PackagePascalName}}

import (
	"fmt"
	"time"
	
    "{{.ServiceName}}/apps/{{.AppName}}/internal/errorCode"
    "{{.ServiceName}}/apps/{{.AppName}}/internal/model"

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

type {{.StructName}}Repo struct {
	model.Base
}

func New{{.StructName}}Repo() *{{.StructName}}Repo {
	return &{{.StructName}}Repo{}
}

func (repo *{{.StructName}}Repo) TableName() string {
	return repo.TableName{{.StructName}}Repo
}

func (repo *{{.StructName}}Repo) WithTx(db *gorm.DB) *{{.StructName}}Repo {
	return &{{.StructName}}Repo{
		Base: model.Base{Tx: db},
	}
}

func (repo *{{.StructName}}Repo) Insert(c *gin.Context, entity *{{.StructName}}Entity) error {
	db := repo.Db(c).Model(&{{.StructName}}Entity{})
	db = db.Table(repo.TableName())
	if err := db.Create(entity).Error; err != nil {
		return errorCode.ErrorDbInsert.Wrapf(err, "[{{.StructName}}Repo] Insert fail, entity:%s", gutils.ToJsonString(entity))
	}
	return nil
}

func (repo *{{.StructName}}Repo) BatchInsert(c *gin.Context, entityList {{.StructName}}EntityList) error {
	if len(entityList) == 0 {
		return errorCode.ErrorDbInsert.Wrapf(nil, "[{{.StructName}}Repo] BatchInsert fail, entityList is empty")
	}

	db := repo.Db(c).Model(&{{.StructName}}Entity{})
	db = db.Table(repo.TableName())
	if err := db.Create(entityList).Error; err != nil {
		return errorCode.ErrorDbInsert.Wrapf(err, "[{{.StructName}}Repo] BatchInsert fail, entityList:%s", gutils.ToJsonString(entityList))
	}
	return nil
}

func (repo *{{.StructName}}Repo) Update(c *gin.Context, entity *{{.StructName}}Entity) error {
	db := repo.Db(c).Model(&{{.StructName}}Entity{})
	db = db.Table(repo.TableName())
	if err := db.Where("id = ?", entity.ID).Updates(entity).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[{{.StructName}}Repo] Update fail, entity:%s", gutils.ToJsonString(entity))
	}
	return nil
}

func (repo *{{.StructName}}Repo) UpdateMap(c *gin.Context, id uint64, updateMap map[string]interface{}) error {
	db := repo.Db(c).Model(&{{.StructName}}Entity{})
	db = db.Table(repo.TableName())
	if err := db.Where("id = ?", id).Updates(updateMap).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[{{.StructName}}Repo] UpdateMap fail, id:%d, updateMap:%s", id, gutils.ToJsonString(updateMap))
	}
	return nil
}

func (repo *{{.StructName}}Repo) Delete(c *gin.Context, id, deletedBy uint64) error {
	db := repo.Db(c).Model(&{{.StructName}}Entity{})
	db = db.Table(repo.TableName())
	updatedField := map[string]interface{}{
		"deleted_time": time.Now(),
		"deleted_by":   deletedBy,
	}
	if err := db.Where("id = ?", id).Updates(updatedField).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[{{.StructName}}Repo] Delete fail, id:%d, deletedBy:%d", id, deletedBy)
	}
	return nil
}

func (repo *{{.StructName}}Repo) GetById(c *gin.Context, id uint64) (*{{.StructName}}Entity, error) {
	var entity {{.StructName}}Entity
	db := repo.Db(c).Model(&{{.StructName}}Entity{})
	db = db.Table(repo.TableName())
	if err := db.Where("id = ?", id).Find(&entity).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[{{.StructName}}Repo] GetById fail, id:%d", id)
	}
	return &entity, nil
}

func (repo *{{.StructName}}Repo) GetByCond(c *gin.Context, cond *{{.StructName}}Cond) (*{{.StructName}}Entity, error) {
	var entity {{.StructName}}Entity
	db := repo.Db(c).Model(&{{.StructName}}Entity{})
	db = db.Table(repo.TableName())

	repo.BuildCondition(db, cond)

	if err := db.Find(&entity).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[{{.StructName}}Repo] GetById fail, cond:%s", gutils.ToJsonString(cond))
	}
	return &entity, nil
}

func (repo *{{.StructName}}Repo) GetListByCond(c *gin.Context, cond *{{.StructName}}Cond) ({{.StructName}}EntityList, error) {
	var entityList {{.StructName}}EntityList
	db := repo.Db(c).Model(&{{.StructName}}Entity{})
	db = db.Table(repo.TableName())

	repo.BuildCondition(db, cond)

	if err := db.Find(&entityList).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[{{.StructName}}Repo] GetListByCond fail, cond:%s", gutils.ToJsonString(cond))
	}
	return entityList, nil
}

func (repo *{{.StructName}}Repo) GetPageListByCond(c *gin.Context, cond *{{.StructName}}Cond) ({{.StructName}}EntityList, int64, error) {
	db := repo.Db(c).Model(&{{.StructName}}Entity{})
	db = db.Table(repo.TableName())

	repo.BuildCondition(db, cond)

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, errorCode.ErrorDbFind.Wrapf(err, "[{{.StructName}}Repo] GetPageListByCond count fail, cond:%s", gutils.ToJsonString(cond))
	}
	if cond.PageSize > 0 && cond.Page > 0 {
		db.Offset((cond.Page - 1) * cond.PageSize).Limit(cond.PageSize)
	}
	var list {{.StructName}}EntityList
	if err := db.Find(&list).Error; err != nil {
		return nil, 0, errorCode.ErrorDbFind.Wrapf(err, "[{{.StructName}}Repo] GetPageListByCond find fail, cond:%s", gutils.ToJsonString(cond))
	}
	return list, count, nil
}

func (repo *{{.StructName}}Repo) CountByCond(c *gin.Context, cond *{{.StructName}}Cond) (int64, error) {
	db := repo.Db(c).Model(&{{.StructName}}Entity{})
	db = db.Table(repo.TableName())

	repo.BuildCondition(db, cond)
	var count int64
	if err := db.Count(&count).Error; err != nil {
		return 0, errorCode.ErrorDbFind.Wrapf(err, "[{{.StructName}}Repo] CountByCond fail, cond:%s", gutils.ToJsonString(cond))
	}
	return count, nil
}

func (repo *{{.StructName}}Repo) BuildCondition(db *gorm.DB, cond *{{.StructName}}Cond) {
	if cond.ID > 0 {
		query := fmt.Sprintf("%s.id = ?", repo.TableName())
		db.Where(query, cond.ID)
	}
	if len(cond.IDs) > 0 {
		query := fmt.Sprintf("%s.id in (?)", repo.TableName())
		db.Where(query, cond.IDs)
	}
	if cond.CreatedAtStart > 0 {
		query := fmt.Sprintf("%s.created_at >= ?", repo.TableName())
		db.Where(query, time.Unix(cond.CreatedAtStart, 0))
	}
	if cond.CreatedAtEnd > 0 {
		query := fmt.Sprintf("%s.created_at <= ?", repo.TableName())
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