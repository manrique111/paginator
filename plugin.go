package paginator

import (
	"fmt"
	"gorm.io/gorm"
	"reflect"
)

type Pages struct {
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalPages int64         `json:"total_pages"`
	TotalCount int64         `json:"total_count"`
	Data       []interface{} `json:"data"`
}

type Order struct {
	Column string
	Dir    string
}

type QueryParams struct {
	Where        map[string]interface{}
	OrWhere      []string
	Associations []string
	Joins        []string
	Order        []Order
}

// Constructor
func PagesInit(page int, pageSize int) *Pages {
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
	pages := &Pages{
		Page:     page,
		PageSize: pageSize,
	}
	// Devolver los parámetros de paginación y ordenamiento
	return pages
}

func (p *Pages) SetRecord(db *gorm.DB, modelo interface{}, params QueryParams) error {
	var totalCount int64

	query := db.Model(modelo)
	// Apply joins
	for _, join := range params.Joins {
		query = query.Joins(join)
	}
	//Apply wheres
	for key, value := range params.Where {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}
	// Apply OrWheres
	for _, condition := range params.OrWhere {
		query = query.Or(condition)
	}
	// Pre-cargar asociaciones si se especifican
	for _, association := range params.Associations {
		query = query.Preload(association)
	}

	// Aplicar orden si se especifica
	for _, order := range params.Order {
		orderDir := "ASC"
		if order.Dir != "" {
			orderDir = order.Dir
		}
		query = query.Order(fmt.Sprintf("%s %s", order.Column, orderDir))
	}

	// Obtener el total de registros que cumplen las condiciones
	query.Count(&totalCount)

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
