package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	rand "math/rand"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

// path to the cli executable
var cli string = "./raven-cli"

// number of connected peers triggering the clean-up
var minConnections int64 = 110

// maximum number of peers that can be disconnected in one run
var maxDisconnect int = 32

func main() {

	// getblockchaininfo
	_, outBuf, errBuf, err := runOSCommand(cli, []string{"getblockchaininfo"}, true)
	if err != nil {
		fmt.Println(err)
		fmt.Println(string(errBuf.Bytes()))
	}

	result := make(map[string]interface{})
	err = json.Unmarshal(outBuf.Bytes(), &result)
	if err != nil {
		log.Fatalln(err)
	}
	// check the host is synched
	var hostBlock int64 = -2
	synced := false
	if result["blocks"] == result["headers"] {
		synced = true
		hostBlock, _ = strconv.ParseInt(fmt.Sprintf("%.0f", result["blocks"]), 10, 64)
		fmt.Println("Host is synched at block " + fmt.Sprint(hostBlock))
	}

	if !synced {
		log.Println("Host is not synced, skipping peers clean-up")
		os.Exit(0)
	}

	// getpeerinfo
	_, outBuf, errBuf, err = runOSCommand(cli, []string{"getpeerinfo"}, true)
	if err != nil {
		fmt.Println(err)
		fmt.Println(string(errBuf.Bytes()))
	}
	var peersResult []map[string]interface{}
	err = json.Unmarshal(outBuf.Bytes(), &peersResult)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Peers number: " + fmt.Sprint(len(peersResult)))
	if len(peersResult) >= int(minConnections) {
		n := 0
		x := 0
		for _, r := range peersResult {
			peerBlock, _ := strconv.ParseInt(fmt.Sprintf("%.0f", r["synced_blocks"]), 10, 64)
			// analyse peers blockcount vs host blockcount
			if hostBlock == peerBlock && fmt.Sprint(r["inbound"]) == "true" {
				n++
				// roll the dice: disconnect or not
				rand.NewSource(time.Now().UnixNano())
				z := rand.Intn(99)
				if x < maxDisconnect && z >= 50 {
					// disconnect the peer
					_, outBuf, errBuf, err = runOSCommand(cli, []string{"disconnectnode", fmt.Sprint(r["addr"])}, false)
					if err != nil {
						fmt.Println(err)
						fmt.Println(string(errBuf.Bytes()))
					} else {
						fmt.Println("Disconnected node: " +
							"id=" + fmt.Sprintf("%5v", r["id"]) +
							"; addr=" + fmt.Sprint(r["addr"]))
						x++
					}
				}

			}
		}

		log.Println("Number of nodes that can be disconnected:  " + fmt.Sprint(n))
		log.Println("Number of nodes successfully disconnected: " + fmt.Sprint(x))
	} else {
		log.Println("Number of connected nodes is below " + fmt.Sprint(minConnections) + ", skipping disconnections...")
	}

}

// RunOSCommandSilent running os command without user interaction
func runOSCommand(command string, params []string, fatal bool) (exitCode int, stdoutBuf bytes.Buffer, stderrBuf bytes.Buffer, err error) {
	exitCode = -1
	cmd := exec.Command(command, params...)

	cmd.Stdout = io.MultiWriter(&stdoutBuf)
	cmd.Stderr = io.MultiWriter(&stderrBuf)

	err = cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		}
		if fatal {
			fmt.Println()
			log.Fatalf("Command execution failed with %s\n", err)
		}
	} else {
		// success, exitCode should be 0 if all goes ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}

	return
}
