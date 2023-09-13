package filelib

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"

	"github.com/pavlo67/common/common/errors"
)

const onSearchByRegexp = "on filelib.SearchByRegexp()"

func SearchByRegexp(path string, re regexp.Regexp, getFirst bool) ([]string, error) {
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return nil, errors.Wrap(err, onSearchByRegexp)
	}

	names := make([]string, len(dirEntries))
	for i, dirEntry := range dirEntries {
		names[i] = dirEntry.Name()
	}

	slices.Sort(names)
	// sort.Slice(names, func(i, j int) bool { return names[i] < names[j] })

	var matches []string
	for _, name := range names {
		if matches = re.FindStringSubmatch(name); getFirst && len(matches) > 0 {
			return matches, nil
		}
	}

	return matches, nil
}

func ListByRegexp(path string, re regexp.Regexp) ([]string, error) {
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return nil, errors.Wrap(err, onSearchByRegexp)
	}

	names := make([]string, len(dirEntries))
	for i, dirEntry := range dirEntries {
		names[i] = dirEntry.Name()
	}

	slices.Sort(names)
	// sort.Slice(names, func(i, j int) bool { return names[i] < names[j] })

	var namesListed []string
	for _, name := range names {
		if re.MatchString(name) {
			namesListed = append(namesListed, filepath.Join(path, name))
		}
	}

	return namesListed, nil
}

const onFileExists = "on filelib.FileExists()"

func FileExists(path string, isDir bool) (bool, error) {
	fileInfo, _ := os.Stat(path)
	if fileInfo == nil {
		return false, nil
	}

	if fileInfo.IsDir() {
		if isDir {
			return true, nil
		}
		return false, fmt.Errorf("%s is not a directory / "+onFileExists, path)
	}

	if isDir {
		return false, fmt.Errorf("%s is a directory / "+onFileExists, path)
	}

	return true, nil
}
