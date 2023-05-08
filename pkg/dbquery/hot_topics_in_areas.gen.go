// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dbquery

import (
	"context"
	"whale/pkg/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"
)

func newHotTopicsInArea(db *gorm.DB, opts ...gen.DOOption) hotTopicsInArea {
	_hotTopicsInArea := hotTopicsInArea{}

	_hotTopicsInArea.hotTopicsInAreaDo.UseDB(db, opts...)
	_hotTopicsInArea.hotTopicsInAreaDo.UseModel(&models.HotTopicsInArea{})

	tableName := _hotTopicsInArea.hotTopicsInAreaDo.TableName()
	_hotTopicsInArea.ALL = field.NewAsterisk(tableName)
	_hotTopicsInArea.CityID = field.NewString(tableName, "city_id")
	_hotTopicsInArea.TopicMetrics = field.NewField(tableName, "topic_metrics")
	_hotTopicsInArea.UpdatedAt = field.NewTime(tableName, "updated_at")
	_hotTopicsInArea.CreatedAt = field.NewTime(tableName, "created_at")

	_hotTopicsInArea.fillFieldMap()

	return _hotTopicsInArea
}

type hotTopicsInArea struct {
	hotTopicsInAreaDo hotTopicsInAreaDo

	ALL          field.Asterisk
	CityID       field.String
	TopicMetrics field.Field
	UpdatedAt    field.Time
	CreatedAt    field.Time

	fieldMap map[string]field.Expr
}

func (h hotTopicsInArea) Table(newTableName string) *hotTopicsInArea {
	h.hotTopicsInAreaDo.UseTable(newTableName)
	return h.updateTableName(newTableName)
}

func (h hotTopicsInArea) As(alias string) *hotTopicsInArea {
	h.hotTopicsInAreaDo.DO = *(h.hotTopicsInAreaDo.As(alias).(*gen.DO))
	return h.updateTableName(alias)
}

func (h *hotTopicsInArea) updateTableName(table string) *hotTopicsInArea {
	h.ALL = field.NewAsterisk(table)
	h.CityID = field.NewString(table, "city_id")
	h.TopicMetrics = field.NewField(table, "topic_metrics")
	h.UpdatedAt = field.NewTime(table, "updated_at")
	h.CreatedAt = field.NewTime(table, "created_at")

	h.fillFieldMap()

	return h
}

func (h *hotTopicsInArea) WithContext(ctx context.Context) IHotTopicsInAreaDo {
	return h.hotTopicsInAreaDo.WithContext(ctx)
}

func (h hotTopicsInArea) TableName() string { return h.hotTopicsInAreaDo.TableName() }

func (h hotTopicsInArea) Alias() string { return h.hotTopicsInAreaDo.Alias() }

func (h *hotTopicsInArea) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := h.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (h *hotTopicsInArea) fillFieldMap() {
	h.fieldMap = make(map[string]field.Expr, 4)
	h.fieldMap["city_id"] = h.CityID
	h.fieldMap["topic_metrics"] = h.TopicMetrics
	h.fieldMap["updated_at"] = h.UpdatedAt
	h.fieldMap["created_at"] = h.CreatedAt
}

func (h hotTopicsInArea) clone(db *gorm.DB) hotTopicsInArea {
	h.hotTopicsInAreaDo.ReplaceConnPool(db.Statement.ConnPool)
	return h
}

func (h hotTopicsInArea) replaceDB(db *gorm.DB) hotTopicsInArea {
	h.hotTopicsInAreaDo.ReplaceDB(db)
	return h
}

type hotTopicsInAreaDo struct{ gen.DO }

