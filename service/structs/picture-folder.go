package structs

import (
	"crypto/tls" // Import the crypto/tls package for configuring TLS
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
	// Configure a transport that disables certificate verification to bypass the docker run error.
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Create an HTTP client with the custom transport
	client := &http.Client{Transport: tr}

	err := os.MkdirAll(PicFolder, os.ModeDir|os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating pictures folder: \n\tPicFolder: %s \n\tError: %w", PicFolder, err)
	}

	fullFilePath := filepath.Join(PicFolder, FileName)

	// Make an HTTP GET request to the URL using the custom client
	resp, err := client.Get(URL)
	if err != nil {
		return fmt.Errorf("error downloading default profile photo at App startup: \n\tURL: %s \n\tError: %w", URL, err)
	}
	defer resp.Body.Close()

	file, err := os.Create(fullFilePath)
	if err != nil {
		return fmt.Errorf("error creating the file path: \n\tFilename: %s \n\tError: %w", FileName, err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("error with io.Copy(...): \n\tError: %w", err)
	}

	return nil
}
