package constz

type Color struct {
	Name       string
	R, G, B, A uint8
}

func (c *Color) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R)
	r |= r << 8
	g = uint32(c.G)
	g |= g << 8
	b = uint32(c.B)
	b |= b << 8
	a = uint32(c.A)
	a |= a << 8
	return
}

func GraphColor(i int) *Color {
	n := len(colors)
	if i < 0 {
		return colors[i%n+n]
	}
	return colors[i%n]
}

// NOTE: カラーユニバーサルデザイン推奨配色セット https://jfly.uni-koeln.de/colorset/
// nolint: gochecknoglobals
var colors = []*Color{
	// ベースカラー
	{"LightRed", 255, 202, 191, 255},
	{"Cream", 255, 255, 128, 255},
	{"LightYellowGreen", 216, 242, 85, 255},
	{"LightSky", 191, 228, 255, 255},
	{"Beige", 255, 202, 128, 255},
	{"LightGreen", 119, 217, 168, 255},
	{"LightPurple", 201, 172, 230, 255},
	// アクセントカラー
	{"Red", 255, 75, 0, 255},
	{"Yellow", 255, 241, 0, 255},
	{"Green", 3, 175, 122, 255},
	{"Blue", 0, 90, 255, 255},
	{"Sky", 77, 196, 255, 255},
	{"Pink", 255, 128, 130, 255},
	{"Orange", 246, 170, 0, 255},
	{"Purple", 153, 0, 153, 255},
	{"Brown", 128, 64, 0, 255},
}
