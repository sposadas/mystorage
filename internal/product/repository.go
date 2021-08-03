package product

import (
	"context"
	"database/sql"
	database "github.com/sposadas/mystorage/db"
	"github.com/sposadas/mystorage/internal/domain"
	"log"
)

type Repository interface {
	Store(domain.Product) (domain.Product, error)
	GetOne(id int) domain.Product
	GetByName(name string) (domain.Product, error)
	Update(ctx context.Context, product domain.Product) (domain.Product, error)
	GetAll() ([]domain.Product, error)
	Delete(id int) error
	GetFullData(id int) domain.Product
	GetOneWithContext(ctx context.Context, id int) (domain.Product, error)
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

const (
	SelectProducts               = "SELECT p.id, p.name, p.type, p.count, p.price FROM products p"
	InsertProduct                = "INSERT INTO products(name, type, count, price) VALUES ( ?, ?, ?, ?)"
	GetOneProduct                = "SELECT p.id, p.name, p.type, p.count, p.price FROM products p WHERE id = ?"
	GetProductByName             = "SELECT p.id, p.name, p.type, p.count, p.price FROM products p WHERE name = ?"
	UpdateProduct                = "UPDATE products SET name = ?, type = ?, count = ?, price = ? WHERE id = ?"
	DeleteProduct                = "DELETE FROM products WHERE id = ?"
	SelectProductsWarehousesJoin = "SELECT p.id, p.name, p.type, p.count, p.price, w.name, w.address FROM products p INNER JOIN warehouses w ON p.id_warehouse = w.id where p.id = ?"
)

func (r *repository) Store(product domain.Product) (domain.Product, error) {
	db := database.StorageDB

	stmt, err := db.Prepare(InsertProduct)
	if err != nil {
		log.Println(err.Error())
		return domain.Product{}, err
	}
	defer stmt.Close()
	var result sql.Result
	result, err = stmt.Exec(product.Name, product.Type, product.Count, product.Price)
	if err != nil {
		return domain.Product{}, err
	}
	insertId, _ := result.LastInsertId()
	product.ID = int(insertId)

	return product, nil
}

func (r *repository) GetOne(id int) domain.Product {
	var product domain.Product
	db := database.StorageDB

	rows, err := db.Query(GetOneProduct, id)
	if err != nil {
		log.Println(err.Error())
		return product
	}

	for rows.Next() {
		if err := rows.Scan(&product.ID, &product.Name, &product.Type, &product.Count, &product.Price); err != nil {
			log.Println(err.Error())
			return product
		}
	}
	return product
}

func (r *repository) Update(ctx context.Context, product domain.Product) (domain.Product, error) {
	db := database.StorageDB
	result, err := db.ExecContext(ctx, UpdateProduct, product.Name, product.Type, product.Count, product.Price, product.ID)
	if err != nil {
		log.Println(err.Error())
		return domain.Product{}, err
	}
	_, err = result.RowsAffected()
	if err != nil {
		log.Println(err.Error())
		return domain.Product{}, err
	}
	return product, nil
}

func (r *repository) GetAll() ([]domain.Product, error) {
	var products []domain.Product
	db := database.StorageDB
	rows, err := db.Query(SelectProducts)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	for rows.Next() {
		var product domain.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Type, &product.Count, &product.Price); err != nil {
			log.Fatal(err.Error())
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (r *repository) Delete(id int) error {
	db := database.StorageDB
	stmt, err := db.Prepare(DeleteProduct)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (r *repository) GetFullData(id int) domain.Product {
	var product domain.Product
	db := database.StorageDB
	rows, err := db.Query(SelectProductsWarehousesJoin, id)
	if err != nil {
		log.Println(err.Error())
		return product
	}
	for rows.Next() {
		if err := rows.Scan(&product.ID, &product.Name, &product.Type, &product.Count, &product.Price, &product.Warehouse, &product.WarehouseAddress); err != nil {
			log.Println(err.Error())
			return product
		}
	}
	return product
}

func (r *repository) GetOneWithContext(ctx context.Context, id int) (domain.Product, error) {
	var product domain.Product
	db := database.StorageDB

	rows, err := db.QueryContext(ctx, GetOneProduct, id)
	if err != nil {
		log.Println(err.Error())
		return domain.Product{}, err
	}
	for rows.Next() {
		if err := rows.Scan(&product.ID, &product.Name, &product.Type, &product.Count, &product.Price); err != nil {
			log.Println(err.Error())
			return product, err
		}
	}
	return product, nil
}

func (r *repository) GetByName(name string) (domain.Product, error) {
	var product domain.Product
	db := database.StorageDB

	rows, err := db.Query(GetProductByName, name)
	if err != nil {
		log.Println(err)
		return domain.Product{}, err
	}
	for rows.Next() {
		if err := rows.Scan(&product.ID, &product.Name, &product.Type, &product.Count, &product.Price); err != nil {
			log.Println(err.Error())
			return domain.Product{}, err
		}
	}
	return product, nil
}
