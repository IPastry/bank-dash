package models

type LookupBankPeer struct {
    BankID    int    `json:"bank_id"`                      // Corresponds to "bank_id" in the table
    PeerGroup string `json:"peer_group"`                  // Corresponds to "peer_group" in the table
    IsDefault bool   `json:"is_default" default:"false"`  // Corresponds to "is_default" in the table
}
