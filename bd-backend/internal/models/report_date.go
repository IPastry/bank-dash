package models

// ReportDate represents a row in the report_date table
type ReportDate struct {
    Date CustomDate `json:"date"`
}
