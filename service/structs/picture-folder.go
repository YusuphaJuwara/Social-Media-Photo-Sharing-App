package structs

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const PicFolder = "/tmp/pictures"
const URL = "https://w7.pngwing.com/pngs/269/897/png-transparent-family-cosmetic-dentistry-family-smile-happy-family-child-photography-people-thumbnail.png"
const FileName = "ffffffff-ffff-ffff-ffff-ffffffffffff.png"

func PicFolderDownloadAndSavePhoto() error {

	err := os.MkdirAll(PicFolder, os.ModeDir|os.ModePerm)
	if err != nil {

		// If path is already a directory, MkdirAll does nothing and returns nil.
		// if !os.IsExist(err) {
		// 	logger.WithError(err).Error("error regarding the pictures folder")
		// 	return fmt.Errorf("error creating/checking pictures folder: %w", err)
		// }

		return fmt.Errorf("error creating pictures folder: \n\tPicFolder: %s \n\tError: %w", PicFolder, err)
	}

	// Create the full file path
	fullFilePath := filepath.Join(PicFolder, FileName)

	// Make an HTTP GET request to the URL
	resp, err := http.Get(URL)
	if err != nil {
		return fmt.Errorf("error downloading default profile photo at App startup: \n\tURL: %s \n\tError: %w", URL, err)
	}
	defer resp.Body.Close()

	// Open the file for writing
	file, err := os.Create(fullFilePath)
	if err != nil {
		return fmt.Errorf("error creating the file path: \n\tFilename: %s \n\tError: %w", FileName, err)
	}
	defer file.Close()

	// Copy the response body to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("error with io.Copy(...): \n\tError: %w", err)
	}

	return nil
}
