package handlers

import (
	"encoding/json"
	"testing"
	"time"

	rs_models "github.com/nutrixpos/pos/modules/reservations/models"
)

func TestReservation_Serialization(t *testing.T) {
	r := rs_models.Reservation{
		Id:              "res-1",
		CustomerName:    "Janez",
		CustomerPhone:   "041123456",
		GuestCount:      4,
		ReservationDate: "2024-02-14",
		ReservationTime: "19:00",
		Status:          "pending",
		CreatedAt:       time.Now(),
	}
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var d rs_models.Reservation
	if err := json.Unmarshal(data, &d); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if d.CustomerName != "Janez" {
		t.Errorf("name=%s", d.CustomerName)
	}
}

func TestReservation_Confirmed(t *testing.T) {
	r := rs_models.Reservation{
		Id:     "res-2",
		Status: "confirmed",
		TableAssignment: "T5",
	}
	data, _ := json.Marshal(r)
	var d rs_models.Reservation
	json.Unmarshal(data, &d)
	if d.TableAssignment != "T5" {
		t.Errorf("table=%s", d.TableAssignment)
	}
}

func TestReservationSlot_Serialization(t *testing.T) {
	s := rs_models.ReservationSlot{
		Date:      "2024-02-14",
		Time:      "19:00",
		Available: 3,
		Total:     10,
	}
	data, _ := json.Marshal(s)
	var d rs_models.ReservationSlot
	json.Unmarshal(data, &d)
	if d.Available != 3 {
		t.Errorf("available=%d", d.Available)
	}
}

func TestReservation_Delete(t *testing.T) {
	r := rs_models.Reservation{Id: "res-3"}
	data, _ := json.Marshal(r)
	var d rs_models.Reservation
	json.Unmarshal(data, &d)
	if d.Id != "res-3" {
		t.Errorf("id=%s", d.Id)
	}
}
