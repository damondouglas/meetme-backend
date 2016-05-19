package meetme

import (
	"github.com/jmcvetta/randutil"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/api/calendar/v3"
	"log"
	"os"
	"testing"
	"time"
)

var (
	srv       *calendar.Service
	calendars []string
)

func TestGetClient(t *testing.T) {
	t.Skip("Skipping TestGetClient")
	ctx := context.TODO()
	client := GetClient(ctx)
	head, _ := client.Head("https://www.googleapis.com/calendar/v3/users/me/calendarList")
	assert.Equal(t, 200, head.StatusCode, "Status code should be 200")
}

func TestListCalendars(t *testing.T) {
	t.Skip("Skipping TestListCalendars")
	name, _ := randutil.AlphaString(10)
	testCalendar := CreateCalendar(srv, name)
	calendars = append(calendars, testCalendar.Id)

	calList := ListCalendars(srv)
	containsCalendar := false
	for _, calendar := range calList {
		containsCalendar = calendar.Id == testCalendar.Id
		if containsCalendar {
			break
		}
	}
	assert.Equal(t, containsCalendar, true, "ListCalendars should list calendars in service account.")
}

func TestCreateCalendar(t *testing.T) {
	t.Skip("Skipping TestCreateCalendar")
	name, _ := randutil.AlphaString(10)
	testCalendar := CreateCalendar(srv, name)
	calendars = append(calendars, testCalendar.Id)
	assert.Equal(t, testCalendar.Summary, name, "CreateCalendar should create calendar with specified name")
}

func TestInsertEvent(t *testing.T) {
	t.Skip("Skipping TestInsertEvent")
	calendarName, _ := randutil.AlphaString(10)
	testCalendar := CreateCalendar(srv, calendarName)
	calendars = append(calendars, testCalendar.Id)

	start, end := RandDateRange()
	startStr := start.Format(time.RFC3339)
	endStr := end.Format(time.RFC3339)
	dateRange := new(DateRange)
	dateRange.Start = startStr
	dateRange.End = endStr
	testEvent := InsertEvent(srv, testCalendar, dateRange)

	assert.Equal(t, testEvent.Start.DateTime, startStr, "InsertEvent start time should be startStr")
	assert.Equal(t, testEvent.End.DateTime, endStr, "InsertEvent end time should be endStr")
}

func TestListEvents(t *testing.T) {
	t.Skip("Skipping TestListEvents")
	calendarName, _ := randutil.AlphaString(10)
	testCalendar := CreateCalendar(srv, calendarName)
	calendars = append(calendars, testCalendar.Id)

	start, end := RandDateRange()
	startStr := start.Format(time.RFC3339)
	endStr := end.Format(time.RFC3339)
	dateRange := new(DateRange)
	dateRange.Start = startStr
	dateRange.End = endStr
	testEvent := InsertEvent(srv, testCalendar, dateRange)

	eventList := ListEvents(srv, testCalendar.Id)
	containsEvent := false
	for _, event := range eventList {
		containsEvent = event.Id == testEvent.Id
		if containsEvent {
			break
		}
	}

	assert.True(t, containsEvent, "ListEvents should list events of given calendar")
}

func TestDeleteCalendar(t *testing.T) {
	t.Skip("Skipping TestDeleteCalendar")
	calendarName, _ := randutil.AlphaString(10)
	testCalendar := CreateCalendar(srv, calendarName)

	calList := ListCalendars(srv)
	containsCalendar := false
	for _, calendar := range calList {
		containsCalendar = calendar.Id == testCalendar.Id
		if containsCalendar {
			break
		}
	}
	assert.True(t, containsCalendar, "")

	DeleteCalendar(srv, testCalendar.Id)
	calList = ListCalendars(srv)

	for _, calendar := range calList {
		containsCalendar = calendar.Id == testCalendar.Id
		if containsCalendar {
			break
		}
	}
	assert.False(t, containsCalendar, "DeleteCalendar should remove calendar of given Id")
}

