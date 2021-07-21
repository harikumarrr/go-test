package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/gin-gonic/gin"

	"go-test/pkg/routes"

	flag "github.com/spf13/pflag"
)

var opts = godog.Options{Output: colors.Colored(os.Stdout)}

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

func TestMain(m *testing.M) {
	flag.Parse()
	opts.Paths = flag.Args()

	status := godog.TestSuite{
		Name:                 "godogs",
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
		Options:              &opts,
	}.Run()

	os.Exit(status)
}

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() { Godogs = 0 })
}

type apiReceiver struct {
	engine           *gin.Engine
	responseRecorder *httptest.ResponseRecorder
}

type stepAdderFunction func(expr interface{}, stepFunc interface{})

func (api *apiReceiver) setEngin(*godog.Scenario) {
	api.engine = routes.SetupRouter()
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.BeforeScenario(func(*godog.Scenario) {
		Godogs = 0 // clean the state before every scenario
	})
	api := &apiReceiver{}
	ctx.BeforeScenario(api.setEngin)
	addStep := func(expr interface{}, stepFunc interface{}) {
		ctx.Step(expr, stepFunc)
	}
	addGroupOfSteps := func(groupAdd func(_ *apiReceiver, _ stepAdderFunction)) {
		groupAdd(api, addStep)
	}

	addGroupOfSteps(addRequestSteps)
}

func addRequestSteps(api *apiReceiver, addStep stepAdderFunction) {
	addStep(`^I send a "([^"]*)" request to "([^"]*)"$`, api.iSendARequestTo)
	addStep(`^the response code should be (\d+)$`, api.theResponseCodeShouldBe)
}

func (api *apiReceiver) iSendARequestTo(method, endpoint string) error {
	return api.sendRequest(method, endpoint, "", nil)
}

func (api *apiReceiver) sendRequest(method, endpoint string, body string, header textproto.MIMEHeader) error {
	//return fmt.Errorf(fmt.Sprintf("endpoint: %s, method: %s", endpoint, method))
	httptest.NewRecorder()
	reqURL, err := url.Parse(endpoint)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(method, reqURL.String(), strings.NewReader(body))
	if err != nil {
		return err
	}

	if body != "" {
		request.Header.Add("Content-Type", "application/json")
	}

	api.responseRecorder = httptest.NewRecorder()
	fmt.Printf("Sending %s to %s with body %s\n", method, reqURL.String(), body)
	api.engine.ServeHTTP(api.responseRecorder, request)
	return nil
}

func (api *apiReceiver) theResponseCodeShouldBe(expectedCode int) error {
	if api.responseRecorder.Code != expectedCode {
		return fmt.Errorf("Expected: %d, Actual: %d, response body: %s",
			expectedCode,
			api.responseRecorder.Code,
			api.responseRecorder.Body.String())
	}
	return nil
}
