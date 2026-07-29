package handlers

import (
	"encoding/json"
	"testing"

	fp_models "github.com/nutrixpos/pos/modules/floorplan/models"
)

func TestFloorTable(t *testing.T) {
	tbl := fp_models.FloorTable{Id: "t-1", Label: "T1", Zone: "Main", Capacity: 4, Shape: "rect", X: 100, Y: 200, Status: "available"}
	data, _ := json.Marshal(tbl)
	var d fp_models.FloorTable
	json.Unmarshal(data, &d)
	if d.Label != "T1" {
		t.Errorf("label=%s", d.Label)
	}
}
func TestFloorZone(t *testing.T) {
	z := fp_models.FloorZone{Id: "z-1", Name: "Terrace", W: 800, H: 600}
	data, _ := json.Marshal(z)
	var d fp_models.FloorZone
	json.Unmarshal(data, &d)
	if d.Name != "Terrace" {
		t.Errorf("name=%s", d.Name)
	}
}
func TestTableOccupied(t *testing.T) {
	tbl := fp_models.FloorTable{Status: "occupied"}
	if tbl.Status != "occupied" {
		t.Error("should be occupied")
	}
}
