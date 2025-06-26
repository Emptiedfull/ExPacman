package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	fails := 0
	logFile, err := os.OpenFile("process.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Could not open log file:", err)
		return
	}
	defer logFile.Close()

	for {
		cmd := exec.Command("./ExPacman.exe")
		cmd.Stdout = logFile
		cmd.Stderr = logFile

		err := cmd.Run()
		if err != nil {
			fmt.Fprintln(logFile, "Process exited with error:", err)
			fmt.Println("Process exited with error:", err)
		} else {
			fmt.Fprintln(logFile, "Process exited normally.")
			fmt.Println("Process exited normally.")
		}
		fails += 1
		fmt.Fprintln(logFile, "Restarting process...", fails)
		fmt.Println("Restarting process...", fails)
	}
}
