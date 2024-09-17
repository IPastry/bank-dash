package db

import (
	"bd-backend/internal/auth"
	"bd-backend/internal/db/users"
	"bd-backend/internal/models"
	"bd-backend/internal/utils"
	"context"
	"os"
	"testing"
	// "time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var dbPool *pgxpool.Pool

func TestMain(m *testing.M) {
	err := os.Chdir("D:/projects/go/bd-backend")
	if err != nil {
		panic("Error changing directory: " + err.Error())
	}

	err = godotenv.Load()
	if err != nil {
		panic("Error loading .env file: " + err.Error())
	}

	dbPool, err = ConnectDB(context.Background())
	if err != nil {
		panic("Unable to connect to database: " + err.Error())
	}

	auth.Init()

	// Run the tests
	code := m.Run()

	// // Clean up
	// err = CleanupExpiredTokens(dbPool)
	// if err != nil {
	// 	panic("Error cleaning up expired tokens: " + err.Error())
	// }

	CloseDB(dbPool)
	os.Exit(code)
}

func mockContextWithUser(userID int, role models.Role) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, utils.UserIDKey, userID)
	ctx = context.WithValue(ctx, utils.RoleKey, role)
	return ctx
}

func TestUserWorkflow(t *testing.T) {
    // Transaction 1: Create an elevated user
    tx1, err := dbPool.Begin(context.Background())
    if err != nil {
        t.Fatalf("Error starting transaction 1: %v", err)
    }
    defer tx1.Rollback(context.Background()) // Ensure rollback on error or completion

    managerEmail := "manager@example.com"
    managerPassword := "managerPassword1!"
    err = users.CreateElevatedUser(context.Background(), tx1, managerEmail, managerPassword)
    if err != nil {
        t.Fatalf("Error creating elevated user: %v", err)
    }

    // Commit the transaction for creating the elevated user
    if err := tx1.Commit(context.Background()); err != nil {
        t.Fatalf("Error committing transaction 1: %v", err)
    }

    // Transaction 2: Retrieve manager ID and update profile and bank ID
    tx2, err := dbPool.Begin(context.Background())
    if err != nil {
        t.Fatalf("Error starting transaction 2: %v", err)
    }
    defer tx2.Rollback(context.Background()) // Ensure rollback on error or completion

    // Retrieve manager ID
    var managerID int
    row := tx2.QueryRow(context.Background(), `SELECT user_id FROM users WHERE email = $1`, managerEmail)
    err = row.Scan(&managerID)
    if err != nil {
        t.Fatalf("Error retrieving manager user ID: %v", err)
    }

    // Set mock context with manager ID
    mockCtx := mockContextWithUser(managerID, models.Elevated)

    // Update manager's profile and bank ID
    managerBankID := 37
    firstName := "John"
    lastName := "Doe"
    phoneNumber := "1234567890"
    err = users.UpdateProfile(mockCtx, tx2, &firstName, &lastName, &phoneNumber)
    if err != nil {
        t.Fatalf("Error updating manager profile: %v", err)
    }
    err = users.UpdateUserBankID(mockCtx, tx2, managerBankID)
    if err != nil {
        t.Fatalf("Error updating manager bank ID: %v", err)
    }

    // Commit the transaction for profile and bank ID updates
    if err := tx2.Commit(context.Background()); err != nil {
        t.Fatalf("Error committing transaction 2: %v", err)
    }

    // Transaction 3: Create shared account
    tx3, err := dbPool.Begin(mockCtx) // Use mockCtx here
    if err != nil {
        t.Fatalf("Error starting transaction 3: %v", err)
    }
    defer tx3.Rollback(mockCtx) // Ensure rollback on error or completion

    sharedEmail := "shared@example.com"
    sharedPassword := "sharedPassword1!"
    phone := "9876543210"
    sharedPhoneNumber := &phone
    err = users.CreateSharedAccount(mockCtx, tx3, sharedEmail, sharedPassword, sharedPhoneNumber)
    if err != nil {
        t.Fatalf("Error creating shared account: %v", err)
    }

    // Commit the transaction for shared account creation
    if err := tx3.Commit(mockCtx); err != nil {
        t.Fatalf("Error committing transaction 3: %v", err)
    }

    // Transaction 4: Verify shared account creation
    tx4, err := dbPool.Begin(mockCtx) // Use mockCtx here
    if err != nil {
        t.Fatalf("Error starting transaction 4: %v", err)
    }
    defer tx4.Rollback(mockCtx) // Ensure rollback on error or completion

    row = tx4.QueryRow(mockCtx, `SELECT email, manager_id, phone_number, is_verified FROM users WHERE email = $1`, sharedEmail)
    var email string
    var retrievedManagerID int
    var retrievedPhoneNumber string
    var isVerified bool

    err = row.Scan(&email, &retrievedManagerID,  &retrievedPhoneNumber, &isVerified)
    if err != nil {
        t.Fatalf("Expected user to be created, got %v", err)
    }

    assert.Equal(t, sharedEmail, email)
    assert.Equal(t, managerID, retrievedManagerID)
    assert.Equal(t, phone, retrievedPhoneNumber)
    assert.True(t, isVerified)

    // Commit the transaction for verifying shared account creation
    if err := tx4.Commit(mockCtx); err != nil {
        t.Fatalf("Error committing transaction 4: %v", err)
    }
}


