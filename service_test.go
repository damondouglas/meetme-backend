package meetme

import (
	"github.com/jmcvetta/randutil"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/api/calendar/v3"
	"log"
	"testing"
	"time"
)

var (
	testService *calendar.Service
	calendars   []string
	serviceImpl = new(serviceimpl)
)

func TestGetClient(t *testing.T) {
	t.Skip("Skipping TestGetClient")
	ctx := context.TODO()
	client := serviceImpl.GetClient(ctx)
	head, _ := client.Head("https://www.googleapiserviceImpl.com/calendar/v3/users/me/calendarList")
	assert.Equal(t, 200, head.StatusCode, "Status code should be 200")
}

func TestListCalendars(t *testing.T) {
	t.Skip("Skipping TestListCalendars")
	name, _ := randutil.AlphaString(10)
	testCalendar, err := serviceImpl.CreateCalendar(testService, name)
	if err != nil {
		t.Fatal(err)
	}
	calendars = append(calendars, testCalendar.Id)

	calList, err := serviceImpl.ListCalendars(testService)
	if err != nil {
		t.Fatal(err)
	}

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
	testCalendar, err := serviceImpl.CreateCalendar(testService, name)
	if err != nil {
		t.Fatal(err)
	}

	calendars = append(calendars, testCalendar.Id)
	assert.Equal(t, testCalendar.Summary, name, "CreateCalendar should create calendar with specified name")
}

func TestInsertEvent(t *testing.T) {
	t.Skip("Skipping TestInsertEvent")
	calendarName, _ := randutil.AlphaString(10)
	testCalendar, err := serviceImpl.CreateCalendar(testService, calendarName)
	if err != nil {
		t.Fatal(err)
	}
	calendars = append(calendars, testCalendar.Id)

	start, end := RandDateRange()
	startStr := start.Format(time.RFC3339)
	endStr := end.Format(time.RFC3339)
	dateRange := new(DateRange)
	dateRange.Start = startStr
	dateRange.End = endStr
	testEvent, err := serviceImpl.InsertEvent(testService, testCalendar, dateRange)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, testEvent.Start.DateTime, startStr, "InsertEvent start time should be startStr")
	assert.Equal(t, testEvent.End.DateTime, endStr, "InsertEvent end time should be endStr")
}

