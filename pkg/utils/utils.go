package utils

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/term"
)

func ElapsedTime(sinceUnix int64) string {
	// Convert the Unix timestamp to a time.Time
	since := time.Unix(sinceUnix, 0)

	// Get the current time
	now := time.Now()

	// Calculate the difference in minutes, hours, and days
	diff := now.Sub(since)
	seconds := int(diff.Seconds())
	minutes := int(diff.Minutes())
	hours := int(diff.Hours())
	days := int(diff.Hours() / 24)

	// Return a single unit based on the elapsed time
	if seconds < 60 {
		return fmt.Sprintf("%ds ago", seconds)
	} else if minutes < 60 {
		return fmt.Sprintf("%dm ago", minutes)
	} else if hours < 24 {
		return fmt.Sprintf("%dh ago", hours)
	} else {
		return fmt.Sprintf("%dd ago", days)
	}
}

func Contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func ToggleState(states []string, state string) []string {
	// Check if the state exists in the slice
	for i, s := range states {
		if s == state {
			// State exists, remove it
			return append(states[:i], states[i+1:]...)
		}
	}

	// State doesn't exist, add it
	return append(states, state)
}

func GetHomeDir() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return dir, nil
}

func GetRootDir() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}
	return dir, nil
}

func Reader() *bufio.Reader {
	// Initiate user input reader
	reader := bufio.NewReader(os.Stdin)
	return reader
}

func UserInput(reader *bufio.Reader, msg string, secret bool) (string, error) {
	fmt.Printf("%s:", msg)
	var key string
	var err error
	if secret {
		key, err = readSecretInput()
	} else {
		key, err = readStandardInput(reader)
	}
	if err != nil {
		return "", err
	}
	if key == "" {
		return "", errors.New("no value provided")
	}
	return key, nil
}

func readStandardInput(reader *bufio.Reader) (string, error) {
	key, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return key, nil
}

func readSecretInput() (string, error) {
	secret, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	return string(secret), err
}
