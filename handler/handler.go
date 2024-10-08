// Package usecase contains the business logic for the application
package usecase

import (
	"fmt"
	"net/http"
	"os"
	"test/lambda/interfaces"
	mysqlservice "test/lambda/services/mysql-service"
	uploaderservice "test/lambda/services/uploader-service"
	"test/lambda/utils"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Error   string      `json:"error,omitempty"`
	Request interface{} `json:"request,omitempty"`
	Context interface{} `json:"context,omitempty"`
}

// HandlePostRequest handles POST requests to upload tabloid data.
// It parses the multipart form data, validates the request event,
// performs database operations to insert tabloid data, uploads images,
// and commits the transaction.
func HandlePostRequest(c *gin.Context) {
	// Parse the multipart form data
	if err := c.Request.ParseMultipartForm(0); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	// Bind the form data to the RequestEvent struct
	var request interfaces.RequestEvent
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	// Validate the request event struct
	if err := utils.ValidateStruct(request); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	// Initialize MySQL service for tabloid repository
	mysqlService := mysqlservice.NewMysqlTabloideRepository()
	transaction, _ := mysqlService.GetTransaction()

	// Initialize upload service for uploading images
	uploadService, err := uploaderservice.NewUploaderAdapter()
	if err != nil {
		fmt.Println("err de uploadService", err)
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	// Parse multipart form data
	err = c.Request.ParseMultipartForm(10 << 20) // 10mb de tamanho máximo TODO VERIFICAR TAMANHO MÁXIMO
	if err != nil {
		fmt.Println("err de ParseMultipartFomr", err)
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	// Parse form data into RequestEvent object
	formData, err := utils.ParseFormData(c)
	if err != nil {
		fmt.Println("err de ParseFormData", err)
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	// Retrieve region by ID from MySQL service
	region, err := mysqlService.GetRegionById(formData.RegionID)
	if err != nil || region == nil {
		fmt.Println("err de GetRegionById", err)
		c.JSON(http.StatusBadRequest, Response{Error: "Region not found"})
		return
	}

	// Insert tabloid data into database
	tabloidID, err := mysqlService.InsertTabloid(formData.Name, formData.RegionID, formData.StartValidityDate, formData.EndValidityDate, transaction)
	if err != nil {
		fmt.Println("err de InsertTabloid", err)
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	// Read file content and upload image
	convertedImageContent, err := utils.ReadFileContent(formData.File.Data)
	if err != nil {
		fmt.Println("ReadFileContent", err)
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}
	var order = 0
	imageUrl, err := uploadService.UploadImage(convertedImageContent, tabloidID, order)
	if err != nil {
		fmt.Println("UploadImage", err)
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	// Format image URL
	formatedImageUrl := fmt.Sprintf("%s%s", os.Getenv("CDN_URL"), imageUrl)
	/* fmt.Println(formatedImageUrl) */

	// Insert tabloid image into database
	err = mysqlService.InsertTabloidImage(formatedImageUrl, tabloidID, order, transaction)
	if err != nil {
		fmt.Println("Error uploading image:", err)
		return
	}

	// Commit the transaction
	transaction.Commit()

	// Respond with success
	c.JSON(http.StatusOK, formData)
}
