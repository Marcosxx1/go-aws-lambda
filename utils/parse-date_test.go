package utils

import (
	"testing"
	"time"
)

func TestParseDate(t *testing.T) {
	dateSrt := "2022-04-08"
	expectedDate, _ := time.Parse("2006-01-02", dateSrt)
	parsedDate, err := parseDate(dateSrt)
	if err != nil {
		t.Errorf("parseDate(%s) returned an error: %v", dateSrt, err)
	}

	if !parsedDate.Equal(expectedDate) {
		t.Errorf("parseDate(%s) returned %v, expected %v", dateSrt, parsedDate, expectedDate)
	}
}

func TestParseDate_InvalidDateString(t *testing.T) {
	invalidDateStr := "invalid-date"
	parsedDate, err := parseDate(invalidDateStr)
	if err == nil {
		t.Errorf("parseDate(%s) did not return an error", invalidDateStr)
	}
	if !parsedDate.IsZero() {
		t.Errorf("parseDate(%s) returned a non-zero date: %v", invalidDateStr, parsedDate)
	}
}

func TestParseDate_EmptyString(t *testing.T) {
	emptyStr := ""
	parsedDate, err := parseDate(emptyStr)
	if err == nil {
		t.Errorf("parseDate(%s) did not return an error", emptyStr)
	}
	if !parsedDate.IsZero() {
		t.Errorf("parseDate(%s) returned a non-zero date: %v", emptyStr, parsedDate)
	}
}

func TestParseDate_Nil(t *testing.T) {
	var nilStr string
	parsedDate, err := parseDate(nilStr)
	if err == nil {
		t.Errorf("parseDate(%v) did not return an error", nilStr)
	}
	if !parsedDate.IsZero() {
		t.Errorf("parseDate(%v) returned a non-zero date: %v", nilStr, parsedDate)
	}
}

func TestParseDate_LeapYear(t *testing.T) {
	dateSrt := "2020-02-29"
	expectedDate, _ := time.Parse("2006-01-02", dateSrt)
	parsedDate, err := parseDate(dateSrt)
	if err != nil {
		t.Errorf("parseDate(%s) returned an error: %v", dateSrt, err)
	}

	if !parsedDate.Equal(expectedDate) {
		t.Errorf("parseDate(%s) returned %v, expected %v", dateSrt, parsedDate, expectedDate)
	}
}

func TestParseDate_BoundaryDates(t *testing.T) {
	earliestDateStr := "0001-01-01"
	latestDateStr := "9999-12-31"
	_, err := parseDate(earliestDateStr)
	if err != nil {
		t.Errorf("parseDate(%s) returned an error: %v", earliestDateStr, err)
	}
	_, err = parseDate(latestDateStr)
	if err != nil {
		t.Errorf("parseDate(%s) returned an error: %v", latestDateStr, err)
	}
}
