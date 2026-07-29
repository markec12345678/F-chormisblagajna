package handlers

import (
	"encoding/json"
	"testing"

	q_models "github.com/nutrixpos/pos/modules/queue/models"
)

func TestQueueEntry(t *testing.T) {
	e := q_models.QueueEntry{Id: "q-1", CustomerName: "Ana", PartySize: 3, Position: 1, Status: "waiting"}
	data, _ := json.Marshal(e)
	var d q_models.QueueEntry
	json.Unmarshal(data, &d)
	if d.Position != 1 {
		t.Errorf("pos=%d", d.Position)
	}
}
func TestQueueSeated(t *testing.T) {
	e := q_models.QueueEntry{Id: "q-2", Status: "seated"}
	data, _ := json.Marshal(e)
	var d q_models.QueueEntry
	json.Unmarshal(data, &d)
	if d.Status != "seated" {
		t.Errorf("status=%s", d.Status)
	}
}
func TestQueueCancel(t *testing.T) {
	e := q_models.QueueEntry{Id: "q-3", Status: "cancelled"}
	data, _ := json.Marshal(e)
	var d q_models.QueueEntry
	json.Unmarshal(data, &d)
	if d.Status != "cancelled" {
		t.Errorf("status=%s", d.Status)
	}
}
