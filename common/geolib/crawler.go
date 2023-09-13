package geolib

type Crawler struct {
	LastXT, LastYT int
	XYRanges
}

func (crawler *Crawler) Next(dx, dy int) bool {
	if crawler == nil {
		return false
	}

	var x, y int

	if crawler.LastYT < crawler.XYRanges.YT[0] {
		y = crawler.XYRanges.YT[0]
	} else {
		y = crawler.LastYT
	}
	if crawler.LastXT < crawler.XYRanges.XT[0] {
		x = crawler.XYRanges.XT[0]
	} else if x = crawler.LastXT + dx; x > crawler.XYRanges.XT[1] {
		y += dy
		x = crawler.XYRanges.XT[0]
	}

	if x > crawler.XYRanges.XT[1] || y > crawler.XYRanges.YT[1] {
		return false
	}

	crawler.LastXT, crawler.LastYT = x, y

	return true
}
