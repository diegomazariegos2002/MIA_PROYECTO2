package comandos

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"miapp/singleton"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type Rep struct {
	Path                           string
	Name                           string
	Id                             string
	Ruta                           string
	Directorio                     string
	Extension                      string
	SuperBloqueGlobal              SuperBloque
	ContadorBloques_ReporteBloques int
	singleton                      *singleton.Singleton
	MountList                      *MountList
}

func NewRep() *Rep {
	return &Rep{
		Path:                           "",
		Name:                           "",
		Id:                             "",
		Ruta:                           "",
		Directorio:                     "",
		Extension:                      "",
		SuperBloqueGlobal:              SuperBloque{},
		ContadorBloques_ReporteBloques: 0,
		singleton:                      singleton.GetInstance(),
		MountList:                      NewMountList(),
	}
}
func (r *Rep) GetCarpetas(ruta string) string {
	dir, _ := filepath.Abs(filepath.Dir(ruta))
	dir += "/"
	return dir
}

func (r *Rep) GetExtensionFile(path string) string {
	i := strings.LastIndex(path, ".")
	extension := path[i+1:]
	return extension
}

func (r *Rep) Generate() {
	if r.Id != " " {
		if r.Path != " " {
			switch r.Name {
			case "mbr":
				//r.ejecutarReporte_mbr()
			case "disk":
				r.ejecutarReporte_disk()
			case "inode":
				//r.ejecutarReporte_inode()
			case "block":
				//r.ejecutarReporte_block()
			case "bm_inode":
				//r.ejecutarReporte_bm_inode()
			case "bm_block":
				//r.ejecutarReporte_bm_block()
			case "sb":
				//r.ejecutarReporte_sb()
			case "journaling":
				//r.ejecutarReporte_Journaling()
			case "file":
				if r.Ruta != " " {
					//r.ejecutarReporte_file()
				} else {
					r.singleton.AddSalidaConsola("EL PARAMETRO RUTA PARA EL REPORTE FILE ES OBLIGATORIO\n")
				}
			case "ls":
				if r.Ruta != " " {
					//r.ejecutarReporte_ls()
				} else {
					r.singleton.AddSalidaConsola("EL PARAMETRO RUTA PARA EL REPORTE LS ES OBLIGATORIO\n")
				}
			case "tree":
				//r.ejecutarReporte_tree()
			default:
				r.singleton.AddSalidaConsola("EL NOMBRE ASIGNADO PARA EL REPORTE ES INVALIDO\n")
			}
		} else {
			r.singleton.AddSalidaConsola("EL PARAMETRO DE LA UBICACION DEL REPORTE ES OBLIGATORIO\n")
		}
	} else {
		r.singleton.AddSalidaConsola("EL ID DE LA PARTICION ES OBLIGATORIO\n")
	}
}

/*
*metodo para imprimir reporte disk
 */