// func TestGenerateAndParseToken(t *testing.T) {
// 	// Generate a token
// 	userID := 123
// 	role := models.Elevated // assuming Elevated is a role alias
// 	token, err := auth.GenerateToken(userID, role)
// 	if err != nil {
// 		t.Fatalf("Error generating token: %v", err)
// 	}

// 	// Parse the token
// 	_, claims, err := auth.ParseToken(token)
// 	if err != nil {
// 		t.Fatalf("Error parsing token: %v", err)
// 	}

// 	// Extract and verify claims
// 	parsedUserIDFloat, ok := claims["user_id"].(float64)
// 	if !ok {
// 		t.Fatalf("Error retrieving user ID from claims")
// 	}
// 	parsedUserID := int(parsedUserIDFloat)

// 	parsedRoleStr, ok := claims["role"].(string)
// 	if !ok {
// 		t.Fatalf("Error retrieving role from claims")
// 	}
// 	parsedRole := models.Role(parsedRoleStr)

// 	// Assertions
// 	assert.Equal(t, userID, parsedUserID)
// 	assert.Equal(t, role, parsedRole) // Role comparison
// }

// func TestGenerateAndParseRefreshToken(t *testing.T) {
// 	userID := 12345

// 	// Generate refresh token
// 	tokenStr, err := auth.GenerateRefreshToken(userID)
// 	assert.NoError(t, err, "Error generating refresh token")

// 	// Parse refresh token
// 	_, claims, err := auth.ParseRefreshToken(tokenStr)
// 	assert.NoError(t, err, "Error parsing refresh token")

// 	// Validate claims
// 	assert.Equal(t, float64(userID), claims["user_id"], "User ID does not match")

// 	// Check expiration
// 	exp := claims["exp"].(float64)
// 	expTime := time.Unix(int64(exp), 0)
// 	assert.True(t, time.Now().Before(expTime), "Refresh token has expired")
// }

// func TestCreateElevatedUser(t *testing.T) {
//     ctx := context.Background() // Create a context for this test

//     // Begin a new transaction for the test
//     tx, err := dbPool.Begin(ctx)
//     if err != nil {
//         t.Fatalf("Error starting transaction: %v", err)
//     }
//     defer tx.Rollback(ctx) // Ensure the transaction is rolled back after the test

//     email := "elevated@example.com"
//     password := "elevatedPassword1!"

