package chls

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"
)

func Test_UnmarshalSession(t *testing.T) {

	file, err := os.Open("sample_session.chlsj")
	if err != nil {
		t.Fatal(err)
	}

	var session Session
	decoder := json.NewDecoder(bufio.NewReader(file))
	err = decoder.Decode(&session)
	if err != nil {
		t.Fatal(err)
	}

	if len(session) < 2 {
		t.Errorf("unexpected session size: %d", len(session))
	}

	pnum := session[0].ActualPort
	if pnum != 443 {
		t.Errorf("unexpected port number: %d", pnum)
	}

	// TODO: perform more validation
}
