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

func newCityTopics(db *gorm.DB, opts ...gen.DOOption) cityTopics {
	_cityTopics := cityTopics{}

	_cityTopics.cityTopicsDo.UseDB(db, opts...)
	_cityTopics.cityTopicsDo.UseModel(&models.CityTopics{})

	tableName := _cityTopics.cityTopicsDo.TableName()
	_cityTopics.ALL = field.NewAsterisk(tableName)
	_cityTopics.CityID = field.NewString(tableName, "city_id")
	_cityTopics.TopicIDs = field.NewField(tableName, "topic_ids")
	_cityTopics.UpdatedAt = field.NewTime(tableName, "updated_at")

	_cityTopics.fillFieldMap()

	return _cityTopics
}

type cityTopics struct {
	cityTopicsDo cityTopicsDo

	ALL       field.Asterisk
	CityID    field.String
	TopicIDs  field.Field
	UpdatedAt field.Time

	fieldMap map[string]field.Expr
}

func (c cityTopics) Table(newTableName string) *cityTopics {
	c.cityTopicsDo.UseTable(newTableName)
	return c.updateTableName(newTableName)
}

func (c cityTopics) As(alias string) *cityTopics {
	c.cityTopicsDo.DO = *(c.cityTopicsDo.As(alias).(*gen.DO))
	return c.updateTableName(alias)
}

func (c *cityTopics) updateTableName(table string) *cityTopics {
	c.ALL = field.NewAsterisk(table)
	c.CityID = field.NewString(table, "city_id")
	c.TopicIDs = field.NewField(table, "topic_ids")
	c.UpdatedAt = field.NewTime(table, "updated_at")

	c.fillFieldMap()

	return c
}

func (c *cityTopics) WithContext(ctx context.Context) ICityTopicsDo {
	return c.cityTopicsDo.WithContext(ctx)
}

func (c cityTopics) TableName() string { return c.cityTopicsDo.TableName() }

func (c cityTopics) Alias() string { return c.cityTopicsDo.Alias() }

func (c *cityTopics) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := c.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (c *cityTopics) fillFieldMap() {
	c.fieldMap = make(map[string]field.Expr, 3)
	c.fieldMap["city_id"] = c.CityID
	c.fieldMap["topic_ids"] = c.TopicIDs
	c.fieldMap["updated_at"] = c.UpdatedAt
}

func (c cityTopics) clone(db *gorm.DB) cityTopics {
	c.cityTopicsDo.ReplaceConnPool(db.Statement.ConnPool)
	return c
}

func (c cityTopics) replaceDB(db *gorm.DB) cityTopics {
	c.cityTopicsDo.ReplaceDB(db)
	return c
}

type cityTopicsDo struct{ gen.DO }

type ICityTopicsDo interface {
	gen.SubQuery
	Debug() ICityTopicsDo
	WithContext(ctx context.Context) ICityTopicsDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() ICityTopicsDo
	WriteDB() ICityTopicsDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) ICityTopicsDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) ICityTopicsDo
	Not(conds ...gen.Condition) ICityTopicsDo
	Or(conds ...gen.Condition) ICityTopicsDo
	Select(conds ...field.Expr) ICityTopicsDo
	Where(conds ...gen.Condition) ICityTopicsDo
	Order(conds ...field.Expr) ICityTopicsDo
	Distinct(cols ...field.Expr) ICityTopicsDo
	Omit(cols ...field.Expr) ICityTopicsDo
	Join(table schema.Tabler, on ...field.Expr) ICityTopicsDo
	LeftJoin(table schema.Tabler, on ...field.Expr) ICityTopicsDo
	RightJoin(table schema.Tabler, on ...field.Expr) ICityTopicsDo
	Group(cols ...field.Expr) ICityTopicsDo
	Having(conds ...gen.Condition) ICityTopicsDo
	Limit(limit int) ICityTopicsDo
	Offset(offset int) ICityTopicsDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) ICityTopicsDo
	Unscoped() ICityTopicsDo
	Create(values ...*models.CityTopics) error
	CreateInBatches(values []*models.CityTopics, batchSize int) error
	Save(values ...*models.CityTopics) error
	First() (*models.CityTopics, error)
	Take() (*models.CityTopics, error)
	Last() (*models.CityTopics, error)
	Find() ([]*models.CityTopics, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.CityTopics, err error)
	FindInBatches(result *[]*models.CityTopics, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*models.CityTopics) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) ICityTopicsDo
	Assign(attrs ...field.AssignExpr) ICityTopicsDo
	Joins(fields ...field.RelationField) ICityTopicsDo
	Preload(fields ...field.RelationField) ICityTopicsDo
	FirstOrInit() (*models.CityTopics, error)
	FirstOrCreate() (*models.CityTopics, error)
	FindByPage(offset int, limit int) (result []*models.CityTopics, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) ICityTopicsDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (c cityTopicsDo) Debug() ICityTopicsDo {
	return c.withDO(c.DO.Debug())
}

func (c cityTopicsDo) WithContext(ctx context.Context) ICityTopicsDo {
	return c.withDO(c.DO.WithContext(ctx))
}

func (c cityTopicsDo) ReadDB() ICityTopicsDo {
	return c.Clauses(dbresolver.Read)
}

func (c cityTopicsDo) WriteDB() ICityTopicsDo {
	return c.Clauses(dbresolver.Write)
}

func (c cityTopicsDo) Session(config *gorm.Session) ICityTopicsDo {
	return c.withDO(c.DO.Session(config))
}

func (c cityTopicsDo) Clauses(conds ...clause.Expression) ICityTopicsDo {
	return c.withDO(c.DO.Clauses(conds...))
}

func (c cityTopicsDo) Returning(value interface{}, columns ...string) ICityTopicsDo {
	return c.withDO(c.DO.Returning(value, columns...))
}

