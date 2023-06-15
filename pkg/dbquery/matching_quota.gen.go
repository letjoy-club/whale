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

func newMatchingQuota(db *gorm.DB, opts ...gen.DOOption) matchingQuota {
	_matchingQuota := matchingQuota{}

	_matchingQuota.matchingQuotaDo.UseDB(db, opts...)
	_matchingQuota.matchingQuotaDo.UseModel(&models.MatchingQuota{})

	tableName := _matchingQuota.matchingQuotaDo.TableName()
	_matchingQuota.ALL = field.NewAsterisk(tableName)
	_matchingQuota.UserID = field.NewString(tableName, "user_id")
	_matchingQuota.Remain = field.NewInt(tableName, "remain")
	_matchingQuota.Total = field.NewInt(tableName, "total")
	_matchingQuota.MatchingNum = field.NewInt(tableName, "matching_num")
	_matchingQuota.InvitationNum = field.NewInt(tableName, "invitation_num")
	_matchingQuota.CreatedAt = field.NewTime(tableName, "created_at")
	_matchingQuota.UpdatedAt = field.NewTime(tableName, "updated_at")

	_matchingQuota.fillFieldMap()

	return _matchingQuota
}

type matchingQuota struct {
	matchingQuotaDo matchingQuotaDo

	ALL           field.Asterisk
	UserID        field.String
	Remain        field.Int
	Total         field.Int
	MatchingNum   field.Int
	InvitationNum field.Int
	CreatedAt     field.Time
	UpdatedAt     field.Time

	fieldMap map[string]field.Expr
}

func (m matchingQuota) Table(newTableName string) *matchingQuota {
	m.matchingQuotaDo.UseTable(newTableName)
	return m.updateTableName(newTableName)
}

func (m matchingQuota) As(alias string) *matchingQuota {
	m.matchingQuotaDo.DO = *(m.matchingQuotaDo.As(alias).(*gen.DO))
	return m.updateTableName(alias)
}

func (m *matchingQuota) updateTableName(table string) *matchingQuota {
	m.ALL = field.NewAsterisk(table)
	m.UserID = field.NewString(table, "user_id")
	m.Remain = field.NewInt(table, "remain")
	m.Total = field.NewInt(table, "total")
	m.MatchingNum = field.NewInt(table, "matching_num")
	m.InvitationNum = field.NewInt(table, "invitation_num")
	m.CreatedAt = field.NewTime(table, "created_at")
	m.UpdatedAt = field.NewTime(table, "updated_at")

	m.fillFieldMap()

	return m
}

func (m *matchingQuota) WithContext(ctx context.Context) IMatchingQuotaDo {
	return m.matchingQuotaDo.WithContext(ctx)
}

func (m matchingQuota) TableName() string { return m.matchingQuotaDo.TableName() }

func (m matchingQuota) Alias() string { return m.matchingQuotaDo.Alias() }

func (m *matchingQuota) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := m.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (m *matchingQuota) fillFieldMap() {
	m.fieldMap = make(map[string]field.Expr, 7)
	m.fieldMap["user_id"] = m.UserID
	m.fieldMap["remain"] = m.Remain
	m.fieldMap["total"] = m.Total
	m.fieldMap["matching_num"] = m.MatchingNum
	m.fieldMap["invitation_num"] = m.InvitationNum
	m.fieldMap["created_at"] = m.CreatedAt
	m.fieldMap["updated_at"] = m.UpdatedAt
}

func (m matchingQuota) clone(db *gorm.DB) matchingQuota {
	m.matchingQuotaDo.ReplaceConnPool(db.Statement.ConnPool)
	return m
}

func (m matchingQuota) replaceDB(db *gorm.DB) matchingQuota {
	m.matchingQuotaDo.ReplaceDB(db)
	return m
}

type matchingQuotaDo struct{ gen.DO }

