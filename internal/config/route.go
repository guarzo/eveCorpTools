package config

// DataMode defines the type of data fetching
type DataMode int

const (
	Unset DataMode = iota
	YearToDate
	MonthToDate
	PreviousMonth
)

var DataModeToString = map[DataMode]string{
	YearToDate:    "ytd",
	MonthToDate:   "mtd",
	PreviousMonth: "lastM",
}

var StringToDataMode = map[string]DataMode{
	"ytd":   YearToDate,
	"mtd":   MonthToDate,
	"lastM": PreviousMonth,
}

// Route defines the type of route
type Route int

const (
	All Route = iota
	Config
	Snippets
)

var RouteToString = map[Route]string{
	All:      "All",
	Config:   "Config",
	Snippets: "Snippets",
}
