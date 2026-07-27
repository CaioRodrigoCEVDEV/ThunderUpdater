package repository

type Release struct {
	Version     int
	FileName    string
	DownloadURL string
	Exists      bool
}

type ReleaseCatalog struct {
	Releases      []Release
	LatestVersion int
}

func (c *ReleaseCatalog) LatestRelease() *Release {
	if len(c.Releases) == 0 {
		return nil
	}
	return &c.Releases[len(c.Releases)-1]
}

func (c *ReleaseCatalog) FindByVersion(version int) *Release {
	for i := range c.Releases {
		if c.Releases[i].Version == version {
			return &c.Releases[i]
		}
	}
	return nil
}