type IMatchingQuotaDo interface {
	gen.SubQuery
	Debug() IMatchingQuotaDo
	WithContext(ctx context.Context) IMatchingQuotaDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IMatchingQuotaDo
	WriteDB() IMatchingQuotaDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IMatchingQuotaDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IMatchingQuotaDo
	Not(conds ...gen.Condition) IMatchingQuotaDo
	Or(conds ...gen.Condition) IMatchingQuotaDo
	Select(conds ...field.Expr) IMatchingQuotaDo
	Where(conds ...gen.Condition) IMatchingQuotaDo
	Order(conds ...field.Expr) IMatchingQuotaDo
	Distinct(cols ...field.Expr) IMatchingQuotaDo
	Omit(cols ...field.Expr) IMatchingQuotaDo
	Join(table schema.Tabler, on ...field.Expr) IMatchingQuotaDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IMatchingQuotaDo
	RightJoin(table schema.Tabler, on ...field.Expr) IMatchingQuotaDo
	Group(cols ...field.Expr) IMatchingQuotaDo
	Having(conds ...gen.Condition) IMatchingQuotaDo
	Limit(limit int) IMatchingQuotaDo
	Offset(offset int) IMatchingQuotaDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IMatchingQuotaDo
	Unscoped() IMatchingQuotaDo
	Create(values ...*models.MatchingQuota) error
	CreateInBatches(values []*models.MatchingQuota, batchSize int) error
	Save(values ...*models.MatchingQuota) error
	First() (*models.MatchingQuota, error)
	Take() (*models.MatchingQuota, error)
	Last() (*models.MatchingQuota, error)
	Find() ([]*models.MatchingQuota, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.MatchingQuota, err error)
	FindInBatches(result *[]*models.MatchingQuota, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*models.MatchingQuota) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IMatchingQuotaDo
	Assign(attrs ...field.AssignExpr) IMatchingQuotaDo
	Joins(fields ...field.RelationField) IMatchingQuotaDo
	Preload(fields ...field.RelationField) IMatchingQuotaDo
	FirstOrInit() (*models.MatchingQuota, error)
	FirstOrCreate() (*models.MatchingQuota, error)
	FindByPage(offset int, limit int) (result []*models.MatchingQuota, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IMatchingQuotaDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (m matchingQuotaDo) Debug() IMatchingQuotaDo {
	return m.withDO(m.DO.Debug())
}

func (m matchingQuotaDo) WithContext(ctx context.Context) IMatchingQuotaDo {
	return m.withDO(m.DO.WithContext(ctx))
}

func (m matchingQuotaDo) ReadDB() IMatchingQuotaDo {
	return m.Clauses(dbresolver.Read)
}

func (m matchingQuotaDo) WriteDB() IMatchingQuotaDo {
	return m.Clauses(dbresolver.Write)
}

func (m matchingQuotaDo) Session(config *gorm.Session) IMatchingQuotaDo {
	return m.withDO(m.DO.Session(config))
}

func (m matchingQuotaDo) Clauses(conds ...clause.Expression) IMatchingQuotaDo {
	return m.withDO(m.DO.Clauses(conds...))
}

func (m matchingQuotaDo) Returning(value interface{}, columns ...string) IMatchingQuotaDo {
	return m.withDO(m.DO.Returning(value, columns...))
}

func (m matchingQuotaDo) Not(conds ...gen.Condition) IMatchingQuotaDo {
	return m.withDO(m.DO.Not(conds...))
}

func (m matchingQuotaDo) Or(conds ...gen.Condition) IMatchingQuotaDo {
	return m.withDO(m.DO.Or(conds...))
}

func (m matchingQuotaDo) Select(conds ...field.Expr) IMatchingQuotaDo {
	return m.withDO(m.DO.Select(conds...))
}

func (m matchingQuotaDo) Where(conds ...gen.Condition) IMatchingQuotaDo {
	return m.withDO(m.DO.Where(conds...))
}

func (m matchingQuotaDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) IMatchingQuotaDo {
	return m.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (m matchingQuotaDo) Order(conds ...field.Expr) IMatchingQuotaDo {
	return m.withDO(m.DO.Order(conds...))
}

func (m matchingQuotaDo) Distinct(cols ...field.Expr) IMatchingQuotaDo {
	return m.withDO(m.DO.Distinct(cols...))
}

func (m matchingQuotaDo) Omit(cols ...field.Expr) IMatchingQuotaDo {
	return m.withDO(m.DO.Omit(cols...))
}

func (m matchingQuotaDo) Join(table schema.Tabler, on ...field.Expr) IMatchingQuotaDo {
	return m.withDO(m.DO.Join(table, on...))
}

func (m matchingQuotaDo) LeftJoin(table schema.Tabler, on ...field.Expr) IMatchingQuotaDo {
	return m.withDO(m.DO.LeftJoin(table, on...))
}

func (m matchingQuotaDo) RightJoin(table schema.Tabler, on ...field.Expr) IMatchingQuotaDo {
	return m.withDO(m.DO.RightJoin(table, on...))
}

func (m matchingQuotaDo) Group(cols ...field.Expr) IMatchingQuotaDo {
	return m.withDO(m.DO.Group(cols...))
}

func (m matchingQuotaDo) Having(conds ...gen.Condition) IMatchingQuotaDo {
	return m.withDO(m.DO.Having(conds...))
}

func (m matchingQuotaDo) Limit(limit int) IMatchingQuotaDo {
	return m.withDO(m.DO.Limit(limit))
}

func (m matchingQuotaDo) Offset(offset int) IMatchingQuotaDo {
	return m.withDO(m.DO.Offset(offset))
}

func (m matchingQuotaDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IMatchingQuotaDo {
	return m.withDO(m.DO.Scopes(funcs...))
}

func (m matchingQuotaDo) Unscoped() IMatchingQuotaDo {
	return m.withDO(m.DO.Unscoped())
}

func (m matchingQuotaDo) Create(values ...*models.MatchingQuota) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Create(values)
}

func (m matchingQuotaDo) CreateInBatches(values []*models.MatchingQuota, batchSize int) error {
	return m.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (m matchingQuotaDo) Save(values ...*models.MatchingQuota) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Save(values)
}

func (m matchingQuotaDo) First() (*models.MatchingQuota, error) {
	if result, err := m.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*models.MatchingQuota), nil
	}
}

func (m matchingQuotaDo) Take() (*models.MatchingQuota, error) {
	if result, err := m.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*models.MatchingQuota), nil
	}
}

