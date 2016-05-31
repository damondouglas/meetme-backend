package meetme

// CreateCalendarReq is the request type for CalendarService.CreateCalendar
type CreateCalendarReq struct {
	Name string `json:"name"`
}

// InsertEventReq is the request type for CalendarService.InsertEvent
type InsertEventReq struct {
	CalendarID string `json:"calendarID"`
	Name       string `json:"name"`
	Start      string `json:"start"`
	End        string `json:"end"`
}

// ListEventsReq is the request type for CalendarService.ListEvents
type ListEventsReq struct {
	CalendarID string `json:"calendarID"`
}

// ListCalendarsReq is the request type for CalendarService.ListCalendars
type ListCalendarsReq struct{}

// DeleteCalendarReq is the request type for CalendarService.DeleteCalendar
type DeleteCalendarReq struct {
	CalendarID string `json:"calendarID"`
}

// GetCalendarReq is the request type for CalendarService.GetCalendar
type GetCalendarReq struct {
	CalendarID string `json:"calendarID"`
}

// UpdateEventReq is the request type for CalendarService.GetCalendar
type UpdateEventReq struct {
	CalendarID string `json:"calendarID"`
	EventID    string `json:"eventID"`
	Name       string `json:"name"`
	Start      string `json:"start"`
	End        string `json:"end"`
}

// UpdateCalendarReq is the request type for CalendarService.UpdateCalendar
type UpdateCalendarReq struct {
	CalendarID string `json:"calendarID"`
	Name       string `json:"name"`
}
