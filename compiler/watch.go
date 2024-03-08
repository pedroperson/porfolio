package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
)

// Watcher encapsulates the file system watcher and the action to be performed on file changes.
type Watcher struct {
	watcher *fsnotify.Watcher
	action  func()
	path    string
}

// NewWatcher creates and initializes a new Watcher.
func NewWatcher(action func(), path string) (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &Watcher{
		watcher: fsWatcher,
		action:  action,
		path:    path,
	}, nil
}

// Start begins watching the file system and executes the provided action on changes.
func (w *Watcher) Start() error {
	done := make(chan bool)

	// Initial action call
	w.action()

	go func() {
		for {
			select {
			case event, ok := <-w.watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
					w.action()
				}
			case err, ok := <-w.watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err := recursiveAdd(w.watcher, w.path)
	if err != nil {
		return err
	}

	<-done // Keep running until program is terminated
	return nil
}

// Cleanup performs necessary cleanup actions.
func (w *Watcher) Cleanup() {
	w.watcher.Close()
}

func main() {
	// Paths and action function
	const layoutPath = "templates/layout.html"
	const templatesDir = "templates"
	const pagesDir = "templates/pages"
	const publicDir = "public"

	colorBlue := color.New(color.FgBlue).SprintFunc()
	colorRed := color.New(color.FgRed).SprintFunc()

	var serverWG sync.WaitGroup
	serverWG.Add(1)
	// Start the dev server
	go executeScript("yarn dev", "[vercel]", &serverWG, colorRed)

	action := func() {
		fmt.Println("lets compile")
		CompileTemplates(layoutPath, templatesDir, pagesDir, publicDir)
		fmt.Println("compiled")

		var wg sync.WaitGroup
		wg.Add(1)
		// go executeScript("go run compiler/watch.go", "[html]", &wg, colorRed)
		go executeScript("yarn tailwindcss -i templates/main.css -o public/style.css", "[tailwind]", &wg, colorBlue)

		wg.Wait() // Wait for first to complete

		fmt.Println("css done as well")

	}

	watcher, err := NewWatcher(action, "templates")
	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Cleanup()

	err = watcher.Start()
	if err != nil {
		log.Fatal(err)
	}

}

// executeScript runs a given shell script and outputs its content with a colored tag.
func executeScript(script string, tag string, wg *sync.WaitGroup, colorFunc func(a ...interface{}) string) {
	defer wg.Done()

	cmd := exec.Command("bash", "-c", script) // Use "bash" for Linux/macOS. Use "cmd", "/C", script for Windows.
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error obtaining stdout:", err)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting command:", err)
		return
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Println(colorFunc(tag), scanner.Text()) // Print each output line with the colored tag
	}
	fmt.Println("IT HAS ENDED")

	if err := cmd.Wait(); err != nil {
		fmt.Println("Error waiting for command:", err)
		return
	}

}

// recursiveAdd watches all subdirectories of the given directory path
func recursiveAdd(watcher *fsnotify.Watcher, basePath string) error {
	return filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}

		return watcher.Add(path)
	})
}

type PageData struct {
	Title   string
	Content template.HTML // Allows for HTML content
}

// compileAllTemplates compiles the layout and all other component templates.
func compileAllTemplates(baseLayout string, dirs ...string) (*template.Template, error) {
	tmpl := template.New("layout")
	var err error

	// Parse the base layout first
	tmpl, err = tmpl.ParseFiles(baseLayout)
	if err != nil {
		return nil, err
	}

	// Parse all other templates from specified directories
	for _, dir := range dirs {
		pattern := filepath.Join(dir, "*.html")
		tmpl, err = tmpl.ParseGlob(pattern)
		if err != nil {
			return nil, err
		}
	}

	return tmpl, nil
}

