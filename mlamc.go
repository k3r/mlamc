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
		return filepath, " cannot be sent are you sure the file exists and is available ?"
	}

	if fileSize > MAX_SIZE {
		return filepath, " not sent because the file is too big."
	}
	
	return filepath, submitFile(filepath, vt)
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
	filesList = append(filesList, strings.Split(files, ";")...)

	directoriesToBrowse := strings.Split(directories, ";")

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

	for i := 0; i< len(filesList); i++ {
		if filesList[i] != "" {
			filepath, result := testFile(filesList[i], false)
			if strings.Contains(result, "malware") {
				malwaresList = append(malwaresList, filepath)
			}
			if *verbose {
				fmt.Println("[", i, "/", len(filesList), "] - ", filepath, " = ", result)

			}
		}
	}

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
			filepath, result := testFile(malwaresList[i], true)
			if result != "Your file seems clean." {
				result = strings.Replace(result, "<a href=\"", "", -1)
				result = strings.Replace(result, "\">virustotal.com</a>.", " .", -1)
				fmt.Println(filepath, " = " ,result)
			}
			if *verbose {
				fmt.Println("[", i, "/", len(malwaresList), "] - ", filepath, " = ", result)
			}
			time.Sleep(30 * time.Second)
		}

	}
	
}