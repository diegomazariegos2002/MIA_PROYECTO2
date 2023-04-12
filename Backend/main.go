package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Route struct {
	Path    string
	Handler gin.HandlerFunc
	Method  string
}

func main() {
	router := gin.Default()

	// Agregar middleware CORS
	config := cors.DefaultConfig()      // la configuración por default permite todo
	config.AllowOrigins = []string{"*"} // Cuidado, esto implica un riesgo de seguridad
	router.Use(cors.New(config))

	// Definiendo rutas del array routes
	routes := []Route{
		{"/", indexHandler, "GET"},
	}

	// Crea cada ruta del array routes
	for _, r := range routes {
		switch r.Method {
		case "GET":
			router.GET(r.Path, r.Handler)
		case "POST":
			router.POST(r.Path, r.Handler)
		case "PUT":
			router.PUT(r.Path, r.Handler)
		case "DELETE":
			router.DELETE(r.Path, r.Handler)
		default:
			panic(fmt.Sprintf("unsupported HTTP method %s for route %s", r.Method, r.Path))
		}
	}

	// Crear instancia de http.Server con el router de Gin
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Iniciar el servidor HTTP en una goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
	}()

	// Esperar señales de interrupción para detener el servidor
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	// Detener el servidor y liberar el puerto
	if err := server.Shutdown(context.Background()); err != nil {
		fmt.Println(err)
	}
}

// definition of routes
func indexHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hola mundo",
	})
}
