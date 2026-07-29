package handlers

import (
	"encoding/json"
	"testing"

	kk_models "github.com/nutrixpos/pos/modules/kiosk/models"
)

func TestKioskConfig(t *testing.T) {
	k := kk_models.KioskConfig{Id: "k-1", Name: "Entry Kiosk", Theme: "dark", Active: true}
	data, _ := json.Marshal(k)
	var d kk_models.KioskConfig
	json.Unmarshal(data, &d)
	if d.Name != "Entry Kiosk" {
		t.Errorf("name=%s", d.Name)
	}
}
func TestKioskOrder(t *testing.T) {
	o := kk_models.KioskOrder{Id: "ko-1", KioskId: "k-1", Status: "pending", Items: []kk_models.KioskOrderItem{{ProductName: "Burger", Quantity: 2, UnitPrice: 8}}}
	data, _ := json.Marshal(o)
	var d kk_models.KioskOrder
	json.Unmarshal(data, &d)
	if len(d.Items) != 1 {
		t.Errorf("items=%d", len(d.Items))
	}
}
func TestKioskOrderTotal(t *testing.T) {
	o := kk_models.KioskOrder{Total: 16}
	data, _ := json.Marshal(o)
	var d kk_models.KioskOrder
	json.Unmarshal(data, &d)
	if d.Total != 16 {
		t.Errorf("total=%f", d.Total)
	}
}
func TestKioskConfigInactive(t *testing.T) {
	k := kk_models.KioskConfig{Active: false}
	if k.Active {
		t.Error("should be inactive")
	}
}
