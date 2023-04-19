package comandos

import "strconv"

type NodoM struct {
	Path  string
	Name  string
	Id    string
	Num   int
	Pos   int
	Type  byte
	Letra string
	Start int
	Sig   *NodoM
}

func NewNodoM(path string, name string, typeVal byte, num int, letra string, pos int, start int) *NodoM {
	id := "75" + strconv.Itoa(num) + letra
	return &NodoM{
		Path:  path,
		Name:  name,
		Id:    id,
		Num:   num,
		Pos:   pos,
		Type:  typeVal,
		Letra: letra,
		Start: start,
		Sig:   nil,
	}
}