type IHotTopicsInAreaDo interface {
	gen.SubQuery
	Debug() IHotTopicsInAreaDo
	WithContext(ctx context.Context) IHotTopicsInAreaDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IHotTopicsInAreaDo
	WriteDB() IHotTopicsInAreaDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IHotTopicsInAreaDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IHotTopicsInAreaDo
	Not(conds ...gen.Condition) IHotTopicsInAreaDo
	Or(conds ...gen.Condition) IHotTopicsInAreaDo
	Select(conds ...field.Expr) IHotTopicsInAreaDo
	Where(conds ...gen.Condition) IHotTopicsInAreaDo
	Order(conds ...field.Expr) IHotTopicsInAreaDo
	Distinct(cols ...field.Expr) IHotTopicsInAreaDo
	Omit(cols ...field.Expr) IHotTopicsInAreaDo
	Join(table schema.Tabler, on ...field.Expr) IHotTopicsInAreaDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IHotTopicsInAreaDo
	RightJoin(table schema.Tabler, on ...field.Expr) IHotTopicsInAreaDo
	Group(cols ...field.Expr) IHotTopicsInAreaDo
	Having(conds ...gen.Condition) IHotTopicsInAreaDo
	Limit(limit int) IHotTopicsInAreaDo
	Offset(offset int) IHotTopicsInAreaDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IHotTopicsInAreaDo
	Unscoped() IHotTopicsInAreaDo
	Create(values ...*models.HotTopicsInArea) error
	CreateInBatches(values []*models.HotTopicsInArea, batchSize int) error
	Save(values ...*models.HotTopicsInArea) error
	First() (*models.HotTopicsInArea, error)
	Take() (*models.HotTopicsInArea, error)
	Last() (*models.HotTopicsInArea, error)
	Find() ([]*models.HotTopicsInArea, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.HotTopicsInArea, err error)
	FindInBatches(result *[]*models.HotTopicsInArea, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*models.HotTopicsInArea) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IHotTopicsInAreaDo
	Assign(attrs ...field.AssignExpr) IHotTopicsInAreaDo
	Joins(fields ...field.RelationField) IHotTopicsInAreaDo
	Preload(fields ...field.RelationField) IHotTopicsInAreaDo
	FirstOrInit() (*models.HotTopicsInArea, error)
	FirstOrCreate() (*models.HotTopicsInArea, error)
	FindByPage(offset int, limit int) (result []*models.HotTopicsInArea, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IHotTopicsInAreaDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (h hotTopicsInAreaDo) Debug() IHotTopicsInAreaDo {
	return h.withDO(h.DO.Debug())
}

func (h hotTopicsInAreaDo) WithContext(ctx context.Context) IHotTopicsInAreaDo {
	return h.withDO(h.DO.WithContext(ctx))
}

func (h hotTopicsInAreaDo) ReadDB() IHotTopicsInAreaDo {
	return h.Clauses(dbresolver.Read)
}

func (h hotTopicsInAreaDo) WriteDB() IHotTopicsInAreaDo {
	return h.Clauses(dbresolver.Write)
}

func (h hotTopicsInAreaDo) Session(config *gorm.Session) IHotTopicsInAreaDo {
	return h.withDO(h.DO.Session(config))
}

func (h hotTopicsInAreaDo) Clauses(conds ...clause.Expression) IHotTopicsInAreaDo {
	return h.withDO(h.DO.Clauses(conds...))
}

func (h hotTopicsInAreaDo) Returning(value interface{}, columns ...string) IHotTopicsInAreaDo {
	return h.withDO(h.DO.Returning(value, columns...))
}

func (h hotTopicsInAreaDo) Not(conds ...gen.Condition) IHotTopicsInAreaDo {
	return h.withDO(h.DO.Not(conds...))
}

func (h hotTopicsInAreaDo) Or(conds ...gen.Condition) IHotTopicsInAreaDo {
	return h.withDO(h.DO.Or(conds...))
}

func (h hotTopicsInAreaDo) Select(conds ...field.Expr) IHotTopicsInAreaDo {
	return h.withDO(h.DO.Select(conds...))
}

func (h hotTopicsInAreaDo) Where(conds ...gen.Condition) IHotTopicsInAreaDo {
	return h.withDO(h.DO.Where(conds...))
}

func (h hotTopicsInAreaDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) IHotTopicsInAreaDo {
	return h.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (h hotTopicsInAreaDo) Order(conds ...field.Expr) IHotTopicsInAreaDo {
	return h.withDO(h.DO.Order(conds...))
}

func (h hotTopicsInAreaDo) Distinct(cols ...field.Expr) IHotTopicsInAreaDo {
	return h.withDO(h.DO.Distinct(cols...))
}

func (h hotTopicsInAreaDo) Omit(cols ...field.Expr) IHotTopicsInAreaDo {
	return h.withDO(h.DO.Omit(cols...))
}

func (h hotTopicsInAreaDo) Join(table schema.Tabler, on ...field.Expr) IHotTopicsInAreaDo {
	return h.withDO(h.DO.Join(table, on...))
}

func (h hotTopicsInAreaDo) LeftJoin(table schema.Tabler, on ...field.Expr) IHotTopicsInAreaDo {
	return h.withDO(h.DO.LeftJoin(table, on...))
}

func (h hotTopicsInAreaDo) RightJoin(table schema.Tabler, on ...field.Expr) IHotTopicsInAreaDo {
	return h.withDO(h.DO.RightJoin(table, on...))
}

func (h hotTopicsInAreaDo) Group(cols ...field.Expr) IHotTopicsInAreaDo {
	return h.withDO(h.DO.Group(cols...))
}

func (h hotTopicsInAreaDo) Having(conds ...gen.Condition) IHotTopicsInAreaDo {
	return h.withDO(h.DO.Having(conds...))
}

func (h hotTopicsInAreaDo) Limit(limit int) IHotTopicsInAreaDo {
	return h.withDO(h.DO.Limit(limit))
}

func (h hotTopicsInAreaDo) Offset(offset int) IHotTopicsInAreaDo {
	return h.withDO(h.DO.Offset(offset))
}

func (h hotTopicsInAreaDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IHotTopicsInAreaDo {
	return h.withDO(h.DO.Scopes(funcs...))
}

func (h hotTopicsInAreaDo) Unscoped() IHotTopicsInAreaDo {
	return h.withDO(h.DO.Unscoped())
}

func (h hotTopicsInAreaDo) Create(values ...*models.HotTopicsInArea) error {
	if len(values) == 0 {
		return nil
	}
	return h.DO.Create(values)
}

func (h hotTopicsInAreaDo) CreateInBatches(values []*models.HotTopicsInArea, batchSize int) error {
	return h.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (h hotTopicsInAreaDo) Save(values ...*models.HotTopicsInArea) error {
	if len(values) == 0 {
		return nil
	}
	return h.DO.Save(values)
}

func (h hotTopicsInAreaDo) First() (*models.HotTopicsInArea, error) {
	if result, err := h.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*models.HotTopicsInArea), nil
	}
}

func (h hotTopicsInAreaDo) Take() (*models.HotTopicsInArea, error) {
	if result, err := h.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*models.HotTopicsInArea), nil
	}
}

func (h hotTopicsInAreaDo) Last() (*models.HotTopicsInArea, error) {
	if result, err := h.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*models.HotTopicsInArea), nil
	}
}

