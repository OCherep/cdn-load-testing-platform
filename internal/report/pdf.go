package report

import (
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type SLAEvidence struct {
	TestID    string
	TargetURL string
	StartTime time.Time
	EndTime   time.Time

	AvgLatencyMs    float64
	P95LatencyMs    float64
	ErrorRate       float64
	StickinessRatio float64

	LatencySLAms  float64
	ErrorRateSLA  float64
	StickinessSLA float64
}

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

func ExportSLAEvidencePDF(e SLAEvidence) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(0, 10, "CDN SLA Evidence Report")
	pdf.Ln(14)

	pdf.SetFont("Arial", "", 12)

	writeRow := func(label string, value string) {
		pdf.CellFormat(60, 8, label, "", 0, "", false, 0, "")
		pdf.CellFormat(0, 8, value, "", 1, "", false, 0, "")
	}

	writeRow("Test ID:", e.TestID)
	writeRow("Target URL:", e.TargetURL)
	writeRow("Start Time:", e.StartTime.Format(time.RFC3339))
	writeRow("End Time:", e.EndTime.Format(time.RFC3339))
	pdf.Ln(4)

	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Measured Metrics")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	writeRow("Avg Latency (ms):", fmt.Sprintf("%.2f", e.AvgLatencyMs))
	writeRow("P95 Latency (ms):", fmt.Sprintf("%.2f", e.P95LatencyMs))
	writeRow("Error Rate:", fmt.Sprintf("%.4f", e.ErrorRate))
	writeRow("Stickiness Ratio:", fmt.Sprintf("%.2f", e.StickinessRatio))
	pdf.Ln(4)

	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "SLA Thresholds")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	writeRow("Latency SLA (ms):", fmt.Sprintf("%.2f", e.LatencySLAms))
	writeRow("Error Rate SLA:", fmt.Sprintf("%.4f", e.ErrorRateSLA))
	writeRow("Stickiness SLA:", fmt.Sprintf("%.2f", e.StickinessSLA))
	pdf.Ln(6)

	pdf.SetFont("Arial", "B", 14)
	if e.AvgLatencyMs <= e.LatencySLAms &&
		e.ErrorRate <= e.ErrorRateSLA &&
		e.StickinessRatio >= e.StickinessSLA {
		pdf.SetTextColor(0, 150, 0)
		pdf.Cell(0, 10, "SLA RESULT: PASSED")
	} else {
		pdf.SetTextColor(200, 0, 0)
		pdf.Cell(0, 10, "SLA RESULT: FAILED")
	}

	file := fmt.Sprintf("sla_evidence_%s.pdf", e.TestID)
	return file, pdf.OutputFileAndClose(file)
}
