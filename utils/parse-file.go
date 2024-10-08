// Package usecase contains utilities for parsing the ginContext
package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"test/lambda/interfaces"
)

// parseFile parses the given multipart file header.
// It returns a File object containing file metadata or an error if parsing fails.
//
// Example:
//
//	file, err := parseFile(fileHeader)
//	if err != nil {
//	    fmt.Println("Error:", err)
//	    return
//	}
//	fmt.Println("Parsed file:", file)
func parseFile(file *multipart.FileHeader) (interfaces.File, error) {
	fileReader, err := file.Open()
	if err != nil {
		return interfaces.File{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer fileReader.Close()

	return interfaces.File{
		Name:        file.Filename,
		ContentType: file.Header.Get("Content-Type"),
		Size:        file.Size,
		Data:        file,
	}, nil
}

// ReadFileContent reads the content of the file from the provided *multipart.FileHeader and returns it as a []byte
//
// Example:
//
//	fileContent, err := ReadFileContent(fileHeader)
//	if err != nil {
//	    fmt.Println("Error:", err)
//	    return
//	}
//	fmt.Println("File content:", string(fileContent))
func ReadFileContent(fileHeader *multipart.FileHeader) ([]byte, error) {
	// Oppen the file from the multipart form data
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return fileContent, nil
}
