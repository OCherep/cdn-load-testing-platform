package report

import (
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
)

/*
=========================
LEGACY SIMPLE EXPORT
=========================
DO NOT REMOVE
*/
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

/*
=========================
SLA EVIDENCE REPORT
=========================
*/

type SLAEvidence struct {
	TestID          string
	TargetURL       string
	StartTime       time.Time
	EndTime         time.Time
	AvgLatencyMs    float64
	P95LatencyMs    float64
	ErrorRate       float64
	StickinessRatio float64

	LatencySLAms  float64
	ErrorRateSLA  float64
	StickinessSLA float64
}

func (e SLAEvidence) Passed() bool {
	return e.AvgLatencyMs <= e.LatencySLAms &&
		e.ErrorRate <= e.ErrorRateSLA &&
		e.StickinessRatio >= e.StickinessSLA
}

func ExportSLAEvidencePDF(e SLAEvidence) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(0, 12, "CDN SLA Evidence Report")
	pdf.Ln(14)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, fmt.Sprintf("Test ID: %s", e.TestID))
	pdf.Ln(8)
	pdf.Cell(0, 8, fmt.Sprintf("Target URL: %s", e.TargetURL))
	pdf.Ln(8)
	pdf.Cell(0, 8, fmt.Sprintf("Start: %s", e.StartTime.Format(time.RFC3339)))
	pdf.Ln(8)
	pdf.Cell(0, 8, fmt.Sprintf("End: %s", e.EndTime.Format(time.RFC3339)))
	pdf.Ln(12)

	// Metrics
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Measured Metrics")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, fmt.Sprintf("Average Latency: %.2f ms", e.AvgLatencyMs))
	pdf.Ln(8)
	pdf.Cell(0, 8, fmt.Sprintf("P95 Latency: %.2f ms", e.P95LatencyMs))
	pdf.Ln(8)
	pdf.Cell(0, 8, fmt.Sprintf("Error Rate: %.4f", e.ErrorRate))
	pdf.Ln(8)
	pdf.Cell(0, 8, fmt.Sprintf("Stickiness Ratio: %.4f", e.StickinessRatio))
	pdf.Ln(12)

	// SLA
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "SLA Thresholds")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, fmt.Sprintf("Latency SLA: ≤ %.2f ms", e.LatencySLAms))
	pdf.Ln(8)
	pdf.Cell(0, 8, fmt.Sprintf("Error Rate SLA: ≤ %.4f", e.ErrorRateSLA))
	pdf.Ln(8)
	pdf.Cell(0, 8, fmt.Sprintf("Stickiness SLA: ≥ %.4f", e.StickinessSLA))
	pdf.Ln(14)

	// Verdict
	pdf.SetFont("Arial", "B", 16)
	result := "FAIL"
	if e.Passed() {
		result = "PASS"
	}
	pdf.Cell(0, 12, "SLA RESULT: "+result)

	file := "sla_evidence_" + e.TestID + ".pdf"
	return file, pdf.OutputFileAndClose(file)
}
