package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type PageData struct {
	Title   string
	Content template.HTML // Allows for HTML content
}

func main() {
	// Define the base directories
	baseTemplateDir := "templates/pages"
	baseLayoutDir := "templates" // Adjust if your directory structure is different
	baseOutputDir := "public"

	// Combine all layout related files into a single template object
	// This should match your layout, header, and footer files
	layoutPattern := filepath.Join(baseLayoutDir, "*.html")

	layoutTemplates, err := template.ParseGlob(layoutPattern)
	if err != nil {
		panic(err)
	}

	// Process each page template
	err = filepath.Walk(baseTemplateDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil // Skip directories
		}

		// Read and prepare the page-specific content
		pageContent, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Set up the full page data
		fullPageData := PageData{
			Title:   "Your Page Title",          // Change as needed
			Content: template.HTML(pageContent), // Directly use read content
		}

		// Construct the output file path
		relPath, err := filepath.Rel(baseTemplateDir, path)
		if err != nil {
			return err
		}
		outputFilePath := filepath.Join(baseOutputDir, strings.Replace(relPath, filepath.Ext(relPath), ".html", -1))

		// Create and open the output file
		file, err := os.Create(outputFilePath)
		if err != nil {
			return err
		}
		defer file.Close()

		// Execute the layout template with the full page content
		err = layoutTemplates.ExecuteTemplate(file, "layout.html", fullPageData)
		if err != nil {
			return err
		}

		fmt.Println("Compiled to html:", path)

		return nil
	})

	if err != nil {
		panic(err)
	}
}
