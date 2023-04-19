package singleton

import (
	"miapp/comandos"
	"sync"
)

type Singleton struct {
	salidaConsola string
	MountList     *comandos.MountList
}

var instance *Singleton
var once sync.Once

/*
*Constructor de única instancia
 */
func GetInstance() *Singleton {
	once.Do(func() {
		instance = &Singleton{
			salidaConsola: "",
			MountList:     comandos.NewMountList(),
		}
	})
	return instance
}

// Parte de métodos para salidaConsola

func (s *Singleton) SalidaConsola() string {
	return s.salidaConsola
}
func (s *Singleton) AddSalidaConsola(str string) {
	s.salidaConsola += str
}
func (s *Singleton) ResetSalidaConsola() {
	s.salidaConsola = ""
}
