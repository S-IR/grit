package js

import (
	"fmt"
	"log"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

func Recompile(css *[]byte, jsBundle *[]byte) {
	options := api.BuildOptions{
		EntryPoints:    []string{"./src/.Entry.tsx"},
		Outfile:        "./public/bundle.js",
		Bundle:         true,
		AllowOverwrite: true,
		JSX:            api.JSXAutomatic,
		Target:         api.Target(1),
		Loader: map[string]api.Loader{
			".tsx": api.LoaderTSX,
			".ts":  api.LoaderTS,

			".svg": api.LoaderFile, // Or use some other loader
		},
	}

	result := api.Build(options)
	if len(result.Errors) > 0 {
		// Handle errors
		log.Println("Build failed.")
		for _, err := range result.Errors {
			fmt.Printf("Error in file: %s at line: %d, column: %d\n", err.Location.File, err.Location.Line, err.Location.Column)
			fmt.Println("Error text:", err.Text)
		}
		panic("errors with compiling")
	}

	for _, file := range result.OutputFiles {
		// If you have multiple outputs, you'd switch based on file path or some other logic
		if strings.HasSuffix(file.Path, ".js") {
			*jsBundle = file.Contents
		} else if strings.HasSuffix(file.Path, ".css") {
			*css = file.Contents
		}
	}

}
