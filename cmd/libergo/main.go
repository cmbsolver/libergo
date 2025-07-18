package main

import (
	"flag"
	"fmt"
	"liberdatabase"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
	"titler"

	"config"
)

// main is the entry point for the application, handling command-line flags and executing the appropriate functionality.
func main() {
	titler.PrintTitle("Liber Go Configurator")

	initFlag := flag.Bool("init", false, "Initialize the default configuration")
	listFlag := flag.Bool("list", false, "List the current configuration")
	workersFlag := flag.Int("workers", 0, "Set the number of workers")
	initDBServerFlag := flag.Bool("initdbserver", false, "Initialize the podman database server")
	showHashesFlag := flag.Bool("showhashes", false, "Show all found hashes")
	hideTitleFlag := flag.Bool("hidetitle", false, "Hide the title")
	flag.Parse()

	if !*initFlag && !*listFlag && *workersFlag <= 0 && !*initDBServerFlag && !*showHashesFlag && !*hideTitleFlag {
		fmt.Println("Usage:")
		flag.PrintDefaults()
	}

	if *initFlag {
		err := config.CreateDefaultConfig()
		if err != nil {
			_, err := fmt.Fprintf(os.Stderr, "Error creating default config: %v\n", err)
			if err != nil {
				return
			}
			os.Exit(1)
		}
		fmt.Println("Default configuration created successfully.")
		os.Exit(0)
	}

	if *listFlag {
		cfg, err := config.LoadConfig()
		if err != nil {
			_, err := fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			if err != nil {
				return
			}
			os.Exit(1)
		}
		fmt.Println("Current configuration:")
		fmt.Printf("  NumWorkers:              %d\n", cfg.NumWorkers)
		fmt.Printf("  ExistingHash:            %s\n", cfg.ExistingHash)
		fmt.Printf("  AdminConnectionString:   %s\n", cfg.AdminConnectionString)
		fmt.Printf("  GeneralConnectionString: %s\n", cfg.GeneralConnectionString)
		fmt.Printf("  MaxPermutationsPerLine:  %d\n", cfg.MaxPermutationsPerLine)
		fmt.Printf("  MaxRangesPerSegment:     %d\n", cfg.MaxRangesPerSegment)
		fmt.Printf("  MaxSegmentsPerPackage:   %d\n", cfg.MaxSegmentsPerPackage)
		fmt.Printf("  HideTitle:               %t\n", cfg.HideTitle)
		os.Exit(0)
	}

	if *initDBServerFlag {
		executeScript()
		fmt.Println("Tables loaded successfully.")
		os.Exit(0)
	}

	if *showHashesFlag {
		err := liberdatabase.GetAllFoundHashes()
		if err != nil {
			_, err := fmt.Fprintf(os.Stderr, "Error showing found hashes: %v\n", err)
			if err != nil {
				return
			}
			os.Exit(1)
		}
		os.Exit(0)
	}

	if *hideTitleFlag {
		err := config.UpdateConfig("HideTitle", *hideTitleFlag)
		if err != nil {
			_, err := fmt.Fprintf(os.Stderr, "Error updating HideTitle: %v\n", err)
			if err != nil {
				return
			}
			os.Exit(1)
		}
		fmt.Println("HideTitle updated successfully.")
		os.Exit(0)
	}

	if *workersFlag > 0 {
		err := config.UpdateConfig("NumWorkers", *workersFlag)
		if err != nil {
			_, err := fmt.Fprintf(os.Stderr, "Error updating NumWorkers: %v\n", err)
			if err != nil {
				return
			}
			os.Exit(1)
		}
		fmt.Println("NumWorkers updated successfully.")
		os.Exit(0)
	}

	fmt.Println("This program is used to init or update the configuration for the toolset.")
}

// executeScript executes a shell script to initialize a database, waits for its completion, and then initializes the database connection.
func executeScript() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin", "linux":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting user home directory: %v\n", err)
			return
		}
		cmdPath := filepath.Join(homeDir, ".libergo/create_podman_db.sh")
		cmd = exec.Command("sh", cmdPath)
	default:
		fmt.Println("Unsupported operating system")
		return
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error executing script: %v\n", err)
	}

	fmt.Println("Sleeping for 2 minutes to fully initialize...")
	time.Sleep(2 * time.Minute)
	fmt.Println("Awake!")

	_, dbError := liberdatabase.InitDatabase()
	if dbError != nil {
		return
	} else {
		fmt.Println("Database initialized successfully.")
	}
}
