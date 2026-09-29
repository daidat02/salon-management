package product

import (
	"errors"
	"time"
)

const (
	ProductTypeRetail   = "retail"
	ProductTypeMaterial = "material"
	ProductTypeBoth     = "both"
)

type Product struct {
	ID             string
	OrganizationID string
	CategoryID     *string
	SKU            string
	Name           string
	Unit           string
	NetUnit        *string
	NetAmount      *float64
	ProductType    string
	ShowOnWeb      bool
	CostPrice      float64
	SellPrice      float64
	StockQuantity  float64
	MinStock       float64
	IsActive       bool
	CreatedAt      *time.Time
	UpdatedAt      *time.Time
}

func (p *Product) Validate() error {
	if p.SKU == "" {
		return errors.New("Mã SKU không được để trống")
	}
	if p.Name == "" {
		return errors.New("Tên sản phẩm không được để trống")
	}
	if p.ProductType != ProductTypeRetail && p.ProductType != ProductTypeMaterial && p.ProductType != ProductTypeBoth {
		return errors.New("Loại sản phẩm không hợp lệ (retail | material | both)")
	}
	if p.CostPrice < 0 {
		return errors.New("Giá vốn không được âm")
	}
	if p.SellPrice < 0 {
		return errors.New("Giá bán không được âm")
	}
	if p.StockQuantity < 0 {
		return errors.New("Tồn kho không được âm")
	}
	if p.MinStock < 0 {
		return errors.New("Tồn kho tối thiểu không được âm")
	}
	return nil
}
