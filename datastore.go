package meetme

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"log"
)

const meetieDataType = "Meetie"
const userLimit = 10

// Greeting is a datastore entity that represents a single greeting.
// It also serves as (a part of) a response of GreetingService.
type Greeting struct {
	Key     *datastore.Key `json:"id" datastore:"-"`
	Author  string         `json:"author"`
	Content string         `json:"content" datastore:",noindex"`
}

// Meetie is the data type stored in the datastore to associate users with calendars
type Meetie struct {
	Email      string `json:"email"`
	CalendarID string `json:"calendarID"`
}

func getUserCalendarQuery(userEmail string) *datastore.Query {
	return datastore.NewQuery(meetieDataType).Filter("Email =", userEmail)
}

func getUserCalendarCount(c context.Context, userEmail string) int {
	q := getUserCalendarQuery(userEmail)
	count, _ := q.Count(c)
	return count
}

func getUserReachedLimit(c context.Context, userEmail string) bool {
	count := getUserCalendarCount(c, userEmail)
	return count >= userLimit
}

func associateUserWithCalendar(c context.Context, userEmail string, calendarID string) error {

	meetie := &Meetie{
		Email:      userEmail,
		CalendarID: calendarID,
	}

	// key := datastore.NewKey(c, meetieDataType, "", 1, nil)
	key := datastore.NewKey(c, meetieDataType, calendarID, 0, nil)
	if _, err := datastore.Put(c, key, meetie); err != nil {
		log.Fatal(err)
	}

	m := Meetie{}
	if err := datastore.Get(c, key, &m); err != nil {
		log.Fatal(err)
	}

	return nil
}

func getCalendarsFromUser(c context.Context, userEmail string) ([]string, error) {
	calendarIDList := []string{}
	q := getUserCalendarQuery(userEmail)
	itr := q.Run(c)
	for {
		var m Meetie
		_, err := itr.Next(&m)
		if err == datastore.Done {
			break
		}

		if err != nil {
			return calendarIDList, err
		}

		calendarIDList = append(calendarIDList, m.CalendarID)
	}

	return calendarIDList, nil
}
