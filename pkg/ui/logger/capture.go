package logger

import (
	"bufio"
	"io"
	"os"
	"sync"

	"github.com/colecrouter/gameboy-go/pkg/display/debug"
)

// OriginalOutput holds the original stdout.
var OriginalOutput io.Writer

// RedirectOutput redirects stdout and stderr so that any output is also sent to logMenu.
func RedirectOutput(logMenu *debug.LogMenu) {
	// Save original stdout/stderr.
	OriginalOutput = os.Stdout

	// Create a pipe.
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	// Redirect stdout & stderr.
	os.Stdout = w
	os.Stderr = w

	// Use a wait group for concurrent copying.
	var wg sync.WaitGroup
	wg.Add(1)

	// Start a goroutine that reads from pipe and forwards output.
	go func(out io.Writer) {
		defer wg.Done()
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			line := scanner.Text()
			// Send to log menu.
			logMenu.AddLog(line)
			// Also write to original output so the user still sees it.
			out.Write([]byte(line + "\n"))
		}
	}(OriginalOutput)

}
