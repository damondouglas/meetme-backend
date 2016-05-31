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

type service interface {
	GetClient(ctx context.Context) *http.Client
	BuildService(client *http.Client) *calendar.Service
	CreateCalendar(srv *calendar.Service, name string) (*calendar.Calendar, error)
	InsertEvent(srv *calendar.Service, cal *calendar.Calendar, dateRange *DateRange) (*calendar.Event, error)
	ListEvents(srv *calendar.Service, calendarID string) ([]*calendar.Event, error)
	ListCalendars(srv *calendar.Service) ([]*calendar.CalendarListEntry, error)
	DeleteCalendar(srv *calendar.Service, calendarID string) error
	GetCalendar(srv *calendar.Service, calendarID string) (*calendar.Calendar, error)
	UpdateEvent(srv *calendar.Service, calendarID string, eventID string, name string, dateRange ...*DateRange) (*calendar.Event, error)
	UpdateCalendar(srv *calendar.Service, calendarID string, summary string) (*calendar.Calendar, error)
}

type serviceimpl struct{}

// DateRange is a range from [Start] to [End]
type DateRange struct {
	Start string `json:"dateTime,omitempty"`
	End   string `json:"dateTime,omitempty"`
}

// GetClient builds OAuth2 Client
func (s serviceimpl) GetClient(ctx context.Context) *http.Client {

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
func (s serviceimpl) BuildService(client *http.Client) *calendar.Service {
	srv, err := calendar.New(client)
	if err != nil {
		log.Fatal(err)
	}
	return srv
}

// CreateCalendar in service account
func (s serviceimpl) CreateCalendar(srv *calendar.Service, name string) (*calendar.Calendar, error) {
	cal := &calendar.Calendar{Summary: name}
	calRes, err := srv.Calendars.Insert(cal).Do()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return calRes, nil
}

// InsertEvent in calendar
// [start] and [end] must be rfc3339 string format
func (s serviceimpl) InsertEvent(srv *calendar.Service, cal *calendar.Calendar, dateRange *DateRange) (*calendar.Event, error) {
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
		return nil, err
	}
	return evt, nil
}

// ListEvents lists calendar events
func (s serviceimpl) ListEvents(srv *calendar.Service, calendarID string) ([]*calendar.Event, error) {
	events, err := srv.Events.List(calendarID).Do()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return events.Items, nil
}

// ListCalendars lists calendars in service account
func (s serviceimpl) ListCalendars(srv *calendar.Service) ([]*calendar.CalendarListEntry, error) {
	calendarList, err := srv.CalendarList.List().Do()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return calendarList.Items, nil
}

// DeleteCalendar deletes calendar in service account
func (s serviceimpl) DeleteCalendar(srv *calendar.Service, calendarID string) error {
	err := srv.Calendars.Delete(calendarID).Do()
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

// GetCalendar loads calendar from service account given calendarID
func (s serviceimpl) GetCalendar(srv *calendar.Service, calendarID string) (*calendar.Calendar, error) {
	cal, err := srv.Calendars.Get(calendarID).Do()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return cal, nil
}

// UpdateEvent updates event data
func (s serviceimpl) UpdateEvent(srv *calendar.Service, calendarID string, eventID string, name string, dateRange ...*DateRange) (*calendar.Event, error) {
	evt, err := srv.Events.Get(calendarID, eventID).Do()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	evt.Summary = name
	if len(dateRange) > 0 {
		evt.Start.DateTime = dateRange[0].Start
		evt.End.DateTime = dateRange[0].End
	}

	evt, err = srv.Events.Update(calendarID, eventID, evt).Do()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return evt, nil
}

// UpdateCalendar updates calendar from service account given calendarID and new Summary
func (s serviceimpl) UpdateCalendar(srv *calendar.Service, calendarID string, summary string) (*calendar.Calendar, error) {
	cal, err := s.GetCalendar(srv, calendarID)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if cal != nil {
		cal.Summary = summary
		updatedCal, err := srv.Calendars.Update(calendarID, cal).Do()
		cal = updatedCal
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		eventList, err := s.ListEvents(srv, cal.Id)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		for _, evt := range eventList {
			_, err = s.UpdateEvent(srv, cal.Id, evt.Id, summary)
			if err != nil {
				log.Fatal(err)
				return nil, err
			}
		}
	}

	return cal, nil
}
