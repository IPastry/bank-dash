package data

import (
	"bd-backend/internal/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// FetchRecentReportData retrieves the 5 most recent quarters of data from the report_data table for a specific bank_id
func FetchRecentReportData(ctx context.Context, db *pgxpool.Pool, bankID int) ([]models.ReportData, error) {
	// Query to get the 5 most recent quarters of data for a specific bank_id
	query := `SELECT date, bank_id, peer_group, name, metric, value, section
              FROM public.report_data
              WHERE bank_id = $1
              ORDER BY date DESC
              LIMIT 5`

	// Execute the query
	rows, err := db.Query(ctx, query, bankID)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var reportData []models.ReportData
	for rows.Next() {
		var data models.ReportData
		if err := rows.Scan(&data.Date, &data.BankID, &data.PeerGroup, &data.Name, &data.Metric, &data.Value, &data.Section); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		reportData = append(reportData, data)
	}

	// Check if there were any errors encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return reportData, nil
}
