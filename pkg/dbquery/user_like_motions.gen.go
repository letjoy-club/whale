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

func newUserLikeMotion(db *gorm.DB, opts ...gen.DOOption) userLikeMotion {
	_userLikeMotion := userLikeMotion{}

	_userLikeMotion.userLikeMotionDo.UseDB(db, opts...)
	_userLikeMotion.userLikeMotionDo.UseModel(&models.UserLikeMotion{})

	tableName := _userLikeMotion.userLikeMotionDo.TableName()
	_userLikeMotion.ALL = field.NewAsterisk(tableName)
	_userLikeMotion.ID = field.NewInt(tableName, "id")
	_userLikeMotion.ToMotionID = field.NewString(tableName, "to_motion_id")
	_userLikeMotion.ToUserID = field.NewString(tableName, "to_user_id")
	_userLikeMotion.UserID = field.NewString(tableName, "user_id")
	_userLikeMotion.CreatedAt = field.NewTime(tableName, "created_at")

	_userLikeMotion.fillFieldMap()

	return _userLikeMotion
}

type userLikeMotion struct {
	userLikeMotionDo userLikeMotionDo

	ALL        field.Asterisk
	ID         field.Int
	ToMotionID field.String
	ToUserID   field.String
	UserID     field.String
	CreatedAt  field.Time

	fieldMap map[string]field.Expr
}

func (u userLikeMotion) Table(newTableName string) *userLikeMotion {
	u.userLikeMotionDo.UseTable(newTableName)
	return u.updateTableName(newTableName)
}

func (u userLikeMotion) As(alias string) *userLikeMotion {
	u.userLikeMotionDo.DO = *(u.userLikeMotionDo.As(alias).(*gen.DO))
	return u.updateTableName(alias)
}

func (u *userLikeMotion) updateTableName(table string) *userLikeMotion {
	u.ALL = field.NewAsterisk(table)
	u.ID = field.NewInt(table, "id")
	u.ToMotionID = field.NewString(table, "to_motion_id")
	u.ToUserID = field.NewString(table, "to_user_id")
	u.UserID = field.NewString(table, "user_id")
	u.CreatedAt = field.NewTime(table, "created_at")

	u.fillFieldMap()

	return u
}

func (u *userLikeMotion) WithContext(ctx context.Context) IUserLikeMotionDo {
	return u.userLikeMotionDo.WithContext(ctx)
}

func (u userLikeMotion) TableName() string { return u.userLikeMotionDo.TableName() }

func (u userLikeMotion) Alias() string { return u.userLikeMotionDo.Alias() }

func (u *userLikeMotion) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := u.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (u *userLikeMotion) fillFieldMap() {
	u.fieldMap = make(map[string]field.Expr, 5)
	u.fieldMap["id"] = u.ID
	u.fieldMap["to_motion_id"] = u.ToMotionID
	u.fieldMap["to_user_id"] = u.ToUserID
	u.fieldMap["user_id"] = u.UserID
	u.fieldMap["created_at"] = u.CreatedAt
}

func (u userLikeMotion) clone(db *gorm.DB) userLikeMotion {
	u.userLikeMotionDo.ReplaceConnPool(db.Statement.ConnPool)
	return u
}

func (u userLikeMotion) replaceDB(db *gorm.DB) userLikeMotion {
	u.userLikeMotionDo.ReplaceDB(db)
	return u
}

type userLikeMotionDo struct{ gen.DO }

