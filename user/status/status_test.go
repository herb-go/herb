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
	ok, err = NoramlOrBannedService.IsAvailable(StatusNormal)
	if ok != true || err != nil {
		t.Fatal()
	}
	label, err = NoramlOrBannedService.Label(StatusNormal)
	if label != StatusLabelNormal || err != nil {
		t.Fatal()
	}
	ok, err = NoramlOrBannedService.IsAvailable(StatusBanned)
	if ok != false || err != nil {
		t.Fatal()
	}
	label, err = NoramlOrBannedService.Label(StatusBanned)
	if label != StatusLabelBanned || err != nil {
		t.Fatal()
	}
	ok, err = NoramlOrBannedService.IsAvailable(999)
	if ok != false || err != nil {
		t.Fatal()
	}
	label, err = NoramlOrBannedService.Label(999)
	if label != "" || err != nil {
		t.Fatal()
	}
}
