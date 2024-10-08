// Package awsconfig contains structs related with the flow of the application.
package interfaces

import (
	"fmt"
	"mime/multipart"
	"time"
)

// File represents a file uploaded via HTTP.
type File struct {
	Name        string                `json:"name"`                                                                  // Name of the file.
	ContentType string                `json:"content_type" validate:"required,oneof=image/png image/jpg image/jpeg"` // Content type of the file.
	Size        int64                 `json:"size"`                                                                  // Size of the file in bytes.
	Data        *multipart.FileHeader `json:"-"`                                                                     // File data.
}

// RequestEvent represents an event request.
type RequestEvent struct {
	Name              string    `json:"name" validate:"required"`                                        // Name of the event.
	RegionID          int       `json:"region_id" validate:"required,min=1"`                             // ID of the region where the event occurs.
	StartValidityDate time.Time `json:"start_validity_date" validate:"required"`                         // Start date of the event's validity.
	EndValidityDate   time.Time `json:"end_validity_date" validate:"required,gtfield=StartValidityDate"` // End date of the event's validity.
	File              File      `json:"file"`                                                            // Uploaded file associated with the event.
}

// This method uses fmt.Sprintf() to format a string containing all of File's attributes
func (f File) String() string {
	return fmt.Sprintf("Name: %s, ContentType: %s, Size: %d, Data: %s", f.Name, f.ContentType, f.Size, f.Data.Filename)
}
