package sorter

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/alde/plexsorter/parser"
	"github.com/sirupsen/logrus"
)

// Sort watched videos into their proper place
func Sort(watched []parser.PlexVideo, target parser.PlexSection) {
	for _, video := range watched {
		handle(video, target)
	}
}

func handle(video parser.PlexVideo, target parser.PlexSection) {
	bestMatch := isArchived(target.Location.Path, video.Title)
	if (bestMatch == match{}) {
		logrus.WithField("directory", video.Title).Debug("show not archived - skipping")
		return
	}

	season, err := parser.ExtractSeason(video.Title)
	if err != nil {
		logrus.WithError(err).Error("unable to determine season - skipping")
		return
	}

	targetFolder := fmt.Sprintf("%s/%s/Season %d", target.Location.Path, bestMatch.show, season)
	if _, err := os.Stat(targetFolder); err != nil {
		err := os.MkdirAll(targetFolder, 0755)
		if err != nil {
			logrus.WithError(err).Error("unable to create directory season - skipping")
			return
		}
	}
	parts := strings.Split(video.Media.Part.File, "/")
	filename := parts[len(parts)-1]
	targetFile := fmt.Sprintf("%s/%s", targetFolder, filename)
	logrus.WithFields(logrus.Fields{
		"from": video.Media.Part.File,
		"to":   targetFile,
	}).Info("moving video")
	if err := os.Rename(video.Media.Part.File, targetFile); err != nil {
		logrus.WithError(err).Error("unable to move directory")
	}
}

type match struct {
	count int
	show  string
}

func isArchived(target, filename string) match {
	filename = strings.ToUpper(filename)
	files, err := ioutil.ReadDir(target)
	if err != nil {
		logrus.WithError(err).WithField("target", target).Fatal("unable to read target directory")
	}
	bestMatches := []match{}
	for _, directory := range files {
		logrus.WithField("entry", directory.Name()).Debug("checking entry")
		if !directory.IsDir() {
			continue
		}
		parts := strings.Split(strings.ToUpper(directory.Name()), " ")
		matches := 0
		for _, part := range parts {
			if part == "-" || part == "_" {
				continue
			}
			fileParts := strings.Split(filename, ".")
			if contains(fileParts, part) {
				logrus.Debugf(">>> %s contains %s", filename, part)
				matches++
			}
		}
		if matches > 0 {
			bestMatches = append(bestMatches, match{
				count: matches,
				show:  directory.Name(),
			})
		}
	}
	sort.Slice(bestMatches, func(i, j int) bool {
		return bestMatches[i].count > bestMatches[j].count
	})
	logrus.WithField("filename", filename).Debug("file name")
	logrus.WithField("best matches", bestMatches).Debug("best matches sorted")
	return bestMatches[0]
}

func contains(slice []string, candidate string) bool {
	for _, s := range slice {
		if s == candidate {
			return true
		}
	}
	return false
}