func TestUpdateEvent(t *testing.T) {
	t.Skip("Skipping TestUpdateEvent")
	calendarName, _ := randutil.AlphaString(10)
	testCalendar := CreateCalendar(srv, calendarName)
	calendars = append(calendars, testCalendar.Id)
	start, end := RandDateRange()
	startStr := start.Format(time.RFC3339)
	endStr := end.Format(time.RFC3339)
	dateRange := new(DateRange)
	dateRange.Start = startStr
	dateRange.End = endStr
	testEvent := InsertEvent(srv, testCalendar, dateRange)

	a := assert.New(t)
	a.Equal(testEvent.Summary, testCalendar.Summary, "Event should be same name as calendar")
	// change name without supplying dateRange arg
	calendarName, _ = randutil.AlphaString(10)
	evt := UpdateEvent(srv, testCalendar.Id, testEvent.Id, calendarName)
	a.NotEqual(testEvent.Summary, evt.Summary, "Updated event should not have the same name")
	// change dateRange
	start, end = RandDateRange()
	startStr = start.Format(time.RFC3339)
	endStr = end.Format(time.RFC3339)
	dateRange = new(DateRange)
	dateRange.Start = startStr
	dateRange.End = endStr

	evt = UpdateEvent(srv, testCalendar.Id, testEvent.Id, calendarName, dateRange)
	a.Equal(evt.Summary, calendarName, "Updated event should have the same name")
	a.NotEqual(testEvent.Start.DateTime, evt.Start.DateTime, "Updated event should not have the same start datetime")
	a.NotEqual(testEvent.End.DateTime, evt.End.DateTime, "Updated event should not have the same end datetime")
}

func TestUpdateCalendar(t *testing.T) {
	t.Skip("Skipping TestUpdateCalendar")
	calendarName, _ := randutil.AlphaString(10)
	testCalendar := CreateCalendar(srv, calendarName)
	calendars = append(calendars, testCalendar.Id)
	start, end := RandDateRange()
	startStr := start.Format(time.RFC3339)
	endStr := end.Format(time.RFC3339)
	dateRange1 := new(DateRange)
	dateRange1.Start = startStr
	dateRange1.End = endStr

	start, end = RandDateRange()
	startStr = start.Format(time.RFC3339)
	endStr = end.Format(time.RFC3339)
	dateRange2 := new(DateRange)
	dateRange2.Start = startStr
	dateRange2.End = endStr

	start, end = RandDateRange()
	startStr = start.Format(time.RFC3339)
	endStr = end.Format(time.RFC3339)
	dateRange3 := new(DateRange)
	dateRange3.Start = startStr
	dateRange3.End = endStr

	testEvent1 := InsertEvent(srv, testCalendar, dateRange1)
	testEvent2 := InsertEvent(srv, testCalendar, dateRange2)
	testEvent3 := InsertEvent(srv, testCalendar, dateRange3)

	a := assert.New(t)
	a.Equal(testCalendar.Summary, testEvent1.Summary, "Created event name should equal parent calendar name")
	a.Equal(testCalendar.Summary, testEvent2.Summary, "Created event name should equal parent calendar name")
	a.Equal(testCalendar.Summary, testEvent3.Summary, "Created event name should equal parent calendar name")

	calendarName, _ = randutil.AlphaString(10)
	updatedCalendar := UpdateCalendar(srv, testCalendar.Id, calendarName)

	a.Equal(updatedCalendar.Summary, calendarName, "Updated calendar should reflect changed name")

	eventList := ListEvents(srv, updatedCalendar.Id)
	for _, evt := range eventList {
		a.Equal(evt.Summary, updatedCalendar.Summary, "Calendar events should reflect changed name")
	}
}

func CleanUp() {
	log.Println("Cleaning up...")
	for _, calendarID := range calendars {
		DeleteCalendar(srv, calendarID)
		log.Printf("Deleted %s", calendarID)

	}
}

func TestMain(m *testing.M) {
	ctx := context.TODO()
	client := GetClient(ctx)
	srv = BuildService(client)
	result := m.Run()
	CleanUp()
	os.Exit(result)
}

func RandDateRange() (time.Time, time.Time) {
	duration := time.Duration(60) * time.Minute
	start := RandDate()
	end := start.Add(duration)
	return start, end
}

func RandDate() time.Time {
	year := randbetween(2000, 2020)
	month := time.Month(randbetween(1, 12))
	day := randbetween(1, 15)
	hour := randbetween(1, 23)
	min := randbetween(0, 3) * 15
	loc := time.UTC
	return time.Date(year, month, day, hour, min, 0, 0, loc)
}
