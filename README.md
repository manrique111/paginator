# Paginador en Gorm

28-Junio-2024
Hola, ¿qué tal? Adjunto esta investigación para la comunidad. Actualmente, no tengo mucho tiempo trabajando con este lenguaje. Todavía me considero un desarrollador Junior y no sé si esto puede ser de ayuda o si estoy reinventando la rueda.

Paginar en GORM es un poco complicado. Busqué en internet y encontré algunas soluciones, pero ninguna se adecuaba a lo que se requería.

Pido disculpas por la documentación en español, pero aún no domino bien el inglés. ¡Saludos!

Actualizacion
Debido a los usos actualizo la funcionalidad de pages ahora los filtros y las condiciones sera por fuera
para aprovechar la funcionalida de go

la funcionalida tendra por objetivo regresar tu objeto page formado con la data para que lo pueda usar en tu front

## Instalación

    go get github.com/manrique111/paginator

    import (
	"github.com/manrique111/paginator"
	)


## Forma de usarlo
    //	@Router			/GetMyFunction/order?supplierID=&status=&page=&pageSize=&filter=
    func GetMyFunction(context *gin.Context) {
        Para esta parte anexare un ejemolo de como se debera de usar
        db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
        if err != nil {
        panic("failed to connect database")
        }
        var purchaseOrders []models.PurchaseOrders
        suplierID := context.Query("suplierID")
        status := context.Query("status")
        page, _ := strconv.Atoi(context.Query("page"))
        pageSize, _ := strconv.Atoi(context.Query("page_size"))
        filter := context.Query("filter")
    
        paginator := paginator.NewPaginator(page, pageSize)
        paginator.SetDebug(true)    // anexar si se requiere que imprima los sql por consola opcional
        // Añadir condiciones de búsqueda usando parámetros nativos de GORM
        query := db.Model(&models.PurchaseOrders{})
        if suplierID != "" {
            query = query.Where("supplier_id = ?", suplierID)
        }
        if status != "TODOS" && status != "" {
            query = query.Where("status = ?", status)
        }
        if filter != "" {
            if shared.IsValidDate(filter) {
                query = query.Or("DATE_FORMAT(CreatedAt, '%Y-%m-%d') LIKE ?", "%"+filter+"%")
            } else {
                query = query.Where("cve_doc_sae LIKE ?", "%"+filter+"%")
                query = query.Or("ID LIKE ?", "%"+filter+"%")
                query = query.Or("quantity LIKE ?", "%"+filter+"%")
                query = query.Or("quantity_received LIKE ?", "%"+filter+"%")
            }
        }
        query = query.Order("ID DESC")
    
        if err := paginator.SetRecord(query, &purchaseOrders); err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                context.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Spanish": "No existen registros"})
            } else {
                context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Spanish": err.Error()})
            }
        }
    
        context.JSON(http.StatusOK, paginator)
        fmt.Printf("Total Count: %d, Total Pages: %d\n", pages.TotalCount, pages.TotalPages)
        context.JSON(http.StatusOK, paginator)
    }





## **La estructura que devolvera sera**

    type Pages struct {  
      Page       int `json:"page"`  
      PageSize   int `json:"page_size"`  
      TotalPages int64 `json:"total_pages"`  
      TotalCount int64 `json:"total_count"`  
      Data       []interface{} `json:"data"`  
	}




