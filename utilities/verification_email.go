package utilities

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mail.v2"
)

// Function to generate a verification token
func GenerateVerificationToken(email string, contentCode string) (string, error) {
	// Ambil kunci rahasia dari environment variables
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return "", errors.New("SECRET_KEY is not set")
	}

	// Buat data token dengan email dan timestamp
	data := fmt.Sprintf("%s:%d:%s", email, time.Now().Unix(), contentCode)

	// Buat HMAC hash menggunakan kunci rahasia
	hmacHash := hmac.New(sha256.New, []byte(secretKey))
	hmacHash.Write([]byte(data))
	hash := hmacHash.Sum(nil)

	// Encode data dan hash menggunakan Base64
	encodedData := base64.StdEncoding.EncodeToString([]byte(data))
	encodedHash := base64.StdEncoding.EncodeToString(hash)

	// Gabungkan data dan hash, lalu percent-encode token
	rawToken := fmt.Sprintf("%s:%s", encodedData, encodedHash)
	percentEncodedToken := url.QueryEscape(rawToken)

	return percentEncodedToken, nil
}

func DecodeToken(token string) (email string, code string, contentCode string, err error) {
	secretKey := os.Getenv("SECRET_KEY")
	parts := strings.Split(token, ":")
	fmt.Println("part ", parts)
	if len(parts) != 2 {
		return "", "", "", errors.New("invalid token format")
	}

	// Decode the base64-encoded data and hash
	encodedData, encodedHash := parts[0], parts[1]
	data, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to decode data: %v", err)
	}
	hash, err := base64.StdEncoding.DecodeString(encodedHash)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to decode hash: %v", err)
	}

	// Split the decoded data into email and code
	dataParts := strings.Split(string(data), ":")
	if len(dataParts) != 3 {
		return "", "", "", errors.New("invalid data format")
	}
	email, code, contentCode = dataParts[0], dataParts[1], dataParts[2]

	// Recompute HMAC hash using the same secret key
	hmacHash := hmac.New(sha256.New, []byte(secretKey))
	hmacHash.Write([]byte(string(data)))
	recomputedHash := hmacHash.Sum(nil)

	// Compare the recomputed hash with the hash from the token
	if !hmac.Equal(hash, recomputedHash) {
		return "", "", "", errors.New("invalid token: hash mismatch")
	}

	return email, code, contentCode, nil
}

func ValidateTime(timestamp int64) (string, bool) {
	// Convert Unix timestamp to time.Time object
	t := time.Unix(timestamp, 0)

	// Format the time to a human-readable string
	formattedTime := t.Format("2006-01-02 15:04:05")

	// Check if the time has passed 1 hour from the current time
	now := time.Now()
	oneHourAgo := now.Add(-1 * time.Hour)
	hasPassed := t.Before(oneHourAgo)

	return formattedTime, hasPassed
}

// Function to send a verification email
func SendVerificationEmail(email string, link string, subject string, body string) error {
	// Implement email sending logic using a library like gomail or smtp
	err := SendEmail(email, subject, body)
	if err != nil {
		return err
	}
	return nil
}

func SendVerificationEmailWithAttachment(email, link, subject, body, attachmentPath string) error {
	err := SendEmailWithAttachment(email, subject, body, attachmentPath)
	fmt.Println(err)
	if err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}
	return nil
}

