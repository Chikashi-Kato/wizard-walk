package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/hajimehoshi/ebiten/v2"
)

type Wizard struct {
	up    *ebiten.Image
	left  *ebiten.Image
	down  *ebiten.Image
	right *ebiten.Image
	loc   location
}

func getWizard(id uint) (*Wizard, error) {
	w := Wizard{}

	if err := w.downloadWizardImages(id); err != nil {
		return nil, fmt.Errorf("failed to download wizard image: %w", err)
	}

	return &w, nil
}

func (w *Wizard) downloadWizardImages(id uint) error {
	url := fmt.Sprintf("https://nftz.forgottenrunes.com/wizard/%d.zip", id)

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download wizard file (%d): %w", id, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// unzip
	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		log.Fatal(err)
	}

	// Read all the files from zip archive
	var wizardFilePattern = regexp.MustCompile(`^50\/wizard-.+-nobg.png$`)
	for _, zipFile := range zipReader.File {
		if wizardFilePattern.MatchString(zipFile.Name) {
			fmt.Println("Reading file:", zipFile.Name)
			unzippedFileBytes, err := readZipFile(zipFile)
			if err != nil {
				return fmt.Errorf("failed to read wizard file in zip: %w", err)
			}
			wizImg, err := createEbitenImage(unzippedFileBytes)
			if err != nil {
				return fmt.Errorf("failed to create ebiten image: %w", err)
			}
			w.up = wizImg
			w.down = wizImg
			w.left = wizImg
			w.right = wizImg
			return nil
		}
	}

	return fmt.Errorf("wizard png file not found in zip")
}
