package meetme

import (
	// "google.golang.org/appengine/aetest"
	"github.com/jmcvetta/randutil"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"log"
	"testing"
	// "google.golang.org/appengine/datastore"
	"google.golang.org/appengine/aetest"
)

var (
	datastoreTestContext context.Context
)

func SetupDatastoreTest() {
	var err error
	datastoreTestContext, _, err = aetest.NewContext()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Setup Datastore Test")
}

func TeardownDatastoreTest() {
	log.Println("Teardown Datastore Test")
}

func TestAssociateUserWithCalendar(t *testing.T) {
	// t.Skip("Skipping TestAssociateUserWithCalendar")
	a := assert.New(t)

	calendarID := generateMockCalendarID()
	email := generateMockEmail()
	_ = associateUserWithCalendar(datastoreTestContext, email, calendarID)
	a.False(getUserReachedLimit(datastoreTestContext, email), "User limit should not yet be reached")

	calendarID = generateMockCalendarID()
	_ = associateUserWithCalendar(datastoreTestContext, email, calendarID)
	a.False(getUserReachedLimit(datastoreTestContext, email), "User limit should not yet be reached")

	calendarID = generateMockCalendarID()
	_ = associateUserWithCalendar(datastoreTestContext, email, calendarID)
	a.False(getUserReachedLimit(datastoreTestContext, email), "User limit should not yet be reached")

	calendarID = generateMockCalendarID()
	_ = associateUserWithCalendar(datastoreTestContext, email, calendarID)
	a.False(getUserReachedLimit(datastoreTestContext, email), "User limit should not yet be reached")

	calendarID = generateMockCalendarID()
	_ = associateUserWithCalendar(datastoreTestContext, email, calendarID)
	a.False(getUserReachedLimit(datastoreTestContext, email), "User limit should not yet be reached")

	calendarID = generateMockCalendarID()
	_ = associateUserWithCalendar(datastoreTestContext, email, calendarID)
	a.False(getUserReachedLimit(datastoreTestContext, email), "User limit should not yet be reached")

	calendarID = generateMockCalendarID()
	_ = associateUserWithCalendar(datastoreTestContext, email, calendarID)
	a.False(getUserReachedLimit(datastoreTestContext, email), "User limit should not yet be reached")

	calendarID = generateMockCalendarID()
	_ = associateUserWithCalendar(datastoreTestContext, email, calendarID)
	a.False(getUserReachedLimit(datastoreTestContext, email), "User limit should not yet be reached")

	calendarID = generateMockCalendarID()
	_ = associateUserWithCalendar(datastoreTestContext, email, calendarID)
	a.False(getUserReachedLimit(datastoreTestContext, email), "User limit should be reached")

	calendarID = generateMockCalendarID()
	_ = associateUserWithCalendar(datastoreTestContext, email, calendarID)
	a.True(getUserReachedLimit(datastoreTestContext, email), "User limit should not yet be reached")

	email = generateMockEmail()
	_ = associateUserWithCalendar(datastoreTestContext, email, calendarID)
	a.False(getUserReachedLimit(datastoreTestContext, email), "User limit should be reached")
}

func generateMockCalendarID() string {
	calID, _ := randutil.AlphaString(16)
	return calID
}

func generateMockEmail() string {
	username, _ := randutil.AlphaString(8)
	return username + "@example.com"
}
