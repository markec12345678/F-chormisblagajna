package handlers

import (
	"encoding/json"
	"testing"
	"time"

	dv_models "github.com/nutrixpos/pos/modules/delivery/models"
)

func TestDeliveryZone(t *testing.T) {
	z := dv_models.DeliveryZone{Id: "dz-1", Name: "Center", Fee: 3.50, MinOrder: 10, Active: true}
	data, _ := json.Marshal(z)
	var d dv_models.DeliveryZone
	json.Unmarshal(data, &d)
	if d.Name != "Center" {
		t.Errorf("name=%s", d.Name)
	}
}
func TestDeliveryOrder(t *testing.T) {
	o := dv_models.DeliveryOrder{Id: "do-1", OrderId: "o-1", CustomerName: "Miha", Status: "pending", PlacedAt: time.Now()}
	data, _ := json.Marshal(o)
	var d dv_models.DeliveryOrder
	json.Unmarshal(data, &d)
	if d.Status != "pending" {
		t.Errorf("status=%s", d.Status)
	}
}
func TestDeliveryOrderDelivered(t *testing.T) {
	now := time.Now()
	o := dv_models.DeliveryOrder{Id: "do-2", Status: "delivered", DeliveredAt: &now}
	data, _ := json.Marshal(o)
	var d dv_models.DeliveryOrder
	json.Unmarshal(data, &d)
	if d.DeliveredAt == nil {
		t.Error("expected delivered_at")
	}
}
func TestZoneFee(t *testing.T) {
	z := dv_models.DeliveryZone{Id: "dz-2", Fee: 0}
	if z.Fee != 0 {
		t.Error("fee should be 0")
	}
}
