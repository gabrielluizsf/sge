package town

import "github.com/gabrielluizsf/sge/terrain"

// Fences stand beside a road in runs, stretch by stretch: a stretch of
// stretchLen cells has a run or none, and a run leaves a cell or two free
// at both its ends, so a gap is never far.
const (
	stretchLen = 12
	// fenceShare is how many of a hundred stretches have a run.
	fenceShare = 65
	// endGap is how many cells at a road's ends, at its junctions and
	// beside the towns' streets stay without a fence.
	endGap = 3
	// apart is how near another road a fence does not stand.
	apart = 2
)

// Fence reports whether a fence stands on cell col, row: one cell off a
// road, in a run, away from the road's ends and from any other road.
func (l *Layout) Fence(col, row int) bool {
	_, ok := l.fenceRoad(col, row)
	return ok
}

// FenceDown reports whether the fence on cell col, row runs down the rows,
// as the road it stands beside does, and not across the columns.
func (l *Layout) FenceDown(col, row int) bool {
	r, ok := l.fenceRoad(col, row)
	return ok && r.C0 == r.C1
}

// fenceRoad is the road a fence on cell col, row stands beside, if one does.
func (l *Layout) fenceRoad(col, row int) (Rect, bool) {
	rects := l.roadsOf(floorDiv(col, l.Region), floorDiv(row, l.Region))
	for i, r := range rects {
		if beside(r, col, row) && inRun(l.Seed, r, col, row) && !crosses(rects, i, col, row) {
			return r, true
		}
	}
	return Rect{}, false
}

// crosses reports whether a road other than rects[i] passes within apart
// cells of cell col, row.
func crosses(rects []Rect, i, col, row int) bool {
	for j, o := range rects {
		if j != i && o.Dist(col, row) <= apart {
			return true
		}
	}
	return false
}

// beside reports whether cell col, row is one cell off the side of road
// r, clear of its ends.
func beside(r Rect, col, row int) bool {
	if r.R0 == r.R1 && r.C1-r.C0 >= 2*endGap {
		return abs(row-r.R0) == 1 && col >= r.C0+endGap && col <= r.C1-endGap
	}
	if r.C0 == r.C1 && r.R1-r.R0 >= 2*endGap {
		return abs(col-r.C0) == 1 && row >= r.R0+endGap && row <= r.R1-endGap
	}
	return false
}

// inRun reports whether the stretch of road r that cell col, row is
// beside has a fence there.
func inRun(seed uint32, r Rect, col, row int) bool {
	along, side, line := col-r.C0, row-r.R0, r.C0*1009+r.R0
	salt := seed + 3
	if r.C0 == r.C1 {
		along, side, salt = row-r.R0, col-r.C0, seed+4
	}
	k := floorDiv(along, stretchLen)
	h := terrain.Hash(salt, k*2+(side+1)/2, line)
	if int(h%100) >= fenceShare {
		return false
	}
	i := along - k*stretchLen
	return i >= 1+int(h>>8)%3 && i < stretchLen-1-int(h>>12)%3
}
