# EFT_DLV_Grabber-GO
 Rewrite of Mao's DLV Grabber tool to support ETS and provide a bit more information

## What does this tool do?
 `AppData\Local\Battlestate Games\BsgLauncher\Logs` contain logs that provide direct links to updates and full client zips that you have downloaded.
 
 This tool scans through those log files and pulls the Version, GUID, File Size and Direct Download Link for these files to allow you to download these files directly, instead of through the sluggish launcher


## To Build:
`go build -ldflags="-s -w" -o .` in your terminal
