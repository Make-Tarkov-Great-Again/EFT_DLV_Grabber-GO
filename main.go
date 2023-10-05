package main

import (
	"bufio"
	"fmt"
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
	var start, end int

	var clientInfo, URL, GUID, Version string
	var versionSplit, guidSplit, splitClientInfo []string

	fmt.Println("--------------------")
	fmt.Println("Detected 'Game Versions' on your machine:")
	fmt.Println()

	for _, file := range logFiles {
		fileName = file.Name()

		if !strings.Contains(fileName, "BSG_Launcher_") {
			continue
		}
		filePath = filepath.Join(logDirectory, fileName)

		data, err := readLines(filePath)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Error in splitting file into []string, exiting!")
			os.Exit(0)
		}

		for _, line := range data {
			if strings.Contains(line, "Starting download") {

				start = strings.Index(line, "/client")

				if strings.Contains(line, ".update") {
					end = strings.Index(line, ".update")
				} else if strings.Contains(line, ".zip") {
					end = strings.Index(line, ".zip")
				} else {
					fmt.Println(err)
					fmt.Println(".update or .zip was not found on line, exiting")
					os.Exit(0)
				}

				clientInfo = line[start:end]
				splitClientInfo = strings.Split(clientInfo, "/")

				URL = cdn + clientInfo + ".zip"

				guidSplit = strings.Split(splitClientInfo[4], "_")
				GUID = guidSplit[1]

				if strings.Contains(splitClientInfo[5], "-") {
					versionSplit = strings.Split(splitClientInfo[5], "-")
					Version = versionSplit[1]
				} else {
					clientCut, _ := strings.CutPrefix(splitClientInfo[5], "Client.")
					zipCut, _ := strings.CutSuffix(clientCut, ".zip")

					Version = zipCut
				}

				fmt.Println("Version:", Version, "GUID:", GUID)
				fmt.Println("Download Link:", URL)
				fmt.Println()
			}
		}
	}

	fmt.Println("Press Enter to exit...")
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
