package renamer

import (
	"testing"
)

func TestGetEpisodeName(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"Iron man (720).mkv", "Iron man"},
		{"Spiderman ().avi", "Spiderman"},
		{"Ant man [HDTV].avi", "Ant man"},
		{"[2018] Thor.mkv", "Thor"},
		{"masterchef.australia.s10e52.mkv", "masterchef australia"},
		{"()", ""},
		{"", ""},
		{"the.grand.tour.s03e01.720p.web.h264-skgtv.mkv", "the grand tour"},
		{"The.Legend.of.Korra.S03E03.The.Earth.Queen.REPACK.720p.WEBRip.x264.AAC.mp4", "The Legend of Korra"},
	}

	for _, c := range cases {
		episodeName := GetEpisodeName(c.in)
		if episodeName != c.want {
			t.Errorf("GetEpisodeName(%v) == %v, want %v", c.in, episodeName, c.want)
		}
	}
}

func TestGetEpisodeNumber(t *testing.T) {
	cases := []struct {
		in   string
		want EpisodeNumber
	}{
		{"DoctorWho S3E5", EpisodeNumber{3, 5}},
		{"Spiderman S1e4 SomethingInteresting", EpisodeNumber{1, 4}},
		{"TestEpidose s3E4", EpisodeNumber{3, 4}},
		{"Captain America", EpisodeNumber{-1, -1}},
	}

	for _, c := range cases {
		episodeName := GetEpisodeNumber(c.in)
		if episodeName != c.want {
			t.Errorf("GetEpisodeNumber(%v) == %v, want %v", c.in, episodeName, c.want)
		}
	}
}
