package handlers

import (
	"encoding/json"
	"testing"

	pc_models "github.com/nutrixpos/pos/modules/purchase/models"
)

func TestPurchaseOrder(t *testing.T) {
	po := pc_models.PurchaseOrder{Id: "po-1", SupplierName: "Dobavitelj d.o.o.", Status: "pending",
		Items: []pc_models.POItem{{MaterialName: "Moka", Quantity: 10, UnitPrice: 1.5, TotalPrice: 15}}}
	data, _ := json.Marshal(po)
	var d pc_models.PurchaseOrder
	json.Unmarshal(data, &d)
	if d.SupplierName != "Dobavitelj d.o.o." {
		t.Errorf("supplier=%s", d.SupplierName)
	}
}
func TestPOItem(t *testing.T) {
	item := pc_models.POItem{MaterialName: "Sladkor", Quantity: 5, UnitPrice: 2, TotalPrice: 10}
	data, _ := json.Marshal(item)
	var d pc_models.POItem
	json.Unmarshal(data, &d)
	if d.TotalPrice != 10 {
		t.Errorf("total=%f", d.TotalPrice)
	}
}
func TestPOReceived(t *testing.T) {
	po := pc_models.PurchaseOrder{Id: "po-2", Status: "received"}
	data, _ := json.Marshal(po)
	var d pc_models.PurchaseOrder
	json.Unmarshal(data, &d)
	if d.Status != "received" {
		t.Errorf("status=%s", d.Status)
	}
}
func TestPOCancelled(t *testing.T) {
	po := pc_models.PurchaseOrder{Id: "po-3", Status: "cancelled"}
	data, _ := json.Marshal(po)
	var d pc_models.PurchaseOrder
	json.Unmarshal(data, &d)
	if d.Status != "cancelled" {
		t.Errorf("status=%s", d.Status)
	}
}
