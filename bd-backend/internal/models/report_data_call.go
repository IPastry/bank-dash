package models

// ReportDataCall represents a row in the report_data_call table
type ReportDataCall struct {
    Date      CustomDate `json:"date"`        // Custom date type for "date"
    BankID    int        `json:"bank_id"`     // Corresponds to "bank_id" in the table
    Name      string     `json:"name"`        // Corresponds to "name" in the table
    Metric    string     `json:"metric"`      // Corresponds to "metric" in the table
    Value     *string    `json:"value,omitempty"` // Corresponds to "value" in the table, nullable
    Section   Section    `json:"section"`     // Custom type for "section" enum
}
