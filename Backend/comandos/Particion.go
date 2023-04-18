package comandos

type Particion struct {
	S             int
	Add           int
	U, T, F, Flag byte
	D, P, Name    string
}

func NewParticion() *Particion {
	return &Particion{
		S:    0,
		Add:  0,
		U:    'k',
		P:    " ",
		T:    'p',
		F:    'w',
		D:    " ",
		Name: " ",
		Flag: 'n',
	}
}

func (p *Particion) fdisk() {}

func (p *Particion) primaryPartition() {}

func (p *Particion) extendPartition() {}

func (p *Particion) LogicPartition() {}

func (p *Particion) deleteFullPartition() {}

func (p *Particion) reducePartition() {}

func (p *Particion) incrementPartition() {}

func (p *Particion) deleteMegaByte(path string, pos int64, posicionFinal int) {}

func (p *Particion) deleteKiloByte(path string, posicionInicial int64, posicionFinal int) {}

func (p *Particion) deleteByte(path string, posicionInicial int64, posicionFinal int) {}

func (p *Particion) freeSpace(s int, path string, u byte, address int) bool {
	return false
}

func (p *Particion) existeParticionExtendida(path string) bool {
	return false
}

func (p *Particion) existeParticion(path string, name string) bool {
	return false
}
