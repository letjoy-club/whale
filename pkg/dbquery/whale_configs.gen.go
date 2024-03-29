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

func newWhaleConfig(db *gorm.DB, opts ...gen.DOOption) whaleConfig {
	_whaleConfig := whaleConfig{}

	_whaleConfig.whaleConfigDo.UseDB(db, opts...)
	_whaleConfig.whaleConfigDo.UseModel(&models.WhaleConfig{})

	tableName := _whaleConfig.whaleConfigDo.TableName()
	_whaleConfig.ALL = field.NewAsterisk(tableName)
	_whaleConfig.ID = field.NewInt(tableName, "id")
	_whaleConfig.Name = field.NewString(tableName, "name")
	_whaleConfig.Desc = field.NewString(tableName, "desc")
	_whaleConfig.Enable = field.NewBool(tableName, "enable")
	_whaleConfig.StartAt = field.NewTime(tableName, "start_at")
	_whaleConfig.EndAt = field.NewTime(tableName, "end_at")
	_whaleConfig.Content = field.NewBytes(tableName, "content")
	_whaleConfig.CreatedAt = field.NewTime(tableName, "created_at")
	_whaleConfig.UpdatedAt = field.NewTime(tableName, "updated_at")
	_whaleConfig.IsDeleted = field.NewUint(tableName, "is_deleted")

	_whaleConfig.fillFieldMap()

	return _whaleConfig
}

type whaleConfig struct {
	whaleConfigDo whaleConfigDo

	ALL       field.Asterisk
	ID        field.Int
	Name      field.String
	Desc      field.String
	Enable    field.Bool
	StartAt   field.Time
	EndAt     field.Time
	Content   field.Bytes
	CreatedAt field.Time
	UpdatedAt field.Time
	IsDeleted field.Uint

	fieldMap map[string]field.Expr
}

func (w whaleConfig) Table(newTableName string) *whaleConfig {
	w.whaleConfigDo.UseTable(newTableName)
	return w.updateTableName(newTableName)
}

func (w whaleConfig) As(alias string) *whaleConfig {
	w.whaleConfigDo.DO = *(w.whaleConfigDo.As(alias).(*gen.DO))
	return w.updateTableName(alias)
}

func (w *whaleConfig) updateTableName(table string) *whaleConfig {
	w.ALL = field.NewAsterisk(table)
	w.ID = field.NewInt(table, "id")
	w.Name = field.NewString(table, "name")
	w.Desc = field.NewString(table, "desc")
	w.Enable = field.NewBool(table, "enable")
	w.StartAt = field.NewTime(table, "start_at")
	w.EndAt = field.NewTime(table, "end_at")
	w.Content = field.NewBytes(table, "content")
	w.CreatedAt = field.NewTime(table, "created_at")
	w.UpdatedAt = field.NewTime(table, "updated_at")
	w.IsDeleted = field.NewUint(table, "is_deleted")

	w.fillFieldMap()

	return w
}

func (w *whaleConfig) WithContext(ctx context.Context) IWhaleConfigDo {
	return w.whaleConfigDo.WithContext(ctx)
}

func (w whaleConfig) TableName() string { return w.whaleConfigDo.TableName() }

func (w whaleConfig) Alias() string { return w.whaleConfigDo.Alias() }

func (w *whaleConfig) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := w.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (w *whaleConfig) fillFieldMap() {
	w.fieldMap = make(map[string]field.Expr, 10)
	w.fieldMap["id"] = w.ID
	w.fieldMap["name"] = w.Name
	w.fieldMap["desc"] = w.Desc
	w.fieldMap["enable"] = w.Enable
	w.fieldMap["start_at"] = w.StartAt
	w.fieldMap["end_at"] = w.EndAt
	w.fieldMap["content"] = w.Content
	w.fieldMap["created_at"] = w.CreatedAt
	w.fieldMap["updated_at"] = w.UpdatedAt
	w.fieldMap["is_deleted"] = w.IsDeleted
}

func (w whaleConfig) clone(db *gorm.DB) whaleConfig {
	w.whaleConfigDo.ReplaceConnPool(db.Statement.ConnPool)
	return w
}

func (w whaleConfig) replaceDB(db *gorm.DB) whaleConfig {
	w.whaleConfigDo.ReplaceDB(db)
	return w
}

type whaleConfigDo struct{ gen.DO }