type IUserLikeMotionDo interface {
	gen.SubQuery
	Debug() IUserLikeMotionDo
	WithContext(ctx context.Context) IUserLikeMotionDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IUserLikeMotionDo
	WriteDB() IUserLikeMotionDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IUserLikeMotionDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IUserLikeMotionDo
	Not(conds ...gen.Condition) IUserLikeMotionDo
	Or(conds ...gen.Condition) IUserLikeMotionDo
	Select(conds ...field.Expr) IUserLikeMotionDo
	Where(conds ...gen.Condition) IUserLikeMotionDo
	Order(conds ...field.Expr) IUserLikeMotionDo
	Distinct(cols ...field.Expr) IUserLikeMotionDo
	Omit(cols ...field.Expr) IUserLikeMotionDo
	Join(table schema.Tabler, on ...field.Expr) IUserLikeMotionDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IUserLikeMotionDo
	RightJoin(table schema.Tabler, on ...field.Expr) IUserLikeMotionDo
	Group(cols ...field.Expr) IUserLikeMotionDo
	Having(conds ...gen.Condition) IUserLikeMotionDo
	Limit(limit int) IUserLikeMotionDo
	Offset(offset int) IUserLikeMotionDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IUserLikeMotionDo
	Unscoped() IUserLikeMotionDo
	Create(values ...*models.UserLikeMotion) error
	CreateInBatches(values []*models.UserLikeMotion, batchSize int) error
	Save(values ...*models.UserLikeMotion) error
	First() (*models.UserLikeMotion, error)
	Take() (*models.UserLikeMotion, error)
	Last() (*models.UserLikeMotion, error)
	Find() ([]*models.UserLikeMotion, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.UserLikeMotion, err error)
	FindInBatches(result *[]*models.UserLikeMotion, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*models.UserLikeMotion) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IUserLikeMotionDo
	Assign(attrs ...field.AssignExpr) IUserLikeMotionDo
	Joins(fields ...field.RelationField) IUserLikeMotionDo
	Preload(fields ...field.RelationField) IUserLikeMotionDo
	FirstOrInit() (*models.UserLikeMotion, error)
	FirstOrCreate() (*models.UserLikeMotion, error)
	FindByPage(offset int, limit int) (result []*models.UserLikeMotion, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IUserLikeMotionDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (u userLikeMotionDo) Debug() IUserLikeMotionDo {
	return u.withDO(u.DO.Debug())
}

func (u userLikeMotionDo) WithContext(ctx context.Context) IUserLikeMotionDo {
	return u.withDO(u.DO.WithContext(ctx))
}

func (u userLikeMotionDo) ReadDB() IUserLikeMotionDo {
	return u.Clauses(dbresolver.Read)
}

func (u userLikeMotionDo) WriteDB() IUserLikeMotionDo {
	return u.Clauses(dbresolver.Write)
}

func (u userLikeMotionDo) Session(config *gorm.Session) IUserLikeMotionDo {
	return u.withDO(u.DO.Session(config))
}

func (u userLikeMotionDo) Clauses(conds ...clause.Expression) IUserLikeMotionDo {
	return u.withDO(u.DO.Clauses(conds...))
}

func (u userLikeMotionDo) Returning(value interface{}, columns ...string) IUserLikeMotionDo {
	return u.withDO(u.DO.Returning(value, columns...))
}

func (u userLikeMotionDo) Not(conds ...gen.Condition) IUserLikeMotionDo {
	return u.withDO(u.DO.Not(conds...))
}

func (u userLikeMotionDo) Or(conds ...gen.Condition) IUserLikeMotionDo {
	return u.withDO(u.DO.Or(conds...))
}

func (u userLikeMotionDo) Select(conds ...field.Expr) IUserLikeMotionDo {
	return u.withDO(u.DO.Select(conds...))
}

func (u userLikeMotionDo) Where(conds ...gen.Condition) IUserLikeMotionDo {
	return u.withDO(u.DO.Where(conds...))
}

func (u userLikeMotionDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) IUserLikeMotionDo {
	return u.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (u userLikeMotionDo) Order(conds ...field.Expr) IUserLikeMotionDo {
	return u.withDO(u.DO.Order(conds...))
}

func (u userLikeMotionDo) Distinct(cols ...field.Expr) IUserLikeMotionDo {
	return u.withDO(u.DO.Distinct(cols...))
}

func (u userLikeMotionDo) Omit(cols ...field.Expr) IUserLikeMotionDo {
	return u.withDO(u.DO.Omit(cols...))
}

func (u userLikeMotionDo) Join(table schema.Tabler, on ...field.Expr) IUserLikeMotionDo {
	return u.withDO(u.DO.Join(table, on...))
}

func (u userLikeMotionDo) LeftJoin(table schema.Tabler, on ...field.Expr) IUserLikeMotionDo {
	return u.withDO(u.DO.LeftJoin(table, on...))
}

func (u userLikeMotionDo) RightJoin(table schema.Tabler, on ...field.Expr) IUserLikeMotionDo {
	return u.withDO(u.DO.RightJoin(table, on...))
}

func (u userLikeMotionDo) Group(cols ...field.Expr) IUserLikeMotionDo {
	return u.withDO(u.DO.Group(cols...))
}

func (u userLikeMotionDo) Having(conds ...gen.Condition) IUserLikeMotionDo {
	return u.withDO(u.DO.Having(conds...))
}

func (u userLikeMotionDo) Limit(limit int) IUserLikeMotionDo {
	return u.withDO(u.DO.Limit(limit))
}

func (u userLikeMotionDo) Offset(offset int) IUserLikeMotionDo {
	return u.withDO(u.DO.Offset(offset))
}

func (u userLikeMotionDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IUserLikeMotionDo {
	return u.withDO(u.DO.Scopes(funcs...))
}

func (u userLikeMotionDo) Unscoped() IUserLikeMotionDo {
	return u.withDO(u.DO.Unscoped())
}

func (u userLikeMotionDo) Create(values ...*models.UserLikeMotion) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Create(values)
}

func (u userLikeMotionDo) CreateInBatches(values []*models.UserLikeMotion, batchSize int) error {
	return u.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (u userLikeMotionDo) Save(values ...*models.UserLikeMotion) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Save(values)
}

func (u userLikeMotionDo) First() (*models.UserLikeMotion, error) {
	if result, err := u.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*models.UserLikeMotion), nil
	}
}

func (u userLikeMotionDo) Take() (*models.UserLikeMotion, error) {
	if result, err := u.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*models.UserLikeMotion), nil
	}
}

