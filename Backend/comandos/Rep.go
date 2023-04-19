package comandos

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math"
	"miapp/singleton"
	"os"
	"os/exec"
	"path/filepath"
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
				//r.ejecutarReporte_disk()
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
	nodoMontura := r.singleton.MountList.Buscar(r.Id)
	if nodoMontura == nil {
		r.singleton.AddSalidaConsola("No se encontró el ID de montaje especificado")
		return
	}

	var mbr MBR
	r.Extension = r.GetExtensionFile(r.Path)
	r.Directorio = r.GetCarpetas(r.Path)
	os.MkdirAll(r.Directorio, 0777)
	os.Chmod(r.Directorio, 0777)

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
					porcentaje := ((mbr.Mbr_partition[i].Part_start - int(start)) / (size * 1.0)) * 100.0
					fmt.Fprintln(fileDot, fmt.Sprintf("<td bgcolor=\"lavender\" rowspan=\"2\">LIBRE <br/>%.0f</td>\n", math.Round(float64(porcentaje))))
					break
				}
				i++
			}
			if i == 4 {
				porcentaje := float64(size-mbr.Mbr_partition[i-1].Part_start-mbr.Mbr_partition[i-1].Part_s) / float64(size) * 100.0
				r.singleton.AddSalidaConsola(fmt.Sprintf("<td bgcolor=\"lavender\" rowspan=\"2\">LIBRE <br/>%.0f</td>\n", math.Round(porcentaje)))
				goto salida1
			}
			i--
		} else { // ESTA OCUPADA
			if mbr.Mbr_partition[i].Part_type == 'e' {
				contadorBloquesExtendida := 0
				var ebr EBR
				if _, err := fileReporte.Seek(int64(mbr.Mbr_partition[i].Part_start), 0); err != nil {
					r.singleton.AddSalidaConsola("ERROR AL LEER EL ARCHIVO: " + err.Error() + "\n")
					return
				}
				if err := binary.Read(fileReporte, binary.LittleEndian, &ebr); err != nil {
					r.singleton.AddSalidaConsola("ERROR AL LEER EL ARCHIVO: " + err.Error() + "\n")
					return
				}
				if !(ebr.Part_s == -1 && ebr.Part_next == -1) {
					if ebr.Part_s > -1 {
						contadorBloquesExtendida += 2
					} else {
						contadorBloquesExtendida += 2
					}
					if _, err := fileReporte.Seek(int64(ebr.Part_next), 0); err != nil {
						r.singleton.AddSalidaConsola("ERROR AL LEER EL ARCHIVO: " + err.Error() + "\n")
						return
					}
					if err := binary.Read(fileReporte, binary.LittleEndian, &ebr); err != nil {
						r.singleton.AddSalidaConsola("ERROR AL LEER EL ARCHIVO: " + err.Error() + "\n")
						return
					}
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
						if _, err := fileReporte.Seek(int64(ebr.Part_next), 0); err != nil {
							r.singleton.AddSalidaConsola("ERROR AL LEER EL ARCHIVO: " + err.Error() + "\n")
							return
						}
						if err := binary.Read(fileReporte, binary.LittleEndian, &ebr); err != nil {
							r.singleton.AddSalidaConsola("NO SE PUDO LEER EL EBR DEL DISCO: " + err.Error())
							return
						}
					}
				}
				fmt.Fprintln(fileDot, "<td bgcolor=\"darkolivegreen1\" colspan=\""+string(contadorBloquesExtendida)+"\">EXTENDIDA</td>")
			} else if mbr.Mbr_partition[i].Part_type == 'p' {
				p1 := float32(mbr.Mbr_partition[i].Part_s) / float32(size)
				porcentaje := p1 * 100.0
				name1 := mbr.Mbr_partition[i].Part_name
				fmt.Fprintf(fileDot, "<td bgcolor=violet rowspan=2>%s <br/>%d</td>", name1, int(math.Round(float64(porcentaje))))
				if i != 3 {
					if mbr.Mbr_partition[i].Part_start+mbr.Mbr_partition[i].Part_s < mbr.Mbr_partition[i+1].Part_start {
						porcentaje = (float32(mbr.Mbr_partition[i+1].Part_start-(mbr.Mbr_partition[i].Part_start+mbr.Mbr_partition[i].Part_s)) / float32(size)) * 100
						fmt.Fprintf(fileDot, "<td bgcolor=lavender rowspan=2>LIBRE <br/>%d</td>", int(math.Round(float64(porcentaje))))
					}
				} else if mbr.Mbr_partition[i].Part_start+mbr.Mbr_partition[i].Part_s < size {
					porcentaje = (float32(size-(mbr.Mbr_partition[i].Part_start+mbr.Mbr_partition[i].Part_s)) / float32(size)) * 100
					fmt.Fprintf(fileDot, "<td bgcolor=lavender rowspan=2>LIBRE <br/>%d</td>", int(math.Round(float64(porcentaje))))
				}
			}
			start = int64(mbr.Mbr_partition[i].Part_start + mbr.Mbr_partition[i].Part_s)
		}
		i++
	}

