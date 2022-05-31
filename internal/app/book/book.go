package book

import "gorm.io/gorm"

type Book struct {
	gorm.Model

	Name    string `gorm:"column:name,index"`
	IBSN    string `gorm:"column:ibsn"`
	Auth    string `gorm:"column:auth"`
	Version string `gorm:"column:version"`
	Kind    int    `gorm:"column:kind"`
}
