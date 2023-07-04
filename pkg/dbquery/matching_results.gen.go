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

func newMatchingResult(db *gorm.DB, opts ...gen.DOOption) matchingResult {
	_matchingResult := matchingResult{}

	_matchingResult.matchingResultDo.UseDB(db, opts...)
	_matchingResult.matchingResultDo.UseModel(&models.MatchingResult{})

	tableName := _matchingResult.matchingResultDo.TableName()
	_matchingResult.ALL = field.NewAsterisk(tableName)
	_matchingResult.ID = field.NewInt(tableName, "id")
	_matchingResult.MatchingIDs = field.NewField(tableName, "matching_ids")
	_matchingResult.TopicID = field.NewString(tableName, "topic_id")
	_matchingResult.UserIDs = field.NewField(tableName, "user_ids")
	_matchingResult.ConfirmStates = field.NewField(tableName, "confirm_states")
	_matchingResult.ChatGroupState = field.NewString(tableName, "chat_group_state")
	_matchingResult.ChatGroupID = field.NewString(tableName, "chat_group_id")
	_matchingResult.Closed = field.NewBool(tableName, "closed")
	_matchingResult.MatchingScore = field.NewInt(tableName, "matching_score")
	_matchingResult.CreatedBy = field.NewString(tableName, "created_by")
	_matchingResult.FinishedAt = field.NewTime(tableName, "finished_at")
	_matchingResult.ChatGroupCreatedAt = field.NewTime(tableName, "chat_group_created_at")
	_matchingResult.CreatedAt = field.NewTime(tableName, "created_at")
	_matchingResult.UpdatedAt = field.NewTime(tableName, "updated_at")

	_matchingResult.fillFieldMap()

	return _matchingResult
}

type matchingResult struct {
	matchingResultDo matchingResultDo

	ALL                field.Asterisk
	ID                 field.Int
	MatchingIDs        field.Field
	TopicID            field.String
	UserIDs            field.Field
	ConfirmStates      field.Field
	ChatGroupState     field.String
	ChatGroupID        field.String
	Closed             field.Bool
	MatchingScore      field.Int
	CreatedBy          field.String
	FinishedAt         field.Time
	ChatGroupCreatedAt field.Time
	CreatedAt          field.Time
	UpdatedAt          field.Time

	fieldMap map[string]field.Expr
}

func (m matchingResult) Table(newTableName string) *matchingResult {
	m.matchingResultDo.UseTable(newTableName)
	return m.updateTableName(newTableName)
}

func (m matchingResult) As(alias string) *matchingResult {
	m.matchingResultDo.DO = *(m.matchingResultDo.As(alias).(*gen.DO))
	return m.updateTableName(alias)
}

func (m *matchingResult) updateTableName(table string) *matchingResult {
	m.ALL = field.NewAsterisk(table)
	m.ID = field.NewInt(table, "id")
	m.MatchingIDs = field.NewField(table, "matching_ids")
	m.TopicID = field.NewString(table, "topic_id")
	m.UserIDs = field.NewField(table, "user_ids")
	m.ConfirmStates = field.NewField(table, "confirm_states")
	m.ChatGroupState = field.NewString(table, "chat_group_state")
	m.ChatGroupID = field.NewString(table, "chat_group_id")
	m.Closed = field.NewBool(table, "closed")
	m.MatchingScore = field.NewInt(table, "matching_score")
	m.CreatedBy = field.NewString(table, "created_by")
	m.FinishedAt = field.NewTime(table, "finished_at")
	m.ChatGroupCreatedAt = field.NewTime(table, "chat_group_created_at")
	m.CreatedAt = field.NewTime(table, "created_at")
	m.UpdatedAt = field.NewTime(table, "updated_at")

	m.fillFieldMap()

	return m
}

func (m *matchingResult) WithContext(ctx context.Context) IMatchingResultDo {
	return m.matchingResultDo.WithContext(ctx)
}

func (m matchingResult) TableName() string { return m.matchingResultDo.TableName() }

func (m matchingResult) Alias() string { return m.matchingResultDo.Alias() }

