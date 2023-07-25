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

func newMatchingView(db *gorm.DB, opts ...gen.DOOption) matchingView {
	_matchingView := matchingView{}

	_matchingView.matchingViewDo.UseDB(db, opts...)
	_matchingView.matchingViewDo.UseModel(&models.MatchingView{})

	tableName := _matchingView.matchingViewDo.TableName()
	_matchingView.ALL = field.NewAsterisk(tableName)
	_matchingView.MatchingID = field.NewString(tableName, "matching_id")
	_matchingView.ViewCount = field.NewInt(tableName, "view_count")
	_matchingView.UpdatedAt = field.NewTime(tableName, "updated_at")

	_matchingView.fillFieldMap()

	return _matchingView
}

type matchingView struct {
	matchingViewDo matchingViewDo

	ALL        field.Asterisk
	MatchingID field.String
	ViewCount  field.Int
	UpdatedAt  field.Time

	fieldMap map[string]field.Expr
}

func (m matchingView) Table(newTableName string) *matchingView {
	m.matchingViewDo.UseTable(newTableName)
	return m.updateTableName(newTableName)
}

func (m matchingView) As(alias string) *matchingView {
	m.matchingViewDo.DO = *(m.matchingViewDo.As(alias).(*gen.DO))
	return m.updateTableName(alias)
}

func (m *matchingView) updateTableName(table string) *matchingView {
	m.ALL = field.NewAsterisk(table)
	m.MatchingID = field.NewString(table, "matching_id")
	m.ViewCount = field.NewInt(table, "view_count")
	m.UpdatedAt = field.NewTime(table, "updated_at")

	m.fillFieldMap()

	return m
}

func (m *matchingView) WithContext(ctx context.Context) IMatchingViewDo {
	return m.matchingViewDo.WithContext(ctx)
}

func (m matchingView) TableName() string { return m.matchingViewDo.TableName() }

func (m matchingView) Alias() string { return m.matchingViewDo.Alias() }

func (m *matchingView) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := m.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (m *matchingView) fillFieldMap() {
	m.fieldMap = make(map[string]field.Expr, 3)
	m.fieldMap["matching_id"] = m.MatchingID
	m.fieldMap["view_count"] = m.ViewCount
	m.fieldMap["updated_at"] = m.UpdatedAt
}

func (m matchingView) clone(db *gorm.DB) matchingView {
	m.matchingViewDo.ReplaceConnPool(db.Statement.ConnPool)
	return m
}

func (m matchingView) replaceDB(db *gorm.DB) matchingView {
	m.matchingViewDo.ReplaceDB(db)
	return m
}

type matchingViewDo struct{ gen.DO }

