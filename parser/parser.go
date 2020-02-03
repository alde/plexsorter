package parser

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strconv"
)

// PlexSections is the toplevel wrapping object
type PlexSections struct {
	Sections []PlexSection `xml:"Directory"`
}

// PlexSection represents a library section
type PlexSection struct {
	Key      int    `xml:"key,attr"`
	Title    string `xml:"title,attr"`
	Location struct {
		ID   int    `xml:"id,attr"`
		Path string `xml:"path,attr"`
	} `xml:"Location"`
}

// FindByTitle will return a section based on it's title
func (p PlexSections) FindByTitle(title string) *PlexSection {
	for _, section := range p.Sections {
		if section.Title == title {
			return &section
		}
	}
	return nil
}

// ParseSections will parse an xml and return the struct representation
func ParseSections(xmlValue []byte) (PlexSections, error) {
	var p PlexSections
	err := xml.Unmarshal(xmlValue, &p)
	return p, err
}

// PlexVideos is the top level object of the video information
type PlexVideos struct {
	Videos []PlexVideo `xml:"Video"`
}

// PlexVideo represents a video item
type PlexVideo struct {
	Title     string `xml:"title,attr"`
	ViewCount int    `xml:"viewCount,attr,omitifempty"`
	Media     struct {
		Part struct {
			File      string `xml:"file,attr"`
			Container string `xml:"container,attr"`
		} `xml:"Part"`
	} `xml:"Media"`
}

// FindSeen is used to find watched videos
func (v PlexVideos) FindSeen() []PlexVideo {
	seen := []PlexVideo{}
	for _, video := range v.Videos {
		if video.ViewCount > 0 {
			seen = append(seen, video)
		}
	}
	return seen
}

// ParseVideos will parse an xml and return the struct representation
func ParseVideos(xmlValue []byte) (PlexVideos, error) {
	var v PlexVideos
	err := xml.Unmarshal(xmlValue, &v)
	return v, err
}

// ExtractSeason will get the season number from the filename
func ExtractSeason(filename string) (int, error) {
	patterns := []string{
		`.*s(\d+)e(\d+).*`, // abc.s01e01.whatever
		`.*S(\d+)E(\d+).*`, // abc.s01e01.whatever
		`(\d+)x(\d+).*`,    // abc.01x01.whatever
	}
	for _, pattern := range patterns {
		p, err := regexp.Compile(pattern)
		if err != nil {
			return 0, err
		}
		match := p.FindStringSubmatch(filename)
		if match != nil {
			return strconv.Atoi(match[1])
		}
	}

	return 0, fmt.Errorf("unable to find a season in %s", filename)
}