//     err = users.CreateElevatedUser(ctx, tx, email, password)
//     if err != nil {
//         t.Fatalf("Error creating elevated user: %v", err)
//     }

//     // Commit the transaction to make sure the data is persisted
//     if err := tx.Commit(ctx); err != nil {
//         t.Fatalf("Error committing transaction: %v", err)
//     }

//     // Verify that the user was created
//     user, err := users.GetUserByEmail(ctx, dbPool, email)
//     if err != nil {
//         t.Fatalf("Error retrieving user: %v", err)
//     }
//     if user["email"] != email {
//         t.Fatalf("Expected email %v, got %v", email, user["email"])
//     }
// }

// func TestCreateSharedAccount(t *testing.T) {
//     ctx := context.Background() // Create a context for this test

//     // Begin a new transaction for the test
//     tx, err := dbPool.Begin(ctx)
//     if err != nil {
//         t.Fatalf("Error starting transaction: %v", err)
//     }
//     defer tx.Rollback(ctx) // Ensure the transaction is rolled back after the test

//     // Create an elevated user (manager)
//     managerEmail := "manager@example.com"
//     managerPassword := "managerPassword1!"
//     err = users.CreateElevatedUser(ctx, tx, managerEmail, managerPassword)
//     if err != nil {
//         t.Fatalf("Error creating elevated user: %v", err)
//     }

//     // Update manager's profile to include bank_id
//     managerBankID := 37
//     firstName := "John"
//     lastName := "Doe"
//     phoneNumber := "1234567890"
//     err = users.UpdateProfile(ctx, tx, managerEmail, &firstName, &lastName, &phoneNumber)
//     if err != nil {
//         t.Fatalf("Error updating manager profile: %v", err)
//     }

//     // Retrieve manager ID and bank ID
//     var managerID int
//     var retrievedManagerBankID int
//     row := tx.QueryRow(ctx, `SELECT user_id, bank_id FROM users WHERE email = $1`, managerEmail)
//     err = row.Scan(&managerID, &retrievedManagerBankID)
//     if err != nil {
//         t.Fatalf("Error retrieving manager user ID and bank ID: %v", err)
//     }

//     // Create shared account
//     sharedEmail := "shared@example.com"
//     sharedPassword := "sharedPassword1!"
//     phone := "1234567890"
//     sharedPhoneNumber := &phone
//     err = users.CreateSharedAccount(ctx, tx, sharedEmail, sharedPassword, sharedPhoneNumber)
//     if err != nil {
//         t.Fatalf("Error creating shared account: %v", err)
//     }

//     // Commit the transaction to make sure the data is persisted
//     if err := tx.Commit(ctx); err != nil {
//         t.Fatalf("Error committing transaction: %v", err)
//     }

//     // Verify shared account creation
//     row = dbPool.QueryRow(ctx, `SELECT email, manager_id, bank_id, phone_number, is_verified FROM users WHERE email = $1`, sharedEmail)
//     var email string
//     var retrievedManagerID int
//     var retrievedBankID int
//     var retrievedPhoneNumber string
//     var isVerified bool

//     err = row.Scan(&email, &retrievedManagerID, &retrievedBankID, &retrievedPhoneNumber, &isVerified)
//     if err != nil {
//         t.Fatalf("Expected user to be created, got %v", err)
//     }

//     if email != sharedEmail {
//         t.Fatalf("Expected email %v, got %v", sharedEmail, email)
//     }
//     if retrievedManagerID != managerID {
//         t.Fatalf("Expected manager ID %v, got %v", managerID, retrievedManagerID)
//     }
//     if retrievedBankID != retrievedManagerBankID {
//         t.Fatalf("Expected bank ID %v, got %v", retrievedManagerBankID, retrievedBankID)
//     }
//     if retrievedPhoneNumber != phone {
//         t.Fatalf("Expected phone number %v, got %v", phone, retrievedPhoneNumber)
//     }
//     if !isVerified {
//         t.Fatalf("Expected is_verified to be true")
//     }
// }

