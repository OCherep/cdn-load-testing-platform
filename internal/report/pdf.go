package report

import "github.com/jung-kurt/gofpdf"

func ExportPDF(testID string, summary string) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "CDN Load Test Report")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 8, summary, "", "", false)

	file := "report_" + testID + ".pdf"
	return file, pdf.OutputFileAndClose(file)
}
