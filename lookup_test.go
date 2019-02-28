package renamer

import (
	"reflect"
	"testing"
)

func TestCullSeriesDataNoAliases(t *testing.T) {
	input := []SeriesData{SeriesData{Aliases: []string{}}}

	episodeName := CullSeriesData(input)
	if len(episodeName) != 0 {
		t.Errorf("CullSeriesData(%v) == %v, want %v", input, episodeName, "[]")
	}
}

func TestCullSeriesDataSomeAliases(t *testing.T) {
	input := []SeriesData{SeriesData{Aliases: []string{"Doctor who", "Doctor who (2005)"}},
		SeriesData{Aliases: []string{"Doctor who (2015)"}},
		SeriesData{Aliases: []string{"Doctor who specials"}},
		SeriesData{Aliases: []string{}}}

	episodeName := len(CullSeriesData(input))
	if episodeName != 3 {
		t.Errorf("CullSeriesData(%v) == %v, want %v", input, episodeName, 3)
	}
}

func TestCullSeriesDataOrderByDate(t *testing.T) {
	input := []SeriesData{SeriesData{Aliases: []string{"Doctor who"}, FirstAired: "2010-09-15"},
		SeriesData{Aliases: []string{"Doctor who (2015)"}, FirstAired: "2016-05-10"},
		SeriesData{Aliases: []string{"Doctor who specials"}, FirstAired: "2013-09-15"},
		SeriesData{Aliases: []string{}, FirstAired: "2016-09-15"}}

	want := []SeriesData{SeriesData{Aliases: []string{"Doctor who (2015)"}, FirstAired: "2016-05-10"},
		SeriesData{Aliases: []string{"Doctor who specials"}, FirstAired: "2013-09-15"},
		SeriesData{Aliases: []string{"Doctor who"}, FirstAired: "2010-09-15"}}

	seriesData := CullSeriesData(input)
	if !reflect.DeepEqual(seriesData, want) {
		t.Errorf("CullSeriesData(%v) == %v, want %v", input, seriesData, want)
	}
}

func TestCullSeriesDataOrderByStatus(t *testing.T) {
	input := []SeriesData{SeriesData{Aliases: []string{"Doctor who"}, FirstAired: "2016-05-10", Status: "Ended"},
		SeriesData{Aliases: []string{"Doctor who (2015)"}, FirstAired: "2016-05-10", Status: "Continuing"},
		SeriesData{Aliases: []string{"Doctor who specials"}, FirstAired: "2016-05-10", Status: "Ended"}}

	want := []SeriesData{SeriesData{Aliases: []string{"Doctor who (2015)"}, FirstAired: "2016-05-10", Status: "Continuing"},
		SeriesData{Aliases: []string{"Doctor who specials"}, FirstAired: "2016-05-10", Status: "Ended"},
		SeriesData{Aliases: []string{"Doctor who"}, FirstAired: "2016-05-10", Status: "Ended"}}

	seriesData := CullSeriesData(input)
	if !reflect.DeepEqual(seriesData, want) {
		t.Errorf("CullSeriesData(%v) == %v, want %v", input, seriesData, want)
	}
}
