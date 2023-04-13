package singleton

import (
	"sync"
)

type singleton struct {
	salidaConsola string
}

func (s *singleton) SalidaConsola() string {
	return s.salidaConsola
}

var instance *singleton
var once sync.Once

func GetInstance() *singleton {
	once.Do(func() {
		instance = &singleton{salidaConsola: ""}
	})
	return instance
}

func (s *singleton) AddSalidaConsola(str string) {
	s.salidaConsola += str
}

func (s *singleton) ResetSalidaConsola() {
	s.salidaConsola = ""
}
