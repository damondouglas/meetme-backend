package meetme

import (
	"os"
	"testing"
)

func Setup() {
	SetupServiceTest()
	SetupEndpointTest()
	SetupDatastoreTest()
}

func Teardown() {
	TeardownServiceTest()
	TeardownEndpointTest()
	TeardownDatastoreTest()
}

func TestMain(m *testing.M) {
	Setup()
	result := m.Run()
	Teardown()
	os.Exit(result)
}