func (m *matchingResult) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := m.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (m *matchingResult) fillFieldMap() {
	m.fieldMap = make(map[string]field.Expr, 14)
	m.fieldMap["id"] = m.ID
	m.fieldMap["matching_ids"] = m.MatchingIDs
	m.fieldMap["topic_id"] = m.TopicID
	m.fieldMap["user_ids"] = m.UserIDs
	m.fieldMap["confirm_states"] = m.ConfirmStates
	m.fieldMap["chat_group_state"] = m.ChatGroupState
	m.fieldMap["chat_group_id"] = m.ChatGroupID
	m.fieldMap["closed"] = m.Closed
	m.fieldMap["matching_score"] = m.MatchingScore
	m.fieldMap["created_by"] = m.CreatedBy
	m.fieldMap["finished_at"] = m.FinishedAt
	m.fieldMap["chat_group_created_at"] = m.ChatGroupCreatedAt
	m.fieldMap["created_at"] = m.CreatedAt
	m.fieldMap["updated_at"] = m.UpdatedAt
}

func (m matchingResult) clone(db *gorm.DB) matchingResult {
	m.matchingResultDo.ReplaceConnPool(db.Statement.ConnPool)
	return m
}

func (m matchingResult) replaceDB(db *gorm.DB) matchingResult {
	m.matchingResultDo.ReplaceDB(db)
	return m
}

type matchingResultDo struct{ gen.DO }

type IMatchingResultDo interface {
	gen.SubQuery
	Debug() IMatchingResultDo
	WithContext(ctx context.Context) IMatchingResultDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IMatchingResultDo
	WriteDB() IMatchingResultDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IMatchingResultDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IMatchingResultDo
	Not(conds ...gen.Condition) IMatchingResultDo
	Or(conds ...gen.Condition) IMatchingResultDo
	Select(conds ...field.Expr) IMatchingResultDo
	Where(conds ...gen.Condition) IMatchingResultDo
	Order(conds ...field.Expr) IMatchingResultDo
	Distinct(cols ...field.Expr) IMatchingResultDo
	Omit(cols ...field.Expr) IMatchingResultDo
	Join(table schema.Tabler, on ...field.Expr) IMatchingResultDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IMatchingResultDo
	RightJoin(table schema.Tabler, on ...field.Expr) IMatchingResultDo
	Group(cols ...field.Expr) IMatchingResultDo
	Having(conds ...gen.Condition) IMatchingResultDo
	Limit(limit int) IMatchingResultDo
	Offset(offset int) IMatchingResultDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IMatchingResultDo
	Unscoped() IMatchingResultDo
	Create(values ...*models.MatchingResult) error
	CreateInBatches(values []*models.MatchingResult, batchSize int) error
	Save(values ...*models.MatchingResult) error
	First() (*models.MatchingResult, error)
	Take() (*models.MatchingResult, error)
	Last() (*models.MatchingResult, error)
	Find() ([]*models.MatchingResult, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.MatchingResult, err error)
	FindInBatches(result *[]*models.MatchingResult, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*models.MatchingResult) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IMatchingResultDo
	Assign(attrs ...field.AssignExpr) IMatchingResultDo
	Joins(fields ...field.RelationField) IMatchingResultDo
	Preload(fields ...field.RelationField) IMatchingResultDo
	FirstOrInit() (*models.MatchingResult, error)
	FirstOrCreate() (*models.MatchingResult, error)
	FindByPage(offset int, limit int) (result []*models.MatchingResult, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IMatchingResultDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (m matchingResultDo) Debug() IMatchingResultDo {
	return m.withDO(m.DO.Debug())
}

func (m matchingResultDo) WithContext(ctx context.Context) IMatchingResultDo {
	return m.withDO(m.DO.WithContext(ctx))
}

func (m matchingResultDo) ReadDB() IMatchingResultDo {
	return m.Clauses(dbresolver.Read)
}

func (m matchingResultDo) WriteDB() IMatchingResultDo {
	return m.Clauses(dbresolver.Write)
}

func (m matchingResultDo) Session(config *gorm.Session) IMatchingResultDo {
	return m.withDO(m.DO.Session(config))
}

func (m matchingResultDo) Clauses(conds ...clause.Expression) IMatchingResultDo {
	return m.withDO(m.DO.Clauses(conds...))
}

func (m matchingResultDo) Returning(value interface{}, columns ...string) IMatchingResultDo {
	return m.withDO(m.DO.Returning(value, columns...))
}

func (m matchingResultDo) Not(conds ...gen.Condition) IMatchingResultDo {
	return m.withDO(m.DO.Not(conds...))
}

func (m matchingResultDo) Or(conds ...gen.Condition) IMatchingResultDo {
	return m.withDO(m.DO.Or(conds...))
}

func (m matchingResultDo) Select(conds ...field.Expr) IMatchingResultDo {
	return m.withDO(m.DO.Select(conds...))
}

func (m matchingResultDo) Where(conds ...gen.Condition) IMatchingResultDo {
	return m.withDO(m.DO.Where(conds...))
}

func (m matchingResultDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) IMatchingResultDo {
	return m.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (m matchingResultDo) Order(conds ...field.Expr) IMatchingResultDo {
	return m.withDO(m.DO.Order(conds...))
}

func (m matchingResultDo) Distinct(cols ...field.Expr) IMatchingResultDo {
	return m.withDO(m.DO.Distinct(cols...))
}

func (m matchingResultDo) Omit(cols ...field.Expr) IMatchingResultDo {
	return m.withDO(m.DO.Omit(cols...))
}

func (m matchingResultDo) Join(table schema.Tabler, on ...field.Expr) IMatchingResultDo {
	return m.withDO(m.DO.Join(table, on...))
}

func (m matchingResultDo) LeftJoin(table schema.Tabler, on ...field.Expr) IMatchingResultDo {
	return m.withDO(m.DO.LeftJoin(table, on...))
}

func (m matchingResultDo) RightJoin(table schema.Tabler, on ...field.Expr) IMatchingResultDo {
	return m.withDO(m.DO.RightJoin(table, on...))
}

func (m matchingResultDo) Group(cols ...field.Expr) IMatchingResultDo {
	return m.withDO(m.DO.Group(cols...))
}

func (m matchingResultDo) Having(conds ...gen.Condition) IMatchingResultDo {
	return m.withDO(m.DO.Having(conds...))
}

func (m matchingResultDo) Limit(limit int) IMatchingResultDo {
	return m.withDO(m.DO.Limit(limit))
}

func (m matchingResultDo) Offset(offset int) IMatchingResultDo {
	return m.withDO(m.DO.Offset(offset))
}

func (m matchingResultDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IMatchingResultDo {
	return m.withDO(m.DO.Scopes(funcs...))
}

func (m matchingResultDo) Unscoped() IMatchingResultDo {
	return m.withDO(m.DO.Unscoped())
}

func (m matchingResultDo) Create(values ...*models.MatchingResult) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Create(values)
}

func (m matchingResultDo) CreateInBatches(values []*models.MatchingResult, batchSize int) error {
	return m.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (m matchingResultDo) Save(values ...*models.MatchingResult) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Save(values)
}

func (m matchingResultDo) First() (*models.MatchingResult, error) {
	if result, err := m.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*models.MatchingResult), nil
	}
}

func (m matchingResultDo) Take() (*models.MatchingResult, error) {
	if result, err := m.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*models.MatchingResult), nil
	}
}

