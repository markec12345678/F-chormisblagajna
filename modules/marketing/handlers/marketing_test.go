package handlers

import (
	"encoding/json"
	"testing"

	mk_models "github.com/nutrixpos/pos/modules/marketing/models"
)

func TestCampaign(t *testing.T) {
	c := mk_models.Campaign{Id: "cm-1", Name: "Summer Sale", Type: "discount", DiscountPct: 20, Active: true}
	data, _ := json.Marshal(c)
	var d mk_models.Campaign
	json.Unmarshal(data, &d)
	if d.Name != "Summer Sale" {
		t.Errorf("name=%s", d.Name)
	}
}
func TestCampaignInactive(t *testing.T) {
	c := mk_models.Campaign{Active: false}
	if c.Active {
		t.Error("should be inactive")
	}
}
func TestCampaignStats(t *testing.T) {
	s := mk_models.CampaignStats{Id: "cm-1", TotalSent: 100, TotalRedeemed: 15, Revenue: 450}
	data, _ := json.Marshal(s)
	var d mk_models.CampaignStats
	json.Unmarshal(data, &d)
	if d.Revenue != 450 {
		t.Errorf("revenue=%f", d.Revenue)
	}
}
