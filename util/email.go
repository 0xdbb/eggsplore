package util

import (
	"fmt"
	"github.com/0xdbb/eggsplore/internal/config"
	db "github.com/0xdbb/eggsplore/internal/database/sqlc"
	"github.com/0xdbb/eggsplore/token"

	"github.com/resend/resend-go/v2"
)

// SendVerificationEmail sends the 2FA OTP email using Resend
func SendVerificationEmail(email, otp, key string) error {
	client := resend.NewClient(key)

	params := &resend.SendEmailRequest{
		From:    "GADE <noreply@fvlcon.org>",
		To:      []string{email},
		Html:    fmt.Sprintf("<p>Your 2FA code is: <strong>%s</strong></p>", otp),
		Subject: "Your 2FA Code",
	}

	_, err := client.Emails.Send(params)
	return err
}

func SendGenericEmail(email, message, subject, key string) error {
	client := resend.NewClient(key)

	params := &resend.SendEmailRequest{
		From:    "GADE <noreply@fvlcon.org>",
		To:      []string{email},
		Html:    fmt.Sprintf("<p>%s</p>", message),
		Subject: subject,
	}

	_, err := client.Emails.Send(params)
	return err
}

// SendInvitationEmail sends an invitation email with a setup link to the user
func SendInvitationEmail(email, deparmtment, setupToken, resendApiKey string, appDomain string) error {
	client := resend.NewClient(resendApiKey)

	// Construct the setup link
	setupLink := fmt.Sprintf("%s/setup?token=%s", appDomain, setupToken)

	// Email content
	subject := "Complete Your GADE Account Setup"
	htmlContent := fmt.Sprintf(`
        <p>Dear User,</p>
        <p>Welcome to GADE! An account has been created for you under the %s. To complete your account setup, please set your password by clicking the link below:</p>
        <p><a href="%s">Complete Account Setup</a></p>
        <p>This link will expire in 48 hours. If the link has expired, please contact your administrator to request a new invitation.</p>
        <p>If you did not expect this email, please ignore it or contact <a href="mailto:info@blvcksapphire.com">info@blvcksapphire.com</a> for assistance.</p>
        <p>Best regards,<br>The GADE Team</p>
    `, deparmtment, setupLink)

	params := &resend.SendEmailRequest{
		From:    "GADE <noreply@fvlcon.org>",
		To:      []string{email},
		Subject: subject,
		Html:    htmlContent,
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send invitation email: %w", err)
	}
	return nil
}

