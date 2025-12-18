package report

import (
	"encoding/csv"
	"os"
)

type Row struct {
	Timestamp string
	Edge      string
	IP        string
	P95       string
	RPS       string
	ErrorRate string
}

func ExportCSV(testID string, rows []Row) (string, error) {
	file := "report_" + testID + ".csv"
	f, err := os.Create(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	w.Write([]string{"timestamp", "edge", "ip", "p95", "rps", "error_rate"})
	for _, r := range rows {
		w.Write([]string{
			r.Timestamp, r.Edge, r.IP, r.P95, r.RPS, r.ErrorRate,
		})
	}
	return file, nil
}
