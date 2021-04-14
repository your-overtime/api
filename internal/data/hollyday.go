package data

import "git.goasum.de/jasper/overtime/pkg"

func (d *Db) SaveHollyday(a *pkg.Hollyday) error {
	tx := d.Conn.Save(a)
	return tx.Error
}

func (d *Db) GetHollyday(id uint) (*pkg.Hollyday, error) {
	a := pkg.Hollyday{}
	tx := d.Conn.First(&a, id)
	return &a, tx.Error
}
