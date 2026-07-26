package services

import (
	"testing"

	"github.com/nutrixpos/pos/modules/table/dto"
)

func TestGenerateQRCode(t *testing.T) {
	svc := &TableService{}

	qr := svc.GenerateQRCode("abc123", 5)

	if qr != "http://localhost:8080/order?table=5" {
		t.Errorf("expected http://localhost:8080/order?table=5, got %s", qr)
	}
}

func TestGenerateQRCode_TableNumber1(t *testing.T) {
	svc := &TableService{}

	qr := svc.GenerateQRCode("xyz789", 1)

	if qr != "http://localhost:8080/order?table=1" {
		t.Errorf("expected http://localhost:8080/order?table=1, got %s", qr)
	}
}

func TestGenerateQRCode_TableNumber99(t *testing.T) {
	svc := &TableService{}

	qr := svc.GenerateQRCode("def456", 99)

	if qr != "http://localhost:8080/order?table=99" {
		t.Errorf("expected http://localhost:8080/order?table=99, got %s", qr)
	}
}

func TestCreateTable_StructFields(t *testing.T) {
	req := dto.CreateTableRequest{
		Number:   7,
		Name:     "Window Seat",
		Capacity: 4,
		Zone:     "indoor",
		BranchId: "branch-001",
	}

	if req.Number != 7 {
		t.Errorf("expected number 7, got %d", req.Number)
	}
	if req.Name != "Window Seat" {
		t.Errorf("expected name 'Window Seat', got %s", req.Name)
	}
	if req.Capacity != 4 {
		t.Errorf("expected capacity 4, got %d", req.Capacity)
	}
	if req.Zone != "indoor" {
		t.Errorf("expected zone 'indoor', got %s", req.Zone)
	}
	if req.BranchId != "branch-001" {
		t.Errorf("expected branch_id 'branch-001', got %s", req.BranchId)
	}
}

func TestUpdateTableRequest_Fields(t *testing.T) {
	req := dto.UpdateTableRequest{
		Name:     "New Name",
		Capacity: 6,
		Zone:     "outdoor",
		Status:   "occupied",
	}

	if req.Name != "New Name" {
		t.Errorf("expected name 'New Name', got %s", req.Name)
	}
	if req.Capacity != 6 {
		t.Errorf("expected capacity 6, got %d", req.Capacity)
	}
	if req.Zone != "outdoor" {
		t.Errorf("expected zone 'outdoor', got %s", req.Zone)
	}
	if req.Status != "occupied" {
		t.Errorf("expected status 'occupied', got %s", req.Status)
	}
}

func TestUpdateTableRequest_Omitempty(t *testing.T) {
	req := dto.UpdateTableRequest{}

	if req.Name != "" {
		t.Errorf("expected empty name, got %s", req.Name)
	}
	if req.Capacity != 0 {
		t.Errorf("expected zero capacity, got %d", req.Capacity)
	}
	if req.Zone != "" {
		t.Errorf("expected empty zone, got %s", req.Zone)
	}
	if req.Status != "" {
		t.Errorf("expected empty status, got %s", req.Status)
	}
}

func TestGetTablesParams(t *testing.T) {
	params := GetTablesParams{
		PageNumber: 2,
		PageSize:   25,
	}

	if params.PageNumber != 2 {
		t.Errorf("expected page number 2, got %d", params.PageNumber)
	}
	if params.PageSize != 25 {
		t.Errorf("expected page size 25, got %d", params.PageSize)
	}
}

func TestTableService_ConstructorFields(t *testing.T) {
	svc := &TableService{}

	if svc.Logger != nil {
		t.Error("expected nil logger by default")
	}
	if svc.Config.Databases != nil {
		t.Error("expected nil databases by default")
	}
}
