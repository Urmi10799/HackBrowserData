package main

import (
	"os"
	"github.com/urfave/cli/v2"
	"github.com/moond4rk/hackbrowserdata/browser"
	"github.com/moond4rk/hackbrowserdata/log"
	"github.com/moond4rk/hackbrowserdata/utils/fileutil"
)

// Global variables to store configuration options.
var (
	browserName  string    // Name of the browser to export data from.
	outputDir    string    // Directory to store the exported data.
	outputFormat string    // Format of the exported data (CSV or JSON).
	verbose      bool      // Flag to enable verbose logging.
	compress     bool      // Flag to enable compression of the output data.
	profilePath  string    // Custom browser profile directory path.
	isFullExport bool      // Flag to specify whether to perform a full export of browsing data.
)

func main() {
	// Call the Execute function to run the CLI app.
	Execute()
}

func Execute() {
	// Define the CLI app with flags and commands.
	app := &cli.App{
		Name:      "hack-browser-data", // Name of the application.
		Usage:     "Export passwords|bookmarks|cookies|history|credit cards|download history|localStorage|extensions from browser", // Brief description of the app.
		UsageText: "[hack-browser-data -b chrome -f json --dir results --zip]\nExport all browsing data (passwords/cookies/history/bookmarks) from browser\nGithub Link: https://github.com/moonD4rk/HackBrowserData", // Example usage and link to the GitHub repo.
		Version:   "0.5.0", // Version of the app.
		Flags: []cli.Flag{
			// Flag for verbose logging.
			&cli.BoolFlag{Name: "verbose", Aliases: []string{"vv"}, Destination: &verbose, Value: false, Usage: "verbose"},
			// Flag to compress the result to a zip file.
			&cli.BoolFlag{Name: "compress", Aliases: []string{"zip"}, Destination: &compress, Value: false, Usage: "compress result to zip"},
			// Flag to specify which browser to export data from. Default is 'all'.
			&cli.StringFlag{Name: "browser", Aliases: []string{"b"}, Destination: &browserName, Value: "all", Usage: "available browsers: all|" + browser.Names()},
			// Flag to specify the directory where the exported data will be saved.
			&cli.StringFlag{Name: "results-dir", Aliases: []string{"dir"}, Destination: &outputDir, Value: "results", Usage: "export dir"},
			// Flag to specify the output format (CSV or JSON).
			&cli.StringFlag{Name: "format", Aliases: []string{"f"}, Destination: &outputFormat, Value: "csv", Usage: "output format: csv|json"},
			// Flag to specify a custom profile directory for the browser.
			&cli.StringFlag{Name: "profile-path", Aliases: []string{"p"}, Destination: &profilePath, Value: "", Usage: "custom profile dir path, get with chrome://version"},
			// Flag to specify whether to perform a full export of browsing data.
			&cli.BoolFlag{Name: "full-export", Aliases: []string{"full"}, Destination: &isFullExport, Value: true, Usage: "is export full browsing data"},
		},
		HideHelpCommand: true, // Hide the help command in the CLI output.
		Action: func(c *cli.Context) error {
			// If verbose flag is set, enable verbose logging.
			if verbose {
				log.SetVerbose()
			}
			// Pick browsers based on the specified browser name or custom profile path.
			browsers, err := browser.PickBrowsers(browserName, profilePath)
			if err != nil {
				log.Errorf("pick browsers %v", err)
				return err
			}

			// Loop through the selected browsers and fetch browsing data.
			for _, b := range browsers {
				// Get browsing data for each browser (passwords, cookies, history, etc.).
				data, err := b.BrowsingData(isFullExport)
				if err != nil {
					log.Errorf("get browsing data error %v", err)
					continue // Continue with the next browser if an error occurs.
				}
				// Output the fetched data to the specified directory in the chosen format.
				data.Output(outputDir, b.Name(), outputFormat)
			}

			// If compression is enabled, compress the output directory into a zip file.
			if compress {
				if err = fileutil.CompressDir(outputDir); err != nil {
					log.Errorf("compress error %v", err)
				}
				log.Debug("compress success")
			}
			return nil
		},
	}

	// Run the app with the command line arguments passed to it.
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("run app error %v", err) // Log an error if the app fails to run.
	}
}
