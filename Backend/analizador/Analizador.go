package analizador

import (
	"miapp/comandos"
	"miapp/singleton"
	"strconv"
	"strings"
)

type Analizador struct {
	Entrada   string
	disco     *comandos.Disco
	singleton *singleton.Singleton
	particion *comandos.Particion
	rep       *comandos.Rep
	MountList *comandos.MountList
	Montar    *comandos.Montar
}

/*
*Constructor
 */
func NewAnalizador(entrada string, mountList *comandos.MountList) *Analizador {
	entrada = removeSpace(entrada)
	return &Analizador{
		Entrada:   entrada,
		disco:     comandos.NewDisco(),
		singleton: singleton.GetInstance(),
		particion: comandos.NewParticion(),
		MountList: mountList,
		rep:       comandos.NewRep(),
		Montar:    comandos.NewMontar()}
}

// Metodos principales del Analizador
func (a *Analizador) toLower(cadena string) string {
	cadMinus := ""
	longitud := len(cadena)
	i := 0
	for i < longitud {
		cadMinus += strings.ToLower(string(cadena[i]))
		i++
	}
	return cadMinus
}
func removeSpace(entrada string) string {
	entrada = strings.ReplaceAll(entrada, "\t", " ")
	entrada = strings.ReplaceAll(entrada, "\r", " ")
	entrada = strings.ReplaceAll(entrada, "\n", " ")
	return entrada
}
func (a *Analizador) AnalizarEntrada() {
	if len(a.Entrada) > 0 {
		entradaMinus := a.toLower(a.Entrada)
		if strings.HasPrefix(entradaMinus, " ") {
			i := 1
			// Consumimos de espacios
			for i < len(entradaMinus) && entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
				i++
			}
			entradaMinus = entradaMinus[i:]
			a.Entrada = a.Entrada[i:]
		} else if strings.HasPrefix(entradaMinus, "#") { //COMENTARIO
			// Es un comentario entonces solo se ignora todo.
			return
		} else if strings.HasPrefix(entradaMinus, "mkdisk") { //MKDISK
			i := 6
			// Consumimos de espacios
			for i < len(entradaMinus) && entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
				i++
			}
			entradaMinus = entradaMinus[i:]
			a.Entrada = a.Entrada[i:]
			//parte de verificar parametros
			for len(a.Entrada) > 0 {
				if strings.HasPrefix(entradaMinus, ">size") {
					i := strings.Index(entradaMinus, "=") + 1
					for i < len(entradaMinus) && entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
						i++
					}
					entradaMinus = entradaMinus[i:]
					a.Entrada = a.Entrada[i:]

					i = strings.Index(entradaMinus, " ")
					if i == -1 { // CAMBIO 1 SIZE
						i = len(entradaMinus)
					}
					s, err := strconv.Atoi(entradaMinus[:i])
					if err != nil {
						// manejar el error de conversión
						a.singleton.AddSalidaConsola("ERROR EN EL COMANDO: " + a.Entrada + ", ASIGNACION NUMERO FLOTANTE A PARAMETRO SIZE\n")
						continue
					}
					a.disco.S = s
					for i < len(entradaMinus) && entradaMinus[i] == ' ' && len(entradaMinus) > 0 { //CAMBIO 2 SIZE
						i++
					}
					entradaMinus = entradaMinus[i:]
					a.Entrada = a.Entrada[i:]
				} else if strings.HasPrefix(entradaMinus, ">fit") {
					i := strings.Index(entradaMinus, "=") + 1
					for i < len(entradaMinus) && entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
						i++
					}
					entradaMinus = entradaMinus[i:]
					a.Entrada = a.Entrada[i:]

					i = strings.Index(entradaMinus, " ")
					if i == -1 { // CAMBIO 1 FIT
						i = len(entradaMinus)
					}
					f := entradaMinus[:i]
					a.disco.F = f
					for i < len(entradaMinus) && entradaMinus[i] == ' ' && len(entradaMinus) > 0 { // CAMBIO 2 FIT
						i++
					}
					entradaMinus = entradaMinus[i:]
					a.Entrada = a.Entrada[i:]
				} else if strings.HasPrefix(entradaMinus, ">unit") {
					i := strings.Index(entradaMinus, "=") + 1
					for i < len(entradaMinus) && entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
						i++
					}
					entradaMinus = entradaMinus[i:]
					a.Entrada = a.Entrada[i:]

					i = strings.Index(entradaMinus, " ")
					if i == -1 { // CAMBIO 1 UNIT
						i = len(entradaMinus)
					}
					u := entradaMinus[:i]
					a.disco.U = u
					for i < len(entradaMinus) && entradaMinus[i] == ' ' && len(entradaMinus) > 0 { // CAMBIO 2 UNIT
						i++
					}
					entradaMinus = entradaMinus[i:]
					a.Entrada = a.Entrada[i:]
				} else if strings.HasPrefix(entradaMinus, ">path") { // PARAMETRO PATH
					i = strings.Index(entradaMinus, "=") + 1
					for i < len(entradaMinus) && entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
						i++
					}
					entradaMinus = entradaMinus[i:]
					a.Entrada = a.Entrada[i:]

					if entradaMinus[0] == '"' { // PATH CON COMILLAS
						entradaMinus = entradaMinus[1:]
						a.Entrada = a.Entrada[1:]
						i = strings.Index(entradaMinus, "\"")
						if i == -1 { // CAMBIO 1
							// manejar el error de conversión
							a.singleton.AddSalidaConsola("ERROR EN EL COMANDO: " + a.Entrada + ", SE LE OLVIDO LA COMILLA DEL STRING" + "\n")
							continue
						}
						p := a.Entrada[:i]
						i += 1 // CAMBIO 2
						a.disco.P = p
						for i < len(entradaMinus) && entradaMinus[i] == ' ' && len(entradaMinus) > 0 { // CAMBIO 3
							i++
						}
						entradaMinus = entradaMinus[i:]
						a.Entrada = a.Entrada[i:]
					} else { // PATH SIN COMILLAS
						i = strings.Index(entradaMinus, " ")
						if i == -1 { // CAMBIO 1
							i = len(entradaMinus)
						}
						p := a.Entrada[:i]

						a.disco.P = p
						for i < len(entradaMinus) && entradaMinus[i] == ' ' && len(entradaMinus) > 0 { // CAMBIO 2
							i++
						}
						entradaMinus = entradaMinus[i:]
						a.Entrada = a.Entrada[i:]
					}
				} else if strings.HasPrefix(entradaMinus, "#") {
					// No se opera, ya que entro un comentario
					break
				} else {
					a.singleton.AddSalidaConsola("ERROR EN EL COMANDO: " + a.Entrada + "\n")
					return
				}
			}
			a.disco.Mkdisk()
		} else if strings.HasPrefix(entradaMinus, "rmdisk") { //RMDISK
			i := 6
			for i < len(entradaMinus) && entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
				i++
			}
			entradaMinus = entradaMinus[i:]
			a.Entrada = a.Entrada[i:]

			for len(a.Entrada) > 0 {
				if strings.HasPrefix(entradaMinus, ">path") {
					i = strings.Index(entradaMinus, "=") + 1
					for i < len(entradaMinus) && entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
						i++
					}
					entradaMinus = entradaMinus[i:]
					a.Entrada = a.Entrada[i:]

					if entradaMinus[0] == '"' {
						entradaMinus = entradaMinus[1:]
						a.Entrada = a.Entrada[1:]
						i = strings.Index(entradaMinus, "\"")
						if i == -1 {
							// manejar el error de conversión
							a.singleton.AddSalidaConsola("ERROR EN EL COMANDO: " + a.Entrada + ", SE LE OLVIDO LA COMILLA DEL STRING" + "\n")
							continue
						}
						p := a.Entrada[:i]
						i += 1
						a.disco.P = p
						for i < len(entradaMinus) && entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
							i++
						}
						entradaMinus = entradaMinus[i:]
						a.Entrada = a.Entrada[i:]
					} else {
						i = strings.Index(entradaMinus, " ")
						if i == -1 { // CAMBIO 1
							i = len(entradaMinus)
						}
						p := a.Entrada[:i]
						a.disco.P = p
						for i < len(entradaMinus) && entradaMinus[i] == ' ' && len(entradaMinus) > 0 { // CAMBIO 2
							i++
						}
						entradaMinus = entradaMinus[i:]
						a.Entrada = a.Entrada[i:]
					}
				} else if strings.HasPrefix(entradaMinus, "#") {
					// es un comentario entonces se ignora todo
					break
				} else {
					a.singleton.AddSalidaConsola("ERROR EN EL COMANDO: " + entradaMinus + "\n")
					return
				}
			}
			a.disco.Rmdisk()
		} else if strings.HasPrefix(entradaMinus, "fdisk") { //FDISK
			i := 5
			for i < len(entradaMinus) && entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
				i++
			}
			entradaMinus = entradaMinus[i:]
			a.Entrada = a.Entrada[i:]
			// Verificar qué parámetros trae el comando
			for len(a.Entrada) > 0 {
				if strings.HasPrefix(entradaMinus, ">size") {
					i = strings.Index(entradaMinus, "=") + 1
					for i < len(entradaMinus) && entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
						i++
					}
					entradaMinus = entradaMinus[i:]
					a.Entrada = a.Entrada[i:]

					i = strings.IndexByte(entradaMinus, ' ')
					if i == -1 { // CAMBIO 1 SIZE
						i = len(entradaMinus)
					}
					s, _ := strconv.Atoi(entradaMinus[:i])
					a.particion.S = s
					if a.particion.Flag == 'n' {
						a.particion.Flag = 's'
					}
					for i < len(entradaMinus) && entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
						i++
					}
					entradaMinus = entradaMinus[i:]
					a.Entrada = a.Entrada[i:]

				} else if strings.HasPrefix(entradaMinus, ">unit") {
					i = strings.Index(entradaMinus, "=") + 1
					for i < len(entradaMinus) && entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
						i++
					}
					entradaMinus = entradaMinus[i:]
					a.Entrada = a.Entrada[i:]

					i = strings.IndexByte(entradaMinus, ' ')
					if i == -1 { // CAMBIO 1 SIZE
						i = len(entradaMinus)
					}
					u := entradaMinus[:i]
					a.particion.U = u[0]
					for i < len(entradaMinus) && entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
						i++
					}
					entradaMinus = entradaMinus[i:]
					a.Entrada = a.Entrada[i:]

				} else if strings.HasPrefix(entradaMinus, ">path") {
					i = strings.Index(entradaMinus, "=") + 1
					for i < len(entradaMinus) && entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
						i++
					}
					entradaMinus = entradaMinus[i:]
					a.Entrada = a.Entrada[i:]

					if entradaMinus[0] == '"' {
						entradaMinus = entradaMinus[1:]
						a.Entrada = a.Entrada[1:]
						i = strings.Index(entradaMinus, "\"")
						if i == -1 { // CAMBIO 1
							// manejar el error de conversión
							a.singleton.AddSalidaConsola("ERROR EN EL COMANDO: " + a.Entrada + ", SE LE OLVIDO LA COMILLA DEL STRING" + "\n")
							continue
						}
						p := a.Entrada[:i]
						i += 1 // CAMBIO 2
						a.particion.P = p
						for i < len(entradaMinus) && entradaMinus[i] == ' ' && len(entradaMinus) > 0 { // CAMBIO 3
							i++
						}
						entradaMinus = entradaMinus[i:]
						a.Entrada = a.Entrada[i:]
					} else {
						i = strings.IndexByte(entradaMinus, ' ')
						if i == -1 { // CAMBIO 1
							i = len(entradaMinus)
						}
						p := a.Entrada[:i]
						a.particion.P = p
						for i < len(entradaMinus) && entradaMinus[i] == ' ' && len(entradaMinus) > 0 { // CAMBIO 2
							i++
						}
						entradaMinus = entradaMinus[i:]
						a.Entrada = a.Entrada[i:]
					}

				} else if strings.HasPrefix(entradaMinus, ">type") {
					i := strings.Index(entradaMinus, "=") + 1
					for i < len(entradaMinus) && entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
						i++
					}
					entradaMinus = entradaMinus[i:]
					a.Entrada = a.Entrada[i:]

					i = strings.IndexByte(entradaMinus, ' ')
					t := entradaMinus[:i]
					a.particion.T = t[0]
					for i < len(entradaMinus) && entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
						i++
					}
					entradaMinus = entradaMinus[i:]
					a.Entrada = a.Entrada[i:]
				} else if strings.HasPrefix(entradaMinus, ">fit") {
					i = strings.Index(entradaMinus, "=") + 1
					for i < len(entradaMinus) && entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
						i++
					}
					entradaMinus = entradaMinus[i:]
					a.Entrada = a.Entrada[i:]

					i = strings.IndexByte(entradaMinus, ' ')
					f := entradaMinus[:i]
					if !(f == "bf" || f == "ff" || f == "wf") {
						a.singleton.AddSalidaConsola("OPCION INVALIDA PARA -f\n")
						return
					}
					a.particion.F = f[0]
					for i < len(entradaMinus) && entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
						i++
					}
					entradaMinus = entradaMinus[i:]
					a.Entrada = a.Entrada[i:]

				} else if strings.HasPrefix(entradaMinus, ">delete") {
					i = strings.Index(entradaMinus, "=") + 1
					for i < len(entradaMinus) && entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
						i++
					}
					entradaMinus = entradaMinus[i:]
					a.Entrada = a.Entrada[i:]

					i = strings.IndexByte(entradaMinus, ' ')
					d := entradaMinus[:i]
					if d != "full" {
						a.singleton.AddSalidaConsola("OPCION " + d + " INVALIDA PARA -delete\n")
						return
					}
					a.particion.D = d
					for i < len(entradaMinus) && entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
						i++
					}
					entradaMinus = entradaMinus[i:]
					a.Entrada = a.Entrada[i:]
				} else if strings.HasPrefix(entradaMinus, ">name") {
					i := strings.Index(entradaMinus, "=") + 1
					for i < len(entradaMinus) && entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
						i++
					}
					entradaMinus = entradaMinus[i:]
					a.Entrada = a.Entrada[i:]

					if entradaMinus[0] == '"' {
						entradaMinus = entradaMinus[1:]
						a.Entrada = a.Entrada[1:]
						i := strings.Index(entradaMinus, "\"")
						n := a.Entrada[:i]
						i += 2
						a.particion.Name = n
						for i < len(entradaMinus) && entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
							i++
						}
						entradaMinus = entradaMinus[i:]
						a.Entrada = a.Entrada[i:]
					} else {
						i := strings.Index(entradaMinus, " ")
						n := a.Entrada[:i]
						a.particion.Name = n
						for i < len(entradaMinus) && entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
							i++
						}
						entradaMinus = entradaMinus[i:]
						a.Entrada = a.Entrada[i:]
					}

				} else if strings.HasPrefix(entradaMinus, ">add") {
					i := strings.Index(entradaMinus, "=") + 1
					for i < len(entradaMinus) && entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
						i++
					}
					entradaMinus = entradaMinus[i:]
					a.Entrada = a.Entrada[i:]

					i = strings.Index(entradaMinus, " ")
					add, _ := strconv.Atoi(entradaMinus[:i])
					a.particion.Add = add
					if a.particion.Flag == 'n' {
						a.particion.Flag = 'a'
					}
					for i < len(entradaMinus) && entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
						i++
					}
					entradaMinus = entradaMinus[i:]
					a.Entrada = a.Entrada[i:]
				} else if strings.HasPrefix(entradaMinus, "#") {
					// es un comentario entonces no se hace nada
					break
				} else {
					a.singleton.AddSalidaConsola("ERROR EN EL COMANDO DE ENTRADA: " + entradaMinus + "\n")
					return
				}
			}
			a.particion.Fdisk()
		} else if strings.HasPrefix(entradaMinus, "rep") {
			entradaMinus = strings.TrimSpace(entradaMinus[3:])
			a.Entrada = strings.TrimSpace(a.Entrada[3:])

			for len(entradaMinus) > 0 {
				if strings.HasPrefix(entradaMinus, ">path") {
					i := strings.Index(entradaMinus, "=") + 1
					for entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
						i++
					}
					entradaMinus = entradaMinus[i:]
					a.Entrada = a.Entrada[i:]

					if entradaMinus[0] == '"' {
						entradaMinus = entradaMinus[1:]
						a.Entrada = a.Entrada[1:]
						i = strings.Index(entradaMinus, "\"")
						p := a.Entrada[:i]
						i += 2
						a.rep.Path = p
						for entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
							i++
						}
						entradaMinus = entradaMinus[i:]
						a.Entrada = a.Entrada[i:]
					} else {
						i := strings.Index(entradaMinus, " ")
						p := a.Entrada[:i]
						a.rep.Path = p
						for entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
							i++
						}
						entradaMinus = entradaMinus[i:]
						a.Entrada = a.Entrada[i:]
					}

				} else if strings.HasPrefix(entradaMinus, ">ruta") {
					i := strings.Index(entradaMinus, "=") + 1
					for entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
						i++
					}
					entradaMinus = entradaMinus[i:]
					a.Entrada = a.Entrada[i:]

					if entradaMinus[0] == '"' {
						entradaMinus = entradaMinus[1:]
						a.Entrada = a.Entrada[1:]
						i := strings.Index(entradaMinus, "\"")
						a.rep.Ruta = a.Entrada[:i]
						i += 2
						for entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
							i++
						}
						entradaMinus = entradaMinus[i:]
						a.Entrada = a.Entrada[i:]
					} else {
						i := strings.Index(entradaMinus, " ")
						a.rep.Ruta = a.Entrada[:i]
						for entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
							i++
						}
						entradaMinus = entradaMinus[i:]
						a.Entrada = a.Entrada[i:]
					}

				} else if strings.HasPrefix(entradaMinus, ">name") {
					i := strings.Index(entradaMinus, "=") + 1
					for entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
						i++
					}
					entradaMinus = entradaMinus[i:]
					a.Entrada = a.Entrada[i:]

					i = strings.Index(entradaMinus, " ")
					n := entradaMinus[:i]
					a.rep.Name = n
					for entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
						i++
					}
					entradaMinus = entradaMinus[i:]
					a.Entrada = a.Entrada[i:]
				} else if strings.HasPrefix(entradaMinus, ">id") {
					i := strings.Index(entradaMinus, "=") + 1
					for entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
						i++
					}
					entradaMinus = entradaMinus[i:]
					a.Entrada = a.Entrada[i:]

					if entradaMinus[0] == '"' {
						entradaMinus = entradaMinus[1:]
						a.Entrada = a.Entrada[1:]
						i := strings.Index(entradaMinus, "\"")
						id := a.Entrada[:i]
						i += 2
						a.rep.Id = id
						for entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
							i++
						}
						entradaMinus = entradaMinus[i:]
						a.Entrada = a.Entrada[i:]
					} else {
						i := strings.Index(entradaMinus, " ")
						id := a.Entrada[:i]
						a.rep.Id = id
						for entradaMinus[i] == ' ' && len(entradaMinus) > 0 {
							i++
						}
						entradaMinus = entradaMinus[i:]
						a.Entrada = a.Entrada[i:]
					}

				} else if strings.HasPrefix(entradaMinus, "#") {
					// No se opera, ya que entro un comentario
					break
				} else {
					a.singleton.AddSalidaConsola("ERROR EN EL COMANDO: " + entradaMinus)
					return
				}
			}
			a.rep.MountList = a.MountList
			a.rep.Generate()
			a.MountList = a.rep.MountList
		} else if strings.HasPrefix(entradaMinus, "mount") {
			i := 5
			for i < len(entradaMinus) && i < len(entradaMinus) && entradaMinus[i] == ' ' {
				i++
			}
			entradaMinus = entradaMinus[i:]
			a.Entrada = a.Entrada[i:]

			for len(entradaMinus) > 0 {
				if strings.HasPrefix(entradaMinus, ">path") {
					i = strings.Index(entradaMinus, "=") + 1
					for i < len(entradaMinus) && entradaMinus[i] == ' ' {
						i++
					}
					entradaMinus = entradaMinus[i:]
					a.Entrada = a.Entrada[i:]

					if entradaMinus[0] == '"' {
						entradaMinus = entradaMinus[1:]
						a.Entrada = a.Entrada[1:]
						i = strings.Index(entradaMinus, "\"")
						p := a.Entrada[:i]
						i += 2
						a.Montar.P = p
						for i < len(entradaMinus) && entradaMinus[i] == ' ' {
							i++
						}
						entradaMinus = entradaMinus[i:]
						a.Entrada = a.Entrada[i:]
					} else {
						i = strings.IndexByte(entradaMinus, ' ')
						p := a.Entrada[:i]
						a.Montar.P = p
						for i < len(entradaMinus) && entradaMinus[i] == ' ' {
							i++
						}
						entradaMinus = entradaMinus[i:]
						a.Entrada = a.Entrada[i:]
					}
				} else if strings.HasPrefix(entradaMinus, ">name") {
					i = strings.Index(entradaMinus, "=") + 1
					for i < len(entradaMinus) && entradaMinus[i] == ' ' {
						i++
					}
					entradaMinus = entradaMinus[i:]
					a.Entrada = a.Entrada[i:]

					if entradaMinus[0] == '"' {
						entradaMinus = entradaMinus[1:]
						a.Entrada = a.Entrada[1:]
						i = strings.Index(entradaMinus, "\"")
						n := a.Entrada[:i]
						i += 2
						a.Montar.Name = n
						for i < len(entradaMinus) && entradaMinus[i] == ' ' {
							i++
						}
						entradaMinus = entradaMinus[i:]
						a.Entrada = a.Entrada[i:]
					} else {
						i = strings.IndexByte(entradaMinus, ' ')
						n := a.Entrada[:i]
						a.Montar.Name = n
						for i < len(entradaMinus) && entradaMinus[i] == ' ' {
							i++
						}
						entradaMinus = entradaMinus[i:]
						a.Entrada = a.Entrada[i:]
					}
				} else if strings.HasPrefix(entradaMinus, "#") {
					break
				} else {
					a.singleton.AddSalidaConsola("ERROR EN: " + entradaMinus + "\n")
					return
				}
			}
			a.Montar.MountList = a.MountList
			a.Montar.Mount()
			a.MountList = a.Montar.MountList
		} else {
			a.singleton.AddSalidaConsola(">> COMANDO INVALIDO ASEGURESE DE ESCRIBIR BIEN TODO\n")
		}
	}
	return
}