func (r *Rep) ejecutarReporte_disk() {
	nodoMontura := r.MountList.Buscar(r.Id)
	if nodoMontura == nil {
		r.singleton.AddSalidaConsola("No se encontró el ID de montaje especificado")
		return
	}

	var mbr MBR
	r.Extension = r.GetExtensionFile(r.Path)
	r.Directorio = r.GetCarpetas(r.Path)
	err := os.MkdirAll(r.Directorio, 0777) // Crea el directorio
	if err != nil {                        // Comprueba si hay un error
		fmt.Println(err) // Imprime el error
		return           // Sale de la función
	}
	os.Chmod(r.Directorio, 0777) // Cambia los permisos del directorio

	fileReporte, err := os.OpenFile(nodoMontura.Path, os.O_RDWR, 0777)
	if err != nil {
		r.singleton.AddSalidaConsola("No se pudo abrir el archivo del disco para lectura/escritura: " + err.Error())
		return
	}
	defer fileReporte.Close()

	if err := binary.Read(fileReporte, binary.LittleEndian, &mbr); err != nil {
		r.singleton.AddSalidaConsola("No se pudo leer el MBR del disco: " + err.Error())
		return
	}

	size := int(mbr.Mbr_tamano)

	fileDot, err := os.Create("disk.dot")
	if err != nil {
		r.singleton.AddSalidaConsola("No se pudo crear el archivo para el reporte: " + err.Error())
		return
	}
	defer fileDot.Close()

	fmt.Fprintln(fileDot, "digraph G {")
	fmt.Fprintln(fileDot, "node[shape=none]")
	fmt.Fprintln(fileDot, "start[label=<<table CELLSPACING=\"0\"><tr>")
	fmt.Fprintln(fileDot, "<td bgcolor=\"khaki1\" rowspan=\"2\">MBR</td>")

	start := int64(binary.Size(mbr))
	i := 0
	for i < 4 {
		if mbr.Mbr_partition[i].Part_start == -1 { // ESTA LIBRE
			i++
			for i < 4 {
				if mbr.Mbr_partition[i].Part_start != -1 {
					porcentaje := ((int(mbr.Mbr_partition[i].Part_start) - int(start)) / (size * 1.0)) * 100.0
					fmt.Fprintln(fileDot, fmt.Sprintf("<td bgcolor=\"lavender\" rowspan=\"2\">LIBRE <br/>%.0f</td>\n", math.Round(float64(porcentaje))))
					break
				}
				i++
			}
			if i == 4 {
				porcentaje := float64(size-int(start)) / float64(size) * 100.0
				fmt.Fprintf(fileDot, "<td bgcolor=\"lavender\" rowspan=\"2\">LIBRE <br/> %.0f </td>\n\n", math.Round(porcentaje))
				goto salida1
			}
			i--
		} else { // ESTA OCUPADA
			if mbr.Mbr_partition[i].Part_type == 'e' {
				contadorBloquesExtendida := 0
				var ebr EBR
				fileReporte.Seek(mbr.Mbr_partition[i].Part_start, 0)
				binary.Read(fileReporte, binary.LittleEndian, &ebr)
				if !(ebr.Part_s == -1 && ebr.Part_next == -1) {
					if ebr.Part_s > -1 {
						contadorBloquesExtendida += 2
					} else {
						contadorBloquesExtendida += 2
					}
					fileReporte.Seek(ebr.Part_next, 0)
					binary.Read(fileReporte, binary.LittleEndian, &ebr)
					for {
						contadorBloquesExtendida += 2
						if ebr.Part_next == -1 {
							if ebr.Part_start+ebr.Part_s < mbr.Mbr_partition[i].Part_s {
								contadorBloquesExtendida++
							}
							break
						} else {
							if ebr.Part_start+ebr.Part_s < ebr.Part_next {
								contadorBloquesExtendida++
							}
						}
						fileReporte.Seek(ebr.Part_next, 0)
						binary.Read(fileReporte, binary.LittleEndian, &ebr)
					}
				}
				fmt.Fprintln(fileDot, "<td bgcolor=\"darkolivegreen1\" colspan=\""+strconv.Itoa(contadorBloquesExtendida)+"\">EXTENDIDA</td>")
			} else if mbr.Mbr_partition[i].Part_type == 'p' {
				p1 := float64(mbr.Mbr_partition[i].Part_s) / float64(size)
				porcentaje := float64(p1) * 100.0
				name1 := mbr.Mbr_partition[i].Part_name
				fmt.Fprintf(fileDot, "<td bgcolor=\"violet\" rowspan=\"2\">%s <br/>%d</td>\n", string(bytes.TrimRight(name1[:], "\x00")), int(math.Round(float64(porcentaje))))
				if i != 3 {
					if mbr.Mbr_partition[i].Part_start+mbr.Mbr_partition[i].Part_s < mbr.Mbr_partition[i+1].Part_start {
						porcentaje = (float64(mbr.Mbr_partition[i+1].Part_start-(mbr.Mbr_partition[i].Part_start+mbr.Mbr_partition[i].Part_s)) / float64(size)) * 100
						fmt.Fprintf(fileDot, "<td bgcolor=\"lavender\" rowspan=\"2\">LIBRE <br/>%d</td>\n", int(math.Round(float64(porcentaje))))
					}
				} else if int(mbr.Mbr_partition[i].Part_start+mbr.Mbr_partition[i].Part_s) < size {
					porcentaje = (float64(size-(int(mbr.Mbr_partition[i].Part_start+mbr.Mbr_partition[i].Part_s))) / float64(size)) * 100
					fmt.Fprintf(fileDot, "<td bgcolor=\"lavender\" rowspan=\"2\">LIBRE <br/>%d</td>\n", int(math.Round(float64(porcentaje))))
				}
			}
			start = mbr.Mbr_partition[i].Part_start + mbr.Mbr_partition[i].Part_s
		}
		i++
	}

salida1:
	fmt.Fprintln(fileDot, "</tr>")

	// Por si hay extendida
	i = 0
	for i < 4 {
		if mbr.Mbr_partition[i].Part_start != -1 {
			if mbr.Mbr_partition[i].Part_type == 'e' {
				fmt.Fprintln(fileDot, "<tr>")
				porcentaje := (float64(mbr.Mbr_partition[i].Part_s) / float64(size)) * 100.0
				var ebr EBR
				fileReporte.Seek(mbr.Mbr_partition[i].Part_start, 0)
				binary.Read(fileReporte, binary.LittleEndian, &ebr)
				if !(ebr.Part_s == -1 && ebr.Part_next == -1) {
					nombre1 := strings.Trim(string(ebr.Part_name[:]), "\x00")
					if ebr.Part_s > -1 {
						fmt.Fprintf(fileDot, "<td bgcolor=\"steelblue1\" rowspan=\"1\">EBR <br/> %s </td>\n", string(bytes.TrimRight([]byte(nombre1[:]), "\x00")))
						porcentaje = (float64(ebr.Part_s) / float64(size)) * 100.0
						fmt.Fprintf(fileDot, "<td bgcolor=\"tan1\" rowspan=\"1\">Logica <br/> %d </td>\n", int(math.Round(porcentaje)))
					} else {
						fmt.Fprintln(fileDot, "<td bgcolor=\"steelblue1\" rowspan=\"1\">EBR</td>")
						porcentaje = (float64(ebr.Part_next-ebr.Part_start) / float64(size)) * 100.0
						fmt.Fprintf(fileDot, "<td bgcolor=\"lavender\" rowspan=\"1\">Libre <br/> %d </td>\n", int(math.Round(porcentaje)))
					}
					fileReporte.Seek(ebr.Part_next, 0)
					binary.Read(fileReporte, binary.LittleEndian, &ebr)
					for {
						name1 := ebr.Part_name
						fmt.Fprintf(fileDot, "<td bgcolor=\"steelblue1\" rowspan=\"1\">EBR <br/>%s</td>\n", string(bytes.TrimRight(name1[:], "\x00")))
						porcentaje := float64(ebr.Part_s) / float64(size) * 100.0
						fmt.Fprintf(fileDot, "<td bgcolor=\"tan1\" rowspan=\"1\">Logica <br/>%d</td>\n", int(math.Round(porcentaje)))
						if ebr.Part_next == -1 {
							if (ebr.Part_start + ebr.Part_s) < mbr.Mbr_partition[i].Part_s {
								porcentaje = float64(mbr.Mbr_partition[i].Part_s-(ebr.Part_start+ebr.Part_s)) / float64(size) * 100.0
								fmt.Fprintf(fileDot, "<td bgcolor=\"lavender\" rowspan=\"1\">Libre <br/>%d</td>\n", int(math.Round(porcentaje)))
							}
							break
						} else {
							if (ebr.Part_start + ebr.Part_s) < ebr.Part_next {
								porcentaje = float64(ebr.Part_next-(ebr.Part_start+ebr.Part_s)) / float64(size) * 100.0
								fmt.Fprintf(fileDot, "<td bgcolor=\"lavender\" rowspan=\"1\">Libre <br/>%d</td>\n", int(math.Round(porcentaje)))
							}
						}
						fileReporte.Seek(ebr.Part_next, 0)
						binary.Read(fileReporte, binary.LittleEndian, &ebr)
					}
				}
				fmt.Fprintln(fileDot, "</tr>\n")
			}
		}
		i++
	}

	fmt.Fprintln(fileDot, "</table>>];")
	fmt.Fprintln(fileDot, "}")

	if err := fileDot.Close(); err != nil {
		r.singleton.AddSalidaConsola("NO SE PUDO CERRAR EL ARCHIVO PARA EL REPORTE: " + err.Error() + "\n")
		return
	}

	// Crea el comando para ejecutar graphviz usando el archivo .dot y la imagen
	command := []string{"dot", "-T" + r.Extension, "disk.dot", "-o", r.Path}

	// Crea el objeto cmd con la función Command
	cmd := exec.Command("sudo", "-S", command[0], command[1], command[2], command[3], command[4])

	// Ejecuta el comando y obtiene la salida combinada
	_, err = cmd.CombinedOutput()
	if err != nil {
		r.singleton.AddSalidaConsola("NO SE PUDO GENERAR EL REPORTE: " + err.Error() + "\n")
		return
	}

	r.singleton.AddSalidaConsola("REPORTE GENERADO CON EXITO: DISK\n")
}
