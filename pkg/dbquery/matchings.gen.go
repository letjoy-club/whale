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

func newMatching(db *gorm.DB, opts ...gen.DOOption) matching {
	_matching := matching{}

	_matching.matchingDo.UseDB(db, opts...)
	_matching.matchingDo.UseModel(&models.Matching{})

	tableName := _matching.matchingDo.TableName()
	_matching.ALL = field.NewAsterisk(tableName)
	_matching.ID = field.NewString(tableName, "id")
	_matching.TopicID = field.NewString(tableName, "topic_id")
	_matching.UserID = field.NewString(tableName, "user_id")
	_matching.AreaIDs = field.NewField(tableName, "area_ids")
	_matching.CityID = field.NewString(tableName, "city_id")
	_matching.Gender = field.NewString(tableName, "gender")
	_matching.MyGender = field.NewString(tableName, "my_gender")
	_matching.RejectedUserIDs = field.NewField(tableName, "rejected_user_ids")
	_matching.InChatGroup = field.NewBool(tableName, "in_chat_group")
	_matching.State = field.NewString(tableName, "state")
	_matching.ChatGroupState = field.NewString(tableName, "chat_group_state")
	_matching.ResultID = field.NewInt(tableName, "result_id")
	_matching.Remark = field.NewString(tableName, "remark")
	_matching.DayRange = field.NewField(tableName, "day_range")
	_matching.PreferredPeriods = field.NewField(tableName, "preferred_periods")
	_matching.Properties = field.NewField(tableName, "properties")
	_matching.StartMatchingAt = field.NewTime(tableName, "start_matching_at")
	_matching.RelatedMotionID = field.NewString(tableName, "related_motion_id")
	_matching.FinishedAt = field.NewTime(tableName, "finished_at")
	_matching.MatchedAt = field.NewTime(tableName, "matched_at")
	_matching.Deadline = field.NewTime(tableName, "deadline")
	_matching.CreatedAt = field.NewTime(tableName, "created_at")
	_matching.UpdatedAt = field.NewTime(tableName, "updated_at")

	_matching.fillFieldMap()

	return _matching
}

type matching struct {
	matchingDo matchingDo

	ALL              field.Asterisk
	ID               field.String
	TopicID          field.String
	UserID           field.String
	AreaIDs          field.Field
	CityID           field.String
	Gender           field.String
	MyGender         field.String
	RejectedUserIDs  field.Field
	InChatGroup      field.Bool
	State            field.String
	ChatGroupState   field.String
	ResultID         field.Int
	Remark           field.String
	DayRange         field.Field
	PreferredPeriods field.Field
	Properties       field.Field
	StartMatchingAt  field.Time
	RelatedMotionID  field.String
	FinishedAt       field.Time
	MatchedAt        field.Time
	Deadline         field.Time
	CreatedAt        field.Time
	UpdatedAt        field.Time

	fieldMap map[string]field.Expr
}

func (m matching) Table(newTableName string) *matching {
	m.matchingDo.UseTable(newTableName)
	return m.updateTableName(newTableName)
}

func (m matching) As(alias string) *matching {
	m.matchingDo.DO = *(m.matchingDo.As(alias).(*gen.DO))
	return m.updateTableName(alias)
}

func (m *matching) updateTableName(table string) *matching {
	m.ALL = field.NewAsterisk(table)
	m.ID = field.NewString(table, "id")
	m.TopicID = field.NewString(table, "topic_id")
	m.UserID = field.NewString(table, "user_id")
	m.AreaIDs = field.NewField(table, "area_ids")
	m.CityID = field.NewString(table, "city_id")
	m.Gender = field.NewString(table, "gender")
	m.MyGender = field.NewString(table, "my_gender")
	m.RejectedUserIDs = field.NewField(table, "rejected_user_ids")
	m.InChatGroup = field.NewBool(table, "in_chat_group")
	m.State = field.NewString(table, "state")
	m.ChatGroupState = field.NewString(table, "chat_group_state")
	m.ResultID = field.NewInt(table, "result_id")
	m.Remark = field.NewString(table, "remark")
	m.DayRange = field.NewField(table, "day_range")
	m.PreferredPeriods = field.NewField(table, "preferred_periods")
	m.Properties = field.NewField(table, "properties")
	m.StartMatchingAt = field.NewTime(table, "start_matching_at")
	m.RelatedMotionID = field.NewString(table, "related_motion_id")
	m.FinishedAt = field.NewTime(table, "finished_at")
	m.MatchedAt = field.NewTime(table, "matched_at")
	m.Deadline = field.NewTime(table, "deadline")
	m.CreatedAt = field.NewTime(table, "created_at")
	m.UpdatedAt = field.NewTime(table, "updated_at")

	m.fillFieldMap()

	return m
}