func (m matchingQuotaDo) Last() (*models.MatchingQuota, error) {
	if result, err := m.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*models.MatchingQuota), nil
	}
}

func (m matchingQuotaDo) Find() ([]*models.MatchingQuota, error) {
	result, err := m.DO.Find()
	return result.([]*models.MatchingQuota), err
}

func (m matchingQuotaDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.MatchingQuota, err error) {
	buf := make([]*models.MatchingQuota, 0, batchSize)
	err = m.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (m matchingQuotaDo) FindInBatches(result *[]*models.MatchingQuota, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return m.DO.FindInBatches(result, batchSize, fc)
}

func (m matchingQuotaDo) Attrs(attrs ...field.AssignExpr) IMatchingQuotaDo {
	return m.withDO(m.DO.Attrs(attrs...))
}

func (m matchingQuotaDo) Assign(attrs ...field.AssignExpr) IMatchingQuotaDo {
	return m.withDO(m.DO.Assign(attrs...))
}

func (m matchingQuotaDo) Joins(fields ...field.RelationField) IMatchingQuotaDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Joins(_f))
	}
	return &m
}

func (m matchingQuotaDo) Preload(fields ...field.RelationField) IMatchingQuotaDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Preload(_f))
	}
	return &m
}

func (m matchingQuotaDo) FirstOrInit() (*models.MatchingQuota, error) {
	if result, err := m.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*models.MatchingQuota), nil
	}
}

func (m matchingQuotaDo) FirstOrCreate() (*models.MatchingQuota, error) {
	if result, err := m.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*models.MatchingQuota), nil
	}
}

func (m matchingQuotaDo) FindByPage(offset int, limit int) (result []*models.MatchingQuota, count int64, err error) {
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

func (m matchingQuotaDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = m.Count()
	if err != nil {
		return
	}

	err = m.Offset(offset).Limit(limit).Scan(result)
	return
}

func (m matchingQuotaDo) Scan(result interface{}) (err error) {
	return m.DO.Scan(result)
}

func (m matchingQuotaDo) Delete(models ...*models.MatchingQuota) (result gen.ResultInfo, err error) {
	return m.DO.Delete(models)
}

func (m *matchingQuotaDo) withDO(do gen.Dao) *matchingQuotaDo {
	m.DO = *do.(*gen.DO)
	return m
}
