package inventory

import "time"

type InventoryDocumentInput struct {
	ProductID string  `json:"product_id"`
	Type      string  `json:"type"`
	Quantity  float64 `json:"quantity"`
}

type CreateInventoryRequest struct {
	ProductID string  `json:"product_id"`
	Type      string  `json:"type"`
	Quantity  float64 `json:"quantity"`
}

type InventoryResponse struct {
	ID             string     `json:"id" db:"id"`
	OrganizationID string     `json:"organization_id" db:"organization_id"`
	ProductID      string     `json:"product_id" db:"product_id"`
	Type           string     `json:"type" db:"type"`
	Quantity       float64    `json:"quantity" db:"quantity"`
	ReferenceType  *string    `json:"reference_type" db:"reference_type"`
	ReferenceID    *string    `json:"reference_id" db:"reference_id"`
	CreatedAt      *time.Time `json:"created_at" db:"created_at"`
}

type DocumentItemInput struct {
	ProductID string  `json:"product_id"`
	Unit string 	   `json:"unit"`
	Quantity  float64 `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

type CreateDocumentRequest struct {
	SupplierID   *string             `json:"supplier_id"`
	OrderID      *string             `json:"order_id"`
	DocumentCode string              `json:"document_code"`
	Type         string              `json:"type"`
	Note         *string             `json:"note"`
	Reason 	 	 *string             `json:"reason"`
	Items        []DocumentItemInput `json:"items"`
	CreatedBy    *string             `json:"created_by"`
}

type DocumentItemResponse struct {
	ID          string   `json:"id" db:"id"`
	DocumentID  string   `json:"document_id" db:"document_id"`
	ProductID   string   `json:"product_id" db:"product_id"`
	ProductName string   `json:"product_name" db:"product_name"`
	SKU         string   `json:"sku" db:"sku"`
	Unit        string   `json:"unit" db:"unit"`
	NetUnit     string  `json:"net_unit" db:"net_unit"`
	NetAmount   float64 `json:"net_amount" db:"net_amount"`
	Quantity    float64  `json:"quantity" db:"quantity"`
	UnitPrice   float64  `json:"unit_price" db:"unit_price"`
	TotalPrice  float64  `json:"total_price" db:"total_price"`
	// Field hiển thị đã quy đổi (tính ở query detail, không lưu DB)
	DisplayQuantity  float64 `json:"display_quantity" db:"display_quantity"`
	DisplayUnit      string  `json:"display_unit" db:"display_unit"`
	DisplayUnitPrice float64 `json:"display_unit_price" db:"display_unit_price"`
}

type DocumentResponse struct {
	ID             string     `json:"id" db:"id"`
	OrganizationID string     `json:"organization_id" db:"organization_id"`
	DocumentCode   string     `json:"document_code" db:"document_code"`
	Type           string     `json:"type" db:"type"`
	TotalAmount    float64    `json:"total_amount" db:"total_amount"`
	Note           *string    `json:"note" db:"note"`
	Reason         *string    `json:"reason" db:"reason"`
	ItemsCount     *int        `json:"items_count" db:"items_count"`

	SupplierID     *string    `json:"supplier_id" db:"supplier_id"`
	SupplierName   *string    `json:"supplier_name" db:"supplier_name"`
	SupplierPhone   *string    `json:"supplier_phone" db:"supplier_phone"`
	SupplierAddress   *string    `json:"supplier_address" db:"supplier_address"`
	OrderID        *string    `json:"order_id" db:"order_id"`
	OrderCode		*string    `json:"order_code" db:"order_code"`

	CreatedBy      *string    `json:"created_by" db:"created_by"`
	CreatedByName  *string    `json:"created_by_name" db:"created_by_name"`
	CreatedAt      *time.Time `json:"created_at" db:"created_at"`
}

type DocumentDetailResponse struct {
	Document *DocumentResponse       `json:"document"`
	Items    []*DocumentItemResponse `json:"items"`
}
type DocumentQueryResult struct {
    DocumentResponse
    ItemsRaw []byte `db:"items"` // Hứng chuỗi JSON từ lệnh jsonb_agg
}


type InventoryTransactionResponse struct {
	ID string `json:"id" db:"id"`
	OrganizationID string `json:"organization_id" db:"organization_id"`
	Type string `json:"type" db:"type"`
	Quantity float64 `json:"quantity" db:"quantity"`
	RefrenceType *string `json:"reference_type" db:"reference_type"`
	ReferenceID *string `json:"reference_id" db:"reference_id"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	ProductID string `json:"product_id" db:"product_id"`
	DocumentID *string `json:"document_id" db:"document_id"`
	DocumentCode *string `json:"document_code" db:"document_code"`
	OrderID *string `json:"order_id" db:"order_id"`
	OrderCode *string `json:"order_code" db:"order_code"`
	SupplierID *string `json:"supplier_id" db:"supplier_id"`
	SupplierName *string `json:"supplier_name" db:"supplier_name"`
}