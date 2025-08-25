package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
)

// FormDataToJSONMiddleware converts multipart/form-data to JSON
func FormDataToJSONMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		contentType := c.GetHeader("Content-Type")

		// Check if content type is multipart/form-data
		if strings.HasPrefix(contentType, "multipart/form-data") {
			// Parse multipart form
			if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 32MB max
				c.JSON(400, gin.H{
					"error":   "Failed to parse form data",
					"message": err.Error(),
				})
				c.Abort()
				return
			}

			// Convert form values to JSON
			formData := make(map[string]interface{})

			// Handle form values
			for key, values := range c.Request.Form {
				if len(values) == 1 {
					formData[key] = values[0]
				} else {
					formData[key] = values
				}
			}

			// Handle multipart form values (including files)
			if c.Request.MultipartForm != nil {
				for key, values := range c.Request.MultipartForm.Value {
					if len(values) == 1 {
						formData[key] = values[0]
					} else {
						formData[key] = values
					}
				}

				// Handle files (if any)
				for key, files := range c.Request.MultipartForm.File {
					if len(files) == 1 {
						formData[key] = files[0].Filename
					} else {
						filenames := make([]string, len(files))
						for i, file := range files {
							filenames[i] = file.Filename
						}
						formData[key] = filenames
					}
				}
			}

			// Convert to JSON
			jsonData, err := json.Marshal(formData)
			if err != nil {
				c.JSON(400, gin.H{
					"error":   "Failed to convert form data to JSON",
					"message": err.Error(),
				})
				c.Abort()
				return
			}

			// Replace request body with JSON
			c.Request.Body = io.NopCloser(bytes.NewReader(jsonData))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Request.Header.Set("Content-Length", string(rune(len(jsonData))))
			c.Request.ContentLength = int64(len(jsonData))
		}

		c.Next()
	})
}
