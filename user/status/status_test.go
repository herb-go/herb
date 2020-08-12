package status

import "testing"

func TestStatus(t *testing.T) {
	var s Service = NewService()
	ok, err := s.IsAvailable(StatusNormal)
	if ok == true || err != nil {
		t.Fatal()
	}
	label, err := s.Label(StatusNormal)
	if label != "" || err != nil {
		t.Fatal()
	}
	ok, err = NormalOrBannedService.IsAvailable(StatusNormal)
	if ok != true || err != nil {
		t.Fatal()
	}
	label, err = NormalOrBannedService.Label(StatusNormal)
	if label != StatusLabelNormal || err != nil {
		t.Fatal()
	}
	ok, err = NormalOrBannedService.IsAvailable(StatusBanned)
	if ok != false || err != nil {
		t.Fatal()
	}
	label, err = NormalOrBannedService.Label(StatusBanned)
	if label != StatusLabelBanned || err != nil {
		t.Fatal()
	}
	ok, err = NormalOrBannedService.IsAvailable(999)
	if ok != false || err != nil {
		t.Fatal()
	}
	label, err = NormalOrBannedService.Label(999)
	if label != "" || err != nil {
		t.Fatal()
	}
}
