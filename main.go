package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const cdn string = "http://cdn-11.eft-store.com"

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
	fmt.Println("Logs folder found")

	logFiles, err := os.ReadDir(logDirectory)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Unable to read", logDirectory)
		os.Exit(0)
	}

	for _, file := range logFiles {
		fileName := file.Name()

		if !strings.Contains(fileName, "BSG_Launcher_") {
			continue
		}
		filePath := filepath.Join(logDirectory, fileName)

		data, err := readLines(filePath)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Error in splitting file into []string, exiting")
			os.Exit(0)
		}

		for _, line := range data {
			if strings.Contains(line, "Starting download") {

				var start int = strings.Index(line, "/client")
				var end int

				if strings.Contains(line, ".update") {
					end = strings.Index(line, ".update")
				} else if strings.Contains(line, ".zip") {
					end = strings.Index(line, ".zip")
				} else {
					log.Fatalln(".update or .zip was not found on line, exiting")
				}

				clientInfo := line[start:end]

				URL := cdn + clientInfo + ".zip"
				splitClientInfo := strings.Split(clientInfo, "/")

				guidSplit := strings.Split(splitClientInfo[4], "_")
				GUID := guidSplit[1]

				var Version string
				var versionSplit []string
				if strings.Contains(splitClientInfo[5], "-") {
					versionSplit = strings.Split(splitClientInfo[5], "-")
					Version = versionSplit[1]
				} else {
					clientCut, _ := strings.CutPrefix(splitClientInfo[5], "Client.")
					zipCut, _ := strings.CutSuffix(clientCut, ".zip")

					Version = zipCut
				}

				fmt.Println(fmt.Sprintf(`
Version: %s
GUID: %s
Download Link: %s`, Version, GUID, URL))
			}
		}
	}

	fmt.Print("\n\nPress Enter to exit...")
	_, err = fmt.Scanln()
	if err != nil {
		return
	}
}

func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