func (c cityTopicsDo) Not(conds ...gen.Condition) ICityTopicsDo {
	return c.withDO(c.DO.Not(conds...))
}

func (c cityTopicsDo) Or(conds ...gen.Condition) ICityTopicsDo {
	return c.withDO(c.DO.Or(conds...))
}

func (c cityTopicsDo) Select(conds ...field.Expr) ICityTopicsDo {
	return c.withDO(c.DO.Select(conds...))
}

func (c cityTopicsDo) Where(conds ...gen.Condition) ICityTopicsDo {
	return c.withDO(c.DO.Where(conds...))
}

func (c cityTopicsDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) ICityTopicsDo {
	return c.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (c cityTopicsDo) Order(conds ...field.Expr) ICityTopicsDo {
	return c.withDO(c.DO.Order(conds...))
}

func (c cityTopicsDo) Distinct(cols ...field.Expr) ICityTopicsDo {
	return c.withDO(c.DO.Distinct(cols...))
}

func (c cityTopicsDo) Omit(cols ...field.Expr) ICityTopicsDo {
	return c.withDO(c.DO.Omit(cols...))
}

func (c cityTopicsDo) Join(table schema.Tabler, on ...field.Expr) ICityTopicsDo {
	return c.withDO(c.DO.Join(table, on...))
}

func (c cityTopicsDo) LeftJoin(table schema.Tabler, on ...field.Expr) ICityTopicsDo {
	return c.withDO(c.DO.LeftJoin(table, on...))
}

func (c cityTopicsDo) RightJoin(table schema.Tabler, on ...field.Expr) ICityTopicsDo {
	return c.withDO(c.DO.RightJoin(table, on...))
}

func (c cityTopicsDo) Group(cols ...field.Expr) ICityTopicsDo {
	return c.withDO(c.DO.Group(cols...))
}

func (c cityTopicsDo) Having(conds ...gen.Condition) ICityTopicsDo {
	return c.withDO(c.DO.Having(conds...))
}

func (c cityTopicsDo) Limit(limit int) ICityTopicsDo {
	return c.withDO(c.DO.Limit(limit))
}

func (c cityTopicsDo) Offset(offset int) ICityTopicsDo {
	return c.withDO(c.DO.Offset(offset))
}

func (c cityTopicsDo) Scopes(funcs ...func(gen.Dao) gen.Dao) ICityTopicsDo {
	return c.withDO(c.DO.Scopes(funcs...))
}

func (c cityTopicsDo) Unscoped() ICityTopicsDo {
	return c.withDO(c.DO.Unscoped())
}

func (c cityTopicsDo) Create(values ...*models.CityTopics) error {
	if len(values) == 0 {
		return nil
	}
	return c.DO.Create(values)
}

func (c cityTopicsDo) CreateInBatches(values []*models.CityTopics, batchSize int) error {
	return c.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (c cityTopicsDo) Save(values ...*models.CityTopics) error {
	if len(values) == 0 {
		return nil
	}
	return c.DO.Save(values)
}

func (c cityTopicsDo) First() (*models.CityTopics, error) {
	if result, err := c.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*models.CityTopics), nil
	}
}

func (c cityTopicsDo) Take() (*models.CityTopics, error) {
	if result, err := c.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*models.CityTopics), nil
	}
}

func (c cityTopicsDo) Last() (*models.CityTopics, error) {
	if result, err := c.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*models.CityTopics), nil
	}
}

func (c cityTopicsDo) Find() ([]*models.CityTopics, error) {
	result, err := c.DO.Find()
	return result.([]*models.CityTopics), err
}

func (c cityTopicsDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.CityTopics, err error) {
	buf := make([]*models.CityTopics, 0, batchSize)
	err = c.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (c cityTopicsDo) FindInBatches(result *[]*models.CityTopics, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return c.DO.FindInBatches(result, batchSize, fc)
}

func (c cityTopicsDo) Attrs(attrs ...field.AssignExpr) ICityTopicsDo {
	return c.withDO(c.DO.Attrs(attrs...))
}

func (c cityTopicsDo) Assign(attrs ...field.AssignExpr) ICityTopicsDo {
	return c.withDO(c.DO.Assign(attrs...))
}

func (c cityTopicsDo) Joins(fields ...field.RelationField) ICityTopicsDo {
	for _, _f := range fields {
		c = *c.withDO(c.DO.Joins(_f))
	}
	return &c
}

func (c cityTopicsDo) Preload(fields ...field.RelationField) ICityTopicsDo {
	for _, _f := range fields {
		c = *c.withDO(c.DO.Preload(_f))
	}
	return &c
}

func (c cityTopicsDo) FirstOrInit() (*models.CityTopics, error) {
	if result, err := c.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*models.CityTopics), nil
	}
}

func (c cityTopicsDo) FirstOrCreate() (*models.CityTopics, error) {
	if result, err := c.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*models.CityTopics), nil
	}
}

func (c cityTopicsDo) FindByPage(offset int, limit int) (result []*models.CityTopics, count int64, err error) {
	result, err = c.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = c.Offset(-1).Limit(-1).Count()
	return
}

func (c cityTopicsDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = c.Count()
	if err != nil {
		return
	}

	err = c.Offset(offset).Limit(limit).Scan(result)
	return
}

func (c cityTopicsDo) Scan(result interface{}) (err error) {
	return c.DO.Scan(result)
}

func (c cityTopicsDo) Delete(models ...*models.CityTopics) (result gen.ResultInfo, err error) {
	return c.DO.Delete(models)
}

func (c *cityTopicsDo) withDO(do gen.Dao) *cityTopicsDo {
	c.DO = *do.(*gen.DO)
	return c
}
