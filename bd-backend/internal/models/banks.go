package models

// Bank represents a row in the banks table
type Bank struct {
    BankID       int    `json:"bank_id"`
    Cert         *int   `json:"cert,omitempty"`
    OccCharter   *int   `json:"occ_charter,omitempty"`
    OtsDocket    *int   `json:"ots_docket,omitempty"`
    Routing      *string `json:"routing,omitempty"`
    Name         *string `json:"name,omitempty"`
    Address      *string `json:"address,omitempty"`
    City         *string `json:"city,omitempty"`
    State        *string `json:"state,omitempty"`
    Zip          *string `json:"zip,omitempty"`
    FilingType   *string `json:"filing_type,omitempty"`
}