// func TestInitialProfileUpdate(t *testing.T) {
//     ctx := context.Background() // Create a context for this test
//     email := "user@example.com"
//     password := "Password!2"

//     // Ensure the elevated user is created successfully
//     err := CreateElevatedUser(ctx, dbPool, email, password)
//     if err != nil {
//         t.Fatalf("Error creating elevated user: %v", err)
//     }

//     // Prepare profile update data
//     firstName := "John"
//     lastName := "Doe"
//     phoneNumber := "0987654321"
//     bankID := 37

//     // Perform the profile update
//     err = InitialProfileUpdate(ctx, dbPool, email, &firstName, &lastName, &phoneNumber, &bankID)
//     if err != nil {
//         t.Fatalf("Error updating profile: %v", err)
//     }

//     // Retrieve and validate updated user info
//     user, err := GetUserInfo(ctx, dbPool, email)
//     if err != nil {
//         t.Fatalf("Error retrieving user info: %v", err)
//     }

//     // Validate user information
//     if user["first_name"] != firstName {
//         t.Fatalf("Expected first name %v, got %v", firstName, user["first_name"])
//     }
//     if user["last_name"] != lastName {
//         t.Fatalf("Expected last name %v, got %v", lastName, user["last_name"])
//     }
//     if user["phone_number"] != phoneNumber {
//         t.Fatalf("Expected phone number %v, got %v", phoneNumber, user["phone_number"])
//     }
//     if user["bank_id"] != bankID {
//         t.Fatalf("Expected bank ID %v, got %v", bankID, user["bank_id"])
//     }
// }

// func TestUpdateBasicProfile(t *testing.T) {
//     ctx := context.Background() // Create a context for this test
//     email := "user@example.com"
//     firstName := "Jane"
//     lastName := "Doe"
//     phoneNumber := "1122334455"

//     err := UpdateBasicProfile(ctx, dbPool, email, &firstName, &lastName, &phoneNumber)
//     if err != nil {
//         t.Fatalf("Error updating basic profile: %v", err)
//     }

//     user, err := GetUserInfo(ctx, dbPool, email)
//     if err != nil {
//         t.Fatalf("Error retrieving user info: %v", err)
//     }
//     if user["first_name"] != firstName {
//         t.Fatalf("Expected first name %v, got %v", firstName, user["first_name"])
//     }
//     if user["last_name"] != lastName {
//         t.Fatalf("Expected last name %v, got %v", lastName, user["last_name"])
//     }
//     if user["phone_number"] != phoneNumber {
//         t.Fatalf("Expected phone number %v, got %v", phoneNumber, user["phone_number"])
//     }
// }

// func TestUpdateEmail(t *testing.T) {
//     ctx := context.Background() // Create a context for this test
//     oldEmail := "user@example.com"
//     newEmail := "newuser@example.com"

//     err := UpdateEmail(ctx, dbPool, oldEmail, newEmail)
//     if err != nil {
//         t.Fatalf("Error updating email: %v", err)
//     }

//     user, err := GetUserByEmail(ctx, dbPool, newEmail)
//     if err != nil {
//         t.Fatalf("Error retrieving new email user: %v", err)
//     }
//     if user["email"] != newEmail {
//         t.Fatalf("Expected email %v, got %v", newEmail, user["email"])
//     }
// }

// func TestUpdatePassword(t *testing.T) {
//     ctx := context.Background() // Create a context for this test
//     email := "user@example.com"
//     newPassword := "newPassword1!"

//     err := UpdatePassword(ctx, dbPool, email, newPassword)
//     if err != nil {
//         t.Fatalf("Error updating password: %v", err)
//     }

//     // Verify that password is updated by attempting to login or similar logic
//     // You might need to implement additional logic to verify password update if applicable.
// }

