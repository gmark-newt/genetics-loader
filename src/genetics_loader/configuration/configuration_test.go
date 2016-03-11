package configuration

import (
	"testing"
)

func TestConfiguration(t *testing.T) {
	conf := GetConfiguration()

	t.Logf("Configuration: %#v", conf)
}

func TestIsDebug(t *testing.T) {
	d := IsDebug()

	if d != true && d != false {
		t.Fatal("Can't check IsDebug()")
	}
}
