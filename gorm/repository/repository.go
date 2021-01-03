package repository

import (
	"fmt"
	"github.com/gogf/gf/os/glog"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"math"
	"regexp"
)

// Holds information about result page
type Page struct {
	Total   uint
	Page    uint
	PerPage uint
	Pages   uint
}

type GormRepository struct {
	logger       *glog.Logger
	db           *gorm.DB
	defaultJoins []string
}

// NewGormRepository returns a new base repository that implements TransactionRepository
func NewGormRepository(db *gorm.DB, logger *glog.Logger, defaultJoins ...string) TransactionRepository {
	return &GormRepository{
		defaultJoins: defaultJoins,
		logger:       logger,
		db:           db,
	}
}

func (r *GormRepository) DB() *gorm.DB {
	return r.DBWithPreloads(nil)
}

func (r *GormRepository) GetAll(target interface{}, preloads ...string) error {
	r.logger.Debugf("Executing GetAll on %T", target)

	res := r.DBWithPreloads(preloads).
		Unscoped().
		Find(target)

	return r.HandleError(res)
}

func (r *GormRepository) GetWhere(target interface{}, condition string, preloads ...string) error {
	r.logger.Debugf("Executing GetWhere on %T with %v ", target, condition)

	res := r.DBWithPreloads(preloads).
		Where(condition).
		Find(target)

	return r.HandleError(res)
}

func (r *GormRepository) GetByField(target interface{}, field string, value interface{}, preloads ...string) error {
	r.logger.Debugf("Executing GetByField on %T with %v = %v", target, field, value)

	res := r.DBWithPreloads(preloads).
		Where(fmt.Sprintf("%v = ?", field), value).
		Find(target)

	return r.HandleError(res)
}

func (r *GormRepository) GetByFields(target interface{}, filters map[string]interface{}, preloads ...string) error {
	r.logger.Debugf("Executing GetByField on %T with filters = %+v", target, filters)

	db := r.DBWithPreloads(preloads)
	for field, value := range filters {
		db = db.Where(fmt.Sprintf("%v = ?", field), value)
	}

	res := db.Find(target)
	return r.HandleError(res)
}

func (r *GormRepository) GetOneByField(target interface{}, field string, value interface{}, preloads ...string) error {
	r.logger.Debugf("Executing GetOneByField on %T with %v = %v", target, field, value)

	res := r.DBWithPreloads(preloads).
		Where(fmt.Sprintf("%v = ?", field), value).
		First(target)

	return r.HandleOneError(res)
}

func (r *GormRepository) GetOneByFields(target interface{}, filters map[string]interface{}, preloads ...string) error {
	r.logger.Debugf("Executing FindOneByField on %T with filters = %+v", target, filters)

	db := r.DBWithPreloads(preloads)
	for field, value := range filters {
		db = db.Where(fmt.Sprintf("%v = ?", field), value)
	}

	res := db.First(target)
	return r.HandleOneError(res)
}

func (r *GormRepository) GetOneByID(target interface{}, id string, preloads ...string) error {
	r.logger.Debugf("Executing GetOneByID on %T with ID %v", target, id)

	res := r.DBWithPreloads(preloads).
		Where("id = ?", id).
		First(target)

	return r.HandleOneError(res)
}

func (r *GormRepository) Create(target interface{}) error {
	r.logger.Debugf("Executing Create on %T", target)

	res := r.db.Create(target)
	return r.HandleError(res)
}

func (r *GormRepository) CreateTx(target interface{}, tx *gorm.DB) error {
	r.logger.Debugf("Executing Create on %T", target)

	res := tx.Create(target)
	return r.HandleError(res)
}

func (r *GormRepository) Save(target interface{}) error {
	r.logger.Debugf("Executing Save on %T", target)

	res := r.db.Save(target)
	return r.HandleError(res)
}

func (r *GormRepository) SaveTx(target interface{}, tx *gorm.DB) error {
	r.logger.Debugf("Executing Save on %T", target)

	res := tx.Save(target)
	return r.HandleError(res)
}

