package models

// UserBankInfo represents the bank information sent to the user
type UserBankInfo struct {
    BankID    int     `json:"bank_id"`
    Cert      *int    `json:"cert,omitempty"`
    Routing   *string `json:"routing,omitempty"`
    Name      *string `json:"name,omitempty"`
    Address   *string `json:"address,omitempty"`
    City      *string `json:"city,omitempty"`
    State     *string `json:"state,omitempty"`
    Zip       *string `json:"zip,omitempty"`
}
