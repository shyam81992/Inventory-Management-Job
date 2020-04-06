package models

// Shipment modal
type Shipment struct {
	ID          int64  `json:"id"`
	OrderId     int64  `json:"inventoryid"`
	UsersId     int64  `json:"userid"`
	WarehouseId int64  `json:"warehouseid"`
	Status      string `json:"status"`
}
