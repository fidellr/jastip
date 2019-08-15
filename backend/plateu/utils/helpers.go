package utils

import (
	"archive/tar"
	"compress/zlib"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

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

// isDirExists : returns whether the given file or directory exists
func isDirExists(path string) (bool, error) {
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
		gzipOldFileName = fmt.Sprintf("%s-%s.tar.gz", img.FileLink, img.Needs)
		fileName = fmt.Sprintf("%s-%s.jpg", img.FileLink, needs)

		isExists, err := isDirExists("../saved_data/profile_data")
		if !isExists {
			if err != nil {
				err = os.MkdirAll("../saved_data/profile_data/pictures", 0755)
				if err != nil {
					return fmt.Errorf("Failed to create folder for profile picture : %s", err.Error())
				}

				return fmt.Errorf("Failed to check folder, folder is not exists : %s", err.Error())
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
	fmt.Printf("File uploaded successfully with file name : %s and with size %d\n", fileName, size)

	writer, err := os.Create(gzipNewFileName)
	if err != nil {
		return fmt.Errorf("Failed to create new file : %s", err.Error())
	}
	defer writer.Close()

	gw, err := zlib.NewWriterLevel(writer, zlib.BestSpeed)
	if err != nil {
		return fmt.Errorf("Failed to initialize gzip with level: %s", err.Error())
	}
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	header.Name = gzipNewFileName
	if err = tw.WriteHeader(header); err != nil {
		return fmt.Errorf("Failed to write header : %s", err.Error())
	}

	if !info.Mode().IsRegular() {
		return nil
	}

	fh, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer fh.Close()

	if _, err := io.Copy(tw, rawFile); err != nil {
		return fmt.Errorf("Failed to copy file data : %s", err.Error())
	}

	defer os.Remove(fileName)
	return nil
}

func DecompressFile(file *os.File) (err error) {
	var fileReader io.ReadCloser = file
	if strings.HasSuffix(file.Name(), ".gz") {
		if fileReader, err = zlib.NewReader(file); err != nil {
			log.Printf("Failed to initialize zlib reader : %s", err.Error())
			return err
		}
		defer fileReader.Close()
	}

	tarReader := tar.NewReader(fileReader)

	for {
		header, err := tarReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Printf("Failed to reading next entry in the tar file : %s", err.Error())
			break
		}

		fullpath := header.Name
		filename := strings.Split(fullpath, "/")[4]
		name := strings.TrimSuffix(filename, ".tar.gz")
		untarredFilePath := fmt.Sprintf("../saved_data/profile_data/pictures/%s.jpg", name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err = os.MkdirAll(fullpath, os.FileMode(header.Mode)); err != nil {
				log.Printf("typeDirMkdirAll err : %s", err.Error())
				return err
			}
		case tar.TypeReg:
			writer, err := os.Create(untarredFilePath)
			if err != nil {
				log.Printf("Failed to create image file for output : %s", err.Error())
				return err
			}

			if _, err = io.Copy(writer, tarReader); err != nil {
				log.Printf("Failed to copy tarred file to writer: %s", err.Error())
				return err
			}

			if err = os.Chmod(untarredFilePath, os.FileMode(header.Mode)); err != nil {
				log.Printf("Failed to change the mode of the named file : %s", err.Error())
				return err
			}

			writer.Close()
		default:
			fmt.Printf("Unable to untar type : %c in file %s", header.Typeflag, fullpath)
		}
	}

	return nil
}
