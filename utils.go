package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"net/http"
	"unicode/utf8"

	"github.com/hajimehoshi/ebiten/v2"
)

func loadEbitenImageFromUrl(url string) (*ebiten.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download map file : %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return createEbitenImage(body)
}

func createEbitenImage(imageBytes []byte) (*ebiten.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, fmt.Errorf("failed create image: %w", err)
	}
	ebiImg := ebiten.NewImageFromImage(img)
	return ebiImg, nil
}

func readZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func trimLastChar(s string) (string, string) {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size], s[len(s)-size:]
}
