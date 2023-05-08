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

func newMatchingReview(db *gorm.DB, opts ...gen.DOOption) matchingReview {
	_matchingReview := matchingReview{}

	_matchingReview.matchingReviewDo.UseDB(db, opts...)
	_matchingReview.matchingReviewDo.UseModel(&models.MatchingReview{})

	tableName := _matchingReview.matchingReviewDo.TableName()
	_matchingReview.ALL = field.NewAsterisk(tableName)
	_matchingReview.ID = field.NewInt(tableName, "id")
	_matchingReview.MatchingResultID = field.NewInt(tableName, "matching_result_id")
	_matchingReview.MatchingID = field.NewString(tableName, "matching_id")
	_matchingReview.UserID = field.NewString(tableName, "user_id")
	_matchingReview.ToMatchingID = field.NewString(tableName, "to_matching_id")
	_matchingReview.ToUserID = field.NewString(tableName, "to_user_id")
	_matchingReview.TopicID = field.NewString(tableName, "topic_id")
	_matchingReview.Score = field.NewInt(tableName, "score")
	_matchingReview.Comment = field.NewString(tableName, "comment")
	_matchingReview.CreateTime = field.NewTime(tableName, "create_time")

	_matchingReview.fillFieldMap()

	return _matchingReview
}

type matchingReview struct {
	matchingReviewDo matchingReviewDo

	ALL              field.Asterisk
	ID               field.Int
	MatchingResultID field.Int
	MatchingID       field.String
	UserID           field.String
	ToMatchingID     field.String
	ToUserID         field.String
	TopicID          field.String
	Score            field.Int
	Comment          field.String
	CreateTime       field.Time

	fieldMap map[string]field.Expr
}

func (m matchingReview) Table(newTableName string) *matchingReview {
	m.matchingReviewDo.UseTable(newTableName)
	return m.updateTableName(newTableName)
}

func (m matchingReview) As(alias string) *matchingReview {
	m.matchingReviewDo.DO = *(m.matchingReviewDo.As(alias).(*gen.DO))
	return m.updateTableName(alias)
}

func (m *matchingReview) updateTableName(table string) *matchingReview {
	m.ALL = field.NewAsterisk(table)
	m.ID = field.NewInt(table, "id")
	m.MatchingResultID = field.NewInt(table, "matching_result_id")
	m.MatchingID = field.NewString(table, "matching_id")
	m.UserID = field.NewString(table, "user_id")
	m.ToMatchingID = field.NewString(table, "to_matching_id")
	m.ToUserID = field.NewString(table, "to_user_id")
	m.TopicID = field.NewString(table, "topic_id")
	m.Score = field.NewInt(table, "score")
	m.Comment = field.NewString(table, "comment")
	m.CreateTime = field.NewTime(table, "create_time")

	m.fillFieldMap()

	return m
}

func (m *matchingReview) WithContext(ctx context.Context) IMatchingReviewDo {
	return m.matchingReviewDo.WithContext(ctx)
}

func (m matchingReview) TableName() string { return m.matchingReviewDo.TableName() }

func (m matchingReview) Alias() string { return m.matchingReviewDo.Alias() }

func (m *matchingReview) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := m.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (m *matchingReview) fillFieldMap() {
	m.fieldMap = make(map[string]field.Expr, 10)
	m.fieldMap["id"] = m.ID
	m.fieldMap["matching_result_id"] = m.MatchingResultID
	m.fieldMap["matching_id"] = m.MatchingID
	m.fieldMap["user_id"] = m.UserID
	m.fieldMap["to_matching_id"] = m.ToMatchingID
	m.fieldMap["to_user_id"] = m.ToUserID
	m.fieldMap["topic_id"] = m.TopicID
	m.fieldMap["score"] = m.Score
	m.fieldMap["comment"] = m.Comment
	m.fieldMap["create_time"] = m.CreateTime
}

func (m matchingReview) clone(db *gorm.DB) matchingReview {
	m.matchingReviewDo.ReplaceConnPool(db.Statement.ConnPool)
	return m
}

func (m matchingReview) replaceDB(db *gorm.DB) matchingReview {
	m.matchingReviewDo.ReplaceDB(db)
	return m
}

type matchingReviewDo struct{ gen.DO }

