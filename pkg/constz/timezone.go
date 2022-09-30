package constz

import "time"

const offsetJST = 9 * 60 * 60

// nolint: gochecknoglobals
var tz = map[string]*time.Location{
	"Asia/Tokyo": time.FixedZone("Asia/Tokyo", offsetJST),
	"JST":        time.FixedZone("Asia/Tokyo", offsetJST),
}

func TimeZone(zone string) *time.Location {
	if loc := tz[zone]; loc != nil {
		return loc
	}

	return time.UTC
}
