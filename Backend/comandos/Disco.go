package comandos

import (
	"encoding/binary"
	"miapp/singleton"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Disco struct {
	S         int
	F         string
	U         string
	P         string
	PathFull  string
	singleton *singleton.Singleton
}

func NewDisco() *Disco {
	return &Disco{
		S:         0,
		F:         "bf",
		U:         "m",
		P:         " ",
		PathFull:  "",
		singleton: singleton.GetInstance()}
}

// validar es un método que verifica si el disco tiene los atributos correctos
func (d *Disco) validar() bool {
	bandera := false
	if d.S > 0 {
		if strings.HasPrefix(d.F, "bf") || strings.HasPrefix(d.F, "ff") || strings.HasPrefix(d.F, "wf") {
			if strings.HasPrefix(d.U, "k") || strings.HasPrefix(d.U, "m") {
				if d.P != " " {
					i := strings.Index(d.P, ".")
					extension := d.P[i+1:]
					if extension == "dsk" {
						bandera = true
					} else {
						d.singleton.AddSalidaConsola(">> EXTENSION INCORRECTA\n")
					}
				} else {
					d.singleton.AddSalidaConsola(">> RUTA INCORRECTA\n")
				}
			} else {
				d.singleton.AddSalidaConsola(">> UNIDADES DEL TAMAÑO DE MEMORIA INVALIDO\n")
			}
		} else {
			d.singleton.AddSalidaConsola(">> AJUSTE INVALIDO\n")
		}
	} else {
		d.singleton.AddSalidaConsola(">> EL TAMAÑO DEL DISCO TIENE QUE SER MAYOR A 0\n")
	}
	return bandera
}

func (d *Disco) Mkdisk() {
	// implementar la lógica de crear disco
	if d.validar() {
		if _, err := os.Stat(d.P); !os.IsNotExist(err) {
			d.singleton.AddSalidaConsola(">> ERROR EL DISCO YA EXISTE")
		} else {
			d.PathFull = d.GetDirectorio(d.P)
			cmd := exec.Command("sudo", "mkdir", "-p", d.PathFull)
			cmd.Run()

			cmd = exec.Command("sudo", "chmod", "-R", "777", d.PathFull)
			cmd.Run()

			var buffer [1024]byte
			size := d.S
			if d.U == "k" {
				size = size * 1024
			} else if d.U == "m" {
				size = size * 1024 * 1024
			}
			file, _ := os.Create(d.P)
			defer file.Close()

			for i := 0; i < 1024; i++ {
				buffer[i] = 0
			}
			for i := 0; i < size; i += 1024 {
				file.Write(buffer[:])
			}
			file.Close()
			file, _ = os.OpenFile(d.P, os.O_RDWR, 0644)
			defer file.Close()

			_, err = file.Seek(0, 0)

			var mbr MBR
			mbr.Mbr_fecha_creacion = time.Now().Unix()
			mbr.Mbr_dsk_signature = int64(int(time.Now().Unix()))
			mbr.Mbr_tamano = int64(size)
			mbr.Disk_fit = d.F[0]
			for j := 0; j < 4; j++ {
				mbr.Mbr_partition[j].Part_start = -1
				mbr.Mbr_partition[j].Part_type = 'P'
			}
			err = binary.Write(file, binary.LittleEndian, &mbr)

			d.singleton.AddSalidaConsola(">> COMANDO EJECUTADO CON EXITO SE CREO EL DISCO EXITOSAMENTE\n")
			file.Close()
			// Abrir el archivo binario en modo lectura
			f, err := os.OpenFile(d.P, os.O_RDWR, 0644)
			if err != nil {
				d.singleton.AddSalidaConsola(err.Error())
				return
			}
			defer f.Close()
			// Crear un struct MBR vacío para almacenar los datos leídos
			var m2 MBR

			// Leer el MBR del archivo binario usando encoding/binary
			err = binary.Read(f, binary.LittleEndian, &m2)
			if err != nil {
				d.singleton.AddSalidaConsola(err.Error())
				return
			}

			// Imprimir los datos del MBR por la consola usando fmt.Println
			d.singleton.AddSalidaConsola("Datos del MBR leídos del disco binario:" + "\n")
			d.singleton.AddSalidaConsola("Mbr_tamano: " + strconv.Itoa(int(m2.Mbr_tamano)) + "\n")
			d.singleton.AddSalidaConsola("Mbr_fecha_creacion: " + strconv.FormatInt(m2.Mbr_fecha_creacion, 10) + "\n")
			d.singleton.AddSalidaConsola("Mbr_dsk_signature: " + strconv.Itoa(int(m2.Mbr_dsk_signature)) + "\n")
			d.singleton.AddSalidaConsola("Disk_fit: " + string(m2.Disk_fit) + "\n")
			for i, p := range m2.Mbr_partition {
				d.singleton.AddSalidaConsola("Mbr_partition[%d]: " + strconv.Itoa(i) + "\n")
				d.singleton.AddSalidaConsola("\tPart_status: " + strconv.Itoa(int(p.Part_status)) + "\n")
				d.singleton.AddSalidaConsola("\tPart_type: " + strconv.Itoa(int(p.Part_type)) + "\n")
				d.singleton.AddSalidaConsola("\tPart_fit: " + strconv.Itoa(int(p.Part_fit)) + "\n")
				d.singleton.AddSalidaConsola("\tPart_start: " + strconv.Itoa(int(p.Part_start)) + "\n")
				d.singleton.AddSalidaConsola("\tPart_s: " + strconv.Itoa(int(p.Part_s)) + "\n")
				d.singleton.AddSalidaConsola("\tPart_name: " + string(p.Part_name[:]) + "\n")
			}
			file.Close()
		}
	}
}

func (d *Disco) Rmdisk() {
	if d.P != " " {
		i := strings.Index(d.P, ".")
		extension := d.P[i+1:]

		if extension == "dsk" {
			_, err := os.Stat(d.P)
			if err == nil {
				os.Remove(d.P)
				d.singleton.AddSalidaConsola(">> COMANDO EJECUTADO CON EXITO DISCO ELIMINADO\n")
			} else {
				d.singleton.AddSalidaConsola(">> EL DISCO NO EXISTE, VERIFIQUE LA RUTA\n")
			}
		} else {
			d.singleton.AddSalidaConsola(">> EXTENSION INCORRECTA\n")
		}
	} else {
		d.singleton.AddSalidaConsola(">> ASEGURESE DE ESCRIBIR UNA RUTA\n")
	}
}

func (d *Disco) GetDirectorio(path string) string {
	directorio := ""
	aux := path
	for i := len(aux) - 1; i >= 0; i-- {
		if aux[i] == '/' {
			directorio = aux[:i+1]
			break
		}
	}
	return directorio
}
