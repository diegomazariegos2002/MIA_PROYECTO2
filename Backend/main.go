package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"miapp/analizador"
	"miapp/comandos"
	"miapp/singleton"
	"net/http"
	"os"
	"os/signal"
	"strings"
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
	mountList := comandos.NewMountList()
	router := gin.Default()

	// Agregar middleware CORS
	config := cors.DefaultConfig()      // la configuración por default permite todo
	config.AllowOrigins = []string{"*"} // Cuidado, esto implica un riesgo de seguridad
	router.Use(cors.New(config))

	// Definiendo rutas del array routes
	routes := []Route{
		{"/", indexHandler, "GET"},
		{"/compilar",
			func(c *gin.Context) {
				compilar(c, mountList)
			}, "POST"},
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
		"message": "Server corriendo",
	})
}

func compilar(c *gin.Context, mountlist *comandos.MountList) {
	var data map[string]interface{}           // un mapa para guardar los datos del cuerpo
	if err := c.BindJSON(&data); err != nil { // intentar leer el cuerpo como JSON
		c.JSON(400, gin.H{"error": err.Error()}) // si hay un error, enviar respuesta 400
		return
	}
	var entryConsole string
	var singleton = singleton.GetInstance()
	singleton.ResetSalidaConsola()
	if s, ok := data["entrada"].(string); ok { // verificar si el campo entrada existe
		entryConsole = s
		reader := bufio.NewReader(strings.NewReader(entryConsole))
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					singleton.AddSalidaConsola("FIN DE LA ENTRADA\n")
					break
				}
				singleton.AddSalidaConsola("Error al leer entrada de usuario: " + err.Error() + "\n")
				continue
			}
			line = strings.TrimSpace(line) // eliminar /t/n/_ de los bordes

			// Aquí va la parte del analyzer
			var analyzer = analizador.NewAnalizador(line+" ", comandos.NewMountList()) // agrego el " " de ultimo para evitar errores
			singleton.AddSalidaConsola("//" + analyzer.Entrada + "\n")
			analyzer.MountList = mountlist
			analyzer.AnalizarEntrada()
			mountlist = analyzer.MountList

			// Verificar si se llegó al final de la entrada
			if _, err := reader.Peek(1); err != nil {
				if err == io.EOF {
					singleton.AddSalidaConsola("=======FIN DE LA ENTRADA=======\n")
					break
				}
				singleton.AddSalidaConsola("Error al leer entrada de usuario: " + err.Error() + "\n")
				continue
			}
		}
		c.JSON(200, gin.H{"message": "Datos recibidos", "salida": singleton.SalidaConsola()}) // si todo está bien, enviar respuesta 200 con los datos
		return
	} else {
		c.JSON(400, gin.H{"error": "error con la entrada"}) // si hay un error, enviar respuesta 400
		return
	}

	c.JSON(200, gin.H{"message": "Datos recibidos", "salida": "vacio"}) // si todo está bien, enviar respuesta 200 con los datos
	return
}
