package WebAssets

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

var (
	cssHash      string
	jsHash       string
	assetsHashes map[string]string
)

var Ctx api.BuildContext

type ErrorDisplayed struct {
	ID       string   `json:"ID"`
	Text     string   `json:"Text"`
	Name     string   `json:"Name"`
	Line     int      `json:"Line,omitempty"`
	Column   int      `json:"Column,omitempty"`
	File     string   `json:"File,omitempty"`
	LineText string   `json:"LineText,omitempty"`
	Notes    []string `json:"Notes,omitempty"`
}

func init() {
	assetsHashes = make(map[string]string)

	SERVER_ENV := os.Getenv("SERVER_ENV")
	var minifySintax bool

	if SERVER_ENV == "development" {
		minifySintax = true
	} else {
		minifySintax = false
	}

	options := api.BuildOptions{
		EntryPoints:       []string{"./src/.Entry.tsx"},
		Outfile:           "./public/bundle.js",
		Bundle:            true,
		AllowOverwrite:    true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      minifySintax,
		JSX:               api.JSXAutomatic,
		Target:            api.Target(1),
		Loader: map[string]api.Loader{
			".tsx": api.LoaderTSX,
			".ts":  api.LoaderTS,
			".svg": api.LoaderFile, // Or use some other loader
		},
	}

	var ctxErr *api.ContextError
	Ctx, ctxErr = api.Context(options)

	if ctxErr != nil {
		panic(ctxErr)
	}

}

func Rebuild() (updatedFiles []string, updatedAssetPaths []string, buildErrors []ErrorDisplayed) {
	res := Ctx.Rebuild()
	updatedFiles = []string{}
	updatedAssetPaths = []string{}

	for _, file := range res.OutputFiles {
		if strings.HasSuffix(file.Path, ".js") {
			if file.Hash == jsHash {
				continue
			}
			*JsBundle = file.Contents
			jsHash = file.Hash
			updatedFiles = append(updatedFiles, "js")
		} else if strings.HasSuffix(file.Path, ".css") {
			if file.Hash == cssHash {
				continue
			}
			*Css = file.Contents
			cssHash = file.Hash
			updatedFiles = append(updatedFiles, "css")
		} else {

			newHash := file.Hash
			if newHash == assetsHashes[file.Path] {
				continue
			}
			fileName := filepath.Base(file.Path)
			updatedAssetPaths = append(updatedAssetPaths, fileName)

			assetsHashes[fileName] = newHash
			if Assets[fileName] == nil {
				Assets[fileName] = &[]byte{}
			}
			*Assets[fileName] = file.Contents

			assetsAlreadyPlaced := false
			for _, item := range updatedFiles {
				if strings.Contains(item, "assets") {
					assetsAlreadyPlaced = true

				}
			}

			if !assetsAlreadyPlaced {
				updatedFiles = append(updatedFiles, "assets")
			}

		}
	}
	if res.Errors == nil {
		return updatedFiles, updatedAssetPaths, nil
	}
	buildErrors = []ErrorDisplayed{}
	for _, err := range res.Errors {
		ed := ErrorDisplayed{Text: err.Text, Name: "Error during Render"}

		if err.Location != nil {
			ed.Line = err.Location.Line
			ed.Column = err.Location.Column
			ed.File = err.Location.File
			ed.LineText = err.Location.LineText
		}
		notesText := []string{}

		if len(err.Notes) > 0 {
			for _, note := range err.Notes {
				notesText = append(notesText, note.Text)
			}
		}
		ed.Notes = notesText

		buildErrors = append(buildErrors, ed)
	}
	return updatedFiles, updatedAssetPaths, buildErrors
}
