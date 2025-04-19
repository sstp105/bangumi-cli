package libs

import (
	"reflect"
	"testing"
)

func TestSplitToSlice(t *testing.T) {
	got := SplitToSlice("CHT,  1080P, AVC, MP4", ",")
	want := []string{"CHT", "1080P", "AVC", "MP4"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got:%v, want:%v", got, want)
	}
}
