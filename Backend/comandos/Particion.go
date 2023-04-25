package comandos

import (
	"encoding/binary"
	"io"
	"log"
	"miapp/singleton"
	"os"
	"strconv"
	"unsafe"
)

type Particion struct {
	S             int
	Add           int
	U, T, F, Flag byte
	D, P, Name    string
	singleton     *singleton.Singleton
}

func NewParticion() *Particion {
	return &Particion{
		S:         0,
		Add:       0,
		U:         'k',
		P:         " ",
		T:         'P',
		F:         'w',
		D:         " ",
		Name:      " ",
		Flag:      'n',
		singleton: singleton.GetInstance()}
}

func (p *Particion) Fdisk() {
	if p.Name != " " {
		if p.P != " " {
			if p.D == "full" {
				p.deleteFullPartition()
			} else if p.Add != 0 && p.Flag == 'a' {
				if p.U == 'b' || p.U == 'k' || p.U == 'm' {
					if p.Add < 0 {
						p.reducePartition()
					} else if p.Add > 0 {
						p.incrementPartition()
					}
				} else {
					p.singleton.AddSalidaConsola("UNIDAD DE ALMACENAMIENTO INCORRECTA\n")
				}
				//validar esto tambien
			} else if p.S > 0 && p.Flag == 's' {
				if p.U == 'b' || p.U == 'k' || p.U == 'm' {
					if p.T == 'p' {
						p.primaryPartition()
					} else if p.T == 'e' {
						p.extendPartition()
					} else if p.T == 'l' {
						p.LogicPartition()
					} else {
						p.singleton.AddSalidaConsola("EL TIPO DE LA PARTICION ES INVALIDO\n")
					}
				} else {
					p.singleton.AddSalidaConsola("UNIDAD DE ALMACENAMIENTO INCORRECTO\n")
				}
			} else {
				p.singleton.AddSalidaConsola("PARA EL PARAMETRO -S SE ACEPTA UNICAMENTE VALORES MAYORES A CERO\n")
			}
		} else {
			p.singleton.AddSalidaConsola("ERROR FATAL EL PARAMETRO -P ES DE CARACTER OBLIGATORIO\n")
		}
	} else {
		p.singleton.AddSalidaConsola("ERROR FATAL EL PARAMETRO -NAME ES DE CARACTER OBLIGATORIO\n")
	}
}

func (p *Particion) primaryPartition() {
	var newPartitionPrimary Partition
	pos := -1
	file, err := os.OpenFile(p.P, os.O_RDWR, 0666)
	if err != nil {
		p.singleton.AddSalidaConsola("NO EXISTE EL DISCO EN LA RUTA ESPECIFICADA\n")
		return
	}
	defer file.Close()

	var mbr MBR
	file.Seek(0, 0)
	err = binary.Read(file, binary.LittleEndian, &mbr)
	if err != nil {
		p.singleton.AddSalidaConsola("ERROR AL LEER EL MBR\n")
		return
	}
	for i := 0; i < 4; i++ {
		if mbr.Mbr_partition[i].Part_start == -1 {
			pos = i
			break
		}
	}
	if p.freeSpace(p.S, p.P, p.U, pos) {
		if !p.existeParticion(p.P, p.Name) {
			newPartitionPrimary.Part_fit = p.F
			newPartitionPrimary.Part_type = p.T
			copy(newPartitionPrimary.Part_name[:], p.Name)
			newPartitionPrimary.Part_status = '0'
			if p.U == 'b' {
				newPartitionPrimary.Part_s = int64(p.S)
			} else if p.U == 'k' {
				newPartitionPrimary.Part_s = int64(p.S * 1024)
			} else if p.U == 'm' {
				newPartitionPrimary.Part_s = int64(p.S * 1024 * 1024)
			}
			// buscando donde ubicar en la dirección de memoria el comienzo de la partición
			if pos == 0 {
				newPartitionPrimary.Part_start = int64(binary.Size(mbr))
			} else {
				newPartitionPrimary.Part_start = mbr.Mbr_partition[pos-1].Part_start + mbr.Mbr_partition[pos-1].Part_s
			}
			mbr.Mbr_partition[pos] = newPartitionPrimary

			_, err := file.Seek(0, 0)
			if err != nil {
				p.singleton.AddSalidaConsola("ERROR AL MOVERSE AL INICIO DEL ARCHIVO\n")
				return
			}
			err = binary.Write(file, binary.LittleEndian, &mbr)
			if err != nil {
				p.singleton.AddSalidaConsola("ERROR AL ESCRIBIR EL MBR\n")
				return
			}

			var mbrVerificador MBR
			file.Seek(0, 0)
			binary.Read(file, binary.LittleEndian, &mbrVerificador)
			file.Close()
			p.singleton.AddSalidaConsola("OPERACION REALIZADA CON EXITO\n")
			p.singleton.AddSalidaConsola("PARTICION " + strconv.Itoa(pos+1) + "\n")
			p.singleton.AddSalidaConsola("NOMBRE: " + string(mbrVerificador.Mbr_partition[pos].Part_name[:]) + "\n")
			p.singleton.AddSalidaConsola("TIPO: PARTICION PRIMARIA\n")
			p.singleton.AddSalidaConsola("INICIO: " + strconv.Itoa(int(mbrVerificador.Mbr_partition[pos].Part_start)) + "\n")
			p.singleton.AddSalidaConsola("SIZE: " + strconv.Itoa(int(mbrVerificador.Mbr_partition[pos].Part_s)) + "\n")
		} else {
			p.singleton.AddSalidaConsola("YA HAY UNA PARTICION CON ESE NOMBRE " + p.Name + "\n")
		}
	} else {
		p.singleton.AddSalidaConsola("NO EXISTE EL ESPACIO NECESARIO PARA REALIZAR ESTE COMANDO\n")
	}
}

