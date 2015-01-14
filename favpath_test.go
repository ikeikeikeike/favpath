package favpath

import (
	"testing"
)

func TestA(t *testing.T) {
	var file string
	f := NewFinder()

	file, err := f.Find("http://tensinyakimeshi.blog98.fc2.com/")
	println(file)
	if err != nil {
		t.Error(err)
	}
}

func TestB(t *testing.T) {
	var file string
	f := NewFinder()

	file, err := f.Find("http://matome.sekurosu.com/")
	println(file)
	if err != nil {
		t.Error(err)
	}
}
