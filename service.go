package meetme

import (
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"log"
	"net/http"
	// "golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	// "google.golang.org/appengine/urlfetch"
	// "google.golang.org/appengine"
)

// DateRange is a range from [Start] to [End]
type DateRange struct {
	Start string `json:"dateTime,omitempty"`
	End   string `json:"dateTime,omitempty"`
}

// GetClient builds OAuth2 Client
func GetClient(ctx context.Context) *http.Client {

	data, err := ioutil.ReadFile("cred.json")
	if err != nil {
		log.Fatal(err)
	}

	conf, err := google.JWTConfigFromJSON(data, calendar.CalendarScope)

	if err != nil {
		log.Fatal(err)
	}

	client := conf.Client(ctx)

	return client
}

// BuildService builds calendar service
func BuildService(client *http.Client) *calendar.Service {
	srv, err := calendar.New(client)
	if err != nil {
		log.Fatal(err)
	}
	return srv
}

// CreateCalendar in service account
func CreateCalendar(srv *calendar.Service, name string) *calendar.Calendar {
	cal := &calendar.Calendar{Summary: name}
	calRes, err := srv.Calendars.Insert(cal).Do()
	if err != nil {
		log.Fatal(err)
	}
	return calRes
}

// InsertEvent in calendar
// [start] and [end] is in rfc3339 string format
func InsertEvent(srv *calendar.Service, cal *calendar.Calendar, dateRange *DateRange) *calendar.Event {
	evt := new(calendar.Event)
	evt.Summary = cal.Summary
	startDateTime := new(calendar.EventDateTime)
	startDateTime.DateTime = dateRange.Start
	endDateTime := new(calendar.EventDateTime)
	endDateTime.DateTime = dateRange.End
	evt.Start = startDateTime
	evt.End = endDateTime
	evt, err := srv.Events.Insert(cal.Id, evt).Do()
	if err != nil {
		log.Fatal(err)
	}
	return evt
}

// ListEvents lists calendar events
func ListEvents(srv *calendar.Service, calendarID string) []*calendar.Event {
	events, err := srv.Events.List(calendarID).Do()
	if err != nil {
		log.Fatal(err)
	}
	return events.Items
}

// ListCalendars lists calendars in service account
func ListCalendars(srv *calendar.Service) []*calendar.CalendarListEntry {
	calendarList, err := srv.CalendarList.List().Do()
	if err != nil {
		log.Fatal(err)
	}
	return calendarList.Items
}

// DeleteCalendar deletes calendar in service account
func DeleteCalendar(srv *calendar.Service, calendarID string) {
	err := srv.Calendars.Delete(calendarID).Do()
	if err != nil {
		log.Fatal(err)
	}
}

// GetCalendar loads calendar from service account given calendarID
func GetCalendar(srv *calendar.Service, calendarID string) *calendar.Calendar {
	cal, err := srv.Calendars.Get(calendarID).Do()
	if err != nil {
		log.Fatal(err)
	}
	return cal
}

// UpdateEvent updates event data
func UpdateEvent(srv *calendar.Service, calendarID string, eventID string, name string, dateRange ...*DateRange) *calendar.Event {
	evt, err := srv.Events.Get(calendarID, eventID).Do()
	if err != nil {
		log.Fatal(err)
	}
	evt.Summary = name
	if len(dateRange) > 0 {
		evt.Start.DateTime = dateRange[0].Start
		evt.End.DateTime = dateRange[0].End
	}

	evt, err = srv.Events.Update(calendarID, eventID, evt).Do()
	if err != nil {
		log.Fatal(err)
	}

	return evt
}

// UpdateCalendar updates calendar from service account given calendarID and new Summary
func UpdateCalendar(srv *calendar.Service, calendarID string, summary string) *calendar.Calendar {
	cal := GetCalendar(srv, calendarID)

	if cal != nil {
		cal.Summary = summary
		updatedCal, err := srv.Calendars.Update(calendarID, cal).Do()
		cal = updatedCal
		if err != nil {
			log.Fatal(err)
		}

		eventList := ListEvents(srv, cal.Id)
		for _, evt := range eventList {
			UpdateEvent(srv, cal.Id, evt.Id, summary)
		}
	}

	return cal
}
