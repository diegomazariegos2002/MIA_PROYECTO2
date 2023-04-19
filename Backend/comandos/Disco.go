package comandos

import (
	"encoding/binary"
	"miapp/singleton"
	"os"
	"os/exec"
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
			cmd := exec.Command("sudo", "-S", "mkdir", "-p", d.PathFull)
			cmd.Run()
			cmd = exec.Command("sudo", "-S", "chmod", "-R", "777", d.PathFull)
			cmd.Run()

			var buffer [1024]byte
			size := d.S
			if d.U == "m" {
				size = size * 1024
			}
			file, _ := os.Create(d.P)
			defer file.Close()

			for i := 0; i < 1024; i++ {
				buffer[i] = 0
			}
			for i := 0; i < size; i += 1024 {
				file.Write(buffer[:])
			}

			file, _ = os.OpenFile(d.P, os.O_RDWR, 0644)
			defer file.Close()

			var mbr MBR
			mbr.Mbr_fecha_creacion = time.Now()
			mbr.Mbr_dsk_signature = int32(int(time.Now().Unix()))
			mbr.Mbr_tamano = int32(size * 1024)
			mbr.Disk_fit = d.F[0]
			for j := 0; j < 4; j++ {
				mbr.Mbr_partition[j].Part_start = -1
				mbr.Mbr_partition[j].Part_type = 'p'
			}
			binary.Write(file, binary.LittleEndian, &mbr)
			d.singleton.AddSalidaConsola(">> COMANDO EJECUTADO CON EXITO SE CREO EL DISCO EXITOSAMENTE\n")
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