func (m *matching) WithContext(ctx context.Context) IMatchingDo { return m.matchingDo.WithContext(ctx) }

func (m matching) TableName() string { return m.matchingDo.TableName() }

func (m matching) Alias() string { return m.matchingDo.Alias() }

func (m *matching) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := m.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (m *matching) fillFieldMap() {
	m.fieldMap = make(map[string]field.Expr, 23)
	m.fieldMap["id"] = m.ID
	m.fieldMap["topic_id"] = m.TopicID
	m.fieldMap["user_id"] = m.UserID
	m.fieldMap["area_ids"] = m.AreaIDs
	m.fieldMap["city_id"] = m.CityID
	m.fieldMap["gender"] = m.Gender
	m.fieldMap["my_gender"] = m.MyGender
	m.fieldMap["rejected_user_ids"] = m.RejectedUserIDs
	m.fieldMap["in_chat_group"] = m.InChatGroup
	m.fieldMap["state"] = m.State
	m.fieldMap["chat_group_state"] = m.ChatGroupState
	m.fieldMap["result_id"] = m.ResultID
	m.fieldMap["remark"] = m.Remark
	m.fieldMap["day_range"] = m.DayRange
	m.fieldMap["preferred_periods"] = m.PreferredPeriods
	m.fieldMap["properties"] = m.Properties
	m.fieldMap["start_matching_at"] = m.StartMatchingAt
	m.fieldMap["related_motion_id"] = m.RelatedMotionID
	m.fieldMap["finished_at"] = m.FinishedAt
	m.fieldMap["matched_at"] = m.MatchedAt
	m.fieldMap["deadline"] = m.Deadline
	m.fieldMap["created_at"] = m.CreatedAt
	m.fieldMap["updated_at"] = m.UpdatedAt
}

func (m matching) clone(db *gorm.DB) matching {
	m.matchingDo.ReplaceConnPool(db.Statement.ConnPool)
	return m
}

func (m matching) replaceDB(db *gorm.DB) matching {
	m.matchingDo.ReplaceDB(db)
	return m
}

type matchingDo struct{ gen.DO }

