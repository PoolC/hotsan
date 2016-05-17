package bot

import (
	"fmt"
	"testing"
)

func newMeu() *Meu {
	ret := &Meu{&BaseBot{nil, nil, "<@meu>", nil}, nil, nil}
	return ret
}

func TestSimpleReactionMeu(t *testing.T) {
	meu := newMeu()

	cases := []struct {
		text     string
		username string
		want     interface{}
	}{
		{"메우, 멱살", "", "사람은 일을 하고 살아야한다. 메우"},
		{"메우메우 펫탄탄", "", `메메메 메메메메 메우메우
메메메 메우메우
펫땅펫땅펫땅펫땅 다이스키`},
		{"메우메우 펫탄", "", nil},
	}

	for _, c := range cases {
		ret := meuMessageProcess(meu, newMessageEvent(c.text, c.username))
		if ret != c.want {
			t.Errorf("Simple reaction test failed. expected \"%q\" but \"%q\"", c.want, ret)
		}
	}
}

func TestCalcMeu(t *testing.T) {
	meu := newMeu()

	ret := meuMessageProcess(meu, newMessageEvent("계산하라 메우 1+1", "meu"))
	switch ret.(type) {
	case string:
		if ret.(string) != fmt.Sprintf("%f", 2.0) {
			t.Errorf("Error calc error")
		}
	default:
		t.Errorf("Error type")
	}
}