type IMatchingReviewDo interface {
	gen.SubQuery
	Debug() IMatchingReviewDo
	WithContext(ctx context.Context) IMatchingReviewDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IMatchingReviewDo
	WriteDB() IMatchingReviewDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IMatchingReviewDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IMatchingReviewDo
	Not(conds ...gen.Condition) IMatchingReviewDo
	Or(conds ...gen.Condition) IMatchingReviewDo
	Select(conds ...field.Expr) IMatchingReviewDo
	Where(conds ...gen.Condition) IMatchingReviewDo
	Order(conds ...field.Expr) IMatchingReviewDo
	Distinct(cols ...field.Expr) IMatchingReviewDo
	Omit(cols ...field.Expr) IMatchingReviewDo
	Join(table schema.Tabler, on ...field.Expr) IMatchingReviewDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IMatchingReviewDo
	RightJoin(table schema.Tabler, on ...field.Expr) IMatchingReviewDo
	Group(cols ...field.Expr) IMatchingReviewDo
	Having(conds ...gen.Condition) IMatchingReviewDo
	Limit(limit int) IMatchingReviewDo
	Offset(offset int) IMatchingReviewDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IMatchingReviewDo
	Unscoped() IMatchingReviewDo
	Create(values ...*models.MatchingReview) error
	CreateInBatches(values []*models.MatchingReview, batchSize int) error
	Save(values ...*models.MatchingReview) error
	First() (*models.MatchingReview, error)
	Take() (*models.MatchingReview, error)
	Last() (*models.MatchingReview, error)
	Find() ([]*models.MatchingReview, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.MatchingReview, err error)
	FindInBatches(result *[]*models.MatchingReview, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*models.MatchingReview) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IMatchingReviewDo
	Assign(attrs ...field.AssignExpr) IMatchingReviewDo
	Joins(fields ...field.RelationField) IMatchingReviewDo
	Preload(fields ...field.RelationField) IMatchingReviewDo
	FirstOrInit() (*models.MatchingReview, error)
	FirstOrCreate() (*models.MatchingReview, error)
	FindByPage(offset int, limit int) (result []*models.MatchingReview, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IMatchingReviewDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (m matchingReviewDo) Debug() IMatchingReviewDo {
	return m.withDO(m.DO.Debug())
}

func (m matchingReviewDo) WithContext(ctx context.Context) IMatchingReviewDo {
	return m.withDO(m.DO.WithContext(ctx))
}

func (m matchingReviewDo) ReadDB() IMatchingReviewDo {
	return m.Clauses(dbresolver.Read)
}

func (m matchingReviewDo) WriteDB() IMatchingReviewDo {
	return m.Clauses(dbresolver.Write)
}

func (m matchingReviewDo) Session(config *gorm.Session) IMatchingReviewDo {
	return m.withDO(m.DO.Session(config))
}

func (m matchingReviewDo) Clauses(conds ...clause.Expression) IMatchingReviewDo {
	return m.withDO(m.DO.Clauses(conds...))
}

func (m matchingReviewDo) Returning(value interface{}, columns ...string) IMatchingReviewDo {
	return m.withDO(m.DO.Returning(value, columns...))
}

func (m matchingReviewDo) Not(conds ...gen.Condition) IMatchingReviewDo {
	return m.withDO(m.DO.Not(conds...))
}

func (m matchingReviewDo) Or(conds ...gen.Condition) IMatchingReviewDo {
	return m.withDO(m.DO.Or(conds...))
}

func (m matchingReviewDo) Select(conds ...field.Expr) IMatchingReviewDo {
	return m.withDO(m.DO.Select(conds...))
}

func (m matchingReviewDo) Where(conds ...gen.Condition) IMatchingReviewDo {
	return m.withDO(m.DO.Where(conds...))
}

func (m matchingReviewDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) IMatchingReviewDo {
	return m.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (m matchingReviewDo) Order(conds ...field.Expr) IMatchingReviewDo {
	return m.withDO(m.DO.Order(conds...))
}

func (m matchingReviewDo) Distinct(cols ...field.Expr) IMatchingReviewDo {
	return m.withDO(m.DO.Distinct(cols...))
}

func (m matchingReviewDo) Omit(cols ...field.Expr) IMatchingReviewDo {
	return m.withDO(m.DO.Omit(cols...))
}

func (m matchingReviewDo) Join(table schema.Tabler, on ...field.Expr) IMatchingReviewDo {
	return m.withDO(m.DO.Join(table, on...))
}

func (m matchingReviewDo) LeftJoin(table schema.Tabler, on ...field.Expr) IMatchingReviewDo {
	return m.withDO(m.DO.LeftJoin(table, on...))
}

func (m matchingReviewDo) RightJoin(table schema.Tabler, on ...field.Expr) IMatchingReviewDo {
	return m.withDO(m.DO.RightJoin(table, on...))
}

func (m matchingReviewDo) Group(cols ...field.Expr) IMatchingReviewDo {
	return m.withDO(m.DO.Group(cols...))
}

func (m matchingReviewDo) Having(conds ...gen.Condition) IMatchingReviewDo {
	return m.withDO(m.DO.Having(conds...))
}

func (m matchingReviewDo) Limit(limit int) IMatchingReviewDo {
	return m.withDO(m.DO.Limit(limit))
}

func (m matchingReviewDo) Offset(offset int) IMatchingReviewDo {
	return m.withDO(m.DO.Offset(offset))
}

func (m matchingReviewDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IMatchingReviewDo {
	return m.withDO(m.DO.Scopes(funcs...))
}

func (m matchingReviewDo) Unscoped() IMatchingReviewDo {
	return m.withDO(m.DO.Unscoped())
}

func (m matchingReviewDo) Create(values ...*models.MatchingReview) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Create(values)
}

func (m matchingReviewDo) CreateInBatches(values []*models.MatchingReview, batchSize int) error {
	return m.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (m matchingReviewDo) Save(values ...*models.MatchingReview) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Save(values)
}

func (m matchingReviewDo) First() (*models.MatchingReview, error) {
	if result, err := m.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*models.MatchingReview), nil
	}
}

func (m matchingReviewDo) Take() (*models.MatchingReview, error) {
	if result, err := m.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*models.MatchingReview), nil
	}
}

