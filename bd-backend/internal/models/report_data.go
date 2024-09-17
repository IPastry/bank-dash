package models

// ReportData represents a row in the report_data table
type ReportData struct {
    Date      CustomDate `json:"date"`        // Custom date type for "date"
    BankID    int        `json:"bank_id"`     // Corresponds to "bank_id" in the table
    PeerGroup string     `json:"peer_group"`  // Corresponds to "peer_group" in the table
    Name      string     `json:"name"`        // Corresponds to "name" in the table
    Metric    string     `json:"metric"`      // Corresponds to "metric" in the table
    Value     *string    `json:"value,omitempty"` // Corresponds to "value" in the table, nullable
    Section   Section    `json:"section"`     // Custom type for "section" enum
}