func SendEmail(to, subject, body string) error {
	// Configure the sender email and SMTP settings
	senderEmail := os.Getenv("MAIL_SENDER")
	senderPassword := os.Getenv("MAIL_PASS")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	// Create a new email message
	m := mail.NewMessage()
	m.SetHeader("From", senderEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body) // to send plain text use text/plain

	// Create a new SMTP dialer
	d := mail.NewDialer(smtpHost, smtpPort, senderEmail, senderPassword)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func SendEmailWithAttachment(to, subject, body, attachmentPath string) error {
	// Configure the sender email and SMTP settings
	senderEmail := os.Getenv("MAIL_SENDER")
	senderPassword := os.Getenv("MAIL_PASS")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	// Create a new email message
	m := mail.NewMessage()
	m.SetHeader("From", senderEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body) // to send plain text use text/plain

	// Attach file if attachmentPath is provided
	if attachmentPath != "" {
		m.Attach(attachmentPath) // tidak ada penanganan error karena tidak ada nilai balik
	}

	// Create a new SMTP dialer
	d := mail.NewDialer(smtpHost, smtpPort, senderEmail, senderPassword)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
func GenerateEmailBodyVerification() string {
	const emailTemplate = `
	<!DOCTYPE html>
<html>
<head>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            margin: 0;
            padding: 0;
            background-color: #f9f9f9;
        }
        .container {
            max-width: 600px;
            margin: 20px auto;
            border: 1px solid #e0e0e0;
            border-radius: 8px;
            padding: 20px;
            text-align: left;
            background-color: #ffffff;
        }
        .logo {
            text-align: left;
            margin-bottom: 20px;
        }
        .header {
            font-size: 18px;
            font-weight: bold;
            margin-bottom: 20px;
        }
        .content {
            text-align: justify;
            margin-bottom: 20px;
        }
       .button {
            width: max-content;
            margin: 20px auto;
            padding: 10px 20px;
            font-size: 16px;
            color: #ffffff;
            background-color:  #4ba3f0;
            border: none;
            border-radius: 5px;
            text-decoration: none;
			display : flex;
			justify-content : center;
			align-items : center;
        }
        .link {
            word-wrap: break-word;
            text-align: justify;
            margin-bottom: 20px;
        }
        .footer {
            font-size: 12px;
            color: #777;
            text-align: center;
            margin-top: 20px;
        }
        .footer a {
            color: #4ba3f0;
            text-decoration: none;
        }
        .signature {
            text-align: left;
            margin-top: 20px;
            font-size: 14px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="logo">
            <img src="{{.SchoolLogo}}" alt="EduMaster Logo" style="width: 150px;">
        </div>
        <div class="header">Halo, {{.UserName}}!</div>
        <div class="content">
            Terima kasih telah mendaftar di EduMaster. Untuk mengaktifkan akun Anda dan mulai menggunakan semua fitur yang kami sediakan, silakan verifikasi email Anda dengan mengklik tombol di bawah ini:
        </div>
        <a class="button" href="{{.ConfirmationLink}}">Aktivasi Akun Saya</a>
        <div class="content">
            Jika tombol di atas tidak berfungsi, Anda dapat menyalin dan menempelkan tautan berikut ini di browser Anda:
        </div>
        <div class="link">
            <a href="{{.ConfirmationLink}}">{{.ConfirmationLink}}</a>
        </div>
        <div class="signature">
            Hormat Kami,<br>
            Edu Master
        </div>
        <div class="footer">
            Email ini dibuat secara otomatis. Mohon untuk tidak membalasnya. Jika Anda memiliki pertanyaan atau memerlukan bantuan, silakan hubungi call center kami melalui email di edumaster@gmail.com.
        </div>
    </div>
</body>
</html>
`
	return emailTemplate
}

func GenerateEmailBodyChangePassword() string {
	const emailTemplate = `
	<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f4f4f4;
            margin: 0;
            padding: 0;
        }
        .email-container {
        
            max-width: 600px;
            margin: 20px auto;
            background-color: #ffffff;
            padding: 20px;
            border-radius: 5px;
            box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
        }
         .logo {
            text-align: left;
            margin-bottom: 20px;
        }
       
        .email-header img {
            max-width: 150px;
        }
        .email-content {
            font-size: 16px;
            color: #333333;
            line-height: 1.6;
        }
        .email-content h1 {
            font-size: 22px;
            color: #333333;
        }
        .email-button {
            text-align: center;
            margin: 20px 0;
        }
        .email-button a {
            display: inline-block;
            background-color: #4ba3f0;
            color: #ffffff;
            text-decoration: none;
            padding: 10px 20px;
            border-radius: 5px;
            font-size: 16px;
        }
        .email-footer {
            font-size: 12px;
            color: #666666;
            text-align: center;
            margin-top: 20px;
        }
    </style>
</head>
<body>
    <div class="email-container">
            <div class="logo">
            <img src="{{.SchoolLogo}}" alt="EduMaster Logo" style="width: 150px;">
        </div>
        <div class="email-content">
            <h1>Hello!</h1>
            <p>Anda menerima email ini karena kami menerima permintaan untuk Password akun Anda di EduMaster.</p>
            <div class="email-button">
                <a href="{{.ConfirmationLink}}" class="email-button">Ubah Password</a>
            </div>
            <p>Tautan Ubah Password pendaftaran ini akan kedaluwarsa dalam 24 Jam.<br>
               Jika Anda tidak melakukan pendaftaran, tidak ada tindakan lebih lanjut yang diperlukan.</p>
            <p>Hormat kami,EduMaster.</p>
        </div>
        <div class="email-footer">
            <p>Email ini dibuat secara otomatis, mohon untuk tidak membalasnya. Jika Anda memiliki pertanyaan atau memerlukan bantuan, silakan hubungi call center kami melalui email di <a href="mailto:edumaster@gmail.com">edumaster@gmail.com</a>.</p>
        </div>
    </div>
</body>
</html>
`
	return emailTemplate
}

func GenerateEmailBodyTransactionSukses() string {

	const transactionTemplate = `
		<!DOCTYPE html>
		<html lang="id">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #f5f5f5;
					color: #333333;
					margin: 0;
					padding: 20px;
				}
				.email-container {
					width: 100%;
					max-width: 600px;
					margin: 0 auto;
					background-color: #ffffff;
					padding: 20px;
					border-radius: 8px;
					box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
				}
				.email-header {
					text-align: left;
					margin-bottom: 20px;
				}
				.email-header img {
					max-width: 80px;
				}
				.email-header h1 {
					font-size: 18px;
					font-weight: 700;
					color: #333333;
					margin-top: 5px;
				}
				.email-content {
					margin-bottom: 20px;
				}
				.email-content h2 {
					color: #333333;
					font-size: 16px;
					font-weight: 400;
					margin-bottom: 10px;
				}
				.email-content h2 span {
					font-weight: bold;
				}
				.invoice-info {
					background-color: #f7f7f7;
					border-radius: 8px;
					padding: 15px;
					margin-bottom: 20px;
					text-align: left;
					width: 70%;
				}
				.invoice-info p {
					margin: 8px 0;
				}
				.invoice-info .label {
					color: #595959;
					font-weight: 500;
				}
				.invoice-info .total-payment {
					color: #388AAF;
					font-size: 24px;
					font-weight: bold;
					padding: 10px;
				}
				.invoice-field {
					display: flex;
					justify-content:space-between;
					padding: 10px;
				}
				.footer {
					text-align: left;
					font-size: 12px;
					color: #666666;
					margin-top: 20px;
				}
				.footer p {
					margin: 2px 0;
				}
			</style>
		</head>
		<body>
			<div class="email-container">
				<!-- Header -->
				<div class="email-header">
					<img src="{{.SchoolLogo}}" alt="Logo Sekolah">
					<h1>Pembayaran Anda Berhasil</h1>
				</div>

				<!-- Content -->
				<div class="email-content">
					<h2>Kepada bapak/ibu wali murid dari <span>{{.StudentName}}</span>,</h2>
					<p>Terima kasih telah melakukan pembayaran. Berikut detail pembayaran Anda:</p>

					<!-- Invoice Information -->
					<div style="display:flex; justify-content:center">
						<div class="invoice-info">
						<div class="invoice-field"><span class="label">No Invoice</span> <b>{{.InvoiceNumber}}</b></div>
						<div class="invoice-field"><span class="label">Tanggal Pembayaran</span> <b>{{.PaymentDate}}</b></div>
						<div class="invoice-field"><span class="label">Total Pembayaran</span></div>
						<p class="total-payment"> {{.TotalPayment}}</p>
					</div>
					</div>
				</div>

				<!-- Footer -->
				<div class="footer">
					<p>{{.SchoolName}}</p>
					<p>©{{.Year}} {{.SchoolName}}. All Rights Reserved</p>
				</div>
			</div>
		</body>
		</html>
		`
	return transactionTemplate
}

func GenerateEmailBodyTransactionWaiting() string {
	const transactionTemplate = `
	<!DOCTYPE html>
	<html lang="id">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<style>
			body {
				font-family: Arial, sans-serif;
				background-color: #f5f5f5;
				color: #333333;
				margin: 0;
				padding: 20px;
			}
			.email-container {
				width: 100%;
				max-width: 600px;
				margin: 0 auto;
				background-color: #ffffff;
				padding: 20px;
				border-radius: 8px;
				box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
			}
			.email-header {
				text-align: left;
				margin-bottom: 20px;
			}
			.email-header img {
				max-width: 80px;
			}
			.email-header h1 {
				font-size: 18px;
				font-weight: 700;
				color: #333333;
				margin-top: 5px;
			}
			.email-content {
				margin-bottom: 20px;
			}
			.email-content h2 {
				color: #333333;
				font-size: 16px;
				font-weight: 400;
				margin-bottom: 10px;
			}
			.email-content h2 span {
				font-weight: bold;
			}
			.invoice-info {
				background-color: #f7f7f7;
				border-radius: 8px;
				padding: 15px;
				margin-bottom: 20px;
				text-align: left;
				width: 70%;
			}
			.invoice-info p {
				margin: 8px 0;
			}
			.invoice-info .label {
				color: #595959;
				font-weight: 500;
			}
			.invoice-info .total-payment {
				color: #388AAF;
				font-size: 24px;
				font-weight: bold;
				padding: 10px;
			}
			.invoice-field {
				display: flex;
				justify-content:space-between;
				padding: 10px;
			}
			.footer {
				text-align: left;
				font-size: 12px;
				color: #666666;
				margin-top: 20px;
			}
			.footer p {
				margin: 2px 0;
			}
		</style>
	</head>
	<body>
		<div class="email-container">
			<!-- Header -->
			<div class="email-header">
				<img src="{{.SchoolLogo}}" alt="Logo Sekolah">
				<h1>Segera Melakukan Pembayaran</h1>
			</div>

			<!-- Content -->
			<div class="email-content">
				<h2>Kepada bapak/ibu wali murid dari <span>{{.StudentName}}</span>,</h2>
				<p>Pembayaran Anda dibawah ini akan segera kadaluwarsa pada {{.ExpireDate}}</p>

				<!-- Invoice Information -->
				<div style="display:flex; justify-content:center">
					<div class="invoice-info">
					<div class="invoice-field"><span class="label">No Invoice</span> <b>{{.InvoiceNumber}}</b></div>
					<div class="invoice-field"><span class="label">Tanggal Pembayaran</span> <b>{{.PaymentDate}}</b></div>
					<div class="invoice-field"><span class="label">Total Pembayaran</span></div>
					<p class="total-payment"> {{.TotalPayment}}</p>
				</div>
				</div>
			</div>

			<!-- Footer -->
			<div class="footer">
				<p>{{.SchoolName}}</p>
				<p>©{{.Year}} {{.SchoolName}}. All Rights Reserved</p>
			</div>
		</div>
	</body>
	</html>
	`
	return transactionTemplate
}

func GenerateEmailBodyTransactionFailedMidtrans() string {
	const transactionTemplate = `
		<!DOCTYPE html>
		<html lang="id">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #f5f5f5;
					color: #333333;
					margin: 0;
					padding: 20px;
				}
				.email-container {
					width: 100%;
					max-width: 600px;
					margin: 0 auto;
					background-color: #ffffff;
					padding: 20px;
					border-radius: 8px;
					box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
				}
				.email-header {
					text-align: left;
					margin-bottom: 20px;
				}
				.email-header img {
					max-width: 80px;
				}
				.email-header h1 {
					font-size: 18px;
					font-weight: 700;
					color: #333333;
					margin-top: 5px;
				}
				.email-content {
					margin-bottom: 20px;
				}
				.email-content h2 {
					color: #333333;
					font-size: 16px;
					font-weight: 400;
					margin-bottom: 10px;
				}
				.email-content h2 span {
					font-weight: bold;
				}
				.invoice-info {
					background-color: #f7f7f7;
					border-radius: 8px;
					padding: 15px;
					margin-bottom: 20px;
					text-align: left;
					width: 70%;
				}
				.invoice-info p {
					margin: 8px 0;
				}
				.invoice-info .label {
					color: #595959;
					font-weight: 500;
				}
				.invoice-info .total-payment {
					color: #388AAF;
					font-size: 24px;
					font-weight: bold;
					padding: 10px;
				}
				.invoice-field {
					display: flex;
					justify-content:space-between;
					padding: 10px;
				}
				.footer {
					text-align: left;
					font-size: 12px;
					color: #666666;
					margin-top: 20px;
				}
				.footer p {
					margin: 2px 0;
				}
			</style>
		</head>
		<body>
			<div class="email-container">
				<!-- Header -->
				<div class="email-header">
					<img src="{{.SchoolLogo}}" alt="Logo Sekolah">
					<h1>Pembayaran Anda Dibatalkan</h1>
				</div>

				<!-- Content -->
				<div class="email-content">
					<h2>Kepada bapak/ibu wali murid dari <span>{{.StudentName}}</span>,</h2>
					<p>Sayang sekali, Anda <b>membatalkan</b> pembayaran dengan detail sebagai berikut:</p>

					<!-- Invoice Information -->
					<div style="display:flex; justify-content:center">
						<div class="invoice-info">
						<div class="invoice-field"><span class="label">No Invoice</span> <b>{{.InvoiceNumber}}</b></div>
						<div class="invoice-field"><span class="label">Tanggal Pembayaran</span> <b>{{.PaymentDate}}</b></div>
						<div class="invoice-field"><span class="label">Total Pembayaran</span></div>
						<p class="total-payment"> {{.TotalPayment}}</p>
					</div>
					</div>
				</div>

				<!-- Footer -->
				<div class="footer">
					<p>{{.SchoolName}}</p>
					<p>©{{.Year}} {{.SchoolName}}. All Rights Reserved</p>
				</div>
			</div>
		</body>
		</html>
	`
	return transactionTemplate
}

func GenerateEmailBodyTransactionFailed() string {
	const transactionTemplate = `
		<!DOCTYPE html>
		<html lang="id">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #f5f5f5;
					color: #333333;
					margin: 0;
					padding: 0;
				}
				.email-container {
					width: 100%;
					max-width: 600px;
					margin: 0 auto;
					background-color: #ffffff;
					padding: 20px;
					border-radius: 8px;
					box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
				}
				.email-header {
					margin-bottom: 20px;
				}
                .email-header h1{
					margin-bottom: 20px;
                    margin-left: 15px;
				}
				.email-header img {
					max-width: 100px;
				}
				.email-content {
					padding: 20px;
					background-color: #f7f7f7;
					border-radius: 8px;
					margin-bottom: 20px;
					text-align: center;
				}
				.email-content h2 {
					color: #333333;
					font-size: 16px;
                    font-weight: 400;
					margin-bottom: 10px;
                    text-align: left;
				}
                .email-content h2 span {
					color: #333333;
					font-size: 18px;
                    font-weight: 700;
					margin-bottom: 10px;
                    text-align: left;
				}
				.email-content p {
					margin: 5px 0;
                    text-align: left;
				}
				.invoice-info {
					display: inline-block;
					text-align: left;
					margin-top: 10px;
				}
                .invoice-info p {
					margin-top: 15px;
				}
				.invoice-info .label {
					font-weight: bold;
				}
				.invoice-info .total-payment {
					color: #388AAF;
					font-size: 24px;
					font-weight: bold;
					margin-top: 15px;
				}	
				
				.footer {
					text-align: left;
                    font-weight: bold;
					font-size: 12px;
					margin-top: 20px;
				}
			</style>
		</head>
		<body>
			<div class="email-container">
				<!-- Header with Logo -->
				<div class="email-header">
					<img src="{{.SchoolLogo}}" alt="School Logo">
					<h1>Pembayaran Anda Dibatalkan</h1>
				</div>

				<!-- Email Content -->
				<div class="email-content">
					<h2>Kepada bapak/ibu wali murid dari <span>{{.StudentName}}</span>,</h2>
					<p>Sayang sekali, Anda <b>membatalkan</b> pembayaran dengan detail sebagai berikut:</p>

					<!-- Invoice Details -->
					<div class="invoice-info">
						<p><span class="label">Virtual Account:</span> {{.VirtualAccount}}</p>
						<p><span class="label">Nama Bank:</span> {{.BankName}}</p>
						<p><span class="label">Total Pembayaran</p>
						<p><span class="total-payment">Rp {{.TotalPayment}}</p>
					</div>

				</div>

				<!-- Footer -->
				<div class="footer">
					<p>{{.SchoolName}}</p>
					<p>©{{.Year}} {{.SchoolName}}. All Rights Reserved</p>
				</div>
			</div>
		</body>
		</html>`
	return transactionTemplate
}

func GenerateEmailBodyBillingReminder() string {
    const transactionTemplate = `
    <!DOCTYPE html>
    <html lang="id">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <style>
            body {
                font-family: Arial, sans-serif;
                background-color: #f5f5f5;
                color: #333333;
                margin: 0;
                padding: 20px;
            }
            .email-container {
                width: 100%;
                max-width: 600px;
                margin: 0 auto;
                background-color: #ffffff;
                padding: 20px;
                border-radius: 8px;
                box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
            }
            .email-header {
                text-align: left;
                margin-bottom: 20px;
            }
            .email-header img {
                max-width: 80px;
            }
            .email-header h1 {
                font-size: 18px;
                font-weight: 700;
                color: #333333;
                margin-top: 10px;
            }
            .email-content {
                margin-bottom: 20px;
            }
            .email-content h2 {
                color: #333333;
                font-size: 16px;
                font-weight: 400;
                margin-bottom: 10px;
            }
            .email-content h2 span {
                font-weight: bold;
            }
            .invoice-info {
                background-color: #f7f7f7;
                border-radius: 8px;
                padding: 15px;
                margin: 0 auto 20px auto;
                text-align: left;
                max-width: 400px;
            }
            .invoice-info p {
                margin: 8px 0;
            }
            .invoice-info .label {
                color: #595959;
                font-weight: 500;
            }
            .invoice-info .total-payment {
                color: #388AAF;
                font-size: 24px;
                font-weight: bold;
                padding: 10px 0;
            }
            .invoice-field {
                display: flex;
                justify-content: space-between;
                padding: 5px 0;
            }
            .footer {
                text-align: left;
                font-size: 12px;
                color: #666666;
                margin-top: 20px;
            }
            .footer p {
                margin: 2px 0;
            }
        </style>
    </head>
    <body>
        <div class="email-container">
            <!-- Header -->
            <div class="email-header">
                <img src="{{.SchoolLogo}}" alt="Logo Sekolah">
                <h1>Pengingat untuk Melakukan Pembayaran</h1>
            </div>

            <!-- Content -->
            <div class="email-content">
                <h2>Kepada bapak/ibu wali murid dari <span>{{.StudentName}}</span>,</h2>
                <p>Kami ingin mengingatkan, pembayaran tagihan sekolah sudah mendekati batas waktu dengan detail sebagai berikut:</p>

                <!-- Invoice Information -->
                <div class="invoice-info">
                    <div class="invoice-field">
                        <span class="label">Batas Waktu Pembayaran</span>
                        <b>{{.ExpireDate}}</b>
                    </div>
                    <div class="invoice-field">
                        <span class="label">Total Pembayaran</span>
                    </div>
                    <p class="total-payment">Rp{{.TotalPayment}}</p>
                </div>
            </div>

            <!-- Footer -->
            <div class="footer">
                <p>{{.SchoolName}}</p>
                <p>©{{.Year}} {{.SchoolName}}. All Rights Reserved</p>
            </div>
        </div>
    </body>
    </html>
    `
    return transactionTemplate
}

