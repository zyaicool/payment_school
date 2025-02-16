package utilities

import (
	"fmt"
	"image"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"schoolPayment/constants"
	response "schoolPayment/dtos/response"
	"schoolPayment/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jung-kurt/gofpdf"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Helper function to render a row with label and value
func renderRow(pdf *gofpdf.Fpdf, label string, value string, labelFontSize, valueFontSize, yOffset, totalRowWidth float64) {
	pdf.SetTextColor(128, 128, 128)

	// Label
	pdf.SetFont("Arial", "", labelFontSize)
	labelWidth := pdf.GetStringWidth(label)
	pdf.SetX(150)
	pdf.CellFormat(labelWidth, 6, label, "0", 0, "R", false, 0, "")

	// Value
	pdf.SetFont("Arial", "B", valueFontSize)
	pdf.SetTextColor(0, 0, 0) // Black text color for values
	pdf.SetX(151)
	pdf.CellFormat(50, 6, value+",00", "0", 1, "R", false, 0, "")

	// Adjust line position for spacing
	pdf.Ln(yOffset)
}

func GeneratePDF(c *fiber.Ctx, dataInvoice []response.RespDataInvoice, school models.School, isPrint bool) (string, error) {
	ext := filepath.Ext(school.SchoolLetterhead)
	schoolLogo := ConvertPath(school.SchoolLetterhead)
	subTotal := formatToIDR(int64(dataInvoice[0].SubTotal))
	dicount := int64(dataInvoice[0].Discount)
	totalAmount := int64(dataInvoice[0].SubTotal - dataInvoice[0].Discount)

	transactionType := dataInvoice[0].TransactionType
	if transactionType == "kasir" {
		transactionType = strings.Title(strings.ToLower(transactionType))
	}

	// Get current date and time
	now := time.Now()
	formattedDate := now.Format("02-01-2006") // Format: dd-MM-yyyy
	formattedTime := now.Format("15.04")      // Format: HH-MM-SS

	// Create the filename without using fmt.Sprintf
	filename := "Riwayat_Pembayaran " + "" + formattedDate + " " + formattedTime + ".pdf"

	// Print the filename
	println(filename)

	resp, err := http.Get(schoolLogo)
	if err != nil {
		return "", fmt.Errorf("Error fetching image: %v", err)
	}
	defer resp.Body.Close()

	// Create a temporary file to save the image
	out, err := os.Create("logo_temp" + ext)
	if err != nil {
		return "", fmt.Errorf("Error creating temporary file: %v", err)
	}
	defer out.Close()

	// Copy the response body (image data) into the temporary file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error saving image: %v", err)
	}

	// Open the image file
	file, err := os.Open("./logo_temp" + ext)
	if err != nil {
		return "", fmt.Errorf("Error opening image:", err)
	}
	defer file.Close()

	// Decode the image to get its dimensions
	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return "", fmt.Errorf("Error decoding image:", err)
	}

	// Get the width and height of the image
	imgWidth := float64(img.Width)
	imgHeight := float64(img.Height)

	// Desired maximum width and height for the PDF
	maxWidth := 100.0 // Adjust as needed
	maxHeight := 50.0 // Adjust as needed

	// Scale the image dimensions proportionally
	width := imgWidth
	height := imgHeight
	if width > maxWidth || height > maxHeight {
		ratio := min(maxWidth/width, maxHeight/height)
		width *= ratio
		height *= ratio
	}

	// Create a new PDF instance
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.ImageOptions(
		"logo_temp"+ext,
		55,     // X position
		10,     // Y position
		width,  // Width
		height, // Height
		false,  // Flow
		gofpdf.ImageOptions{ImageType: ext[1:], ReadDpi: true},
		0,
		"",
	)
	fmt.Println("cek heigth ", height+2)
	pdf.Ln(height + 2)
	pdf.SetDrawColor(89, 89, 89)
	pdf.SetLineWidth(0.3)
	pdf.Line(10, pdf.GetY()+2, 200, pdf.GetY()+2)
	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(190, 10, "BUKTI PEMBAYARAN SISWA", "0", 1, "C", false, 0, "")
	pdf.SetLineWidth(0.5)
	pdf.Line(10, pdf.GetY()+2, 200, pdf.GetY()+2)
	pdf.Ln(10)
	// Invoice details
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(89, 89, 89)
	pdf.CellFormat(95, 6, "No Invoice", "0", 0, "", false, 0, "")
	pdf.CellFormat(95, 6, "NIS", "0", 1, "", false, 0, "")
	pdf.SetTextColor(0, 0, 0) // Kembali ke warna default
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(95, 6, dataInvoice[0].InvoiceNumber, "0", 0, "", false, 0, "")
	pdf.CellFormat(95, 6, dataInvoice[0].Nis, "0", 1, "", false, 0, "")

	// Tambahkan elemen lainnya dengan warna merah
	pdf.Ln(1)
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(89, 89, 89)
	pdf.CellFormat(95, 6, "Tanggal Cetak", "0", 0, "", false, 0, "")
	pdf.CellFormat(95, 6, "Nama Siswa", "0", 1, "", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(95, 6, dataInvoice[0].PrintDate.Format("02/01/2006, 15:04"), "0", 0, "", false, 0, "")
	pdf.CellFormat(95, 6, dataInvoice[0].StudentName, "0", 1, "", false, 0, "")

	// Tipe pembayaran dengan warna merah
	pdf.Ln(1)
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(89, 89, 89)
	pdf.CellFormat(95, 6, "Tanggal Bayar", "0", 0, "", false, 0, "")
	pdf.CellFormat(95, 6, "Kelas", "0", 1, "", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(95, 6, dataInvoice[0].PaymentDate.Format("02/01/2006, 15:04"), "0", 0, "", false, 0, "")
	pdf.CellFormat(95, 6, dataInvoice[0].SchoolClassName, "0", 1, "", false, 0, "")

	pdf.Ln(1)
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(89, 89, 89)
	pdf.CellFormat(95, 6, "Tipe Pembayaran", "0", 1, "", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(95, 6, transactionType, "0", 1, "", false, 0, "")
	pdf.Ln(5)

	// Header tabel tanpa tebal
	pdf.SetFont("Arial", "", 10)
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(2)
	pdf.SetTextColor(89, 89, 89)
	pdf.CellFormat(10, 6, "No", "0", 0, "C", false, 0, "") // Menggunakan titik
	pdf.CellFormat(120, 6, "Keterangan Pembayaran", "0", 0, "", false, 0, "")
	pdf.CellFormat(60, 6, "Jumlah (Rp.)", "0", 1, "R", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(2)

	// Baris detail tabel
	for rw, billingStudent := range dataInvoice[0].BillingStudents {
		strRw := fmt.Sprintf("%d.", rw+1)                                                  // Format angka dengan titik
		formattedBillingStudentAmount := formatToIDR(int64(billingStudent.Amount)) + ",00" // Tambahkan ,00 untuk desimal

		pdf.SetFont("Arial", "B", 10)
		pdf.SetTextColor(0, 0, 0)
		pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
		pdf.Ln(2)
		pdf.CellFormat(10, 6, strRw, "", 0, "C", false, 0, "")
		pdf.CellFormat(120, 6, billingStudent.DetailBillingName, "0", 0, "", false, 0, "")
		pdf.SetFont("Arial", "B", 14)
		pdf.SetTextColor(65, 67, 71)
		pdf.CellFormat(60, 6, formattedBillingStudentAmount, "0", 1, "R", false, 0, "")
		pdf.Ln(2)
	}

	// Garis atas untuk memisahkan
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(2)

	// Total lebar baris
	totalRowWidth := 200.0 // Pastikan sesuai dengan lebar dokumen PDF Anda

	renderRow(pdf, "Subtotal", subTotal, 10, 15, 2, totalRowWidth)
	renderRow(pdf, "Diskon", formatToIDR(dicount), 10, 15, 2, totalRowWidth)
	renderRow(pdf, "Total", formatToIDR(totalAmount), 10, 15, 4, totalRowWidth)

	// Garis bawah (sesuaikan posisi garis agar sejajar)
	lineEndX := totalRowWidth
	pdf.Line(150, pdf.GetY()+2, lineEndX, pdf.GetY()+2)

	pdf.Ln(10)
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(89, 89, 89)
	pdf.CellFormat(95, 6, "Catatan:", "0", 1, "", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "B", 10)
	pdf.MultiCell(0, 5, " - Disimpan sebagai bukti pembayaran yang SAH\n - Uang yang sudah dibayarkan tidak dapat diminta kembali.", "", "L", false)

	// Output the PDF
	// err = pdf.Output(c.Response().BodyWriter())
	// if err != nil {
	// 	fmt.Println(err)
	// 	return fmt.Errorf("Failed to generate PDF")
	// }
	tempPDF, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("Error creating temporary PDF file: %v", err)
	}
	defer tempPDF.Close()
	tempPDFPath := tempPDF.Name()

	// Output the PDF to the temporary file
	err = pdf.OutputFileAndClose(tempPDFPath)
	if err != nil {
		return "", fmt.Errorf("Failed to save generated PDF")
	}

	c.Set(constants.ContentType, "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	pdfFile, err := os.Open(tempPDFPath)
	if err != nil {
		return "", fmt.Errorf("Error opening temporary PDF file: %v", err)
	}
	defer pdfFile.Close()

	_, err = io.Copy(c.Response().BodyWriter(), pdfFile)
	if err != nil {
		return "", fmt.Errorf("Failed to write PDF to response")
	}

	go func() {
		// Wait until the response is completely written, then delete the temporary file
		// (this can be done after the PDF output has been processed)
		if isPrint {
			err := os.Remove(filename)
			if err != nil {
				fmt.Printf("Error deleting temporary file: %v\n", err)
			}
		}

		_ = os.Remove("logo_temp" + ext)
		// if err != nil {
		// 	fmt.Printf("Error deleting temporary file: %v\n", err)
		// }
	}()

	return filename, nil
}

func formatToIDR(amount int64) string {
	p := message.NewPrinter(language.Indonesian)
	return p.Sprintf("%d", amount)
}