// func TestGetUserByEmail(t *testing.T) {
//     ctx := context.Background() // Create a context for this test
//     email := "user@example.com"
//     user, err := GetUserByEmail(ctx, dbPool, email)
//     if err != nil {
//         t.Fatalf("Error retrieving user by email: %v", err)
//     }
//     if user["email"] != email {
//         t.Fatalf("Expected email %v, got %v", email, user["email"])
//     }
// }

// func TestGetUserInfo(t *testing.T) {
//     ctx := context.Background() // Create a context for this test
//     email := "user@example.com"
//     user, err := GetUserInfo(ctx, dbPool, email)
//     if err != nil {
//         t.Fatalf("Error retrieving user info: %v", err)
//     }
//     if user["email"] != email {
//         t.Fatalf("Expected email %v, got %v", email, user["email"])
//     }
// }

// func TestHashPassword(t *testing.T) {
// 	password := "testPassword"
// 	hashedPassword, err := utils.HashPassword(password)

// 	if err != nil {
// 		t.Fatalf("Expected no error, got %v", err)
// 	}
// 	if hashedPassword == "" {
// 		t.Fatal("Expected non-empty hashed password")
// 	}

// 	// To further test, you might want to verify that the hashed password matches the original password
// 	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
// 	if err != nil {
// 		t.Fatalf("Password comparison failed: %v", err)
// 	}
// }

// func CleanupTokens() error {
// 	_, err := dbPool.Exec(context.Background(), `DELETE FROM email_verifications`)
// 	return err
// }

// func TestSaveTokenInDB(t *testing.T) {
// 	err := CleanupTokens()
// 	assert.NoError(t, err, "Error cleaning up tokens before test")

// 	email := "test@example.com"
// 	token := "test-token"
// 	expirationTime := time.Now().Add(24 * time.Hour)

// 	err = SaveTokenInDB(dbPool, email, token, expirationTime)
// 	assert.NoError(t, err, "Error saving token to database")

// 	var savedToken string
// 	var savedExpirationTime time.Time
// 	err = dbPool.QueryRow(context.Background(), `
//         SELECT token, expiration_time
//         FROM email_verifications
//         WHERE email = $1`, email).Scan(&savedToken, &savedExpirationTime)

// 	assert.NoError(t, err, "Error retrieving token from database")
// 	assert.Equal(t, token, savedToken, "Token mismatch")
// 	assert.WithinDuration(t, expirationTime, savedExpirationTime, time.Second, "Expiration time mismatch")
// }

// func TestVerifyToken(t *testing.T) {
// 	err := CleanupTokens()
// 	assert.NoError(t, err, "Error cleaning up tokens before test")

// 	email := "test@example.com"
// 	token := "test-token"
// 	expirationTime := time.Now().Add(24 * time.Hour)

// 	err = SaveTokenInDB(dbPool, email, token, expirationTime)
// 	assert.NoError(t, err, "Error saving token to database")

// 	retrievedEmail, valid, err := VerifyToken(dbPool, token)
// 	assert.NoError(t, err, "Error verifying token")
// 	assert.True(t, valid, "Token should be valid")
// 	assert.Equal(t, email, retrievedEmail, "Email mismatch")

// 	// Test with an expired token
// 	expiredToken := "expired-token"
// 	expiredTime := time.Now().Add(-24 * time.Hour)
// 	err = SaveTokenInDB(dbPool, email, expiredToken, expiredTime)
// 	assert.NoError(t, err, "Error saving expired token to database")

// 	_, valid, err = VerifyToken(dbPool, expiredToken)
// 	assert.NoError(t, err, "Error verifying expired token")
// 	assert.False(t, valid, "Expired token should be invalid")

// 	// Test with a non-existing token
// 	nonExistingToken := "non-existing-token"
// 	_, valid, err = VerifyToken(dbPool, nonExistingToken)
// 	assert.NoError(t, err, "Error verifying non-existing token")
// 	assert.False(t, valid, "Non-existing token should be invalid")
// }

