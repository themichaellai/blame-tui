package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type BlameHeader struct {
	AuthorName      string
	CommitHash      string
	AuthorTimeEpoch int64
	AuthorTimezone  string
}

type BlameLine struct {
	BlameHeader
	Code string
}

var firstLineRegexp = regexp.MustCompile("^([0-9a-f]{40}) ([0-9]+) ([0-9]+) ([0-9]+)$")
var authorRegexp = regexp.MustCompile("^author (.*)$")
var authorTimeRegexp = regexp.MustCompile("^author-time ([0-9]+)$")
var authorTimezoneRegexp = regexp.MustCompile("^author-tz ([0-9]+)$")

func Blame(filename string) ([]BlameLine, error) {
	cmd := exec.Command("git", "blame", "--porcelain", filename)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	commitMap := map[string]BlameHeader{}
	lines := bytes.Split(out.Bytes(), []byte("\n"))
	result := make([]BlameLine, 0)

	var currentHeader BlameHeader
	for _, lineBytes := range lines {
		line := string(lineBytes)

		if firstLineMatch := firstLineRegexp.FindStringSubmatch(line); firstLineMatch != nil {
			commitMap[currentHeader.CommitHash] = currentHeader
			if header, ok := commitMap[firstLineMatch[1]]; ok {
				currentHeader = header
			} else {
				currentHeader = BlameHeader{
					CommitHash: firstLineMatch[1],
				}
			}
		} else if authorMatch := authorRegexp.FindStringSubmatch(line); authorMatch != nil {
			currentHeader.AuthorName = authorMatch[1]
		} else if authorTimeMatch := authorTimeRegexp.FindStringSubmatch(line); authorTimeMatch != nil {
			epoch, err := strconv.ParseInt(authorTimeMatch[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("could not parse author time: %w", err)
			}
			currentHeader.AuthorTimeEpoch = epoch
		} else if authorTimezoneMatch := authorTimezoneRegexp.FindStringSubmatch(line); authorTimezoneMatch != nil {
			currentHeader.AuthorTimezone = authorTimezoneMatch[1]
		} else if strings.HasPrefix(line, "\t") {
			result = append(result, BlameLine{
				BlameHeader: currentHeader,
				Code:        strings.TrimPrefix(line, "\t"),
			})
		}
	}
	return result, nil

}
