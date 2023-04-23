package singleton

import (
	"sync"
)

type Singleton struct {
	salidaConsola string
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