func SendAccountApprovalEmail(email string, user db.Accounts, config *config.Config) error {
	client := resend.NewClient(config.ResendApiKey)
	tokenMaker, err := token.NewJWTMaker(config.TokenSecret)
	if err != nil {
		return err
	}

	// Generate a single token for both actions
	token, _, err := tokenMaker.CreateToken(user.ID, user.Role, config.AccessTokenDuration)
	if err != nil {
		return err
	}
	baseUrl := "https://mdukq3hdsh.us-east-1.awsapprunner.com/api"

	approvalLink := fmt.Sprintf("%s/approval?action=approve&token=%s", baseUrl, token)
	rejectionLink := fmt.Sprintf("%s/approval?action=reject&token=%s", baseUrl, token)

	htmlContent := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<style>
				.button {
					display: inline-block;
					padding: 10px 20px;
					margin: 5px;
					border-radius: 5px;
					color: #fff;
					text-decoration: none;
					font-weight: bold;
				}
				.approve {
					background-color: #28a745;
				}
				.reject {
					background-color: #dc3545;
				}
			</style>
		</head>
		<body>
			<p>Dear Admin,</p>
			<p>A new user has requested account approval:</p>
			<p><strong>User:</strong> %s</p>
			<p>Please take the appropriate action:</p>
			<p>
				<a href="%s" class="button approve">Approve</a>
				<a href="%s" class="button reject">Reject</a>
			</p>
			<p>Thank you,<br>GADE Team</p>
		</body>
		</html>
	`, user.Email, approvalLink, rejectionLink)

	params := &resend.SendEmailRequest{
		From:    "GADE <noreply@fvlcon.org>",
		To:      []string{email},
		Html:    htmlContent,
		Subject: "Account Approval Request",
	}

	_, err = client.Emails.Send(params)
	return err
}

func SendWelcomeEmail(email, firstName, lastName, resendApiKey, appDomain string) error {
	client := resend.NewClient(resendApiKey)

	// Construct login link
	loginLink := fmt.Sprintf("%s/signin", appDomain)

	// Personalize greeting
	name := firstName
	if lastName != "" {
		name = fmt.Sprintf("%s %s", firstName, lastName)
	}

	// Email content
	subject := "Welcome to GADE!"
	htmlContent := fmt.Sprintf(`
        <p>Dear %s,</p>
        <p>Congratulations! Your GADE account has been successfully set up.</p>
        <p>You can now <a href="%s">log in</a> to start exploring GADE.</p>
        <p>If you have any questions or need assistance, please contact our support team at <a href="mailto:info@blvcksapphire.com">info@blvcksapphire.com</a>.</p>
        <p>Welcome aboard!<br>The GADE Team</p>
    `, name, loginLink)

	params := &resend.SendEmailRequest{
		From:    "GADE <noreply@fvlcon.org>",
		To:      []string{email},
		Subject: subject,
		Html:    htmlContent,
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send welcome email: %w", err)
	}
	return nil
}

// SendPasswordResetEmail sends a password reset email with a reset link
func SendPasswordResetEmail(email, firstName, lastName, resetToken, resendApiKey, appDomain string) error {
	client := resend.NewClient(resendApiKey)

	// Construct reset link
	resetLink := fmt.Sprintf("https://%s/reset-password?token=%s", appDomain, resetToken)

	// Personalize greeting
	name := firstName
	if lastName != "" {
		name = fmt.Sprintf("%s %s", firstName, lastName)
	}

	// Email content
	subject := "Reset Your GADE Password"
	htmlContent := fmt.Sprintf(`
        <p>Dear %s,</p>
        <p>We received a request to reset your GADE password. Please click the link below to set a new password:</p>
        <p><a href="%s">Reset Password</a></p>
        <p>This link will expire in 24 hours for security reasons. If you did not request a password reset, please ignore this email or contact <a href="mailto:info@blvcksapphire.com">info@blvcksapphire.com</a> for assistance.</p>
        <p>Best regards,<br>The GADE Team</p>
    `, name, resetLink)

	params := &resend.SendEmailRequest{
		From:    "GADE <noreply@fvlcon.org>",
		To:      []string{email},
		Subject: subject,
		Html:    htmlContent,
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send password reset email: %w", err)
	}
	return nil
}

func SendPasswordResetConfirmationEmail(email, firstName, lastName, resendApiKey, appDomain string) error {
	client := resend.NewClient(resendApiKey)

	// Construct login link
	loginLink := fmt.Sprintf("https://%s/login", appDomain)

	// Personalize greeting
	name := firstName
	if lastName != "" {
		name = fmt.Sprintf("%s %s", firstName, lastName)
	}

	// Email content
	subject := "Your GADE Password Has Been Reset"
	htmlContent := fmt.Sprintf(`
        <p>Dear %s,</p>
        <p>Your GADE password has been successfully reset.</p>
        <p>You can now <a href="%s">log in</a> with your new password.</p>
        <p>If you did not initiate this change, please contact our support team immediately at <a href="mailto:info@blvcksapphire.com">info@blvcksapphire.com</a>.</p>
        <p>Best regards,<br>The GADE Team</p>
    `, name, loginLink)

	params := &resend.SendEmailRequest{
		From:    "GADE <noreply@fvlcon.org>",
		To:      []string{email},
		Subject: subject,
		Html:    htmlContent,
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send password reset confirmation email: %w", err)
	}
	return nil
}

func SendAlertEmail(
	recipients []string,
	metrics map[string]interface{},
	frontendURL string,
	resendApiKey string,
) error {
	client := resend.NewClient(resendApiKey)

	subject := "🚨 New Illegal Mining Segments Detected"

	// Build HTML body
	html := fmt.Sprintf(`<p><strong>%d</strong> new illegal mining segments detected.</p>`, metrics["total_new_segments"].(int))

	if districts, ok := metrics["districts_affected"].([]string); ok && len(districts) > 0 {
		html += "<p><strong>Districts Affected:</strong><ul>"
		for _, d := range districts {
			html += fmt.Sprintf("<li>%s</li>", d)
		}
		html += "</ul></p>"
	}

	if counts, ok := metrics["violation_counts"].(map[string]int); ok && len(counts) > 0 {
		html += "<p><strong>Violation Types:</strong><ul>"
		for vtype, count := range counts {
			html += fmt.Sprintf("<li>%s: %d</li>", vtype, count)
		}
		html += "</ul></p>"
	}

	// Add link to view in frontend
	taskID := metrics["task_id"].(string)
	html += fmt.Sprintf(`<p><a href="%s?task_id=%s">👉 View in GADE Dashboard</a></p>`, frontendURL, taskID)

	params := &resend.SendEmailRequest{
		From:    "GADE <noreply@fvlcon.org>",
		To:      recipients,
		Html:    html,
		Subject: subject,
	}

	_, err := client.Emails.Send(params)
	return err
}