func (u userLikeMotionDo) Last() (*models.UserLikeMotion, error) {
	if result, err := u.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*models.UserLikeMotion), nil
	}
}

func (u userLikeMotionDo) Find() ([]*models.UserLikeMotion, error) {
	result, err := u.DO.Find()
	return result.([]*models.UserLikeMotion), err
}

func (u userLikeMotionDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.UserLikeMotion, err error) {
	buf := make([]*models.UserLikeMotion, 0, batchSize)
	err = u.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (u userLikeMotionDo) FindInBatches(result *[]*models.UserLikeMotion, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return u.DO.FindInBatches(result, batchSize, fc)
}

func (u userLikeMotionDo) Attrs(attrs ...field.AssignExpr) IUserLikeMotionDo {
	return u.withDO(u.DO.Attrs(attrs...))
}

func (u userLikeMotionDo) Assign(attrs ...field.AssignExpr) IUserLikeMotionDo {
	return u.withDO(u.DO.Assign(attrs...))
}

func (u userLikeMotionDo) Joins(fields ...field.RelationField) IUserLikeMotionDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Joins(_f))
	}
	return &u
}

func (u userLikeMotionDo) Preload(fields ...field.RelationField) IUserLikeMotionDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Preload(_f))
	}
	return &u
}

func (u userLikeMotionDo) FirstOrInit() (*models.UserLikeMotion, error) {
	if result, err := u.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*models.UserLikeMotion), nil
	}
}

func (u userLikeMotionDo) FirstOrCreate() (*models.UserLikeMotion, error) {
	if result, err := u.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*models.UserLikeMotion), nil
	}
}

func (u userLikeMotionDo) FindByPage(offset int, limit int) (result []*models.UserLikeMotion, count int64, err error) {
	result, err = u.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = u.Offset(-1).Limit(-1).Count()
	return
}

func (u userLikeMotionDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = u.Count()
	if err != nil {
		return
	}

	err = u.Offset(offset).Limit(limit).Scan(result)
	return
}

func (u userLikeMotionDo) Scan(result interface{}) (err error) {
	return u.DO.Scan(result)
}

func (u userLikeMotionDo) Delete(models ...*models.UserLikeMotion) (result gen.ResultInfo, err error) {
	return u.DO.Delete(models)
}

func (u *userLikeMotionDo) withDO(do gen.Dao) *userLikeMotionDo {
	u.DO = *do.(*gen.DO)
	return u
}
