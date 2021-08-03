package product

import (
	"context"
	"fmt"
	"github.com/sposadas/mystorage/internal/domain"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestStore(t *testing.T) {
	product := domain.Product{
		Name:  "testing",
		Type:  "testing",
		Count: 4,
		Price: 2600.00,
	}
	myRepository := NewRepository()
	productResult, err := myRepository.Store(product)
	if err != nil {
		log.Println(err)
	}
	assert.Equal(t, product.Name, productResult.Name)
}

func TestGetAll(t *testing.T) {
	size := 0
	myRepository := NewRepository()
	products, err := myRepository.GetAll()
	if err != nil {
		log.Println(err)
	}
	assert.NotEqual(t, size, len(products))
}

func TestGetOneWithContext(t *testing.T) {
	id := 1

	product := domain.Product{
		Name: "test",
	}

	myRepository := NewRepository()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	productResult, err := myRepository.GetOneWithContext(ctx, id)
	fmt.Println(err)
	assert.Equal(t, product.Name, productResult.Name)
}

func TestUpdate(t *testing.T) {
	product := domain.Product{
		ID:    1,
		Name:  "test01",
		Type:  "test01",
		Count: 9,
		Price: 1500.00,
	}

	myRepository := NewRepository()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	productResult, err := myRepository.Update(ctx, product)
	if err != nil {
		log.Println(err)
	}

	assert.Equal(t, product, productResult)
}

func TestGetByName(t *testing.T) {
	name := "testing"

	myRepository := NewRepository()
	productResult, err := myRepository.GetByName(name)
	if err != nil {
		log.Println(err)
	}
	assert.Equal(t, name, productResult.Name)
}