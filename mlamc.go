package main

import (
	"fmt"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func testFile(filepath string, vt bool) (string, string) {
	var fileSize = 0
	fi, err := os.Stat(filepath)
	if err == nil {
		fileSize = int(fi.Size())
	} else {
		return "ERROR", filepath + " cannot be sent are you sure the file exists and is readable ?"
	}

	if fileSize > MAX_SIZE {
		return "ERROR", filepath + " not sent because the file is too big."
	}

	// check if the file result is in the cache
	status, message := submitHash(filepath)

	if status != "NOT FOUND" {
		return status, message
	} else {
		// If the file result is not in the cache, send the file to the webservice
		status, message = submitFile(filepath, vt)
		return status, message
	}
}

/*
	in list function
*/
func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

/*
	List all files to check from the cli arguments
*/
func getFilesList(files, directories, extensions string) []string {
	filesList := make([]string, 0)
	// Files list are separated by ;
	filesList = append(filesList, strings.Split(files, ";")...)

	// Directory list are separated by ;
	directoriesToBrowse := strings.Split(directories, ";")

	// Get all files from the directory with selected extension and size <= MAX_SIZE
	if len(directoriesToBrowse) > 0 {
		extensionsToKeep := strings.Split(extensions, ";")
		for _, v := range directoriesToBrowse {

		    filepath.Walk(v, func(path string, info os.FileInfo, err error) error {
		    	fileSize := 0
				fi, err := os.Stat(path)
				if err == nil {
					fileSize += int(fi.Size())
				}

		        if !info.IsDir() && stringInSlice(filepath.Ext(path), extensionsToKeep) && fileSize <= MAX_SIZE {
		            filesList = append(filesList, path)
		        }
		        return nil
		    })
		}
	}

	return filesList
}


/*
	Main function
*/
func main() {
	// CLI
	files := flag.String("f", "", "Files to analyze")
	directories := flag.String("d", "", "Directories to analyze")
	extensions := flag.String("e", ".exe;.dll;.js;.vbs", "Extensions to use for directory analysis")
	disableVT := flag.Bool("disable-vt", false, "Disable Virus Total")
	verbose := flag.Bool("v", false, "Verbose mode")

	flag.Parse()

	// Choose at least one file or one directory
	if *files == "" && *directories == "" {
		fmt.Println("You should use at least -f for files or -d for directories !")
		os.Exit(1)
	}

	// Get the list of files to check
	filesList := getFilesList(*files, *directories, *extensions)
	
	malwaresList := make([]string, 0)

	// First analysis without VirusTotal check to avoid VT threshold
	for i := 0; i< len(filesList); i++ {
		if filesList[i] != "" {
			status, message := testFile(filesList[i], false)
			if status == "MALWARE" || status == "MALWARE_BUT" {
				malwaresList = append(malwaresList, filesList[i])
			}
			if status == "ERROR" {
				fmt.Println("Error: ", message, " for file ", filesList[i])
			}
			if *verbose {
				fmt.Println("[", i, "/", len(filesList), "] - ", filesList[i], " = ", status, " - ", message)
			}
		}
	}

	// If VT is disable display result, else check with vt malwares detected to be sure this is not FP.
	if *disableVT {
		fmt.Println("List of files detected as malicious:")
		for _, v := range malwaresList {
			fmt.Println(v)
		}
	} else {
		if *verbose {
			fmt.Println("Virus total will confirm only malicious files.")
		}
		fmt.Println("List of files detected as malicious:")

		for i := 0; i< len(malwaresList); i++ {
			status, message := testFile(malwaresList[i], true)
			if status == "MALWARE" || status == "MALWARE_BUT" {

				fmt.Println(malwaresList[i], " = ", message)
			}
			if status == "ERROR" {
				fmt.Println("Error: ", message, " for file ", filesList[i])
			}
			if *verbose {
				// Display file sent
				fmt.Println("[", i + 1, "/", len(filesList), "] - ", filesList[i], " = ", status, " - ", message)
			}

			// Wait 30 secondes to avoid VT threshold
			time.Sleep(30 * time.Second)
		}

	}
	
}