func (p *Particion) extendPartition() {
	var newPartitionExtend Partition
	var indice, addressEBR int = -1, -1
	file, err := os.OpenFile(p.P, os.O_RDWR, 0644)
	if err != nil {
		p.singleton.AddSalidaConsola("NO EXISTE UN DISCO CON LA RUTA ESPECIFICADA\n")
		return
	}
	defer file.Close()
	var mbr MBR
	_, err = file.Seek(0, 0)
	err = binary.Read(file, binary.LittleEndian, &mbr)
	if err != nil {
		p.singleton.AddSalidaConsola("ERROR LEYENDO EL MBR\n")
		return
	}
	for i := 0; i < 4; i++ {
		if mbr.Mbr_partition[i].Part_start == -1 {
			indice = i
			break
		}
	}

	if p.freeSpace(p.S, p.P, p.U, indice) {
		if !p.existeParticion(p.P, p.Name) {
			if !p.existeParticionExtendida(p.P) {
				newPartitionExtend.Part_fit = p.F
				newPartitionExtend.Part_type = p.T
				copy(newPartitionExtend.Part_name[:], p.Name)
				newPartitionExtend.Part_status = '0'
				if p.U == 'b' {
					newPartitionExtend.Part_s = int64(p.S)
				} else if p.U == 'k' {
					newPartitionExtend.Part_s = int64(p.S * 1024)
				} else if p.U == 'm' {
					newPartitionExtend.Part_s = int64(p.S * 1024 * 1024)
				}

				if indice == 0 {
					newPartitionExtend.Part_start = int64(binary.Size(mbr))
				} else {
					newPartitionExtend.Part_start = mbr.Mbr_partition[indice-1].Part_start + mbr.Mbr_partition[indice-1].Part_s
				}

				addressEBR = int(newPartitionExtend.Part_start)
				mbr.Mbr_partition[indice] = newPartitionExtend

				_, err = file.Seek(0, 0)
				if err != nil {
					p.singleton.AddSalidaConsola("ERROR SEEKING FILE\n")
					return
				}

				err = binary.Write(file, binary.LittleEndian, &mbr)
				if err != nil {
					p.singleton.AddSalidaConsola("ERROR ESCRIBIENDO MBR\n")
					return
				}

				ebr := EBR{
					Part_next:   -1,
					Part_start:  int64(addressEBR),
					Part_s:      -1,
					Part_status: '0',
				}

				_, err = file.Seek(int64(addressEBR), 0)
				if err != nil {
					p.singleton.AddSalidaConsola("ERROR SEEKING FILE\n")
					return
				}

				err = binary.Write(file, binary.LittleEndian, &ebr)
				if err != nil {
					p.singleton.AddSalidaConsola("ERROR ESCRIBIENDO EBR\n")
					return
				}

				p.singleton.AddSalidaConsola("OPERACION REALIZADA CON EXITO\n")
				p.singleton.AddSalidaConsola("PARTICION " + strconv.Itoa(indice+1) + "\n")
				p.singleton.AddSalidaConsola("NOMBRE: " + string(mbr.Mbr_partition[indice].Part_name[:]) + "\n")
				p.singleton.AddSalidaConsola("TIPO: PARTICION EXTENDIDA\n")
				p.singleton.AddSalidaConsola("INICIO: " + strconv.Itoa(int(mbr.Mbr_partition[indice].Part_start)) + "\n")
				p.singleton.AddSalidaConsola("SIZE: " + strconv.Itoa(int(mbr.Mbr_partition[indice].Part_s)) + "\n")
			} else {
				p.singleton.AddSalidaConsola("YA EXISTE UNA PARTICION EXTENDIDA CON ESE NOMBRE\n")
			}
		} else {
			p.singleton.AddSalidaConsola("YA EXISTE UNA PARTICION CON ESE NOMBRE " + p.Name + "\n")
		}
	} else {
		p.singleton.AddSalidaConsola("NO EXISTE EL ESPACIO NECESARIO PARA EJECUTAR ESTE COMANDO\n")
	}
}

