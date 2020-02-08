package parser

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parse_sections(t *testing.T) {
	f, err := os.Open("testdata/sections.xml")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	s, _ := ParseSections(data)
	assert.Equal(t, 3, len(s.Sections))
	found := s.FindByTitle("Unsorted")
	assert.Equal(t, "Unsorted", found.Title)
	assert.Equal(t, 3, found.Key)
}

func Test_parse_videos(t *testing.T) {
	f, err := os.Open("testdata/all.xml")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	videos, _ := ParseVideos(data)
	assert.Equal(t, 4, len(videos.Videos))
	seen := videos.FindSeen()
	assert.Equal(t, 1, len(seen))
}

func Test_extract_season(t *testing.T) {
	testData := []struct {
		filename string
		season   int
	}{
		{filename: "egg.bacon.s01e02.some.title", season: 1},
		{filename: "egg.bacon.S01E02.some.title", season: 1},
		{filename: "egg.bacon.01x02.some.title", season: 1},
	}
	for _, td := range testData {
		t.Run(fmt.Sprintf("parsing %s", td.filename), func(t *testing.T) {
			season, err := ExtractSeason(td.filename)
			assert.Nil(t, err, "error should be empty")
			assert.Equal(t, season, td.season)
		})
	}
}

func Test_extract_showname(t *testing.T) {
	testData := []struct {
		filename string
		showName string
	}{
		{filename: "egg.bacon.s01e02.some.title", showName: "egg.bacon"},
		{filename: "egg.bacon.S01E02.some.title", showName: "egg.bacon"},
		{filename: "egg.bacon.01x02.some.title", showName: "egg.bacon"},
		{filename: "doctor.who.2005.s01e02e03", showName: "doctor.who.2005"},
	}
	for _, td := range testData {
		t.Run(fmt.Sprintf("parsing %s", td.filename), func(t *testing.T) {
			showName, err := ExtractShowName(td.filename)
			assert.Nil(t, err, "error should be empty")
			assert.Equal(t, showName, td.showName)
		})
	}
}
