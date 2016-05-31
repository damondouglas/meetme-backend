package meetme

// testuser.tr0r@gmail.com
import (
	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/appengine/aetest"
	"log"
	"net/http"
	"testing"
	// "google.golang.org/appengine"
	"github.com/jmcvetta/randutil"
	"google.golang.org/appengine/user"
	"time"
)

var (
	endpointCtx     context.Context
	calendarService = new(CalendarService)
	eventService    = new(EventService)
	mockUser        *user.User
)

func MockCreateCalendar(r CreateCalendarReq) *calendar.Calendar {
	cal := new(calendar.Calendar)
	cal.Summary = r.Name
	return cal
}

func MockLogin(r *http.Request) {
	//example@example.com
	clientIds = append(clientIds, "123456789.apps.googleusercontent.com")
	const jwtValidTokenString = ("eyJhbGciOiAiUlMyNTYiLCAidHlwIjogIkpXVCJ9." +
		"eyJhdWQiOiAibXktY2xpZW50LWlkIiwgImlzcyI6ICJhY2NvdW50cy5nb29nbGUuY29tIiwg" +
		"ImV4cCI6IDEzNzAzNTIyNTIsICJhenAiOiAiaGVsbG8tYW5kcm9pZCIsICJpYXQiOiAxMzcw" +
		"MzQ4NjUyLCAiZW1haWwiOiAiZHVkZUBnbWFpbC5jb20ifQ." +
		"sv7l0v_u6DmVe7s-hg8Q5LOYXNCdUBR7efnvQ4ns6IfBFZ71yPvWfwOYqZuYGQ0a9V5CfR0r" +
		"TfNlXVEpW5NE9rZy8hFiZkHBE30yPDti6PUUtT1bZST1VPFnIvSHobcUj-QPBTRC1Df86Vv0" +
		"Jmx4yowL1z3Yhe0Zh1WcvUUPG9sKJt8_-qKAv9QeeCMveBYpRSh6JvoU_qUKxPTjOOLvQiqV" +
		"4NiNjJ3sDN0P4BHJc3VcqB-SFd7kMRgQy1Fq-NHKN5-T2x4gxPwUy9GvHOftxY47B1NtJ9Q5" +
		"KtSui9uXdyBNJnt0xcIT5CcQYkVLoeCldxpSfwfA2kyfuJKBeiQnSA")

	r.Header.Set("authorization", "oauth"+jwtValidTokenString)
}

func SetupEndpointTest() {
	var err error
	// point endpoint.go endpointService var to mockservice instead of google calendar api
	endpointService = newMockService()
	inst, err := aetest.NewInstance(nil)
	if err != nil {
		log.Fatal(err)
	}

	r, err := inst.NewRequest("GET", "/", nil)
	if err != nil {
		log.Fatal(err)
	}

	endpointCtx = endpoints.NewContext(r)
	MockLogin(r)

	log.Println("Setup Endpoint Test")
}

func TeardownEndpointTest() {
	log.Println("Teardown Endpoint Test")
}

func TestCreateCalendarEndpoint(t *testing.T) {
	t.Skip("Skipping TestCreateCalendarEndpoint")
	req := new(CreateCalendarReq)
	name1, _ := randutil.AlphaString(10)
	name2, _ := randutil.AlphaString(10)
	req.Name = name1 + " " + name2
	cal, err := calendarService.CreateCalendarEndpoint(endpointCtx, req)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, req.Name, cal.Summary, "Created calendar should match name of request")
}

func TestInsertEventEndpoint(t *testing.T) {
	t.Skip("Skipping TestInsertEventEndpoint")
	req := new(CreateCalendarReq)
	req.Name, _ = randutil.AlphaString(10)
	cal, _ := calendarService.CreateCalendarEndpoint(endpointCtx, req)
	insertEvtRequest := new(InsertEventReq)

	insertEvtRequest.CalendarID = cal.Id
	insertEvtRequest.Name = cal.Summary

	start, end := RandDateRange()
	startStr := start.Format(time.RFC3339)
	endStr := end.Format(time.RFC3339)
	insertEvtRequest.Start = startStr
	insertEvtRequest.End = endStr

	evt, err := eventService.InsertEventEndpoint(endpointCtx, insertEvtRequest)
	if err != nil {
		log.Fatal(err)
	}

	a := assert.New(t)
	a.Equal(evt.Summary, cal.Summary, "Event summary should equal calendar summary")
	a.Equal(evt.Start.DateTime, startStr, "Event start should be set to request start")
	a.Equal(evt.End.DateTime, endStr, "Event end should be set to request end")
}
