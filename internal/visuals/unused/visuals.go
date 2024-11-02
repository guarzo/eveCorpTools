package unused

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/render"
)

var colors = []string{"red", "blue", "green", "yellow", "purple", "cyan", "magenta", "lime", "black", "navy", "aqua", "maroon", "olive", "silver", "gray", "fuchsia", "white", "coral", "salmon"}

func newBarChart(title string, tooltip bool) *charts.Bar {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "1200px", // Increase width as needed
			Height: "800px",  // Increase height as needed
		}),
		charts.WithTitleOpts(opts.Title{
			Title: title,
			Left:  "0%",
			TitleStyle: &opts.TextStyle{
				Color: "white",
			},
		}),
		charts.WithYAxisOpts(opts.YAxis{}),
		charts.WithXAxisOpts(opts.XAxis{
			AxisLabel: &opts.AxisLabel{
				Rotate:     45,
				Interval:   "0",
				FontSize:   14,
				FontWeight: "bold",
				Color:      "white",
			},
		}),
		charts.WithGridOpts(opts.Grid{
			Bottom: "30%", // Increase this value to give more space for the labels
		}),
		charts.WithLegendOpts(opts.Legend{Show: opts.Bool(false)}),
	)
	if tooltip {
		bar.SetGlobalOptions(charts.WithTooltipOpts(opts.Tooltip{Show: opts.Bool(tooltip), Formatter: opts.FuncOpts("function (params) {return params.seriesName + ': ' + params.value;}")}))
	} else {
		bar.SetGlobalOptions(charts.WithTooltipOpts(opts.Tooltip{
			Show:      opts.Bool(true),
			Formatter: opts.FuncOpts("function (params) {return params.name + ': ' + params.value;}"),
		}))
	}
	return bar
}

func truncateString(str string, num int) string {
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		str = str[0:num] + "..."
	}
	return str
}

// Helper function to create a new file and render a page to it
func renderBar(bar *charts.Bar, filePath string) error {
	page := components.NewPage()
	page.AddCharts(bar)
	page.PageTitle = "Zoolanders TPS Reports"
	page.AddCustomizedCSSAssets("/custom.css")

	// Save the chart to an HTML file
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()
	return page.Render(f)
}

func renderToHtml(charts ...interface{}) template.HTML {
	var buf bytes.Buffer

	for _, c := range charts {
		buf.WriteString(`<div class="chart-content">`)
		r := c.(render.Renderer)
		err := r.Render(&buf)
		if err != nil {
			log.Printf("Failed to render chart: %s", err)
			return ""
		}
		buf.WriteString(`</div>`)
	}

	return template.HTML(buf.String())
}

// Helper function to prepare data for a bar chart
func prepareBarChartData(data map[string]int) ([]string, []opts.BarData) {
	var names []string
	var counts []opts.BarData
	for name, count := range data {
		names = append(names, name)
		counts = append(counts, opts.BarData{Value: count})
	}
	return names, counts
}
