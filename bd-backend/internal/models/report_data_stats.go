package models

// ReportDataStats represents a row in the report_data_stats table
type ReportDataStats struct {
    Date      CustomDate `json:"date"`        // Custom date type for "date"
    PeerGroup string     `json:"peer_group"`  // Corresponds to "peer_group" in the table
    Name      string     `json:"name"`        // Corresponds to "name" in the table
    Metric    string     `json:"metric"`      // Corresponds to "metric" in the table
    Value     *string    `json:"value,omitempty"` // Corresponds to "value" in the table, nullable
    Section   Section    `json:"section"`     // Custom type for "section" enum
}
