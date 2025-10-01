package verification

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

func GenerateVerificationToken() (token, hash string) {
    bytes := make([]byte, 32)
    rand.Read(bytes)
    token = hex.EncodeToString(bytes)

    hasher := sha256.New()
    hasher.Write([]byte(token))
    hash = hex.EncodeToString(hasher.Sum(nil))

    return token, hash
}

func SendVerificationEmail(email string, token string) error {
    client := resend.NewClient(os.Getenv("RESEND_API_KEY"))

    verifyLink := fmt.Sprintf("%s/verify?token=%s", os.Getenv("FRONTEND_URL"), token)

    html := fmt.Sprintf(`
    <!DOCTYPE html>
    <html>
    <head>
        <link href="https://fonts.googleapis.com/css2?family=Finger+Paint&display=swap" rel="stylesheet">
    </head>
    <body style="margin: 0; padding: 0; background-color: #1E1E24; font-family: 'Finger Paint', cursive;">
        <table width="100%%" cellpadding="0" cellspacing="0" style="background-color: #1E1E24; padding: 40px 0;">
            <tr>
                <td align="center">
                    <table width="600" cellpadding="0" cellspacing="0" style="background-color: #FFF3E5; border-radius: 5px; padding: 40px;">
                        <tr>
                            <td align="center">
                                <h1 style="color: #1E1E24; font-size: 32px; margin-bottom: 20px;">
                                    Welcome to Silly Mini Games!
                                </h1>
                                <p style="color: #1E1E24; font-size: 18px; margin-bottom: 30px;">
                                    Click below to verify your email and start playing online!
                                </p>
                                <a href="%s" style="
                                    background-color: #92140C;
                                    color: #FFF8F0;
                                    padding: 15px 40px;
                                    text-decoration: none;
                                    border-radius: 5px;
                                    font-size: 20px;
                                    display: inline-block;
                                    margin-bottom: 30px;
                                ">
                                    Verify Email & Sign In
                                </a>
                                <p style="color: #666; font-size: 14px; margin-top: 30px;">
                                    Or copy this link:<br>
                                    <span style="color: #92140C;">%s</span>
                                </p>
                                <p style="color: #666; font-size: 12px; margin-top: 20px;">
                                    This link expires in 120 minutes
                                </p>
                            </td>
                        </tr>
                    </table>
                </td>
            </tr>
        </table>
    </body>
    </html>
    `, verifyLink, verifyLink)

    _, err := client.Emails.Send(&resend.SendEmailRequest{
        From: fmt.Sprintf("Silly Mini Games <noreply@%s>", os.Getenv("EMAIL_DOMAIN")),
        To: []string{email},
        Subject: "Verify your account",
        Html: html,
    })

    return err
}
