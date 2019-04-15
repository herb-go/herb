package session

import "testing"

func TestFlag(t *testing.T) {
	store := New()
	s := NewSession("", store)
	s.SetFlag(FlagTemporay, true)
	if s.HasFlag(FlagTemporay) == false {
		t.Fatal(s.HasFlag(FlagTemporay))
	}
	s.SetFlag(FlagTemporay, false)
	if s.HasFlag(FlagTemporay) == true {
		t.Fatal(s.HasFlag(FlagTemporay))
	}

}
