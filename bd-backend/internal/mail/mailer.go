package mail

import (
    "context"
    "github.com/mailgun/mailgun-go/v4"
    "github.com/joho/godotenv"
    "os"
    "log"
)

// Load environment variables from .env file
func init() {
    err := godotenv.Load()
    if err != nil {
        log.Println("Warning: Error loading .env file, proceeding with system environment variables.")
    }
}

// SendVerificationEmail sends a verification email to the user
func SendVerificationEmail(ctx context.Context, email, token string) error {
    mg := mailgun.NewMailgun(
        os.Getenv("MAILGUN_DOMAIN"),
        os.Getenv("MAILGUN_API_KEY"),
    )

    verificationURL := os.Getenv("VERIFICATION_URL")
    if verificationURL == "" {
        verificationURL = "https://DP.com/verify-email"
    }

    message := mg.NewMessage(
        "no-reply@" + os.Getenv("MAILGUN_DOMAIN"),   // From email address
        "Please verify your email",                 // Subject
        "Please verify your email by clicking the following link: " + verificationURL + "?token=" + token, // Body
        email,                                      // Recipient email
    )

    // Send the email using the passed context
    _, _, err := mg.Send(ctx, message)
    return err
}
