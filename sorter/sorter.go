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
	if len(watched) == 0 {
		logrus.Info("no watched videos")
		return
	}
	for _, video := range watched {
		err := handle(video, target)
		if err != nil {
			logrus.WithError(err).WithField("video", video).Error("error processig file")
		}
	}
}

func handle(video parser.PlexVideo, target parser.PlexSection) error {
	bestMatch, err := isArchived(target.Location.Path, video.Title)
	if err != nil {
		logrus.WithError(err).WithField("directory", video.Title).Info("show not archived - skipping")
		return err
	}

	season, err := parser.ExtractSeason(video.Title)
	if err != nil {
		logrus.WithError(err).Error("unable to determine season - skipping")
		return err
	}

	targetFolder := fmt.Sprintf("%s/%s/Season %d", target.Location.Path, bestMatch.show, season)
	if _, err := os.Stat(targetFolder); err != nil {
		err := os.MkdirAll(targetFolder, 0755)
		if err != nil {
			logrus.WithError(err).Error("unable to create season directory - skipping")
			return err
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
		logrus.WithError(err).Error("unable to move file to directory")
		return err
	}
	return nil
}

type match struct {
	count int
	show  string
}

func isArchived(target, filename string) (*match, error) {
	filename = strings.ToUpper(filename)
	files, err := ioutil.ReadDir(target)
	if err != nil {
		logrus.WithError(err).WithField("target", target).Fatal("unable to read target directory")
		return nil, err
	}
	bestMatches := []*match{}
	for _, directory := range files {
		logrus.WithField("entry", directory.Name()).Debug("checking entry")
		if !directory.IsDir() {
			continue
		}
		showName, err := parser.ExtractShowName(filename)
		if err != nil {
			return nil, err
		}
		parts := strings.Split(strings.ToUpper(directory.Name()), " ")
		matches := 0
		for _, part := range parts {
			if part == "-" || part == "_" {
				continue
			}
			fileParts := strings.Split(showName, ".")
			if contains(fileParts, part) {
				logrus.Debugf(">>> %s contains %s", filename, part)
				matches++
			}
		}
		if matches > 0 {
			bestMatches = append(bestMatches, &match{
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
	if len(bestMatches) == 0 {
		return nil, fmt.Errorf("no shows matching %s found", filename)
	}
	return bestMatches[0], nil
}

func contains(slice []string, candidate string) bool {
	for _, s := range slice {
		if s == candidate {
			return true
		}
	}
	return false
}