type IWhaleConfigDo interface {
	gen.SubQuery
	Debug() IWhaleConfigDo
	WithContext(ctx context.Context) IWhaleConfigDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IWhaleConfigDo
	WriteDB() IWhaleConfigDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IWhaleConfigDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IWhaleConfigDo
	Not(conds ...gen.Condition) IWhaleConfigDo
	Or(conds ...gen.Condition) IWhaleConfigDo
	Select(conds ...field.Expr) IWhaleConfigDo
	Where(conds ...gen.Condition) IWhaleConfigDo
	Order(conds ...field.Expr) IWhaleConfigDo
	Distinct(cols ...field.Expr) IWhaleConfigDo
	Omit(cols ...field.Expr) IWhaleConfigDo
	Join(table schema.Tabler, on ...field.Expr) IWhaleConfigDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IWhaleConfigDo
	RightJoin(table schema.Tabler, on ...field.Expr) IWhaleConfigDo
	Group(cols ...field.Expr) IWhaleConfigDo
	Having(conds ...gen.Condition) IWhaleConfigDo
	Limit(limit int) IWhaleConfigDo
	Offset(offset int) IWhaleConfigDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IWhaleConfigDo
	Unscoped() IWhaleConfigDo
	Create(values ...*models.WhaleConfig) error
	CreateInBatches(values []*models.WhaleConfig, batchSize int) error
	Save(values ...*models.WhaleConfig) error
	First() (*models.WhaleConfig, error)
	Take() (*models.WhaleConfig, error)
	Last() (*models.WhaleConfig, error)
	Find() ([]*models.WhaleConfig, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.WhaleConfig, err error)
	FindInBatches(result *[]*models.WhaleConfig, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*models.WhaleConfig) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IWhaleConfigDo
	Assign(attrs ...field.AssignExpr) IWhaleConfigDo
	Joins(fields ...field.RelationField) IWhaleConfigDo
	Preload(fields ...field.RelationField) IWhaleConfigDo
	FirstOrInit() (*models.WhaleConfig, error)
	FirstOrCreate() (*models.WhaleConfig, error)
	FindByPage(offset int, limit int) (result []*models.WhaleConfig, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IWhaleConfigDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (w whaleConfigDo) Debug() IWhaleConfigDo {
	return w.withDO(w.DO.Debug())
}

func (w whaleConfigDo) WithContext(ctx context.Context) IWhaleConfigDo {
	return w.withDO(w.DO.WithContext(ctx))
}

func (w whaleConfigDo) ReadDB() IWhaleConfigDo {
	return w.Clauses(dbresolver.Read)
}

func (w whaleConfigDo) WriteDB() IWhaleConfigDo {
	return w.Clauses(dbresolver.Write)
}

func (w whaleConfigDo) Session(config *gorm.Session) IWhaleConfigDo {
	return w.withDO(w.DO.Session(config))
}

func (w whaleConfigDo) Clauses(conds ...clause.Expression) IWhaleConfigDo {
	return w.withDO(w.DO.Clauses(conds...))
}

func (w whaleConfigDo) Returning(value interface{}, columns ...string) IWhaleConfigDo {
	return w.withDO(w.DO.Returning(value, columns...))
}

func (w whaleConfigDo) Not(conds ...gen.Condition) IWhaleConfigDo {
	return w.withDO(w.DO.Not(conds...))
}

func (w whaleConfigDo) Or(conds ...gen.Condition) IWhaleConfigDo {
	return w.withDO(w.DO.Or(conds...))
}

func (w whaleConfigDo) Select(conds ...field.Expr) IWhaleConfigDo {
	return w.withDO(w.DO.Select(conds...))
}

func (w whaleConfigDo) Where(conds ...gen.Condition) IWhaleConfigDo {
	return w.withDO(w.DO.Where(conds...))
}

func (w whaleConfigDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) IWhaleConfigDo {
	return w.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (w whaleConfigDo) Order(conds ...field.Expr) IWhaleConfigDo {
	return w.withDO(w.DO.Order(conds...))
}

func (w whaleConfigDo) Distinct(cols ...field.Expr) IWhaleConfigDo {
	return w.withDO(w.DO.Distinct(cols...))
}

func (w whaleConfigDo) Omit(cols ...field.Expr) IWhaleConfigDo {
	return w.withDO(w.DO.Omit(cols...))
}

func (w whaleConfigDo) Join(table schema.Tabler, on ...field.Expr) IWhaleConfigDo {
	return w.withDO(w.DO.Join(table, on...))
}

func (w whaleConfigDo) LeftJoin(table schema.Tabler, on ...field.Expr) IWhaleConfigDo {
	return w.withDO(w.DO.LeftJoin(table, on...))
}

func (w whaleConfigDo) RightJoin(table schema.Tabler, on ...field.Expr) IWhaleConfigDo {
	return w.withDO(w.DO.RightJoin(table, on...))
}

func (w whaleConfigDo) Group(cols ...field.Expr) IWhaleConfigDo {
	return w.withDO(w.DO.Group(cols...))
}

func (w whaleConfigDo) Having(conds ...gen.Condition) IWhaleConfigDo {
	return w.withDO(w.DO.Having(conds...))
}

func (w whaleConfigDo) Limit(limit int) IWhaleConfigDo {
	return w.withDO(w.DO.Limit(limit))
}

func (w whaleConfigDo) Offset(offset int) IWhaleConfigDo {
	return w.withDO(w.DO.Offset(offset))
}

func (w whaleConfigDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IWhaleConfigDo {
	return w.withDO(w.DO.Scopes(funcs...))
}

func (w whaleConfigDo) Unscoped() IWhaleConfigDo {
	return w.withDO(w.DO.Unscoped())
}

func (w whaleConfigDo) Create(values ...*models.WhaleConfig) error {
	if len(values) == 0 {
		return nil
	}
	return w.DO.Create(values)
}

func (w whaleConfigDo) CreateInBatches(values []*models.WhaleConfig, batchSize int) error {
	return w.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (w whaleConfigDo) Save(values ...*models.WhaleConfig) error {
	if len(values) == 0 {
		return nil
	}
	return w.DO.Save(values)
}

func (w whaleConfigDo) First() (*models.WhaleConfig, error) {
	if result, err := w.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*models.WhaleConfig), nil
	}
}

func (w whaleConfigDo) Take() (*models.WhaleConfig, error) {
	if result, err := w.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*models.WhaleConfig), nil
	}
}

func (w whaleConfigDo) Last() (*models.WhaleConfig, error) {
	if result, err := w.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*models.WhaleConfig), nil
	}
}

