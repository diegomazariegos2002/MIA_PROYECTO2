package comandos

import (
	"bytes"
	"encoding/binary"
	"miapp/singleton"
	"os"
	"time"
)

type Montar struct {
	P         string
	Name      string
	Id        string
	Fs        string
	MountList *MountList
	singleton *singleton.Singleton
}

func NewMontar() *Montar {
	return &Montar{
		P:         " ",
		Name:      " ",
		Id:        " ",
		Fs:        "2fs",
		MountList: NewMountList(),
		singleton: singleton.GetInstance(),
	}
}

func (m *Montar) Mount() {
	if m.P != "" {
		if m.Name != "" {
			b := []byte(m.Name)
			indiceParticion := -1
			file, notFile := os.OpenFile(m.P, os.O_RDWR, 0644)
			if notFile == nil {
				var mbr MBR
				file.Seek(0, 0)
				binary.Read(file, binary.LittleEndian, &mbr)
				for i := 0; i < 4; i++ {
					if bytes.Equal(b, mbr.Mbr_partition[i].Part_name[:len(b)]) {
						indiceParticion = i
						break
					} else if mbr.Mbr_partition[i].Part_type == 'e' { // ENTRANDO A BUSCAR SI ES LÓGICA
						var ebr EBR
						var superBloque SuperBloque
						file.Seek(mbr.Mbr_partition[i].Part_start, 0)
						notFile = binary.Read(file, binary.LittleEndian, &ebr)
						if !(ebr.Part_s == -1 && ebr.Part_next == -1) {
							if bytes.Equal(b, ebr.Part_name[:len(b)]) {
								if ebr.Part_status == '0' || ebr.Part_status == '1' {
									ebr.Part_status = '1'
								}
								m.MountList.Add(m.P, m.Name, 'l', int(ebr.Part_start), -1)
								file.Seek(int64(ebr.Part_start), 0)
								binary.Write(file, binary.LittleEndian, &ebr)
								if ebr.Part_status == '2' {
									file.Seek(int64(int(ebr.Part_start)+binary.Size(EBR{})), 0)
									binary.Read(file, binary.LittleEndian, &superBloque)
									superBloque.S_mtime = time.Now().Unix()
									superBloque.S_mnt_count++
									file.Seek(int64(int(ebr.Part_start)+(binary.Size(EBR{}))), 0)
									binary.Write(file, binary.LittleEndian, &superBloque)

								}
								file.Close()
								return
							} else if ebr.Part_next != -1 {
								file.Seek(int64(ebr.Part_next), 0)
								binary.Read(file, binary.LittleEndian, &ebr)
								for {
									if bytes.Equal(b, ebr.Part_name[:len(b)]) {
										if ebr.Part_status == '0' || ebr.Part_status == '1' {
											ebr.Part_status = '1'
										}
										m.MountList.Add(m.P, m.Name, 'l', int(ebr.Part_start), -1)

										file.Seek(int64(ebr.Part_start), 0)
										binary.Write(file, binary.LittleEndian, &ebr)

										if ebr.Part_status == '2' {
											file.Seek(int64(ebr.Part_start)+int64(binary.Size(EBR{})), 0)
											binary.Read(file, binary.LittleEndian, &superBloque)
											superBloque.S_mtime = time.Now().Unix()
											superBloque.S_mnt_count++
											file.Seek(int64(ebr.Part_start)+int64(binary.Size(EBR{})), 0)
											binary.Write(file, binary.LittleEndian, &superBloque)
										}

										return
									}
									if ebr.Part_next == -1 {
										break
									}
									file.Seek(int64(ebr.Part_next), 0)
									binary.Read(file, binary.LittleEndian, &ebr)
								}
							}

						}
					}
				}
				if indiceParticion != -1 {
					if mbr.Mbr_partition[indiceParticion].Part_type == 'e' {
						// Corrección al parecer si se montan extendidas.
						m.MountList.Add(m.P, m.Name, 'e', int(mbr.Mbr_partition[indiceParticion].Part_start), indiceParticion)
						m.singleton.AddSalidaConsola("PARTICION EXTENDIDA MONTADA")
						file.Close()
						return
					} else { // Se encontro que fue en un partición primaria
						var superBloque SuperBloque
						if mbr.Mbr_partition[indiceParticion].Part_status == '0' || mbr.Mbr_partition[indiceParticion].Part_status == '1' {
							mbr.Mbr_partition[indiceParticion].Part_status = '1'
						}
						m.MountList.Add(m.P, m.Name, 'p', int(mbr.Mbr_partition[indiceParticion].Part_start), indiceParticion)

						file.Seek(0, 0)
						binary.Write(file, binary.LittleEndian, &mbr)

						if mbr.Mbr_partition[indiceParticion].Part_status == '2' {
							file.Seek(int64(mbr.Mbr_partition[indiceParticion].Part_start), 0)
							binary.Read(file, binary.LittleEndian, &mbr)
							superBloque.S_mtime = time.Now().Unix()
							superBloque.S_mnt_count++
							file.Seek(int64(mbr.Mbr_partition[indiceParticion].Part_start), 0)
							binary.Write(file, binary.LittleEndian, &superBloque)
						}
						file.Close()
						return
					}
				} else {
					m.singleton.AddSalidaConsola("NO SE ENCONTRO UNA PARTICION CON ESE NOMBRE\n")
					file.Close()
					return
				}
			} else {
				m.singleton.AddSalidaConsola("NO EXISTE EL DISCO EN LA RUTA ESPECIFICADA\n")
			}
		} else {
			m.singleton.AddSalidaConsola("ERROR EL PARAMETRO NAME ES OBLIGATORIO\n")
		}
	} else {
		m.singleton.AddSalidaConsola("EL PARAMETRO RUTA ES OBLIGATORIO\n")
	}
}

func (m *Montar) Unmount() {
	// Implementación de la función unmount()
}

func (m *Montar) Mkfs() {
	// Implementación de la función mkfs()
}