func (p *Particion) LogicPartition() {
	file, err := os.OpenFile(p.P, os.O_RDWR, 0644)
	if err != nil {
		p.singleton.AddSalidaConsola("NO SE PUDO ABRIR EL ARCHIVO\n")
		return
	}
	defer file.Close()

	if p.existeParticionExtendida(p.P) {
		if !p.existeParticion(p.P, p.Name) {
			var indice int = -1
			var mbr MBR
			if err := binary.Read(file, binary.LittleEndian, &mbr); err != nil {
				p.singleton.AddSalidaConsola("NO SE PUDO LEER EL MBR\n")
				return
			}

			for i := 0; i < 4; i++ {
				if mbr.Mbr_partition[i].Part_type == 'e' {
					indice = i
					break
				}
			}

			if indice != -1 {
				var ebrAuxiliar EBR
				var fullSpace int = 0
				var filePos int = int(mbr.Mbr_partition[indice].Part_start) // se agrega inicialización

				if _, err := file.Seek(int64(mbr.Mbr_partition[indice].Part_start), io.SeekStart); err != nil {
					p.singleton.AddSalidaConsola("NO SE PUDO MOVER EL PUNTERO EN EL ARCHIVO\n")
					return
				}

				if err := binary.Read(file, binary.LittleEndian, &ebrAuxiliar); err != nil {
					p.singleton.AddSalidaConsola("NO SE PUDO LEER EL EBR AUXILIAR\n")
					return
				}

				if ebrAuxiliar.Part_next != -1 || ebrAuxiliar.Part_s != -1 {
					fullSpace += (int)(unsafe.Sizeof(EBR{})) + (int)(ebrAuxiliar.Part_s)
					for ebrAuxiliar.Part_next != -1 && filePos < (int)(mbr.Mbr_partition[indice].Part_start+mbr.Mbr_partition[indice].Part_s) {

						if _, err := file.Seek(int64(ebrAuxiliar.Part_next), io.SeekStart); err != nil {
							p.singleton.AddSalidaConsola("NO SE PUDO MOVER EL PUNTERO EN EL ARCHIVO\n")
							return
						}

						if err := binary.Read(file, binary.LittleEndian, &ebrAuxiliar); err != nil {
							p.singleton.AddSalidaConsola("NO SE PUDO LEER EL EBR AUXILIAR\n")
							return
						}

						fullSpace += (int)(unsafe.Sizeof(EBR{})) + (int)(ebrAuxiliar.Part_s)
					}

					var newExtend EBR
					newExtend.Part_fit = p.F
					newExtend.Part_start = int64(((int)(ebrAuxiliar.Part_start)) + (int)(unsafe.Sizeof(EBR{})) + ((int)(ebrAuxiliar.Part_s)))
					newExtend.Part_status = '0'
					newExtend.Part_next = -1
					copy(newExtend.Part_name[:], p.Name)

					if p.U == 'b' {
						newExtend.Part_s = int64(p.S)
					} else if p.U == 'k' {
						newExtend.Part_s = int64(p.S * 1024)
					} else if p.U == 'm' {
						newExtend.Part_s = int64(p.S * 1024 * 1024)
					}

					freeSpace := (int)(mbr.Mbr_partition[indice].Part_s) - fullSpace
					espacioNewE := (int)(unsafe.Sizeof(EBR{})) + (int)(newExtend.Part_s)
					ebrAuxiliar.Part_next = newExtend.Part_start
					if freeSpace >= espacioNewE {
						ebrAuxiliar.Part_next = newExtend.Part_start
						file.Seek(ebrAuxiliar.Part_start, 0)
						if err := binary.Write(file, binary.LittleEndian, &ebrAuxiliar); err != nil {
							log.Fatal(err)
						}
						file.Seek(newExtend.Part_start, 0)
						if err := binary.Write(file, binary.LittleEndian, &newExtend); err != nil {
							log.Fatal(err)
						}
						file.Close()

						file, err = os.OpenFile(p.P, os.O_RDWR, 0644)
						var ebrAux EBR
						var ebrNew EBR
						file.Seek(ebrAuxiliar.Part_start, 0)
						binary.Read(file, binary.LittleEndian, &ebrAux)
						file.Seek(ebrAuxiliar.Part_next, 0)
						binary.Read(file, binary.LittleEndian, &ebrNew)
						file.Close()
						p.singleton.AddSalidaConsola("OPERACION REALIZADA CON EXITO\n")
						p.singleton.AddSalidaConsola("Nombre particion: " + string(ebrNew.Part_name[:]) + "\n")
						p.singleton.AddSalidaConsola("Tipo: Logica\n")
						p.singleton.AddSalidaConsola("Inicio: " + strconv.Itoa(int(ebrNew.Part_start)) + "\n")
						p.singleton.AddSalidaConsola("Size: " + strconv.Itoa(int(ebrNew.Part_s)) + "\n")
						p.singleton.AddSalidaConsola("EBR Anterior next: " + strconv.Itoa(int(ebrAux.Part_next)) + "\n")
					} else {
						p.singleton.AddSalidaConsola("NO EXISTE EL ESPACIO NECESARIO PARA EJECUTAR ESTE COMANDO\n")
						file.Close()
					}
				} else {
					ebrAuxiliar.Part_fit = p.F
					ebrAuxiliar.Part_start = mbr.Mbr_partition[indice].Part_start
					ebrAuxiliar.Part_status = '0'
					if p.U == 'b' {
						ebrAuxiliar.Part_s = int64(p.S)
					} else if p.U == 'k' {
						ebrAuxiliar.Part_s = int64(p.S * 1024)
					} else if p.U == 'm' {
						ebrAuxiliar.Part_s = int64(p.S * 1024 * 1024)
					}
					ebrAuxiliar.Part_next = -1
					copy(ebrAuxiliar.Part_name[:], p.Name)

					if int(mbr.Mbr_partition[indice].Part_s) >= (int(ebrAuxiliar.Part_s) + (binary.Size(EBR{}))) {
						file.Seek(0, io.SeekStart)
						binary.Write(file, binary.LittleEndian, &mbr)

						file.Seek(int64(mbr.Mbr_partition[indice].Part_start), io.SeekStart)
						binary.Write(file, binary.LittleEndian, ebrAuxiliar)
						file.Close()

						file, err = os.OpenFile(p.P, os.O_RDWR, 0644)
						var ebr2 EBR
						file.Seek(mbr.Mbr_partition[indice].Part_start, 0)
						binary.Read(file, binary.LittleEndian, &ebr2)

						p.singleton.AddSalidaConsola("OPERACION REALIZADA CON EXITO\n")
						p.singleton.AddSalidaConsola("Nombre particion: " + string(ebr2.Part_name[:]) + "\n")
						p.singleton.AddSalidaConsola("Tipo: Logica\n")
						p.singleton.AddSalidaConsola("Inicio: " + strconv.Itoa(int(ebr2.Part_start)) + "\n")
						p.singleton.AddSalidaConsola("Size: " + strconv.Itoa(int(ebr2.Part_s)) + "\n")
						p.singleton.AddSalidaConsola("Part_next: " + strconv.Itoa(int(ebr2.Part_next)) + "\n")
					} else {
						file.Close()
						p.singleton.AddSalidaConsola("NO EXISTE EL ESPACIO SUFICIENTE PARA EJECUTAR ESTE COMANDO\n")
					}
				}
			} else {
				p.singleton.AddSalidaConsola("YA EXISTE UNA PARTICIÓN CON ESE NOMBRE\n")
			}
		} else {
			p.singleton.AddSalidaConsola("NO HAY PARTICIÓN EXTENDIDA EN EL DISCO\n")
		}
	}
}

