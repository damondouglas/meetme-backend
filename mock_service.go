package meetme

import (
	"github.com/jmcvetta/randutil"
	"golang.org/x/net/context"
	"google.golang.org/api/calendar/v3"
	"log"
	"net/http"
)

type mockservice struct {
	Calendars    map[string]calendar.Calendar
	Events       map[string]calendar.Event
	EventListMap map[string][]string
}

func newMockService() *mockservice {
	s := new(mockservice)
	s.Calendars = make(map[string]calendar.Calendar)
	s.Events = make(map[string]calendar.Event)
	s.EventListMap = make(map[string][]string)
	return s
}

func (s mockservice) GetClient(ctx context.Context) *http.Client {
	return nil
}

func (s mockservice) BuildService(client *http.Client) *calendar.Service {
	service := new(calendar.Service)
	return service
}

func getID() string {
	id, _ := randutil.AlphaString(16)
	return id + "_MOCK_RESOURCE" // So we can confirm we created resource through mock service
}

func (s mockservice) SaveCalendar(c *calendar.Calendar) {
	s.Calendars[c.Id] = *c
}

func (s mockservice) SaveEvent(e *calendar.Event, calendarID string) {
	s.Events[e.Id] = *e
	calendarEventListLength := len(s.EventListMap[calendarID])
	if calendarEventListLength == 0 {
		s.EventListMap[calendarID] = []string{}
	}
	s.EventListMap[calendarID] = append(s.EventListMap[calendarID], e.Id)
}

func (s mockservice) Debug() {
	log.Printf("calendars(%v): %v", len(s.Calendars), s.Calendars)
	log.Printf("events(%v): %v", len(s.Events), s.Events)
	log.Printf("eventlist(%v): %v", len(s.EventListMap), s.EventListMap)
}

func (s mockservice) CreateCalendar(srv *calendar.Service, name string) (*calendar.Calendar, error) {
	cal := new(calendar.Calendar)
	cal.Summary = name
	cal.Id = getID()
	s.SaveCalendar(cal)
	return cal, nil
}

func (s mockservice) InsertEvent(srv *calendar.Service, cal *calendar.Calendar, dateRange *DateRange) (*calendar.Event, error) {
	evt := new(calendar.Event)
	evt.Id = getID()
	evt.Summary = cal.Summary
	start := new(calendar.EventDateTime)
	end := new(calendar.EventDateTime)
	start.DateTime = dateRange.Start
	end.DateTime = dateRange.End
	evt.Start = start
	evt.End = end
	s.SaveEvent(evt, cal.Id)
	return evt, nil
}

func (s mockservice) ListEvents(srv *calendar.Service, calendarID string) ([]*calendar.Event, error) {
	var eventList []*calendar.Event
	return eventList, nil
}

func (s mockservice) ListCalendars(srv *calendar.Service) ([]*calendar.CalendarListEntry, error) {
	var calendarList []*calendar.CalendarListEntry
	return calendarList, nil
}

func (s mockservice) DeleteCalendar(srv *calendar.Service, calendarID string) error {
	return nil
}

func (s mockservice) GetCalendar(srv *calendar.Service, calendarID string) (*calendar.Calendar, error) {
	cal := new(calendar.Calendar)
	return cal, nil
}

func (s mockservice) UpdateEvent(srv *calendar.Service, calendarID string, eventID string, name string, dateRange ...*DateRange) (*calendar.Event, error) {
	evt := new(calendar.Event)
	return evt, nil
}

func (s mockservice) UpdateCalendar(srv *calendar.Service, calendarID string, summary string) (*calendar.Calendar, error) {
	cal := new(calendar.Calendar)
	return cal, nil
}
