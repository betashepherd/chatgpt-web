package archivelib

import (
	"testing"
)

func TestUnzip(t *testing.T) {
	err := Unzip("../../test.zip", "out", func(p float64) {
		t.Log("unzip progress: ", p)
	})
	if err != nil {
		t.Error(err.Error())
	}
	// if err := archiver.NewZip().Unarchive("../../test.zip", "out"); err != nil {
	// 	t.Error(err.Error())
	// }

}
