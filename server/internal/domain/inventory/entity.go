package inventory

import (
	"errors"
	"time"
)

const (
	TypeImport = "import"
	TypeExport = "export"
	TypeAdjust = "adjust"
	TypeSale   = "sale"
	TypeReturn = "return"
)

type InventoryTransaction struct {
	ID             string     `json:"id" db:"id"`
	OrganizationID string     `json:"organization_id" db:"organization_id"`
	ProductID      string     `json:"product_id" db:"product_id"`
	Type           string     `json:"type" db:"type"`
	Quantity       float64    `json:"quantity" db:"quantity"`
	ReferenceType  *string    `json:"reference_type" db:"reference_type"`
	ReferenceID    *string    `json:"reference_id" db:"reference_id"`
	CreatedBy      *string    `json:"created_by" db:"created_by"`
	CreatedAt      *time.Time `json:"created_at" db:"created_at"`
}

func (t *InventoryTransaction) Validate() error {
	if t.ProductID == "" {
		return errors.New("Thiếu product_id của phiếu kho")
	}
	switch t.Type {
	case TypeImport, TypeExport, TypeAdjust, TypeSale, TypeReturn:
	default:
		return errors.New("Loại giao dịch kho không hợp lệ (import | export | adjust | sale | return)")
	}
	if t.Quantity == 0 {
		return errors.New("Số lượng phiếu kho phải khác 0")
	}
	if (t.Type == TypeImport || t.Type == TypeReturn) && t.Quantity < 0 {
		return errors.New("Phiếu nhập/trả hàng phải có số lượng dương")
	}
	if (t.Type == TypeExport || t.Type == TypeSale) && t.Quantity > 0 {
		return errors.New("Phiếu xuất/bán hàng phải có số lượng âm")
	}
	return nil
}

// LowStockItem dùng cho cảnh báo tồn thấp, đọc từ bảng products.
type LowStockItem struct {
	ProductID     string  `json:"product_id" db:"product_id"`
	SKU           string  `json:"sku" db:"sku"`
	Name          string  `json:"name" db:"name"`
	StockQuantity float64 `json:"stock_quantity" db:"stock_quantity"`
	MinStock      float64 `json:"min_stock" db:"min_stock"`
}


// Receipts 

type InventoryDocument struct{
	ID             string
	OrganizationID string
	SupplierID     *string
	OrderID		*string
	DocumentCode   string
	Type		   string
	TotalAmount    float64
	Reason         *string
	Note           *string
	CreatedBy      *string
	CreatedAt      *time.Time
}

func (d *InventoryDocument) Validate() error {
	switch d.Type {
	case TypeImport, TypeExport, TypeAdjust, TypeSale, TypeReturn:
	default:
		return errors.New("Loại phiếu kho không hợp lệ (import | export | adjust | sale | return)")
	}
	if d.DocumentCode == "" {
		return errors.New("Thiếu mã chứng từ kho")
	}
	if d.TotalAmount < 0 {
		return errors.New("Tổng tiền phiếu kho không được âm")
	}
	return nil
}

type InventoryDocumentItem struct{
	ID             string
	OrganizationID string
	ProductID      string
	ProductName    string
	SKU     	   string
	DocumentID     string
	NetUnit       *string
	NetAmount     *float64
	UnitPrice      float64
	Quantity       float64
	TotalPrice     float64
	CreatedAt      *time.Time
}