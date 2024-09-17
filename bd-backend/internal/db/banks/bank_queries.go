package banks

import (
	"bd-backend/internal/models"
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// FindBank retrieves a bank by various criteria
func FindBank(ctx context.Context, dbPool *pgxpool.Pool, criteria string, value interface{}) (*models.UserBankInfo, error) {
	var condition string
	var args []interface{}

	// Determine the condition and arguments based on the criteria
	switch criteria {
	case "name":
		condition = "name ILIKE $1"
		args = append(args, value)
	case "routing":
		condition = "routing = $1"
		args = append(args, removeLeadingZeros(value.(string)))
	case "bank_id":
		condition = "bank_id = $1"
		id, err := convertToInt(value)
		if err != nil {
			return nil, err
		}
		args = append(args, id)
	case "cert":
		condition = "cert = $1"
		certStr, err := convertToInt(value)
		if err != nil {
			return nil, err
		}
		args = append(args, certStr)
	default:
		return nil, fmt.Errorf("invalid criteria: %s", criteria)
	}

	// Perform the database query
	var bank models.Bank
	err := dbPool.QueryRow(ctx, fmt.Sprintf("SELECT * FROM banks WHERE %s", condition), args...).Scan(
		&bank.BankID,
		&bank.Cert,
		&bank.OccCharter,
		&bank.OtsDocket,
		&bank.Routing,
		&bank.Name,
		&bank.Address,
		&bank.City,
		&bank.State,
		&bank.Zip,
		&bank.FilingType,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return scanBank(&bank)
}

// GetPeerGroupFromBankID retrieves the peer group for a given bank ID from the lookup table.
func GetPeerGroupFromBankID(ctx context.Context, dbPool *pgxpool.Pool, bankID interface{}) (string, error) {
	var bankIDInt int

	// Convert bankID to integer if it's a string
	switch v := bankID.(type) {
	case int:
		bankIDInt = v
	case string:
		var err error
		bankIDInt, err = strconv.Atoi(v)
		if err != nil {
			return "", fmt.Errorf("invalid bank_id format: %w", err)
		}
	default:
		return "", fmt.Errorf("unsupported bank_id type")
	}

	// Query the lookup table
	var peerGroup string
	err := dbPool.QueryRow(ctx, `
        SELECT peer_group
        FROM lookup_bank_peer
        WHERE bank_id = $1 AND is_default = true
    `, bankIDInt).Scan(&peerGroup)

	if err != nil {
		return "", fmt.Errorf("error querying peer group: %w", err)
	}

	return peerGroup, nil
}
// // GetAllBanks retrieves all banks from the database
// func GetAllBanks(ctx context.Context, dbPool *pgxpool.Pool) ([]models.UserBankInfo, error) {
// 	rows, err := dbPool.Query(ctx, "SELECT * FROM banks")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var banks []models.UserBankInfo
// 	for rows.Next() {
// 		var bank models.Bank
// 		err := rows.Scan(
// 			&bank.BankID,
// 			&bank.Cert,
// 			&bank.OccCharter,
// 			&bank.OtsDocket,
// 			&bank.Routing,
// 			&bank.Name,
// 			&bank.Address,
// 			&bank.City,
// 			&bank.State,
// 			&bank.Zip,
// 			&bank.FilingType,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}

// 		userBankInfo, err := scanBank(&bank)
// 		if err != nil {
// 			return nil, err
// 		}

// 		banks = append(banks, *userBankInfo)
// 	}

// 	return banks, nil
// }