func (r *GormRepository) Delete(target interface{}) error {
	r.logger.Debugf("Executing Delete on %T", target)

	res := r.db.Delete(target)
	return r.HandleError(res)
}

func (r *GormRepository) DeleteTx(target interface{}, tx *gorm.DB) error {
	r.logger.Debugf("Executing Delete on %T", target)

	res := tx.Delete(target)
	return r.HandleError(res)
}

func (r *GormRepository) HandleError(res *gorm.DB) error {
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		err := fmt.Errorf("Error: %w", res.Error)
		r.logger.Error(err)
		return err
	}

	return nil
}

func (r *GormRepository) HandleOneError(res *gorm.DB) error {
	if err := r.HandleError(res); err != nil {
		return err
	}

	if res.RowsAffected != 1 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *GormRepository) DBWithPreloads(preloads []string) *gorm.DB {
	dbConn := r.db

	for _, join := range r.defaultJoins {
		dbConn = dbConn.Joins(join)
	}

	for _, preload := range preloads {
		dbConn = dbConn.Preload(preload)
	}

	return dbConn
}

// FindPage page finding records
func (r GormRepository) FindPage(page, perPage uint, query *gorm.DB, out interface{}) (Page, error) {
	if perPage > 1000 {
		// cap at 10000 records per call
		perPage = 1000
	}

	if page == 0 {
		// 1 based page index
		page = 1
	}

	results := Page{
		PerPage: perPage,
		Page:    page,
	}

	session := query.Session(&gorm.Session{
		DryRun: true,
		Logger: logger.New(nil, logger.Config{
			LogLevel: logger.Silent,
		}),
	})

	var total int64
	countSqlStr, countArgs := r.buildCountSql(session)
	err := r.db.Raw(countSqlStr, countArgs...).Count(&total).Error
	if err != nil {
		return results, err
	}

	queryStmt := session.
		Offset(int(page - 1*perPage)).
		Limit(int(perPage)).
		Find(nil).Statement

	// for Postgresql, db.Raw() expects '?' as placeholders so replace '$1' placeholders with '?'
	sqlStr := r.replaceNumericPlaceholders(queryStmt.SQL.String())
	vars := queryStmt.Vars

	err = r.db.Raw(sqlStr, vars...).Scan(out).Error
	if err != nil {
		return results, err
	}

	results.Total = uint(total)
	results.Pages = r.calcPageCount(uint64(results.PerPage), uint64(results.Total))

	return results, nil
}

// Count count total number of records for the given query
func (r GormRepository) Count(model, query interface{}, args ...interface{}) (count int64, err error) {
	err = r.db.Model(model).Where(query, args...).Count(&count).Error

	return
}

// AutoMigrateOrWarn create tables for the given models or print a warning message if there's an error
func (r GormRepository) AutoMigrateOrWarn(models ...interface{}) {
	if err := r.db.AutoMigrate(models...); err != nil {
		log.Println("warning:", err.Error())
	}
}

func (r GormRepository) calcPageCount(perPage, total uint64) uint {
	if perPage == 0 || total == 0 {
		return 0
	}
	return uint(math.Ceil(float64(total) / float64(perPage)))
}

func (r GormRepository) buildCountSql(db *gorm.DB) (countSql string, vars []interface{}) {
	if orderByClause, ok := db.Statement.Clauses["ORDER BY"]; ok {
		if _, ok := db.Statement.Clauses["GROUP BY"]; !ok {
			delete(db.Statement.Clauses, "ORDER BY")
			defer func() {
				db.Statement.Clauses["ORDER BY"] = orderByClause
			}()
		}
	}
	var count int64
	countStmt := db.Count(&count).Statement

	// for Postgresql, db.Raw() expects '?' as placeholders so replace '$1' placeholders with '?'
	countSql = r.replaceNumericPlaceholders(countStmt.SQL.String())
	vars = countStmt.Vars

	return
}

func (r GormRepository) replaceNumericPlaceholders(sqlStr string) string {
	var numericPlaceholder = regexp.MustCompile("\\$(\\d+)")

	return numericPlaceholder.ReplaceAllString(sqlStr, "?")
}
