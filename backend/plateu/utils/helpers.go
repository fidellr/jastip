package utils

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/fidellr/jastip_way/backend/plateu/models"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type ErrorHTTPResponse struct {
	Message string `json:"error"`
}

func HandleUncaughtHTTPError(err error, c echo.Context) {
	logrus.Error(err)
	c.JSON(http.StatusInternalServerError, ErrorHTTPResponse{Message: err.Error()})
}

// Exists : returns whether the given file or directory exists
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// CompressFile :
func CompressFile(data interface{}, needs string) (err error) {
	var gzipOldFileName, gzipNewFileName, fileName string

	switch needs {
	case "profile_picture":
		img := data.(*models.Image)
		gzipOldFileName = fmt.Sprintf("%s-%s.gzip", img.FileLink, img.Needs)
		fileName = fmt.Sprintf("%s-%s.jpg", img.FileLink, needs)

		isExists, err := Exists("../saved_data/profile_data")
		if !isExists {
			if err != nil {
				return fmt.Errorf("Failed to check folder, folder is not exists")
			}

			err = os.MkdirAll("../saved_data/profile_data/pictures", 0755)
			if err != nil {
				return fmt.Errorf("Failed to create folder for profile picture")
			}
		}

		gzipNewFileName = fmt.Sprintf("../saved_data/profile_data/pictures/%s", gzipOldFileName)
		break
	default:
		return fmt.Errorf("Failed to read needs while compression, needs is [%s] while its required", needs)
	}

	rawFile, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("Failed to open file with file name : %s : , err : %s", fileName, err.Error())
	}
	defer rawFile.Close()

	info, _ := rawFile.Stat()
	size := info.Size()
	rawBytes := make([]byte, size)
	fmt.Printf("File uploaded successfuly with file name : %s and with size %d\n", fileName, size)

	buffer := bufio.NewReader(rawFile)
	_, err = buffer.Read(rawBytes)
	if err != nil {
		return fmt.Errorf("Failed to read rawBytes : %s", err.Error())
	}

	var buf bytes.Buffer
	gw, err := gzip.NewWriterLevel(&buf, gzip.BestSpeed)
	if err != nil {
		return fmt.Errorf("Failed to initialize gzip with level: %s", err.Error())
	}

	if _, err := gw.Write(rawBytes); err != nil {
		return fmt.Errorf("Failed gzip file : %s", err.Error())
	}

	if err = gw.Flush(); err != nil {
		return fmt.Errorf("Failed to flush pending compressed data : %s", err.Error())
	}
	defer gw.Close()

	err = ioutil.WriteFile(gzipOldFileName, buf.Bytes(), info.Mode())
	if err != nil {
		return fmt.Errorf("Failed to write file : %s", err.Error())
	}

	err = os.Rename(gzipOldFileName, gzipNewFileName)
	if err != nil {
		return fmt.Errorf("Failed to rename file : %s", err.Error())
	}

	gzipSize, err := buf.Read(buf.Bytes())
	fmt.Printf("Gzip File uploaded successfuly with the old file name %s and new file name %s with %d size\n", gzipOldFileName, gzipNewFileName, gzipSize)

	if err = os.Remove(fileName); err != nil {
		return fmt.Errorf("Failed to remove image file : %s", err.Error())
	}

	return nil
}

func DecompressFile(w io.Writer, data []byte) ([]byte, error) {
	gr, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize gzip reader : %s", err.Error())
	}
	defer gr.Close()

	data, err = ioutil.ReadAll(gr)
	if err != nil {
		return nil, fmt.Errorf("Failed to readAll reader : %s", err.Error())
	}

	return data, nil
}
