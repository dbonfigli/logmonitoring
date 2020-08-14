package logstat

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/dbonfigli/logmonitoring/commonlog"
)

// SectionStat is a type that hold statistics about section hits, implements interface LogStat.
// SectionStat must be initialized with the NewSectionStat function.
type SectionStat struct {
	// sections is a maps {} section: counter } to keep track of the seen sections
	sections map[string]int
	// sectionsToShow is the number of top hit sections to show
	sectionsToShow int
}

// NewSectionStat creates and initialize a SectionStat object.
func NewSectionStat(sectionsToShow int) *SectionStat {
	if sectionsToShow < 1 {
		log.Fatal("sectionsToShow cannot be < 1!")
	}
	return &SectionStat{
		sectionsToShow: sectionsToShow,
		sections:       make(map[string]int),
	}
}

// NewInterval clears the statistic for the time interval.
func (s *SectionStat) NewInterval() {
	s.sections = make(map[string]int)
}

// Update updates the statistic with the new log line.
func (s *SectionStat) Update(clog commonlog.CommonLog) error {
	fields := strings.Split(clog.Request, " ")
	if len(fields) != 3 {
		return errors.New("cannot find path in request field: " + clog.Request + ", skipping this log for statistics")
	}
	if fields[1][0] != '/' {
		return errors.New("request path does not start with a slash: " + clog.Request + ", skipping this log for statistics")
	}
	requestPath := fields[1]
	indexSlash := strings.IndexByte(requestPath[1:], '/')
	// if there are no other slashes after the first character, consider the section as "/"
	if indexSlash == -1 {
		s.sections["/"]++
	} else {
		section := requestPath[:indexSlash+1]
		s.sections[section]++
	}
	return nil
}

// String returns a string used to display this statistic.
func (s *SectionStat) String() string {
	if len(s.sections) == 0 {
		return "site hits by section: no traffic seen, no section to show"
	}

	// sort sections by hits
	type kv struct {
		Key   string
		Value int
	}
	var sectionsSlice []kv
	for k, v := range s.sections {
		sectionsSlice = append(sectionsSlice, kv{k, v})
	}
	sort.Slice(sectionsSlice, func(i, j int) bool {
		return sectionsSlice[i].Value > sectionsSlice[j].Value
	})

	// show the first top sections
	sectionsToShow := len(sectionsSlice)
	if s.sectionsToShow < sectionsToShow {
		sectionsToShow = s.sectionsToShow
	}
	plaintextSectionInfoSlice := make([]string, 0, sectionsToShow)
	for _, kv := range sectionsSlice[:sectionsToShow] {
		plaintextSectionInfoSlice = append(plaintextSectionInfoSlice, fmt.Sprintf("%v: %v", kv.Key, kv.Value))
	}
	return "site hits by section: " + strings.Join(plaintextSectionInfoSlice, "; ")
}
