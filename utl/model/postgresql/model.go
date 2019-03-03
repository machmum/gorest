package gorestdb

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Base struct {
	CreatedBy *int       `json:"-"`
	CreatedAt *time.Time `json:"-"`
	UpdatedBy *int       `json:"-"`
	UpdatedAt *time.Time `json:"-"`
	DeletedBy *int       `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

type DB struct {
	*gorm.DB
}

func NewDB(db *gorm.DB) DB {
	// Disable table name's pluralization globally
	db.SingularTable(true)

	return DB{
		DB: db,
	}
}

// return type is *gorm.DB
func (newdb DB) NewTx() DB {
	return NewDB(newdb.Begin())
	// return newdb.Begin()
}

func (newdb DB) CommitTx() DB {
	return NewDB(newdb.Commit())
	// return newdb.Commit()
}

func (newdb DB) RollbackTx() DB {
	return NewDB(newdb.Rollback())
	// return like below if return is *gorm.DB
	// return newdb.RollbackTx()
}