type IMatchingViewDo interface {
	gen.SubQuery
	Debug() IMatchingViewDo
	WithContext(ctx context.Context) IMatchingViewDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IMatchingViewDo
	WriteDB() IMatchingViewDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IMatchingViewDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IMatchingViewDo
	Not(conds ...gen.Condition) IMatchingViewDo
	Or(conds ...gen.Condition) IMatchingViewDo
	Select(conds ...field.Expr) IMatchingViewDo
	Where(conds ...gen.Condition) IMatchingViewDo
	Order(conds ...field.Expr) IMatchingViewDo
	Distinct(cols ...field.Expr) IMatchingViewDo
	Omit(cols ...field.Expr) IMatchingViewDo
	Join(table schema.Tabler, on ...field.Expr) IMatchingViewDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IMatchingViewDo
	RightJoin(table schema.Tabler, on ...field.Expr) IMatchingViewDo
	Group(cols ...field.Expr) IMatchingViewDo
	Having(conds ...gen.Condition) IMatchingViewDo
	Limit(limit int) IMatchingViewDo
	Offset(offset int) IMatchingViewDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IMatchingViewDo
	Unscoped() IMatchingViewDo
	Create(values ...*models.MatchingView) error
	CreateInBatches(values []*models.MatchingView, batchSize int) error
	Save(values ...*models.MatchingView) error
	First() (*models.MatchingView, error)
	Take() (*models.MatchingView, error)
	Last() (*models.MatchingView, error)
	Find() ([]*models.MatchingView, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.MatchingView, err error)
	FindInBatches(result *[]*models.MatchingView, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*models.MatchingView) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IMatchingViewDo
	Assign(attrs ...field.AssignExpr) IMatchingViewDo
	Joins(fields ...field.RelationField) IMatchingViewDo
	Preload(fields ...field.RelationField) IMatchingViewDo
	FirstOrInit() (*models.MatchingView, error)
	FirstOrCreate() (*models.MatchingView, error)
	FindByPage(offset int, limit int) (result []*models.MatchingView, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IMatchingViewDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (m matchingViewDo) Debug() IMatchingViewDo {
	return m.withDO(m.DO.Debug())
}

func (m matchingViewDo) WithContext(ctx context.Context) IMatchingViewDo {
	return m.withDO(m.DO.WithContext(ctx))
}

func (m matchingViewDo) ReadDB() IMatchingViewDo {
	return m.Clauses(dbresolver.Read)
}

func (m matchingViewDo) WriteDB() IMatchingViewDo {
	return m.Clauses(dbresolver.Write)
}

func (m matchingViewDo) Session(config *gorm.Session) IMatchingViewDo {
	return m.withDO(m.DO.Session(config))
}

func (m matchingViewDo) Clauses(conds ...clause.Expression) IMatchingViewDo {
	return m.withDO(m.DO.Clauses(conds...))
}

func (m matchingViewDo) Returning(value interface{}, columns ...string) IMatchingViewDo {
	return m.withDO(m.DO.Returning(value, columns...))
}

func (m matchingViewDo) Not(conds ...gen.Condition) IMatchingViewDo {
	return m.withDO(m.DO.Not(conds...))
}

func (m matchingViewDo) Or(conds ...gen.Condition) IMatchingViewDo {
	return m.withDO(m.DO.Or(conds...))
}

func (m matchingViewDo) Select(conds ...field.Expr) IMatchingViewDo {
	return m.withDO(m.DO.Select(conds...))
}

func (m matchingViewDo) Where(conds ...gen.Condition) IMatchingViewDo {
	return m.withDO(m.DO.Where(conds...))
}

func (m matchingViewDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) IMatchingViewDo {
	return m.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (m matchingViewDo) Order(conds ...field.Expr) IMatchingViewDo {
	return m.withDO(m.DO.Order(conds...))
}

func (m matchingViewDo) Distinct(cols ...field.Expr) IMatchingViewDo {
	return m.withDO(m.DO.Distinct(cols...))
}

func (m matchingViewDo) Omit(cols ...field.Expr) IMatchingViewDo {
	return m.withDO(m.DO.Omit(cols...))
}

func (m matchingViewDo) Join(table schema.Tabler, on ...field.Expr) IMatchingViewDo {
	return m.withDO(m.DO.Join(table, on...))
}

func (m matchingViewDo) LeftJoin(table schema.Tabler, on ...field.Expr) IMatchingViewDo {
	return m.withDO(m.DO.LeftJoin(table, on...))
}

func (m matchingViewDo) RightJoin(table schema.Tabler, on ...field.Expr) IMatchingViewDo {
	return m.withDO(m.DO.RightJoin(table, on...))
}

func (m matchingViewDo) Group(cols ...field.Expr) IMatchingViewDo {
	return m.withDO(m.DO.Group(cols...))
}

func (m matchingViewDo) Having(conds ...gen.Condition) IMatchingViewDo {
	return m.withDO(m.DO.Having(conds...))
}

func (m matchingViewDo) Limit(limit int) IMatchingViewDo {
	return m.withDO(m.DO.Limit(limit))
}

func (m matchingViewDo) Offset(offset int) IMatchingViewDo {
	return m.withDO(m.DO.Offset(offset))
}

func (m matchingViewDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IMatchingViewDo {
	return m.withDO(m.DO.Scopes(funcs...))
}

func (m matchingViewDo) Unscoped() IMatchingViewDo {
	return m.withDO(m.DO.Unscoped())
}

func (m matchingViewDo) Create(values ...*models.MatchingView) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Create(values)
}

func (m matchingViewDo) CreateInBatches(values []*models.MatchingView, batchSize int) error {
	return m.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (m matchingViewDo) Save(values ...*models.MatchingView) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Save(values)
}

func (m matchingViewDo) First() (*models.MatchingView, error) {
	if result, err := m.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*models.MatchingView), nil
	}
}

func (m matchingViewDo) Take() (*models.MatchingView, error) {
	if result, err := m.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*models.MatchingView), nil
	}
}

func (m matchingViewDo) Last() (*models.MatchingView, error) {
	if result, err := m.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*models.MatchingView), nil
	}
}

func (m matchingViewDo) Find() ([]*models.MatchingView, error) {
	result, err := m.DO.Find()
	return result.([]*models.MatchingView), err
}

func (m matchingViewDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.MatchingView, err error) {
	buf := make([]*models.MatchingView, 0, batchSize)
	err = m.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (m matchingViewDo) FindInBatches(result *[]*models.MatchingView, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return m.DO.FindInBatches(result, batchSize, fc)
}

func (m matchingViewDo) Attrs(attrs ...field.AssignExpr) IMatchingViewDo {
	return m.withDO(m.DO.Attrs(attrs...))
}

func (m matchingViewDo) Assign(attrs ...field.AssignExpr) IMatchingViewDo {
	return m.withDO(m.DO.Assign(attrs...))
}

func (m matchingViewDo) Joins(fields ...field.RelationField) IMatchingViewDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Joins(_f))
	}
	return &m
}

func (m matchingViewDo) Preload(fields ...field.RelationField) IMatchingViewDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Preload(_f))
	}
	return &m
}

func (m matchingViewDo) FirstOrInit() (*models.MatchingView, error) {
	if result, err := m.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*models.MatchingView), nil
	}
}

func (m matchingViewDo) FirstOrCreate() (*models.MatchingView, error) {
	if result, err := m.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*models.MatchingView), nil
	}
}

func (m matchingViewDo) FindByPage(offset int, limit int) (result []*models.MatchingView, count int64, err error) {
	result, err = m.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = m.Offset(-1).Limit(-1).Count()
	return
}

func (m matchingViewDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = m.Count()
	if err != nil {
		return
	}

	err = m.Offset(offset).Limit(limit).Scan(result)
	return
}

func (m matchingViewDo) Scan(result interface{}) (err error) {
	return m.DO.Scan(result)
}

func (m matchingViewDo) Delete(models ...*models.MatchingView) (result gen.ResultInfo, err error) {
	return m.DO.Delete(models)
}

func (m *matchingViewDo) withDO(do gen.Dao) *matchingViewDo {
	m.DO = *do.(*gen.DO)
	return m
}
