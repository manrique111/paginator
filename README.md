# Paginador en Gorm

28-Junio-2024
Hola, ¿qué tal? Adjunto esta investigación para la comunidad. Actualmente, no tengo mucho tiempo trabajando con este lenguaje. Todavía me considero un desarrollador Junior y no sé si esto puede ser de ayuda o si estoy reinventando la rueda.

Paginar en GORM es un poco complicado. Busqué en internet y encontré algunas soluciones, pero ninguna se adecuaba a lo que se requería.

Pido disculpas por la documentación en español, pero aún no domino bien el inglés. ¡Saludos!

## Instalación

    go get github.com/manrique111/paginator

    import (
	"github.com/tu-usuario/gorm-plugin-example/paginator"
	)


## Forma de usarlo
Para esta parte anexare un ejemolo de como se debera de usar

    func GetFuncion(context *gin.Context) {
	status := context.Query("status")
    page, _ := strconv.Atoi(context.Query("page"))  
	pageSize, _ := strconv.Atoi(context.Query("page_size"))  
	filter := context.Query("filter")

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    var users []User
	
	// inicializar el constructor
	pages := shared.PagesInit(page, pageSize)
	
	// Amar las condiciones
	var conditions = paginator.QueryParams{  
	    Where: map[string]interface{}{  
	       "supplier_id": suplierID,  
	    },  
	    Joins: []string{  
	       //"INNER JOIN other_table ON other_table.id = main_table.other_table_id",  
	  },  
	    Associations: []string{"Address", "celphones.Compania",}, // Especifica las asociaciones a .Preload  
	  Order: []shared.Order{  
	       {Column: "ID", Dir: "Desc"},  
	       //{Column: orderColumn2, Dir: orderDir2},  
	  },  
	}
    

    //en caso de modificar tu condicion where o anexar nuevas condiciones
    if status != "all" {
		conditions.Where = map[string]interface{}{
			"Status": status,
		}
	}

	//sentencias or
	if filter != "" {
		conditions.OrWhere = []string{
			fmt.Sprintf("ID LIKE '%%%s%%'", filter),
			fmt.Sprintf("name LIKE '%%%s%%'", filter),
		}
	}


	if err := pages.SetRecord(db, &user, conditions); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.AbortWithStatusJSON(http.StatusNotFound, gin.H{"msg": "No existen registros"})
		} else {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		}
	}
	fmt.Printf("Total Count: %d, Total Pages: %d\n", pages.TotalCount, pages.TotalPages)
	context.JSON(http.StatusOK, pages)
}





## **La estructura que devolvera sera**

    type Pages struct {  
	      Page       int `json:"page"`  
		  PageSize   int `json:"page_size"`  
		  TotalPages int64 `json:"total_pages"`  
		  TotalCount int64 `json:"total_count"`  
		  Data       []interface{} `json:"data"`  
	}

## Estructuracion de las clasulas where

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

