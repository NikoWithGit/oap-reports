package repoimpl

import (
	"database/sql"
	"oap-reposts/iface"
	"oap-reposts/model"
	"time"
)

type ReportRepoImpl struct {
	db iface.Idb
}

func NewReportRepoImpl(db iface.Idb) *ReportRepoImpl {
	return &ReportRepoImpl{db}
}

func (rri *ReportRepoImpl) AddOrder(order *model.Order) error {
	tx, err := rri.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO completed_orders(id, date, total) VALUES($1, $2, $3)", order.Id, order.Date, order.Total)
	if err != nil {
		return err
	}
	for i := range order.Products {
		order.Products[i].OrderId = order.Id
		err = rri.AddProduct(tx, &order.Products[i])
		if err != nil {
			return err
		}
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (rri *ReportRepoImpl) AddProduct(tx *sql.Tx, product *model.ProductInOrder) error {
	_, err := tx.Exec(
		"INSERT INTO products_in_orders(uuid, num, price_per_one, order_id) VALUES($1, $2, $3, $4)",
		product.Uuid, product.Num, product.PricePerOne, product.OrderId,
	)
	return err
}

func (rri *ReportRepoImpl) GetAll(from time.Time, to time.Time) ([]model.Order, error) {
	orders, err := rri.getAllRaw(from, to)
	if err != nil {
		return nil, err
	}
	for i := range orders {
		orders[i].Products, err = rri.getProductsByOrderId(orders[i].Id)
		if err != nil {
			return nil, err
		}
	}
	return orders, nil
}

func (rri *ReportRepoImpl) getAllRaw(from time.Time, to time.Time) ([]model.Order, error) {
	rows, err := rri.db.Query("SELECT id, total, date FROM completed_orders WHERE date BETWEEN $1 AND $2", from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []model.Order
	for rows.Next() {
		var order model.Order
		rows.Scan(&order.Id, &order.Total, &order.Date)
		orders = append(orders, order)
	}
	return orders, nil
}

func (rri *ReportRepoImpl) getProductsByOrderId(order_id string) ([]model.ProductInOrder, error) {
	rows, err := rri.db.Query("SELECT uuid, num, price_per_one FROM products_in_orders WHERE order_id = $1", order_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []model.ProductInOrder
	for rows.Next() {
		var product model.ProductInOrder
		rows.Scan(&product.Uuid, &product.Num, &product.PricePerOne)
		products = append(products, product)
	}
	return products, nil
}
