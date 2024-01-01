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
	GameType string
	Version  string
	GUID     string
	Size     string
	URL      string
}

var categories = map[string][]eftDLV{
	"EFT":   make([]eftDLV, 0),
	"ETS":   make([]eftDLV, 0),
	"ARENA": make([]eftDLV, 0),
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

	for name, category := range categories {
		sort.SliceStable(category, func(i, j int) bool {
			return category[i].Version < category[j].Version
		})
		fmt.Println("[", name, "]")
		for _, dlv := range category {
			if dlv.Version != "" && dlv.GUID != "" && dlv.URL != "" {
				fmt.Println("Version:", dlv.Version, "GUID:", dlv.GUID, "FileSize:", dlv.Size)
				fmt.Println("Download Link:", dlv.URL)
				fmt.Println()
			}
		}
	}

	//for name, category := range categories {
	//	fmt.Println("[", name, "]")
	//	for _, dlv := range category {
	//		if dlv.Version != "" && dlv.GUID != "" && dlv.URL != "" {
	//			fmt.Println("Version:", dlv.Version, "GUID:", dlv.GUID, "FileSize:", dlv.Size)
	//			fmt.Println("Download Link:", dlv.URL)
	//			fmt.Println()
	//		}
	//	}
	//}

	fmt.Println("Press Enter to exit...")
	_, err = fmt.Scanln()
	if err != nil {
		return
	}
}

// var output = make([]eftDLV, 0)
var duplicates = make(map[string]struct{})

func extractInfo(line string) {
	var isUpdate bool
	var clientInfo string
	var gameType string

	var filepathSplit = "/client"
	var guidSplit []string

	if strings.Contains(line, "/arena/client") {
		gameType = "ARENA"
		filepathSplit = "/arena"
	} else if strings.Contains(line, "/eft/client") {
		if strings.Contains(line, "ets") {
			gameType = "ETS"
		} else {
			gameType = "EFT"
		}
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

	if filepathSplit == "/eft" || filepathSplit == "/arena" {
		guidSplit = strings.Split(splitClientInfo[5], "_")
	} else {
		guidSplit = strings.Split(splitClientInfo[4], "_")
	}

	if isUpdate {
		updateURL := cdn + clientInfo + ".update"
		if _, ok := duplicates[guidSplit[0]]; ok {
			return
		}

		categories[gameType] = append(categories[gameType], eftDLV{
			GameType: gameType,
			Version:  guidSplit[0],
			GUID:     guidSplit[1],
			URL:      updateURL,
			Size:     line[strings.Index(line, "size of")+8:],
		})
		duplicates[guidSplit[0]] = struct{}{}

		if filepathSplit == "/client" {
			versionSplit := strings.Split(guidSplit[0], "-")[1]
			if _, ok := duplicates[versionSplit]; ok {
				return
			}

			zipURL := cdn + "/" + splitClientInfo[1] + "/" + splitClientInfo[2] + "/distribs/" + versionSplit + "_" + guidSplit[1] + "/Client." + versionSplit + ".zip"
			categories[gameType] = append(categories[gameType], eftDLV{
				GameType: gameType,
				Version:  versionSplit,
				GUID:     guidSplit[1],
				URL:      zipURL,
				Size:     "Unknown",
			})
			duplicates[versionSplit] = struct{}{}
		}
	} else {
		if _, ok := duplicates[guidSplit[0]]; ok {
			return
		}

		categories[gameType] = append(categories[gameType], eftDLV{
			GameType: gameType,
			Version:  guidSplit[0],
			GUID:     guidSplit[1],
			URL:      cdn + clientInfo + ".zip",
			Size:     line[strings.Index(line, "size of")+8:],
		})
		duplicates[guidSplit[0]] = struct{}{}
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