func TestListEvents(t *testing.T) {
	t.Skip("Skipping TestListEvents")
	calendarName, _ := randutil.AlphaString(10)
	testCalendar, err := serviceImpl.CreateCalendar(testService, calendarName)
	if err != nil {
		t.Fatal(err)
	}
	calendars = append(calendars, testCalendar.Id)

	start, end := RandDateRange()
	startStr := start.Format(time.RFC3339)
	endStr := end.Format(time.RFC3339)
	dateRange := new(DateRange)
	dateRange.Start = startStr
	dateRange.End = endStr
	testEvent, err := serviceImpl.InsertEvent(testService, testCalendar, dateRange)
	if err != nil {
		t.Fatal(err)
	}

	eventList, err := serviceImpl.ListEvents(testService, testCalendar.Id)
	if err != nil {
		t.Fatal(err)
	}

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
	testCalendar, err := serviceImpl.CreateCalendar(testService, calendarName)
	if err != nil {
		t.Fatal(err)
	}

	calList, err := serviceImpl.ListCalendars(testService)
	if err != nil {
		t.Fatal(err)
	}
	containsCalendar := false
	for _, calendar := range calList {
		containsCalendar = calendar.Id == testCalendar.Id
		if containsCalendar {
			break
		}
	}
	assert.True(t, containsCalendar, "")

	serviceImpl.DeleteCalendar(testService, testCalendar.Id)
	calList, err = serviceImpl.ListCalendars(testService)
	if err != nil {
		t.Fatal(err)
	}

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
	testCalendar, err := serviceImpl.CreateCalendar(testService, calendarName)
	if err != nil {
		t.Fatal(err)
	}
	calendars = append(calendars, testCalendar.Id)
	start, end := RandDateRange()
	startStr := start.Format(time.RFC3339)
	endStr := end.Format(time.RFC3339)
	dateRange := new(DateRange)
	dateRange.Start = startStr
	dateRange.End = endStr
	testEvent, err := serviceImpl.InsertEvent(testService, testCalendar, dateRange)
	if err != nil {
		t.Fatal(err)
	}

	a := assert.New(t)
	a.Equal(testEvent.Summary, testCalendar.Summary, "Event should be same name as calendar")
	// change name without supplying dateRange arg
	calendarName, _ = randutil.AlphaString(10)
	evt, err := serviceImpl.UpdateEvent(testService, testCalendar.Id, testEvent.Id, calendarName)
	if err != nil {
		t.Fatal(err)
	}
	a.NotEqual(testEvent.Summary, evt.Summary, "Updated event should not have the same name")
	// change dateRange
	start, end = RandDateRange()
	startStr = start.Format(time.RFC3339)
	endStr = end.Format(time.RFC3339)
	dateRange = new(DateRange)
	dateRange.Start = startStr
	dateRange.End = endStr

	evt, err = serviceImpl.UpdateEvent(testService, testCalendar.Id, testEvent.Id, calendarName, dateRange)
	if err != nil {
		t.Fatal(err)
	}
	a.Equal(evt.Summary, calendarName, "Updated event should have the same name")
	a.NotEqual(testEvent.Start.DateTime, evt.Start.DateTime, "Updated event should not have the same start datetime")
	a.NotEqual(testEvent.End.DateTime, evt.End.DateTime, "Updated event should not have the same end datetime")
}

func TestUpdateCalendar(t *testing.T) {
	t.Skip("Skipping TestUpdateCalendar")
	calendarName, _ := randutil.AlphaString(10)
	testCalendar, err := serviceImpl.CreateCalendar(testService, calendarName)
	if err != nil {
		t.Fatal(err)
	}
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

	testEvent1, err := serviceImpl.InsertEvent(testService, testCalendar, dateRange1)
	if err != nil {
		t.Fatal(err)
	}
	testEvent2, err := serviceImpl.InsertEvent(testService, testCalendar, dateRange2)
	if err != nil {
		t.Fatal(err)
	}
	testEvent3, err := serviceImpl.InsertEvent(testService, testCalendar, dateRange3)
	if err != nil {
		t.Fatal(err)
	}

	a := assert.New(t)
	a.Equal(testCalendar.Summary, testEvent1.Summary, "Created event name should equal parent calendar name")
	a.Equal(testCalendar.Summary, testEvent2.Summary, "Created event name should equal parent calendar name")
	a.Equal(testCalendar.Summary, testEvent3.Summary, "Created event name should equal parent calendar name")

	calendarName, _ = randutil.AlphaString(10)
	updatedCalendar, err := serviceImpl.UpdateCalendar(testService, testCalendar.Id, calendarName)
	if err != nil {
		t.Fatal(err)
	}

	a.Equal(updatedCalendar.Summary, calendarName, "Updated calendar should reflect changed name")

	eventList, err := serviceImpl.ListEvents(testService, updatedCalendar.Id)
	if err != nil {
		t.Fatal(err)
	}

	for _, evt := range eventList {
		a.Equal(evt.Summary, updatedCalendar.Summary, "Calendar events should reflect changed name")
	}
}

func SetupServiceTest() {
	ctx := context.TODO()
	client := serviceImpl.GetClient(ctx)
	testService = serviceImpl.BuildService(client)
}

func TeardownServiceTest() {
	log.Println("Cleaning up...")
	for _, calendarID := range calendars {
		serviceImpl.DeleteCalendar(testService, calendarID)
		log.Printf("Deleted %s", calendarID)

	}
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
