package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const cdn string = "http://cdn-11.eft-store.com"

type eftDLV struct {
	Version string
	GUID    string
	Size    string
	URL     string
}

func main() {
	logDirectory := filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local", "Battlestate Games", "BsgLauncher", "Logs")
	_, err := os.Stat(logDirectory)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Folder does not exist, exiting")
			os.Exit(0)
		} else {
			fmt.Println(err)
			fmt.Println("Error checking folder at:", logDirectory)
			os.Exit(0)
		}
	}
	fmt.Println()
	fmt.Println(">> Automatic DownloadLink Creator from BsgLauncher Logs <<")
	fmt.Println("Created by TheMaoci; Rewritten by King")

	logFiles, err := os.ReadDir(logDirectory)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Unable to read", logDirectory)
		os.Exit(0)
	}

	var fileName, filePath string

	fmt.Println("--------------------")
	fmt.Println("Detected 'Game Versions' on your machine:")
	fmt.Println()

	for _, file := range logFiles {
		fileName = file.Name()

		if !strings.Contains(fileName, "BSG_Launcher_") {
			continue
		}
		filePath = filepath.Join(logDirectory, fileName)

		lines := readLines(filePath, "(DWN1) The file", "has a size of")
		if len(lines) == 0 {
			continue
		}

		for _, line := range lines {
			extractInfo(line)
		}
	}

	sort.Slice(output, func(i, j int) bool {
		return output[i].Version < output[j].Version
	})

	for _, dlv := range output {
		if dlv.Version != "" && dlv.GUID != "" && dlv.URL != "" {
			fmt.Println("Version:", dlv.Version, "GUID:", dlv.GUID, "FileSize:", dlv.Size)
			fmt.Println("Download Link:", dlv.URL)
			fmt.Println()
		}
	}

	fmt.Println("Press Enter to exit...")
	_, err = fmt.Scanln()
	if err != nil {
		return
	}
}

var output = make([]eftDLV, 0)

func extractInfo(line string) {
	var isUpdate bool
	var clientInfo string

	var filepathSplit = "/client"
	var guidSplit []string

	if strings.Contains(line, "/eft/client") {
		filepathSplit = "/eft"
	}

	if strings.Contains(line, ".update") {
		isUpdate = true
		clientInfo = line[strings.Index(line, filepathSplit):strings.Index(line, ".update")]
	} else if strings.Contains(line, ".zip") {
		clientInfo = line[strings.Index(line, filepathSplit):strings.Index(line, ".zip")]
	} else {
		fmt.Println(".update or .zip was not found on line, exiting")
		fmt.Println()
		os.Exit(0)
	}

	splitClientInfo := strings.Split(clientInfo, "/")

	if filepathSplit == "/eft" {
		guidSplit = strings.Split(splitClientInfo[5], "_")
	} else {
		guidSplit = strings.Split(splitClientInfo[4], "_")
	}

	if isUpdate {
		updateURL := cdn + clientInfo + ".update"
		output = append(output, eftDLV{
			Version: guidSplit[0],
			GUID:    guidSplit[1],
			URL:     updateURL,
			Size:    line[strings.Index(line, "size of")+8:],
		})

		if filepathSplit == "/client" {
			versionSplit := strings.Split(guidSplit[0], "-")[1]
			zipURL := cdn + "/" + splitClientInfo[1] + "/" + splitClientInfo[2] + "/distribs/" + versionSplit + "_" + guidSplit[1] + "/Client." + versionSplit + ".zip"
			output = append(output, eftDLV{
				Version: versionSplit,
				GUID:    guidSplit[1],
				URL:     zipURL,
				Size:    "Unknown",
			})
		}
	} else {
		output = append(output, eftDLV{
			Version: guidSplit[0],
			GUID:    guidSplit[1],
			URL:     cdn + clientInfo + ".zip",
			Size:    line[strings.Index(line, "size of")+8:],
		})
	}
}

const errorSubstring = "Substrings '%s' and '%s' not found in file '%s'"

var buffer = make([]string, 0, 1024)

func readLines(filename string, firstSubstring string, secondSubstring string) []string {
	buffer = buffer[0:]

	output := make([]string, 0)
	file, err := os.Open(filename)
	if err != nil {
		return output
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, firstSubstring) && strings.Contains(line, secondSubstring) {
			output = append(output, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return output
	}

	return output
}