// CompileTemplates compiles HTML templates from the pages directory using a common layout and additional components.
func CompileTemplates(layoutPath, templatesDir, pagesDir, publicDir string) {
	// Compile layout and all component templates
	tmpl, err := compileAllTemplates(layoutPath, filepath.Join(templatesDir, "components"), filepath.Join(templatesDir, "icons"))
	if err != nil {
		panic(err) // Replace with proper error handling
	}

	// Iterate over each HTML file in the pages directory
	files, err := os.ReadDir(pagesDir)
	if err != nil {
		panic(err) // Replace with proper error handling
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".html") {
			continue
		}

		pageName := file.Name()
		pagePath := filepath.Join(pagesDir, pageName)

		// Parse the page-specific content as a new template
		pageContent, err := os.ReadFile(pagePath)
		if err != nil {
			panic(err)
		}

		pageTmpl, err := tmpl.Clone()
		if err != nil {
			panic(err)
		}

		_, err = pageTmpl.New("content").Parse(string(pageContent))
		if err != nil {
			panic(err)
		}

		// Create the output file in the public directory
		outputPath := filepath.Join(publicDir, pageName)
		outputFile, err := os.Create(outputPath)
		if err != nil {
			panic(err)
		}
		defer outputFile.Close()

		// Execute the combined template with the layout
		err = pageTmpl.ExecuteTemplate(outputFile, "layout", pageData())
		if err != nil {
			panic(err)
		}

		println("Compiled:", pageName)
	}
}

type Project struct {
	Name      string
	URL       string
	Role      string
	Business  string
	Tasks     string
	Image1    string
	Image1Alt string
	Image2    string
	Image2Alt string
	Color     string
}

func pageData() interface{} {
	return struct {
		Projects []Project
		Title    string
		Content  template.HTML
	}{
		Projects: []Project{
			{
				Name:      "bitbu",
				URL:       "https://bitbu.io/",
				Role:      "Co-founder",
				Business:  "Online tools for musician collaboration",
				Tasks:     "everything tech",
				Image1:    "/bitbu_home.jpg",
				Image1Alt: "screenshot of Bitbu's home page",
				Image2:    "/bitbu_short.jpg",
				Image2Alt: "screenshot of Bitbu's playlist sharing page",
				Color:     "[230,230,240]", //"[13,13,13]",
			},
			{
				Name:      "Ludlow Kingsley",
				URL:       "https://ludlowkingsley.com/",
				Role:      "Frontend Developer",
				Business:  "Corporate design agency",
				Tasks:     "built websites for brands designed in-house",
				Image1:    "/ludlow_home.jpg",
				Image1Alt: "screenshot of ludlow kingsley's home page",
				Image2:    "/ludlow_project.jpg",
				Image2Alt: "screenshot of ludlow kingsley's website",
				Color:     "[13,74,27]", // "[253,253,247]"
			},
			{
				Name:      "Jerde",
				URL:       "https://jerde.com/",
				Role:      "Frontend Developer",
				Business:  "Architecture design firm",
				Tasks:     "wrote the frontend for the visually striking Ludlow Kingsley design",
				Image1:    "/jerde_home.jpg",
				Image1Alt: "screenshot of Jerde's home page",
				Image2:    "/jerde_project.jpg",
				Image2Alt: "screenshot of one of Jerde's projects called the 'Hard Rock Seminole Spirit', showing some conceptual art for the project",
				Color:     "[228,233,230]", // "[13,74,27]",
			},
			{
				Name:      "Heloisa Prieto",
				URL:       "https://heloisaprieto.com/?lang=english",
				Role:      "Fullstack developer",
				Business:  "Prolific Brazilian writer",
				Tasks:     "wrote front and backends for eugênia hanitzsch's design",
				Image1:    "/heloisa_home.jpg",
				Image1Alt: "screenshot of home page of Heloisa's website",
				Image2:    "/heloisa_project.jpg",
				Image2Alt: "screenshot of books page of Heloisa's website",
				Color:     "[50,50,50]",
			},
			// {
			// 	Name:      "JM Agency",
			// 	URL:       "https://jm.agency/",
			// 	Role:      "Fullstack developer",
			// 	Business:  "Prolific Brazilian writer",
			// 	Tasks:     "wrote custom wordpress theme for eugênia hanitzsch's design",
			// 	Image1:    "/jm_home.jpg",
			// 	Image1Alt: "screenshot of home page for JM Agency",
			// 	Image2:    "/jm_project.jpg",
			// 	Image2Alt: "screenshot of testimonials page for JM Agency",
			// 	Color:     "[91,11,94]",
			// },
		},
	}
}
