package models

import "time"

type ClockEntry struct {
	Id         string     `json:"id" bson:"id" mapstructure:"id"`
	EmployeeId string     `json:"employee_id" bson:"employee_id" mapstructure:"employee_id"`
	EmployeeName string   `json:"employee_name" bson:"employee_name" mapstructure:"employee_name"`
	ClockIn    time.Time  `json:"clock_in" bson:"clock_in" mapstructure:"clock_in"`
	ClockOut   *time.Time `json:"clock_out" bson:"clock_out" mapstructure:"clock_out"`
	Status     string     `json:"status" bson:"status" mapstructure:"status"` // active, completed
	HoursWorked float64   `json:"hours_worked" bson:"hours_worked" mapstructure:"hours_worked"`
	Notes      string     `json:"notes" bson:"notes" mapstructure:"notes"`
}

type TimeClockSummary struct {
	EmployeeId     string  `json:"employee_id"`
	EmployeeName   string  `json:"employee_name"`
	TotalHours     float64 `json:"total_hours"`
	ShiftCount     int     `json:"shift_count"`
	AvgHoursPerShift float64 `json:"avg_hours_per_shift"`
	OvertimeHours  float64 `json:"overtime_hours"`
}

type TimeClockDashboard struct {
	CurrentlyClockedIn []ClockEntry        `json:"currently_clocked_in"`
	TodaySummary       []TimeClockSummary  `json:"today_summary"`
	WeekSummary        []TimeClockSummary  `json:"week_summary"`
}
