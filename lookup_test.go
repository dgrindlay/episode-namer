package renamer

import (
	"reflect"
	"sort"
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

func TestNormalizeSeriesPriorities(t *testing.T) {
	seriesOne := SeriesData{ID: 1, SeriesName: "Series One"}
	seriesOnePriorities := []float64{20.0, 8.0}

	seriesTwo := SeriesData{ID: 2, SeriesName: "Series Two"}
	seriesTwoPriorities := []float64{15.0, 6.5}

	input := []SeriesPriority{SeriesPriority{series: &seriesOne, priorities: seriesOnePriorities}, SeriesPriority{series: &seriesTwo, priorities: seriesTwoPriorities}}

	normalSeriesOnePriorities := []float64{100.0, 100.0}
	normalSeriesTwoPriorities := []float64{75.0, 81.25}
	want := []SeriesPriority{SeriesPriority{series: &seriesOne, priorities: normalSeriesOnePriorities}, SeriesPriority{series: &seriesTwo, priorities: normalSeriesTwoPriorities}}

	normalizedPriorities, err := normalizeSeriesPriorities(input)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(normalizedPriorities, want) {
		t.Errorf("normalizePriorities(%v) == %v, want %v", input, normalizedPriorities, want)
	}
}

func TestSortNormalizedSeriesPriorities(t *testing.T) {
	seriesOne := SeriesData{ID: 1, SeriesName: "Series One"}
	seriesOnePriorities := []float64{55.0, 4.0}

	seriesTwo := SeriesData{ID: 2, SeriesName: "Series Two"}
	seriesTwoPriorities := []float64{30.0, 9.0}

	orderedSeriesPointers := []*SeriesData{&seriesOne, &seriesTwo}
	input := []SeriesPriority{SeriesPriority{series: &seriesOne, priorities: seriesOnePriorities}, SeriesPriority{series: &seriesTwo, priorities: seriesTwoPriorities}}

	normalizedSeries, err := normalizeSeriesPriorities(input)
	if err != nil {
		t.Error(err)
	}

	seriesPriorityMap := make(map[*SeriesData]int)
	for _, seriesPriority := range normalizedSeries {
		seriesPriorityMap[seriesPriority.series] = int(0.2*seriesPriority.priorities[0] + 0.8*seriesPriority.priorities[1])
	}

	orderedSeriesList := []SeriesData{}
	sort.Sort(byRating{seriesData: orderedSeriesPointers, priorityMap: seriesPriorityMap})
	for _, seriesPointer := range orderedSeriesPointers {
		orderedSeriesList = append(orderedSeriesList, *seriesPointer)
	}

	want := []SeriesData{SeriesData{ID: 2, SeriesName: "Series Two"}, SeriesData{ID: 1, SeriesName: "Series One"}}

	if !reflect.DeepEqual(orderedSeriesList, want) {
		t.Errorf("normalizeSeriesPriorities(%v) == %v, want %v", input, orderedSeriesList, want)
	}
}
