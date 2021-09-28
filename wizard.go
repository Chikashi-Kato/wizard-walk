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

const (
	WizardAngleUp = wizardAngle("up")
	WizardAngleDown = wizardAngle("down")
	WizardAngleLeft = wizardAngle("left")
	WizardAngleRight = wizardAngle("right")
)

type Wizard struct {
	up    *ebiten.Image
	left  *ebiten.Image
	down  *ebiten.Image
	right *ebiten.Image
	loc   location
	angle wizardAngle
}

func getWizard(id uint) (*Wizard, error) {
	w := Wizard{}

	if err := w.downloadWizardImages(id); err != nil {
		return nil, fmt.Errorf("failed to download wizard image: %w", err)
	}

	return &w, nil
}

func (w *Wizard) downloadWizardImages(id uint) error {
	url := fmt.Sprintf("https://www.forgottenrunes.com/api/art/wizards/%d.zip", id)

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
	var wizardFilePattern = regexp.MustCompile(`^50\/turnarounds\/wizards-.+-(.+).png$`)
	for _, zipFile := range zipReader.File {
		if rs := wizardFilePattern.FindStringSubmatch(zipFile.Name); len(rs) != 0 {
			fmt.Println("Reading file:", zipFile.Name)
			unzippedFileBytes, err := readZipFile(zipFile)
			if err != nil {
				return fmt.Errorf("failed to read wizard file in zip: %w", err)
			}
			wizImg, err := createEbitenImage(unzippedFileBytes)
			if err != nil {
				return fmt.Errorf("failed to create ebiten image: %w", err)
			}
			switch rs[1] {
			case "right":
				w.right = wizImg
			case "left":
				w.left = wizImg
			case "back":
				w.up = wizImg
			case "front":
				w.down = wizImg
			}

			if w.right != nil && w.left != nil && w.up != nil && w.down != nil {
				return nil
			}
		}
	}

	return fmt.Errorf("wizard png file not found in zip")
}
