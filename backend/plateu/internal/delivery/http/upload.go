package http

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo"

	"github.com/fidellr/jastip_way/backend/plateu/models"
	"github.com/fidellr/jastip_way/backend/plateu/utils"
)

func UploadFile(c echo.Context, data interface{}, needs string) (err error) {
	var dst *os.File
	var src multipart.File
	var file *multipart.FileHeader
	var fileLink string

	switch needs {
	case "profile_picture":
		img := data.(*models.Image)
		// Get the file from user form
		file, err = c.FormFile("image")
		if err != nil {
			return fmt.Errorf("Failed to get the file from user form : %s", err.Error())
		}

		if img.FileLink != "" {
			fileLink = img.FileLink
		} else {
			fileLink = strings.ToLower(img.PersonName)
			img.FileLink = fileLink
		}

		// Open the file to prepare creating the physical file (Original source file) [Temp]
		src, err = file.Open()
		if err != nil {
			return fmt.Errorf("Failed to open given file : %s", err.Error())
		}
		defer src.Close()

		// Create new file to prepare copying bytes to the new file (Destination file) [Temp]
		dst, err = os.Create(fmt.Sprintf("%s-%s.jpg", img.FileLink, needs))
		if err != nil {
			return fmt.Errorf("Failed to copying bytes to a new file : %s", err.Error())
		}
		defer dst.Close()
		break
	default:
		return fmt.Errorf("Failed to read needs form text with value [%s], while it's required", needs)
	}

	// Copy old original source file to created file as the destination
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatalf("Failed to copy file : %s", err.Error())
		return err
	}

	// If succeed, compress created file to gzip and remove original file
	err = utils.CompressFile(data, needs)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
