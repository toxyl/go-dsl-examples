package main

import (
	"encoding/binary"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/charmbracelet/glamour"
)

/////////////////////////////////////////////////
// Utility functions                           //
/////////////////////////////////////////////////

func renderTemplate(tmplStr string, data interface{}) (string, error) {
	// First render the template
	tmpl, err := template.New("output").Parse(tmplStr)
	if err != nil {
		return "", err
	}

	var markdownOutput strings.Builder
	err = tmpl.Execute(&markdownOutput, data)
	if err != nil {
		return "", err
	}

	// Then convert Markdown to CLI-friendly output
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(140),
	)
	if err != nil {
		return "", err
	}

	cliOutput, err := r.Render(markdownOutput.String())
	if err != nil {
		return "", err
	}

	return cliOutput, nil
}

func parseLastlog() (map[string]time.Time, error) {
	file, err := os.Open("/var/log/lastlog")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Get file size
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// Calculate number of entries (each entry is 292 bytes)
	entrySize := 292
	numEntries := int(stat.Size()) / entrySize

	// Read all entries
	entries := make([]LastlogEntry, numEntries)
	err = binary.Read(file, binary.LittleEndian, &entries)
	if err != nil {
		return nil, err
	}

	// Create map of UID to last login time
	lastLogins := make(map[string]time.Time)
	for uid, entry := range entries {
		if entry.Time != 0 {
			lastLogins[strconv.Itoa(uid)] = time.Unix(int64(entry.Time), 0)
		}
	}

	return lastLogins, nil
}
