package domain

import (
	"image/color"
	"io"
	"strconv"
	"time"

	"github.com/kunitsuinc/ccc/pkg/constz"
	"github.com/kunitsuinc/ccc/pkg/errorz"
	"github.com/kunitsuinc/ccc/pkg/log"
	"github.com/kunitsuinc/util.go/mathz"
	"github.com/kunitsuinc/util.go/slice"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// nolint: cyclop,funlen
func Plot1280x720(
	target io.Writer,
	graphTitle string,
	xLabelText string,
	yLabelText string,
	xAxisPointsCount int,
	from time.Time,
	to time.Time,
	tz *time.Location,
	legends []string,
	legendValuesMap map[string]plotter.Values,
	imageFormat string,
) error {
	mono := font.Font{Typeface: "Liberation", Variant: "Mono"}
	plot.DefaultFont = mono
	plotter.DefaultFont = mono

	p := plot.New()
	p.Title.Text = graphTitle
	p.X.Label.Text = xLabelText
	p.Y.Label.Text = yLabelText

	const graphWidth = (1280 / 4) * 3 // NOTE: 1280 pixel / 4 * 3 = 960
	const graphHight = (graphWidth / 16) * 9
	barChartWidth := vg.Points((graphWidth - 100) / float64(xAxisPointsCount)) // NOTE: グラフの幅から固定長(95)を引いて X 軸の値数で割る

	previousBarChart := (*plotter.BarChart)(nil)
	for i, legend := range legends {
		barChart, err := plotter.NewBarChart(legendValuesMap[legend], barChartWidth)
		if err != nil {
			return errorz.Errorf("plotter.NewBarChart: %w", err)
		}
		barChart.Width = barChartWidth
		barChart.LineStyle.Width = vg.Length(0) // NOTE: グラフの枠線の太さを 0 にする
		barChart.Color = constz.GraphColor(i)
		p.Legend.Add(legend, barChart)

		if previousBarChart != nil {
			barChart.StackOn(previousBarChart)
		}

		p.Add(barChart)

		previousBarChart = barChart
	}

	grid := plotter.NewGrid()
	grid.Horizontal.Color = color.Black
	grid.Horizontal.Dashes = []vg.Length{vg.Length(5)}
	p.Add(grid)

	var xLabels []string
	for i := 0; i <= xAxisPointsCount; i++ {
		var x string
		if i%7 == 1 { // NOTE: 余り 1 -> 1 日前, 1+7 日前, 1+14 日前 1+21 日前にラベルを付与する
			x = to.In(tz).AddDate(0, 0, -i).Format(constz.DateOnly)
			log.Debugf("label: %s", x)
		}
		xLabels = append(xLabels, x)
	}
	xLabels = slice.Reverse(xLabels)
	p.NominalX(xLabels...)

	p.Legend.Top = true
	p.Legend.Left = true
	p.Legend.XOffs = 10
	p.Legend.YOffs = -10
	legendHight := float64(p.Legend.TextStyle.Height("C")) * 8
	legendsHight := legendHight * float64(len(legends))
	log.Debugf("legendHight=%f, legendsHight=%f", legendHight, legendsHight)
	p.Y.Min = 0
	p.Y.Max += legendsHight // NOTE: グラフと Legend が被らないように、 Legend の高さ (文字 C の高さで計算) * Legend 数を足して、 Y 軸の高さを確保している
	p.Y.Tick.Marker = plot.ConstantTicks(func() []plot.Tick {
		var ticks []plot.Tick
		unit := func() int { // NOTE: どの単位で Y 軸グリッドを入れるか。 1, 5, 10, 50, 100, 500, 1000, 5000, 10000, 50000, ... のどれかが入る
			unit := 1
			for i := p.Y.Max; i > 10; i = p.Y.Max / float64(unit) {
				if mathz.IsPow10(float64(unit)) {
					unit *= 5
					continue
				}
				unit *= 2
			}
			return unit
		}()
		numOfYGrid := int(p.Y.Max / float64(unit)) // NOTE: グラフ内に何本 Y 軸 grid を入れるか
		for i := 0; i <= numOfYGrid; i++ {
			value := unit * i
			ticks = append(
				ticks,
				plot.Tick{
					// NOTE: unit==500 の場合は 500, 1000, 1500, 2000, ... に Y 軸グリッドを書き込む
					Label: strconv.Itoa(value),
					Value: float64(value),
				},
			)
		}
		return ticks
	}())

	wt, err := p.WriterTo(graphWidth, graphHight, imageFormat)
	if err != nil {
		return errorz.Errorf("(*plot.Plot).WriterTo: %w", err)
	}

	if _, err := wt.WriteTo(target); err != nil {
		return errorz.Errorf("(io.WriterTo).WriteTo: %w", err)
	}

	return nil
}