// func TestCleanupExpiredTokens(t *testing.T) {
// 	err := CleanupTokens()
// 	assert.NoError(t, err, "Error cleaning up tokens before test")

// 	email := "cleanup@example.com"
// 	token := "cleanup-token"
// 	expirationTime := time.Now().Add(-24 * time.Hour) // Expired token

// 	err = SaveTokenInDB(dbPool, email, token, expirationTime)
// 	assert.NoError(t, err, "Error saving expired token to database")

// 	// Ensure the expired token is in the database
// 	var count int
// 	err = dbPool.QueryRow(context.Background(), `
//         SELECT COUNT(*)
//         FROM email_verifications
//         WHERE email = $1`, email).Scan(&count)
// 	assert.NoError(t, err, "Error counting tokens")
// 	assert.Equal(t, 1, count, "Expired token should be present")

// 	err = CleanupExpiredTokens(dbPool)
// 	assert.NoError(t, err, "Error cleaning up expired tokens")

// 	// Verify the expired token has been removed
// 	err = dbPool.QueryRow(context.Background(), `
//         SELECT COUNT(*)
//         FROM email_verifications
//         WHERE email = $1`, email).Scan(&count)
// 	assert.NoError(t, err, "Error counting tokens")
// 	assert.Equal(t, 0, count, "Expired token should be removed")
// }

// func TestFetchRecentReportData(t *testing.T) {
//     ctx := context.Background()
//     bankID := 999935
//     reportData, err := FetchRecentReportData(ctx, dbPool, bankID)
//     if err != nil {
//         t.Fatalf("Error fetching report data: %v", err)
//     }

//     // Assert that the reportData is not empty (adjust based on your data)
//     assert.NotEmpty(t, reportData, "Expected non-empty report data for bank_id 999935")

//     // Add additional assertions to verify data content if needed
//     for _, data := range reportData {
//         assert.Equal(t, bankID, data.BankID, "Expected bank_id to be 999935")
//     }
// }

// func TestFindBankByName(t *testing.T) {
//     bank, err := FindBankByName(context.Background(), dbPool, "BANK OF HANCOCK COUNTY")
//     assert.NoError(t, err)
//     assert.NotNil(t, bank)
//     assert.Equal(t, 37, bank.BankID)
//     assert.Equal(t, 10057, *bank.Cert)
//     assert.Equal(t, "061107146", *bank.Routing)
//     assert.Equal(t, "Bank Of Hancock County", *bank.Name)
//     assert.Equal(t, "12855 Broad Street", *bank.Address)
//     assert.Equal(t, "Sparta", *bank.City)
//     assert.Equal(t, "GA", *bank.State)
//     assert.Equal(t, "31087", *bank.Zip)
// }

// func TestFindBankByRouting(t *testing.T) {
//     bank, err := FindBankByRouting(context.Background(), dbPool, "61107146")
//     assert.NoError(t, err)
//     assert.NotNil(t, bank)
//     assert.Equal(t, 37, bank.BankID)
//     assert.Equal(t, 10057, *bank.Cert)
//     assert.Equal(t, "061107146", *bank.Routing)
//     assert.Equal(t, "Bank Of Hancock County", *bank.Name)
//     assert.Equal(t, "12855 Broad Street", *bank.Address)
//     assert.Equal(t, "Sparta", *bank.City)
//     assert.Equal(t, "GA", *bank.State)
//     assert.Equal(t, "31087", *bank.Zip)
// }

// func TestFindBankByID(t *testing.T) {
//     bank, err := FindBankByID(context.Background(), dbPool, 37)
//     assert.NoError(t, err)
//     assert.NotNil(t, bank)
//     assert.Equal(t, 37, bank.BankID)
//     assert.Equal(t, 10057, *bank.Cert)
//     assert.Equal(t, "061107146", *bank.Routing)
//     assert.Equal(t, "Bank Of Hancock County", *bank.Name)
//     assert.Equal(t, "12855 Broad Street", *bank.Address)
//     assert.Equal(t, "Sparta", *bank.City)
//     assert.Equal(t, "GA", *bank.State)
//     assert.Equal(t, "31087", *bank.Zip)

