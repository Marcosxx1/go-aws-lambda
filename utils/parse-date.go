package utils

import "time"

// parseDate parses the given date string in the format "YYYY-MM-DD".
// It returns the parsed time.Time or an error if parsing fails.
//
// Example:
//
//	dateStr := "2022-04-08"
//	parsedDate, err := parseDate(dateStr)
//	if err != nil {
//	    fmt.Println("Error:", err)
//	    return
//	}
//	fmt.Println("Parsed date:", parsedDate)
func parseDate(dateString string) (time.Time, error) {
	return time.Parse("2006-01-02", dateString)
}
