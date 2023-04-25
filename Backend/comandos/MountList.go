package comandos

import (
	"miapp/singleton"
	"strings"
)

type MountList struct {
	primero   *NodoM
	ultimo    *NodoM
	singleton *singleton.Singleton
}

func NewMountList() *MountList {
	return &MountList{
		primero:   nil,
		ultimo:    nil,
		singleton: singleton.GetInstance(),
	}
}

func (m *MountList) ExistMount(path string, name string) bool {
	actual := m.primero
	for actual != nil {
		if actual.Path == path && actual.Name == name {
			return true
		}
		actual = actual.Sig
	}
	return false
}

func (m *MountList) GetNum(path string) int {
	num := 0
	actual := m.primero
	for actual != nil {
		if actual.Path == path && actual.Num > num {
			num = actual.Num
		}
		actual = actual.Sig
	}
	return num + 1
}

func (m *MountList) GetName(path string) string {
	aux := path
	p := 0
	name := ""
	for {
		p = strings.Index(aux, "/")
		if p == -1 {
			break
		}
		name += aux[:p+1]
		aux = aux[p+1:]
	}

	if p = strings.Index(aux, "."); p != -1 {
		name = aux[:p]
	}
	return name
}

func (m *MountList) Buscar(id string) *NodoM {
	actual := m.primero
	for actual != nil {
		if actual.Id == id {
			return actual
		}
		actual = actual.Sig
	}
	return nil
}

func (m *MountList) eliminar(id string) bool {
	if m.primero != nil {
		if m.primero == m.ultimo && m.primero.Id == id {
			m.primero, m.ultimo = nil, nil
			return true
		} else if m.primero.Id == id {
			m.primero = m.primero.Sig
			return true
		} else {
			aux, ant := m.primero.Sig, m.primero
			for aux != nil {
				if aux.Id == id {
					ant.Sig = aux.Sig
					aux.Sig = nil
					return true
				}
				ant = aux
				aux = aux.Sig
			}
			m.singleton.AddSalidaConsola("EL ID " + id + " NO REPRESENTA A NINGUNA MONTURA\n")
			return false
		}
	} else {
		m.singleton.AddSalidaConsola("IMPOSIBLE EJECUTAR NO EXISTEN MONTURAS EN EL SISTEMA\n")
		return false
	}
}

func (m *MountList) Add(path string, name string, t byte, start int, pos int) {
	if !m.ExistMount(path, name) {
		num := m.GetNum(path)
		letra := m.GetName(path)
		nuevo := NewNodoM(path, name, t, num, letra, start, pos)
		if m.primero == nil {
			m.primero = nuevo
			m.ultimo = nuevo
		} else {
			m.ultimo.Sig = nuevo
			m.ultimo = nuevo
		}
		m.singleton.AddSalidaConsola("SE MONTO CON EXITO LA PARTICION " + name + " CON ID " + nuevo.Id + "\n")
	} else {
		m.singleton.AddSalidaConsola("IMPOSIBLE EJECUTAR LA PARTICION YA ESTA MONTADA\n")
	}
}
