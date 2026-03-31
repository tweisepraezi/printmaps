package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"syscall"

	"github.com/printmaps/printmaps/pd"
)

/*
runCommand runs a command / program.
*/
func runCommand(command string) (commandExitStatus int, commandOutput []byte, err error) {

	program := "/bin/bash"
	args := []string{"-c", command}
	cmd := exec.Command(program, args...)

	commandOutput, err = cmd.CombinedOutput()

	var waitStatus syscall.WaitStatus
	if err != nil {
		// command was not successful
		if exitError, ok := err.(*exec.ExitError); ok {
			// command fails because of an unsuccessful exit code
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			log.Printf("command exit code = <%d>", waitStatus.ExitStatus())
		}
		log.Printf("error <%v> at cmd.CombinedOutput()", err)
		log.Printf("command (not successful) = <%s>", strings.Join(cmd.Args, " "))
		if len(commandOutput) > 0 {
			log.Printf("command output (stdout, stderr) =\n%s", string(commandOutput))
		}
	} else {
		// command was successful
		waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
		if config.Testmode {
			log.Printf("command (successful) = <%s>", strings.Join(cmd.Args, " "))
			log.Printf("command exit code = <%d>", waitStatus.ExitStatus())
			if len(commandOutput) > 0 {
				log.Printf("command output (stdout, stderr) =\n%s", string(commandOutput))
			}
		}
	}

	commandExitStatus = waitStatus.ExitStatus()
	return
}

/*
nextBuildOrder returns the name of the oldest build order file.
*/
func nextBuildOrder() string {

	path := filepath.Join(pd.PathWorkdir, pd.PathOrders)
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatalf("fatal error <%v> at os.ReadDir(), path = <%v>", err, path)
	}

	if len(files) == 0 {
		return ""
	}

	sort.Slice(files, func(i, j int) bool {
		infoI, errI := files[i].Info()
		infoJ, errJ := files[j].Info()
		if errI != nil || errJ != nil {
			return false
		}
		return infoI.ModTime().Unix() < infoJ.ModTime().Unix()
	})

	for _, fileEntry := range files {
		if !fileEntry.IsDir() {
			return fileEntry.Name()
		}
	}

	return ""
}
