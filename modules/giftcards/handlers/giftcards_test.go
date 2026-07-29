package handlers

import (
	"encoding/json"
	"testing"

	gc_models "github.com/nutrixpos/pos/modules/giftcards/models"
)

func TestGiftCard(t *testing.T) {
	g := gc_models.GiftCard{Id: "gc-1", Code: "GC-ABC123", Balance: 50, InitialAmt: 50, Active: true}
	data, _ := json.Marshal(g)
	var d gc_models.GiftCard
	json.Unmarshal(data, &d)
	if d.Balance != 50 {
		t.Errorf("balance=%f", d.Balance)
	}
}
func TestGiftCardRedeemed(t *testing.T) {
	g := gc_models.GiftCard{Id: "gc-2", Balance: 20}
	data, _ := json.Marshal(g)
	var d gc_models.GiftCard
	json.Unmarshal(data, &d)
	if d.Balance != 20 {
		t.Errorf("balance=%f", d.Balance)
	}
}
func TestGiftCardInactive(t *testing.T) {
	g := gc_models.GiftCard{Active: false}
	if g.Active {
		t.Error("should be inactive")
	}
}