func (m matchingResultDo) Last() (*models.MatchingResult, error) {
	if result, err := m.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*models.MatchingResult), nil
	}
}

func (m matchingResultDo) Find() ([]*models.MatchingResult, error) {
	result, err := m.DO.Find()
	return result.([]*models.MatchingResult), err
}

func (m matchingResultDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.MatchingResult, err error) {
	buf := make([]*models.MatchingResult, 0, batchSize)
	err = m.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (m matchingResultDo) FindInBatches(result *[]*models.MatchingResult, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return m.DO.FindInBatches(result, batchSize, fc)
}

func (m matchingResultDo) Attrs(attrs ...field.AssignExpr) IMatchingResultDo {
	return m.withDO(m.DO.Attrs(attrs...))
}

func (m matchingResultDo) Assign(attrs ...field.AssignExpr) IMatchingResultDo {
	return m.withDO(m.DO.Assign(attrs...))
}

func (m matchingResultDo) Joins(fields ...field.RelationField) IMatchingResultDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Joins(_f))
	}
	return &m
}

func (m matchingResultDo) Preload(fields ...field.RelationField) IMatchingResultDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Preload(_f))
	}
	return &m
}

func (m matchingResultDo) FirstOrInit() (*models.MatchingResult, error) {
	if result, err := m.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*models.MatchingResult), nil
	}
}

func (m matchingResultDo) FirstOrCreate() (*models.MatchingResult, error) {
	if result, err := m.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*models.MatchingResult), nil
	}
}

func (m matchingResultDo) FindByPage(offset int, limit int) (result []*models.MatchingResult, count int64, err error) {
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

func (m matchingResultDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = m.Count()
	if err != nil {
		return
	}

	err = m.Offset(offset).Limit(limit).Scan(result)
	return
}

func (m matchingResultDo) Scan(result interface{}) (err error) {
	return m.DO.Scan(result)
}

func (m matchingResultDo) Delete(models ...*models.MatchingResult) (result gen.ResultInfo, err error) {
	return m.DO.Delete(models)
}

func (m *matchingResultDo) withDO(do gen.Dao) *matchingResultDo {
	m.DO = *do.(*gen.DO)
	return m
}
