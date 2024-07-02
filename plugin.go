package gormplugin

import (
	"fmt"
	"gorm.io/gorm"
	"reflect"
)

type Paginator struct {
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	DebugSql   bool          `json:"-"`
	TotalPages int64         `json:"total_pages"`
	TotalCount int64         `json:"total_count"`
	Data       []interface{} `json:"data"`
}

// Constructor
func NewPaginator(page int, pageSize int) *Paginator {
	// Obtener los parámetros de paginación directamente
	if page <= 0 {
		page = 1
	}
	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}
	// Obtener los parámetros de ordenamiento directamente
	pages := &Paginator{
		Page:     page,
		PageSize: pageSize,
	}
	// Devolver los parámetros de paginación y ordenamiento
	return pages
}

func (p *Paginator) SetDebug(debugSql bool) {
	p.DebugSql = debugSql
}

func (p *Paginator) SetRecord(db *gorm.DB, modelo interface{}) error {
	var totalCount int64

	query := db.Model(modelo)

	// Debugging: Print the SQL query before counting
	if p.DebugSql {
		sqlQueryCount := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
			return query.Session(&gorm.Session{DryRun: true}).Count(&totalCount)
		})
		fmt.Printf("SQL before count: %s\n", sqlQueryCount)
	}

	// Obtener el total de registros que cumplen las condiciones
	if err := query.Count(&totalCount).Find(modelo).Error; err != nil {
		return err
	}

	// Debugging: Print the SQL query after counting
	if p.DebugSql {
		sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
			return query.Session(&gorm.Session{DryRun: true}).Offset(0).Limit(p.PageSize).Find(modelo)
		})
		fmt.Printf("SQL after count: %s\n", sql)
	}

	offset := (p.Page - 1) * p.PageSize
	if err := query.Offset(offset).Limit(p.PageSize).Find(modelo).Error; err != nil {
		return err
	}

	// Verificar que modelo es un puntero a un slice
	val := reflect.ValueOf(modelo)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Slice {
		panic("modelo debe ser un puntero a un slice")
	}

	// Convertir modelo a un slice de interface{}
	modeloVal := val.Elem()
	data := make([]interface{}, modeloVal.Len())
	for i := 0; i < modeloVal.Len(); i++ {
		data[i] = modeloVal.Index(i).Interface()
	}

	p.TotalCount = totalCount
	p.TotalPages = (totalCount + int64(p.PageSize) - 1) / int64(p.PageSize)
	p.Data = data

	return nil
}
