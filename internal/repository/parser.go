package repository

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
)

var zipRegex = regexp.MustCompile(`Thunder_(\d+)\.zip`)

func parseReleases(html, baseURL string) ([]Release, error) {
	matches := zipRegex.FindAllStringSubmatch(html, -1)
	if len(matches) == 0 {
		return nil, fmt.Errorf("nenhuma release encontrada")
	}

	seen := make(map[int]bool)
	var releases []Release

	for _, m := range matches {
		version, err := strconv.Atoi(m[1])
		if err != nil {
			continue
		}
		if seen[version] {
			continue
		}
		seen[version] = true

		releases = append(releases, Release{
			Version:     version,
			FileName:    fmt.Sprintf("Thunder_%d.zip", version),
			DownloadURL: baseURL + fmt.Sprintf("Thunder_%d.zip", version),
			Exists:      true,
		})
	}

	sort.Slice(releases, func(i, j int) bool {
		return releases[i].Version < releases[j].Version
	})

	return releases, nil
}