type IMatchingDo interface {
	gen.SubQuery
	Debug() IMatchingDo
	WithContext(ctx context.Context) IMatchingDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IMatchingDo
	WriteDB() IMatchingDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IMatchingDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IMatchingDo
	Not(conds ...gen.Condition) IMatchingDo
	Or(conds ...gen.Condition) IMatchingDo
	Select(conds ...field.Expr) IMatchingDo
	Where(conds ...gen.Condition) IMatchingDo
	Order(conds ...field.Expr) IMatchingDo
	Distinct(cols ...field.Expr) IMatchingDo
	Omit(cols ...field.Expr) IMatchingDo
	Join(table schema.Tabler, on ...field.Expr) IMatchingDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IMatchingDo
	RightJoin(table schema.Tabler, on ...field.Expr) IMatchingDo
	Group(cols ...field.Expr) IMatchingDo
	Having(conds ...gen.Condition) IMatchingDo
	Limit(limit int) IMatchingDo
	Offset(offset int) IMatchingDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IMatchingDo
	Unscoped() IMatchingDo
	Create(values ...*models.Matching) error
	CreateInBatches(values []*models.Matching, batchSize int) error
	Save(values ...*models.Matching) error
	First() (*models.Matching, error)
	Take() (*models.Matching, error)
	Last() (*models.Matching, error)
	Find() ([]*models.Matching, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.Matching, err error)
	FindInBatches(result *[]*models.Matching, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*models.Matching) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IMatchingDo
	Assign(attrs ...field.AssignExpr) IMatchingDo
	Joins(fields ...field.RelationField) IMatchingDo
	Preload(fields ...field.RelationField) IMatchingDo
	FirstOrInit() (*models.Matching, error)
	FirstOrCreate() (*models.Matching, error)
	FindByPage(offset int, limit int) (result []*models.Matching, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IMatchingDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (m matchingDo) Debug() IMatchingDo {
	return m.withDO(m.DO.Debug())
}

func (m matchingDo) WithContext(ctx context.Context) IMatchingDo {
	return m.withDO(m.DO.WithContext(ctx))
}

func (m matchingDo) ReadDB() IMatchingDo {
	return m.Clauses(dbresolver.Read)
}

func (m matchingDo) WriteDB() IMatchingDo {
	return m.Clauses(dbresolver.Write)
}

func (m matchingDo) Session(config *gorm.Session) IMatchingDo {
	return m.withDO(m.DO.Session(config))
}

func (m matchingDo) Clauses(conds ...clause.Expression) IMatchingDo {
	return m.withDO(m.DO.Clauses(conds...))
}

func (m matchingDo) Returning(value interface{}, columns ...string) IMatchingDo {
	return m.withDO(m.DO.Returning(value, columns...))
}

func (m matchingDo) Not(conds ...gen.Condition) IMatchingDo {
	return m.withDO(m.DO.Not(conds...))
}

func (m matchingDo) Or(conds ...gen.Condition) IMatchingDo {
	return m.withDO(m.DO.Or(conds...))
}

func (m matchingDo) Select(conds ...field.Expr) IMatchingDo {
	return m.withDO(m.DO.Select(conds...))
}

func (m matchingDo) Where(conds ...gen.Condition) IMatchingDo {
	return m.withDO(m.DO.Where(conds...))
}

func (m matchingDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) IMatchingDo {
	return m.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (m matchingDo) Order(conds ...field.Expr) IMatchingDo {
	return m.withDO(m.DO.Order(conds...))
}

func (m matchingDo) Distinct(cols ...field.Expr) IMatchingDo {
	return m.withDO(m.DO.Distinct(cols...))
}

func (m matchingDo) Omit(cols ...field.Expr) IMatchingDo {
	return m.withDO(m.DO.Omit(cols...))
}

func (m matchingDo) Join(table schema.Tabler, on ...field.Expr) IMatchingDo {
	return m.withDO(m.DO.Join(table, on...))
}

func (m matchingDo) LeftJoin(table schema.Tabler, on ...field.Expr) IMatchingDo {
	return m.withDO(m.DO.LeftJoin(table, on...))
}

func (m matchingDo) RightJoin(table schema.Tabler, on ...field.Expr) IMatchingDo {
	return m.withDO(m.DO.RightJoin(table, on...))
}

func (m matchingDo) Group(cols ...field.Expr) IMatchingDo {
	return m.withDO(m.DO.Group(cols...))
}

func (m matchingDo) Having(conds ...gen.Condition) IMatchingDo {
	return m.withDO(m.DO.Having(conds...))
}

func (m matchingDo) Limit(limit int) IMatchingDo {
	return m.withDO(m.DO.Limit(limit))
}

func (m matchingDo) Offset(offset int) IMatchingDo {
	return m.withDO(m.DO.Offset(offset))
}

func (m matchingDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IMatchingDo {
	return m.withDO(m.DO.Scopes(funcs...))
}

func (m matchingDo) Unscoped() IMatchingDo {
	return m.withDO(m.DO.Unscoped())
}

func (m matchingDo) Create(values ...*models.Matching) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Create(values)
}

func (m matchingDo) CreateInBatches(values []*models.Matching, batchSize int) error {
	return m.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (m matchingDo) Save(values ...*models.Matching) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Save(values)
}

func (m matchingDo) First() (*models.Matching, error) {
	if result, err := m.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*models.Matching), nil
	}
}

func (m matchingDo) Take() (*models.Matching, error) {
	if result, err := m.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*models.Matching), nil
	}
}

func (m matchingDo) Last() (*models.Matching, error) {
	if result, err := m.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*models.Matching), nil
	}
}

func (m matchingDo) Find() ([]*models.Matching, error) {
	result, err := m.DO.Find()
	return result.([]*models.Matching), err
}

func (m matchingDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.Matching, err error) {
	buf := make([]*models.Matching, 0, batchSize)
	err = m.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (m matchingDo) FindInBatches(result *[]*models.Matching, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return m.DO.FindInBatches(result, batchSize, fc)
}

func (m matchingDo) Attrs(attrs ...field.AssignExpr) IMatchingDo {
	return m.withDO(m.DO.Attrs(attrs...))
}

func (m matchingDo) Assign(attrs ...field.AssignExpr) IMatchingDo {
	return m.withDO(m.DO.Assign(attrs...))
}

func (m matchingDo) Joins(fields ...field.RelationField) IMatchingDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Joins(_f))
	}
	return &m
}

func (m matchingDo) Preload(fields ...field.RelationField) IMatchingDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Preload(_f))
	}
	return &m
}

func (m matchingDo) FirstOrInit() (*models.Matching, error) {
	if result, err := m.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*models.Matching), nil
	}
}

func (m matchingDo) FirstOrCreate() (*models.Matching, error) {
	if result, err := m.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*models.Matching), nil
	}
}

func (m matchingDo) FindByPage(offset int, limit int) (result []*models.Matching, count int64, err error) {
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

func (m matchingDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = m.Count()
	if err != nil {
		return
	}

	err = m.Offset(offset).Limit(limit).Scan(result)
	return
}

func (m matchingDo) Scan(result interface{}) (err error) {
	return m.DO.Scan(result)
}

func (m matchingDo) Delete(models ...*models.Matching) (result gen.ResultInfo, err error) {
	return m.DO.Delete(models)
}

func (m *matchingDo) withDO(do gen.Dao) *matchingDo {
	m.DO = *do.(*gen.DO)
	return m
}
