package renamer

import (
	"reflect"
	"testing"
)

func TestCullSeriesDataSomeAliases(t *testing.T) {
	input := []SeriesData{SeriesData{Aliases: []string{"Doctor who", "Doctor who (2005)"}},
		SeriesData{Aliases: []string{"Doctor who (2015)"}},
		SeriesData{Aliases: []string{"Doctor who specials"}},
		SeriesData{Aliases: []string{}}}

	episodeName := len(CullSeriesData(input))
	if episodeName != 4 {
		t.Errorf("CullSeriesData(%v) == %v, want %v", input, episodeName, 4)
	}
}

func TestCullSeriesDataOrderByID(t *testing.T) {
	input := []SeriesData{SeriesData{Aliases: []string{"Doctor who"}, FirstAired: "2010-09-15", ID: 2},
		SeriesData{Aliases: []string{"Doctor who (2015)"}, FirstAired: "2016-05-10", ID: 1},
		SeriesData{Aliases: []string{"Doctor who specials"}, FirstAired: "2013-09-15", ID: 3},
		SeriesData{Aliases: []string{}, FirstAired: "2016-09-15", ID: 5}}

	want := []SeriesData{SeriesData{Aliases: []string{"Doctor who (2015)"}, FirstAired: "2016-05-10", ID: 1},
		SeriesData{Aliases: []string{"Doctor who"}, FirstAired: "2010-09-15", ID: 2},
		SeriesData{Aliases: []string{"Doctor who specials"}, FirstAired: "2013-09-15", ID: 3},
		SeriesData{Aliases: []string{}, FirstAired: "2016-09-15", ID: 5}}

	seriesData := CullSeriesData(input)
	if !reflect.DeepEqual(seriesData, want) {
		t.Errorf("CullSeriesData(%v) == %v, want %v", input, seriesData, want)
	}
}
