package database

import (
	"fmt"
	"github.com/ebauman/widgetfactory/types"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB struct {
	db *gorm.DB
}

func New(dsn string) *DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("failed to connect to database: %s", err.Error())
	}

	if err := db.AutoMigrate(&types.Order{}, &types.Widget{}); err != nil {
		logrus.Fatalf("database migration failed: %s", err.Error())
	}

	return &DB{
		db: db,
	}
}

func (d *DB) CreateWidget(w *types.Widget) (*types.Widget, error) {
	tx := d.db.Create(w)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return w, nil
}

func (d *DB) CreateOrder(o *types.Order) (*types.Order, error) {
	tx := d.db.Create(o)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return o, nil
}

func (d *DB) ListWidgets() ([]types.Widget, error) {
	var widgets = make([]types.Widget, 0)

	result := d.db.Find(&widgets)

	if result.Error != nil {
		return nil, result.Error
	}

	return widgets, nil
}

func (d *DB) ListOrders() ([]types.Order, error) {
	var orders = make([]types.Order, 0)

	result := d.db.Find(&orders)

	if result.Error != nil {
		return nil, result.Error
	}

	return orders, nil
}

func (d *DB) DeleteWidget(w *types.Widget) error {
	tx := d.db.Where("ID = ?", w.ID).Delete(w)

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return fmt.Errorf("error deleting widget, no rows affected")
	}

	return nil
}

func (d *DB) DeleteOrder(o *types.Order) error {
	tx := d.db.Where("ID = ?", o.ID).Delete(o)

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return fmt.Errorf("error deleting order, no rows affected")
	}

	return nil
}

func (d *DB) Query(statement string) (any, error) {
	var res map[string]any
	tx := d.db.Raw(statement).Scan(&res)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return res, nil
}
