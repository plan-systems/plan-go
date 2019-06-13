package plan // import "github.com/plan-systems/go-plan/plan"

import (
	//"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestErrors(t *testing.T) {

	_, err := ioutil.ReadFile("FileNotHere.yo")
	if os.IsNotExist(err) {

		num := 55
		perr := Errorf(err, GenericErrorFamily, "custom msg {val:%d}", num)

		actual := perr.Error()

		expected := "custom msg {val:55} {code:5000, err:{open FileNotHere.yo: no such file or directory}}"
		if actual != expected {
			t.Fatalf("got \"%v\", expected \"%v\"", actual, expected)
		}
	} else {
		t.Fatal()
	}

}
