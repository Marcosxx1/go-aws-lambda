package utils

import (
	"fmt"
	"strconv"
	"test/lambda/interfaces"

	"github.com/gin-gonic/gin"
)

// ParseFormData parses form data from the given Gin context.
// It returns a RequestEvent or an error if parsing fails.
//
// Example:
//
//	c.Request.Form.Set("name", "Tabloide Marcos")
//	c.Request.Form.Set("region_id", "144")
//	c.Request.Form.Set("start_validity_date", "2024-04-08")
//	c.Request.Form.Set("end_validity_date", "2024-04-10")
//	file, _ := c.FormFile("file") // Assume you have a file uploaded in the form
//	event, err := ParseFormData(c)
//	if err != nil {
//	    fmt.Println("Error:", err)
//	    return
//	}
//	fmt.Println("Parsed form data:", event)
func ParseFormData(c *gin.Context) (*interfaces.RequestEvent, error) {
	event := &interfaces.RequestEvent{
		Name: c.Request.FormValue("name"),
	}

	regionID, err := strconv.Atoi(c.Request.FormValue("region_id"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse region_id: %w", err)
	}
	event.RegionID = regionID

	startValidityDate, err := parseDate(c.Request.FormValue("start_validity_date"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse start_validity_date: %w", err)
	}
	event.StartValidityDate = startValidityDate

	endValidityDate, err := parseDate(c.Request.FormValue("end_validity_date"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse end_validity_date: %w", err)
	}
	event.EndValidityDate = endValidityDate

	file, err := c.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	event.File, err = parseFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	return event, nil
}