//     // Test with string ID
//     bank, err = FindBankByID(context.Background(), dbPool, "37")
//     assert.NoError(t, err)
//     assert.NotNil(t, bank)
//     assert.Equal(t, 37, bank.BankID)
//     assert.Equal(t, 10057, *bank.Cert)
//     assert.Equal(t, "061107146", *bank.Routing)
//     assert.Equal(t, "Bank Of Hancock County", *bank.Name)
//     assert.Equal(t, "12855 Broad Street", *bank.Address)
//     assert.Equal(t, "Sparta", *bank.City)
//     assert.Equal(t, "GA", *bank.State)
//     assert.Equal(t, "31087", *bank.Zip)
// }

// func TestFindBankByCert(t *testing.T) {
//     bank, err := FindBankByCert(context.Background(), dbPool, 10057)
//     assert.NoError(t, err)
//     assert.NotNil(t, bank)
//     assert.Equal(t, 37, bank.BankID)
//     assert.Equal(t, 10057, *bank.Cert)
//     assert.Equal(t, "061107146", *bank.Routing)
//     assert.Equal(t, "Bank Of Hancock County", *bank.Name)
//     assert.Equal(t, "12855 Broad Street", *bank.Address)
//     assert.Equal(t, "Sparta", *bank.City)
//     assert.Equal(t, "GA", *bank.State)
//     assert.Equal(t, "31087", *bank.Zip)

//     // Test with string cert
//     bank, err = FindBankByCert(context.Background(), dbPool, "10057")
//     assert.NoError(t, err)
//     assert.NotNil(t, bank)
//     assert.Equal(t, 37, bank.BankID)
//     assert.Equal(t, 10057, *bank.Cert)
//     assert.Equal(t, "061107146", *bank.Routing)
//     assert.Equal(t, "Bank Of Hancock County", *bank.Name)
//     assert.Equal(t, "12855 Broad Street", *bank.Address)
//     assert.Equal(t, "Sparta", *bank.City)
//     assert.Equal(t, "GA", *bank.State)
//     assert.Equal(t, "31087", *bank.Zip)
// }

// func TestGetAllBanks(t *testing.T) {
//     Comment out or remove this test if you decide not to use GetAllBanks
//     banks, err := GetAllBanks(context.Background(), dbPool)
//     assert.NoError(t, err)
//     assert.NotEmpty(t, banks, "Expected non-empty bank list")

//     // Check the content of the first bank
//     bank1 := banks[0]
//     assert.Equal(t, 37, bank1.BankID)
//     assert.Equal(t, 10057, *bank1.Cert)
//     assert.Equal(t, "061107146", *bank1.Routing)
//     assert.Equal(t, "Bank Of Hancock County", *bank1.Name)
//     assert.Equal(t, "12855 Broad Street", *bank1.Address)
//     assert.Equal(t, "Sparta", *bank1.City)
//     assert.Equal(t, "GA", *bank1.State)
//     assert.Equal(t, "31087", *bank1.Zip)

//     // Check the content of the second bank
//     bank2 := banks[1]
//     assert.Equal(t, 242, bank2.BankID)
//     assert.Equal(t, 3850, *bank2.Cert)
//     assert.Equal(t, "081220537", *bank2.Routing)
//     assert.Equal(t, "First Community Bank Xenia-Flora", *bank2.Name)
//     assert.Equal(t, "260 Front Street", *bank2.Address)
//     assert.Equal(t, "Xenia", *bank2.City)
//     assert.Equal(t, "IL", *bank2.State)
//     assert.Equal(t, "62899", *bank2.Zip)
// }