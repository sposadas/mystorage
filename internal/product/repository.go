package product

import (
	"database/sql"
	database "github.com/sposadas/mystorage/db"
	"github.com/sposadas/mystorage/internal/domain"
	"log"
)

type Repository interface {
	Store(domain.Product) (domain.Product, error)
	GetOne(id int) domain.Product
	Update(product domain.Product) (domain.Product, error)
	GetAll() ([]domain.Product, error)
	Delete(id int) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

const (
	InsertProduct = "INSERT INTO products(name, type, count, price) VALUES ( ?, ?, ?, ?)"
	GetOneProduct = "SELECT * FROM products WHERE id = ?"
	UpdateProduct = "UPDATE products SET name = ?, type = ?, count = ?, price = ? WHERE id = ?"
)

func (r *repository) Store(product domain.Product) (domain.Product, error) {
	db := database.StorageDB

	stmt, err := db.Prepare(InsertProduct)
	if err != nil {
		log.Fatal(err)
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
		if err := rows.Scan(&product.Name, &product.Type, &product.Count, &product.Price); err != nil {
			log.Println(err.Error())
			return product
		}
	}
	return product
}

func (r *repository) Update(product domain.Product) (domain.Product, error) {
	db := database.StorageDB
	stmt, err := db.Prepare(UpdateProduct)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(product.Name, product.Type, product.Count, product.Price, product.ID)
	if err != nil {
		return domain.Product{}, err
	}
	return product, nil
}

func (r *repository) GetAll() ([]domain.Product, error) {
	return []domain.Product{}, nil
}

func (r *repository) Delete(id int) error {
	return nil
}