func (h hotTopicsInAreaDo) Find() ([]*models.HotTopicsInArea, error) {
	result, err := h.DO.Find()
	return result.([]*models.HotTopicsInArea), err
}

func (h hotTopicsInAreaDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.HotTopicsInArea, err error) {
	buf := make([]*models.HotTopicsInArea, 0, batchSize)
	err = h.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (h hotTopicsInAreaDo) FindInBatches(result *[]*models.HotTopicsInArea, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return h.DO.FindInBatches(result, batchSize, fc)
}

func (h hotTopicsInAreaDo) Attrs(attrs ...field.AssignExpr) IHotTopicsInAreaDo {
	return h.withDO(h.DO.Attrs(attrs...))
}

func (h hotTopicsInAreaDo) Assign(attrs ...field.AssignExpr) IHotTopicsInAreaDo {
	return h.withDO(h.DO.Assign(attrs...))
}

func (h hotTopicsInAreaDo) Joins(fields ...field.RelationField) IHotTopicsInAreaDo {
	for _, _f := range fields {
		h = *h.withDO(h.DO.Joins(_f))
	}
	return &h
}

func (h hotTopicsInAreaDo) Preload(fields ...field.RelationField) IHotTopicsInAreaDo {
	for _, _f := range fields {
		h = *h.withDO(h.DO.Preload(_f))
	}
	return &h
}

func (h hotTopicsInAreaDo) FirstOrInit() (*models.HotTopicsInArea, error) {
	if result, err := h.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*models.HotTopicsInArea), nil
	}
}

func (h hotTopicsInAreaDo) FirstOrCreate() (*models.HotTopicsInArea, error) {
	if result, err := h.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*models.HotTopicsInArea), nil
	}
}

func (h hotTopicsInAreaDo) FindByPage(offset int, limit int) (result []*models.HotTopicsInArea, count int64, err error) {
	result, err = h.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = h.Offset(-1).Limit(-1).Count()
	return
}

func (h hotTopicsInAreaDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = h.Count()
	if err != nil {
		return
	}

	err = h.Offset(offset).Limit(limit).Scan(result)
	return
}

func (h hotTopicsInAreaDo) Scan(result interface{}) (err error) {
	return h.DO.Scan(result)
}

func (h hotTopicsInAreaDo) Delete(models ...*models.HotTopicsInArea) (result gen.ResultInfo, err error) {
	return h.DO.Delete(models)
}

func (h *hotTopicsInAreaDo) withDO(do gen.Dao) *hotTopicsInAreaDo {
	h.DO = *do.(*gen.DO)
	return h
}
