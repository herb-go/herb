package notificationpartmsgpack

import (
	"bytes"
	"testing"

	"github.com/herb-go/herb/notification"
)

func TestPart(t *testing.T) {
	n, err := notification.NewPartedNotificationWithID()
	if err != nil {
		t.Fatal(err)
	}
	teststring := "test"
	var msgstring string
	NotificationPartTitle.Set(n, teststring)
	msgstring, err = NotificationPartTitle.Get(n)
	if err != nil {
		t.Fatal(err)
	}
	if msgstring != teststring {
		t.Error(msgstring)
	}

	NotificationPartSummary.Set(n, teststring)
	msgstring, err = NotificationPartSummary.Get(n)
	if err != nil {
		t.Fatal(err)
	}
	if msgstring != teststring {
		t.Error(msgstring)
	}

	NotificationPartText.Set(n, teststring)
	msgstring, err = NotificationPartText.Get(n)
	if err != nil {
		t.Fatal(err)
	}
	if msgstring != teststring {
		t.Error(msgstring)
	}
	NotificationPartHtml.Set(n, teststring)
	msgstring, err = NotificationPartHtml.Get(n)
	if err != nil {
		t.Fatal(err)
	}
	if msgstring != teststring {
		t.Error(msgstring)
	}

	var testbytes = []byte("testbytes")
	var msgbytes []byte
	NotificationPartAttachment.Set(n, testbytes)
	msgbytes, err = NotificationPartAttachment.Get(n)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(msgbytes, testbytes) != 0 {
		t.Error(msgbytes)
	}
	NotificationPartPicture.Set(n, testbytes)
	msgbytes, err = NotificationPartPicture.Get(n)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(msgbytes, testbytes) != 0 {
		t.Error(msgbytes)
	}
	NotificationPartURL.Set(n, teststring)
	msgstring, err = NotificationPartURL.Get(n)
	if err != nil {
		t.Fatal(err)
	}
	if msgstring != teststring {
		t.Error(msgstring)
	}
	NotificationPartAttachmentURL.Set(n, teststring)
	msgstring, err = NotificationPartAttachmentURL.Get(n)
	if err != nil {
		t.Fatal(err)
	}
	if msgstring != teststring {
		t.Error(msgstring)
	}
	NotificationPartAttachmentFilename.Set(n, teststring)
	msgstring, err = NotificationPartAttachmentFilename.Get(n)
	if err != nil {
		t.Fatal(err)
	}
	if msgstring != teststring {
		t.Error(msgstring)
	}
	NotificationPartPictureFilename.Set(n, teststring)
	msgstring, err = NotificationPartPictureFilename.Get(n)
	if err != nil {
		t.Fatal(err)
	}
	if msgstring != teststring {
		t.Error(msgstring)
	}
	NotificationPartPictureURL.Set(n, teststring)
	msgstring, err = NotificationPartPictureURL.Get(n)
	if err != nil {
		t.Fatal(err)
	}
	if msgstring != teststring {
		t.Error(msgstring)
	}
	var teststringlist = []string{"test1", "test2"}
	var msgteststring []string

	NotificationPartTitleList.Set(n, teststringlist)
	msgteststring, err = NotificationPartTitleList.Get(n)
	if err != nil {
		t.Fatal(err)
	}
	for k := range teststringlist {
		if teststringlist[k] != msgteststring[k] {
			t.Error(k)
		}
	}

	NotificationPartSummaryList.Set(n, teststringlist)
	msgteststring, err = NotificationPartSummaryList.Get(n)
	if err != nil {
		t.Fatal(err)
	}
	for k := range teststringlist {
		if teststringlist[k] != msgteststring[k] {
			t.Error(k)
		}
	}

	NotificationPartTextList.Set(n, teststringlist)
	msgteststring, err = NotificationPartTextList.Get(n)
	if err != nil {
		t.Fatal(err)
	}
	for k := range teststringlist {
		if teststringlist[k] != msgteststring[k] {
			t.Error(k)
		}
	}

	NotificationPartHtmlList.Set(n, teststringlist)
	msgteststring, err = NotificationPartHtmlList.Get(n)
	if err != nil {
		t.Fatal(err)
	}
	for k := range teststringlist {
		if teststringlist[k] != msgteststring[k] {
			t.Error(k)
		}
	}

	var testbyteslist = [][]byte{[]byte("testbytest1"), []byte("testbytest2")}
	var msgbyteslist [][]byte

	NotificationPartAttachmentList.Set(n, testbyteslist)
	msgbyteslist, err = NotificationPartAttachmentList.Get(n)
	if err != nil {
		t.Fatal(err)
	}
	for k := range msgbyteslist {
		if bytes.Compare(msgbyteslist[k], msgbyteslist[k]) != 0 {
			t.Error(k)
		}
	}
	NotificationPartPictureList.Set(n, testbyteslist)
	msgbyteslist, err = NotificationPartPictureList.Get(n)
	if err != nil {
		t.Fatal(err)
	}
	for k := range msgbyteslist {
		if bytes.Compare(msgbyteslist[k], msgbyteslist[k]) != 0 {
			t.Error(k)
		}
	}

	NotificationPartAttachmentURLList.Set(n, teststringlist)
	msgteststring, err = NotificationPartHtmlList.Get(n)
	if err != nil {
		t.Fatal(err)
	}
	for k := range teststringlist {
		if teststringlist[k] != msgteststring[k] {
			t.Error(k)
		}
	}
	NotificationPartPictureURLList.Set(n, teststringlist)
	msgteststring, err = NotificationPartHtmlList.Get(n)
	if err != nil {
		t.Fatal(err)
	}
	for k := range teststringlist {
		if teststringlist[k] != msgteststring[k] {
			t.Error(k)
		}
	}
	NotificationPartURLList.Set(n, teststringlist)
	msgteststring, err = NotificationPartHtmlList.Get(n)
	if err != nil {
		t.Fatal(err)
	}
	for k := range teststringlist {
		if teststringlist[k] != msgteststring[k] {
			t.Error(k)
		}
	}
	var teststringmap = map[string]string{"key1": "test1", "key2": "test2"}
	var msgstringmap map[string]string

	NotificationPartTitleMap.Set(n, teststringmap)
	msgstringmap, err = NotificationPartTitleMap.Get(n)
	if err != nil {
		t.Fatal(err)
	}
	for k := range teststringmap {
		if teststringmap[k] != msgstringmap[k] {
			t.Error(k)
		}
	}
	NotificationPartSummaryMap.Set(n, teststringmap)
	msgstringmap, err = NotificationPartSummaryMap.Get(n)
	if err != nil {
		t.Fatal(err)
	}
	for k := range teststringmap {
		if teststringmap[k] != msgstringmap[k] {
			t.Error(k)
		}
	}

	NotificationPartTextMap.Set(n, teststringmap)
	msgstringmap, err = NotificationPartTextMap.Get(n)
	if err != nil {
		t.Fatal(err)
	}
	for k := range teststringmap {
		if teststringmap[k] != msgstringmap[k] {
			t.Error(k)
		}
	}

	NotificationPartHtmlMap.Set(n, teststringmap)
	msgstringmap, err = NotificationPartHtmlMap.Get(n)
	if err != nil {
		t.Fatal(err)
	}
	for k := range teststringmap {
		if teststringmap[k] != msgstringmap[k] {
			t.Error(k)
		}
	}

	var testbytesmap = map[string][]byte{"key1": []byte("test1"), "key2": []byte("test2")}
	var msgbytesmap map[string][]byte

	NotificationPartAttachmentMap.Set(n, testbytesmap)
	msgbytesmap, err = NotificationPartAttachmentMap.Get(n)
	if err != nil {
		t.Fatal(err)
	}
	for k := range testbytesmap {
		if bytes.Compare(testbytesmap[k], msgbytesmap[k]) != 0 {
			t.Error(k)
		}
	}

	NotificationPartPictureMap.Set(n, testbytesmap)
	msgbytesmap, err = NotificationPartPictureMap.Get(n)
	if err != nil {
		t.Fatal(err)
	}
	for k := range testbytesmap {
		if bytes.Compare(testbytesmap[k], msgbytesmap[k]) != 0 {
			t.Error(k)
		}
	}
	NotificationPartAttachmentURLMap.Set(n, teststringmap)
	msgstringmap, err = NotificationPartAttachmentURLMap.Get(n)
	if err != nil {
		t.Fatal(err)
	}
	for k := range teststringmap {
		if teststringmap[k] != msgstringmap[k] {
			t.Error(k)
		}
	}

	NotificationPartPictureURLMAp.Set(n, teststringmap)
	msgstringmap, err = NotificationPartPictureURLMAp.Get(n)
	if err != nil {
		t.Fatal(err)
	}
	for k := range teststringmap {
		if teststringmap[k] != msgstringmap[k] {
			t.Error(k)
		}
	}

	NotificationPartURLMap.Set(n, teststringmap)
	msgstringmap, err = NotificationPartURLMap.Get(n)
	if err != nil {
		t.Fatal(err)
	}
	for k := range teststringmap {
		if teststringmap[k] != msgstringmap[k] {
			t.Error(k)
		}
	}
	type TestStruct struct {
		DataString string
		DataBytes  []byte
	}
	var TestNotificationInterface = NewInterfacePart("testinterface")
	var testinterface = &TestStruct{
		DataString: "teststring",
		DataBytes:  []byte("testbytes"),
	}
	var msginterface = &TestStruct{}
	err = TestNotificationInterface.Marshal(n, testinterface)
	if err != nil {
		t.Fatal(err)
	}
	err = TestNotificationInterface.Unmarshal(n, msginterface)
	if err != nil {
		t.Fatal(err)
	}
	if testinterface.DataString != msginterface.DataString || bytes.Compare(testinterface.DataBytes, msginterface.DataBytes) != 0 {
		t.Error(testinterface)
	}
}