func (m matchingReviewDo) Last() (*models.MatchingReview, error) {
	if result, err := m.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*models.MatchingReview), nil
	}
}

func (m matchingReviewDo) Find() ([]*models.MatchingReview, error) {
	result, err := m.DO.Find()
	return result.([]*models.MatchingReview), err
}

func (m matchingReviewDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.MatchingReview, err error) {
	buf := make([]*models.MatchingReview, 0, batchSize)
	err = m.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (m matchingReviewDo) FindInBatches(result *[]*models.MatchingReview, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return m.DO.FindInBatches(result, batchSize, fc)
}

func (m matchingReviewDo) Attrs(attrs ...field.AssignExpr) IMatchingReviewDo {
	return m.withDO(m.DO.Attrs(attrs...))
}

func (m matchingReviewDo) Assign(attrs ...field.AssignExpr) IMatchingReviewDo {
	return m.withDO(m.DO.Assign(attrs...))
}

func (m matchingReviewDo) Joins(fields ...field.RelationField) IMatchingReviewDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Joins(_f))
	}
	return &m
}

func (m matchingReviewDo) Preload(fields ...field.RelationField) IMatchingReviewDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Preload(_f))
	}
	return &m
}

func (m matchingReviewDo) FirstOrInit() (*models.MatchingReview, error) {
	if result, err := m.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*models.MatchingReview), nil
	}
}

func (m matchingReviewDo) FirstOrCreate() (*models.MatchingReview, error) {
	if result, err := m.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*models.MatchingReview), nil
	}
}

func (m matchingReviewDo) FindByPage(offset int, limit int) (result []*models.MatchingReview, count int64, err error) {
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

func (m matchingReviewDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = m.Count()
	if err != nil {
		return
	}

	err = m.Offset(offset).Limit(limit).Scan(result)
	return
}

func (m matchingReviewDo) Scan(result interface{}) (err error) {
	return m.DO.Scan(result)
}

func (m matchingReviewDo) Delete(models ...*models.MatchingReview) (result gen.ResultInfo, err error) {
	return m.DO.Delete(models)
}

func (m *matchingReviewDo) withDO(do gen.Dao) *matchingReviewDo {
	m.DO = *do.(*gen.DO)
	return m
}
