package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/alde/plexsorter/parser"
	"github.com/alde/plexsorter/sorter"
	"github.com/sirupsen/logrus"
)

// Command-line Flags
var (
	token      = flag.String("token", "", "API token for plex")
	host       = flag.String("host", "localhost", "Plex server host")
	port       = flag.Int("port", 32401, "Plex API port")
	source     = flag.String("source", "Unsorted", "Title of the plex section to scan")
	sortTarget = flag.String("target", "TV Shows", "Title of the plex section to sort files into")
	debug      = flag.Bool("debug", false, "enable debug logging")
)

func main() {
	flag.Parse()
	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if *token == "" {
		logrus.Fatal("token not provided")
	}

	baseURL := fmt.Sprintf("http://%s:%d", *host, *port)

	sourceSection := findSection(baseURL, *source)
	logrus.WithFields(logrus.Fields{
		"title":    sourceSection.Title,
		"location": sourceSection.Location.Path,
		"key":      sourceSection.Key,
	}).Debug("found source section")

	targetSection := findSection(baseURL, *sortTarget)
	logrus.WithFields(logrus.Fields{
		"title":    targetSection.Title,
		"location": targetSection.Location.Path,
		"key":      targetSection.Key,
	}).Debug("found target section")

	videos := findVideos(baseURL, sourceSection)
	watchedVideos := videos.FindSeen()
	logrus.WithField("count", len(watchedVideos)).Debug("watched videos")
	sorter.Sort(watchedVideos, *targetSection)
}

func findSection(baseURL string, title string) *parser.PlexSection {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/library/sections", baseURL), nil)
	check(err)
	req.Header.Add("X-Plex-Token", *token)
	resp, err := client.Do(req)
	check(err)
	body, err := ioutil.ReadAll(resp.Body)
	check(err)
	sections, err := parser.ParseSections(body)
	check(err)
	return sections.FindByTitle(title)
}

func findVideos(baseURL string, sourceSection *parser.PlexSection) parser.PlexVideos {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/library/sections/%d/all", baseURL, sourceSection.Key), nil)
	check(err)
	req.Header.Add("X-Plex-Token", *token)
	resp, err := client.Do(req)
	check(err)
	body, err := ioutil.ReadAll(resp.Body)
	check(err)
	videos, err := parser.ParseVideos(body)
	check(err)
	return videos

}

func check(err error) {
	if err != nil {
		logrus.WithError(err).Fatal("unrecoverable error")
	}
}
