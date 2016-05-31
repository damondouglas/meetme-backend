package meetme

import (
	"errors"
	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"
	"golang.org/x/net/context"
	"google.golang.org/api/calendar/v3"
	"log"
	// "time"
)

const aetestAppID = "testapp"

var (
	scopes          = []string{endpoints.EmailScope}
	clientIds       = []string{endpoints.APIExplorerClientID}
	audiences       = []string{}
	endpointService service
)

func getEndpointServiceFromContext(c context.Context) service {
	return new(mockservice)
}

// CalendarService offers operations to maintain calendars in the service account
type CalendarService struct{}

// EventService offers operations to maintain events in the service account
type EventService struct{}

func _buildServiceFromContext(c context.Context) *calendar.Service {
	client := endpointService.GetClient(c)
	return endpointService.BuildService(client)
}

func _putMeetie(c context.Context, m Meetie) error {
	log.Println(m)
	return nil
}

// CreateCalendarEndpoint sends request to create a new calendar in service account
func (cs *CalendarService) CreateCalendarEndpoint(c context.Context, r *CreateCalendarReq) (*calendar.Calendar, error) {
	srv := _buildServiceFromContext(c)
	var cal *calendar.Calendar
	u, err := endpoints.CurrentBearerTokenUser(c, scopes, clientIds)
	if err != nil {
		log.Fatal(err)
	}

	userReachedLimit := getUserReachedLimit(c, u.Email)
	if !userReachedLimit {
		cal, _ = endpointService.CreateCalendar(srv, r.Name)
	}
	// var user *user.User
	// var err error
	// if aetestAppID == appengine.AppID(c) {
	// 	cal = GetMockCalendar()
	// 	user = GetMockUser()
	// } else {
	// 	user, err = endpoints.CurrentUser(c, scopes, audiences, clientIds)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	//
	// 	userReachedLimit := getUserReachedLimit(c, user)
	// 	if userReachedLimit {
	// 		return nil, errors.New("User reached limit.")
	// 	}
	//
	// 	srv := _buildServiceFromContext(c)
	// 	cal, err = CreateCalendar(srv, r.Name)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 		return nil, err
	// 	}
	// }
	//
	//
	//
	return cal, nil
}

// InsertEventEndpoint sends request to insert event
func (cs *EventService) InsertEventEndpoint(c context.Context, r *InsertEventReq) (*calendar.Event, error) {
	srv := _buildServiceFromContext(c)
	cal := new(calendar.Calendar)
	cal.Id = r.CalendarID
	cal.Summary = r.Name
	var dateRange *DateRange
	if r.Start != "" && r.End != "" {
		dateRange = new(DateRange)
		dateRange.Start = r.Start
		dateRange.End = r.End
	} else {
		return nil, errors.New("Start and End date must be specified in request.")
	}
	return endpointService.InsertEvent(srv, cal, dateRange)
}

// DeleteCalendarEndpoint sends request to delete calendar in service account
func (cs *CalendarService) DeleteCalendarEndpoint(c context.Context, r *DeleteCalendarReq) error {
	srv := _buildServiceFromContext(c)
	endpointService.DeleteCalendar(srv, r.CalendarID)
	return nil
}

func init() {
	endpointService = new(serviceimpl)
	api, err := endpoints.RegisterService(&CalendarService{},
		"meetme", "v1", "Meetme Calendar API", true)
	if err != nil {
		log.Fatalf("Register service: %v", err)
	}

	create := api.MethodByName("CreateCalendarEndpoint").Info()
	create.Scopes = scopes
	create.ClientIds = clientIds
	create.HTTPMethod = "POST"
	create.Path = "calendar/create"
	create.Name = "create"
	create.Desc = "Create calendar."

	delete := api.MethodByName("DeleteCalendarEndpoint").Info()
	delete.HTTPMethod = "POST"
	delete.Path = "calendar/delete"
	delete.Name = "delete"
	delete.Desc = "Delete calendar."

	endpoints.HandleHTTP()
}