func (w whaleConfigDo) Find() ([]*models.WhaleConfig, error) {
	result, err := w.DO.Find()
	return result.([]*models.WhaleConfig), err
}

func (w whaleConfigDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.WhaleConfig, err error) {
	buf := make([]*models.WhaleConfig, 0, batchSize)
	err = w.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (w whaleConfigDo) FindInBatches(result *[]*models.WhaleConfig, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return w.DO.FindInBatches(result, batchSize, fc)
}

func (w whaleConfigDo) Attrs(attrs ...field.AssignExpr) IWhaleConfigDo {
	return w.withDO(w.DO.Attrs(attrs...))
}

func (w whaleConfigDo) Assign(attrs ...field.AssignExpr) IWhaleConfigDo {
	return w.withDO(w.DO.Assign(attrs...))
}

func (w whaleConfigDo) Joins(fields ...field.RelationField) IWhaleConfigDo {
	for _, _f := range fields {
		w = *w.withDO(w.DO.Joins(_f))
	}
	return &w
}

func (w whaleConfigDo) Preload(fields ...field.RelationField) IWhaleConfigDo {
	for _, _f := range fields {
		w = *w.withDO(w.DO.Preload(_f))
	}
	return &w
}

func (w whaleConfigDo) FirstOrInit() (*models.WhaleConfig, error) {
	if result, err := w.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*models.WhaleConfig), nil
	}
}

func (w whaleConfigDo) FirstOrCreate() (*models.WhaleConfig, error) {
	if result, err := w.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*models.WhaleConfig), nil
	}
}

func (w whaleConfigDo) FindByPage(offset int, limit int) (result []*models.WhaleConfig, count int64, err error) {
	result, err = w.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = w.Offset(-1).Limit(-1).Count()
	return
}

func (w whaleConfigDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = w.Count()
	if err != nil {
		return
	}

	err = w.Offset(offset).Limit(limit).Scan(result)
	return
}

func (w whaleConfigDo) Scan(result interface{}) (err error) {
	return w.DO.Scan(result)
}

func (w whaleConfigDo) Delete(models ...*models.WhaleConfig) (result gen.ResultInfo, err error) {
	return w.DO.Delete(models)
}

func (w *whaleConfigDo) withDO(do gen.Dao) *whaleConfigDo {
	w.DO = *do.(*gen.DO)
	return w
}