func (p *Particion) deleteFullPartition() {}

func (p *Particion) reducePartition() {}

func (p *Particion) incrementPartition() {}

func (p *Particion) deleteMegaByte(path string, posicionInicial int64, posicionFinal int64) {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	_, err = file.Seek(posicionInicial, io.SeekStart)
	if err != nil {
		log.Fatal(err)
	}
	resta := (posicionFinal - posicionInicial) / (1024 * 1024)
	if resta >= 1 {
		buffer := make([]byte, 1024)
		for i := 0; i < 1024; i++ {
			buffer[i] = 0
		}
		for j := int64(0); j < resta*1024; j++ {
			_, err = file.Write(buffer)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	posicionInicial, _ = file.Seek(0, 1)
	p.deleteKiloByte(path, posicionInicial, posicionFinal)
	file.Close()
}

func (p *Particion) deleteKiloByte(path string, posicionInicial int64, posicionFinal int64) {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Seek(posicionInicial, 0)
	if err != nil {
		log.Fatal(err)
	}

	resta := (posicionFinal - posicionInicial) / 1024
	if resta >= 1 {
		buffer := make([]byte, 1024)
		for i := range buffer {
			buffer[i] = 0
		}
		for j := int64(0); j < resta; j++ {
			_, err = file.Write(buffer)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	posicionInicial, _ = file.Seek(0, 1)
	p.deleteByte(path, posicionInicial, posicionFinal)
	file.Close()
}

func (p *Particion) deleteByte(path string, posicionInicial int64, posicionFinal int64) {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	_, err = file.Seek(posicionInicial, 0)
	if err != nil {
		log.Fatal(err)
	}

	resta := posicionFinal - posicionInicial
	if resta >= 1 {
		buffer := make([]byte, resta)
		for i := range buffer {
			buffer[i] = 0
		}
		_, err = file.Write(buffer)
		if err != nil {
			log.Fatal(err)
		}
	}
	file.Close()
}

func (p *Particion) freeSpace(s int, path string, u byte, address int) bool {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return false
	}
	defer file.Close()

	var mbr MBR
	_, err = file.Seek(0, 0)
	err = binary.Read(file, binary.LittleEndian, &mbr)
	if err != nil {
		return false
	}

	if address > -1 {
		size := 0
		if u == 'b' {
			size = s
		} else if u == 'k' {
			size = s * 1024
		} else if u == 'm' {
			size = s * 1024 * 1024
		}

		if size > 0 {
			freeSpace := 0
			if address == 0 { // si es la primera particion
				freeSpace = int(mbr.Mbr_tamano) - binary.Size(mbr)
			} else { // si no es la primera particion
				freeSpace = int(mbr.Mbr_tamano) - int(mbr.Mbr_partition[address-1].Part_start) - int(mbr.Mbr_partition[address-1].Part_s)
			}
			return freeSpace >= size
		}
		return false
	}
	return false
}

func (p *Particion) existeParticionExtendida(path string) bool {
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	defer file.Close()
	if err != nil {
		p.singleton.AddSalidaConsola("ERROR AL ABRIR EL ARCHIVO: " + err.Error())
		return false
	}
	var mbr MBR
	_, err = file.Seek(0, 0)
	err = binary.Read(file, binary.LittleEndian, &mbr)
	if err != nil {
		p.singleton.AddSalidaConsola("ERROR AL LEER EL MBR: " + err.Error())
		return false
	}
	for i := 0; i < 4; i++ {
		if mbr.Mbr_partition[i].Part_type == 'e' {
			return true
		}
	}
	return false
}

func (p *Particion) existeParticion(path string, name string) bool {
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	file.Seek(0, 0)
	if err != nil {
		p.singleton.AddSalidaConsola("ERROR AL LEER EL ARCHIVO: " + err.Error())
	}
	defer file.Close()

	var mbr MBR
	file.Seek(0, 0)
	if err := binary.Read(file, binary.LittleEndian, &mbr); err != nil {
		p.singleton.AddSalidaConsola("ERROR AL LEER EL MBR: " + err.Error())
	}

	for i := 0; i < 4; i++ {
		name1 := string(mbr.Mbr_partition[i].Part_name[:])
		if name1 == name {
			return true
		}

		if mbr.Mbr_partition[i].Part_type == 'e' {
			var ebr EBR
			if _, err := file.Seek(int64(mbr.Mbr_partition[i].Part_start), 0); err != nil {
				p.singleton.AddSalidaConsola("ERROR AL BUSCAR EL INICIO DEL EBR: " + err.Error())
			}
			if err := binary.Read(file, binary.LittleEndian, &ebr); err != nil {
				p.singleton.AddSalidaConsola("ERROR AL LEER EL EBR: " + err.Error())
			}
			if ebr.Part_next != -1 || ebr.Part_s != -1 {
				name1 = string(ebr.Part_name[:])
				if name1 == name {
					return true
				}
				for ebr.Part_next != -1 {
					name1 = string(ebr.Part_name[:])
					if name1 == name {
						return true
					}
					if _, err := file.Seek(int64(ebr.Part_next), 0); err != nil {
						p.singleton.AddSalidaConsola("ERROR AL BUSCAR EL SIGUIENTE EBR: " + err.Error())
					}
					if err := binary.Read(file, binary.LittleEndian, &ebr); err != nil {
						p.singleton.AddSalidaConsola("ERROR AL LEER EL EBR: " + err.Error())
					}
				}
			}
		}
	}
	return false
}