salida1:
	fmt.Fprintln(fileDot, "</tr>\n")
	// Por si hay extendida
	i = 0
	for i < 4 {
		if mbr.Mbr_partition[i].Part_start != -1 {
			if mbr.Mbr_partition[i].Part_type == 'e' {
				fmt.Fprintf(fileDot, "<tr>\n")
				porcentaje := (float64(mbr.Mbr_partition[i].Part_s) / float64(size)) * 100.0
				var ebr EBR
				if _, err := fileDot.Seek(int64(mbr.Mbr_partition[i].Part_start), 0); err != nil {
					r.singleton.AddSalidaConsola("ERROR AL LEER EL ARCHIVO: " + err.Error() + "\n")
					return
				}
				if err := binary.Read(fileDot, binary.LittleEndian, &ebr); err != nil {
					r.singleton.AddSalidaConsola("ERROR AL LEER EL ARCHIVO: " + err.Error() + "\n")
					return
				}
				if !(ebr.Part_s == -1 && ebr.Part_next == -1) {
					nombre_1 := strings.Trim(string(ebr.Part_name[:]), "\x00")
					if ebr.Part_s > -1 {
						fmt.Fprintf(fileDot, "<td bgcolor=\"steelblue1\" rowspan=\"1\">EBR <br/> %s </td>", nombre_1)
						porcentaje = (float64(ebr.Part_s) / float64(size)) * 100.0
						fmt.Fprintf(fileDot, "<td bgcolor=\"tan1\" rowspan=\"1\">Logica <br/> %d </td>", int(math.Round(porcentaje)))
					} else {
						fmt.Fprintf(fileDot, "<td bgcolor=\"steelblue1\" rowspan=\"1\">EBR</td>")
						porcentaje = (float64(ebr.Part_next-ebr.Part_start) / float64(size)) * 100.0
						fmt.Fprintf(fileDot, "<td bgcolor=\"lavender\" rowspan=\"1\">Libre <br/> %d </td>", int(math.Round(porcentaje)))
					}
					if _, err := fileDot.Seek(int64(ebr.Part_next), 0); err != nil {
						r.singleton.AddSalidaConsola("ERROR AL LEER EL ARCHIVO: " + err.Error() + "\n")
						return
					}
					if err := binary.Read(fileDot, binary.LittleEndian, &ebr); err != nil {
						r.singleton.AddSalidaConsola("ERROR AL LEER EL ARCHIVO: " + err.Error() + "\n")
						return
					}
					for {
						name1 := ebr.Part_name
						fmt.Fprintf(fileDot, "<td bgcolor=\"steelblue1\" rowspan=\"1\">EBR <br/>%s</td>", name1)
						porcentaje := float64(ebr.Part_s) / float64(size) * 100.0
						fmt.Fprintf(fileDot, "<td bgcolor=\"tan1\" rowspan=\"1\">Logica <br/>%d</td>", int(math.Round(porcentaje)))
						if ebr.Part_next == -1 {
							if (ebr.Part_start + ebr.Part_s) < mbr.Mbr_partition[i].Part_s {
								porcentaje = float64(mbr.Mbr_partition[i].Part_s-(ebr.Part_start+ebr.Part_s)) / float64(size) * 100.0
								fmt.Fprintf(fileDot, "<td bgcolor=\"lavender\" rowspan=\"1\">Libre <br/>%d</td>", int(math.Round(porcentaje)))
							}
							break
						} else {
							if (ebr.Part_start + ebr.Part_s) < ebr.Part_next {
								porcentaje = float64(ebr.Part_next-(ebr.Part_start+ebr.Part_s)) / float64(size) * 100.0
								fmt.Fprintf(fileDot, "<td bgcolor=\"lavender\" rowspan=\"1\">Libre <br/>%d</td>", int(math.Round(porcentaje)))
							}
						}
						if _, err := fileReporte.Seek(int64(ebr.Part_next), io.SeekStart); err != nil {
							log.Fatal(err)
						}
						if err := binary.Read(fileReporte, binary.LittleEndian, &ebr); err != nil {
							log.Fatal(err)
						}
					}
				}
				fmt.Fprintf(fileDot, "</tr>\n")
			}
		}
		i++
	}
	fmt.Fprintln(fileDot, "</table>>];")
	fmt.Fprintln(fileDot, "}")

	if err := fileDot.Close(); err != nil {
		r.singleton.AddSalidaConsola("NO SE PUDO CERRAR EL ARCHIVO PARA EL REPORTE: " + err.Error())
		return
	}

	command := "dot -T" + r.Extension + " disk.dot -o \"" + r.Path + "\""
	if _, err := exec.Command("sudo", "-S", command).Output(); err != nil {
		r.singleton.AddSalidaConsola("NO SE PUDO GENERAR EL REPORTE: " + err.Error())
		return
	}

	r.singleton.AddSalidaConsola("REPORTE GENERADO CON EXITO: DISK\n")
}
