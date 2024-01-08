package gorm

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

type Payload interface {
	TableName() string
}

type PayloadProduct struct {
	Color string
}

func newPayloadProduct(color string) *PayloadProduct {
	return &PayloadProduct{
		Color: color,
	}
}

func (p *PayloadProduct) TableName() string {
	return "test"
}

type Product struct {
	//unsupported data type
	//Payload
	*PayloadProduct

	Id    int    `gorm:"column(id),primary""`
	Name  string `gorm:"column(name)""`
	Price int    `gorm:"column(price)"`
}

func newProduct(payload *PayloadProduct) *Product {
	return &Product{
		PayloadProduct: payload,
	}
}

func (p Product) TableName() string {
	return "product"
}

func TestCRUD(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err, "sqlite open failed: %v", err)
	db = db.Debug()
	err = db.AutoMigrate(newProduct(newPayloadProduct("red")))
	require.NoError(t, err, "gorm migrate failed: %v", err)
	testCases := []struct {
		name    string
		op      string
		want    Product
		product Product
	}{
		{
			name: "create product",
			op:   "create",
			want: Product{
				PayloadProduct: newPayloadProduct("red"),
				Id:             1, Name: "Phone", Price: 5000},
			product: Product{
				PayloadProduct: newPayloadProduct("red"),
				Id:             1, Name: "Phone", Price: 5000},
		},
		{
			name: "update product",
			op:   "update",
			want: Product{
				PayloadProduct: newPayloadProduct("yellow"),
				Id:             1, Name: "Phone", Price: 3000},
			product: Product{
				PayloadProduct: newPayloadProduct("yellow"),
				Id:             1, Name: "Phone", Price: 3000},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			switch tc.op {
			case "create":
				db.Create(&tc.product)
			case "update":
				db.Updates(&tc.product)
			}
			var product Product
			db.First(&product, "id = ?", tc.product.Id)
			assert.Equal(t, tc.want, product)
		})
	}
}
