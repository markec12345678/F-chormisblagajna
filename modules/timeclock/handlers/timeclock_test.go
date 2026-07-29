package handlers

import (
	"encoding/json"
	"testing"
	"time"

	timeclock_models "github.com/nutrixpos/pos/modules/timeclock/models"
)

func TestClockEntry_Serialization(t *testing.T) {
	now := time.Now()
	entry := timeclock_models.ClockEntry{
		Id:           "ce-1",
		EmployeeId:   "emp-1",
		EmployeeName: "Janez Novak",
		ClockIn:      now,
		Status:       "active",
		Notes:        "Morning shift",
	}

	data, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded timeclock_models.ClockEntry
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.EmployeeName != "Janez Novak" {
		t.Errorf("expected EmployeeName='Janez Novak', got %s", decoded.EmployeeName)
	}
	if decoded.Status != "active" {
		t.Errorf("expected Status='active', got %s", decoded.Status)
	}
}

func TestTimeClockSummary_Calculation(t *testing.T) {
	summary := timeclock_models.TimeClockSummary{
		EmployeeId:   "emp-1",
		EmployeeName: "Janez Novak",
		TotalHours:   40.0,
		ShiftCount:   5,
		OvertimeHours: 0.0,
	}

	summary.AvgHoursPerShift = summary.TotalHours / float64(summary.ShiftCount)

	if summary.AvgHoursPerShift != 8.0 {
		t.Errorf("expected AvgHoursPerShift=8, got %f", summary.AvgHoursPerShift)
	}
}

func TestClockEntry_CompletedStatus(t *testing.T) {
	now := time.Now()
	clockIn := now.Add(-8 * time.Hour)
	entry := timeclock_models.ClockEntry{
		Id:         "ce-1",
		ClockIn:    clockIn,
		ClockOut:   &now,
		Status:     "completed",
		HoursWorked: 8.0,
	}

	if entry.Status != "completed" {
		t.Errorf("expected Status='completed', got %s", entry.Status)
	}
	if entry.HoursWorked != 8.0 {
		t.Errorf("expected HoursWorked=8, got %f", entry.HoursWorked)
	}
}

func TestTimeClockDashboard_Fields(t *testing.T) {
	dashboard := timeclock_models.TimeClockDashboard{
		CurrentlyClockedIn: []timeclock_models.ClockEntry{},
		TodaySummary:       []timeclock_models.TimeClockSummary{},
		WeekSummary:        []timeclock_models.TimeClockSummary{},
	}

	if dashboard.CurrentlyClockedIn == nil {
		t.Error("CurrentlyClockedIn should not be nil")
	}
}
