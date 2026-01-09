package steps

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/vatvit/go-demo/internal/handler"
)

type apiFeature struct {
	server   *httptest.Server
	response *http.Response
	body     string
}

func (a *apiFeature) theServerIsRunning() error {
	h := handler.New()
	a.server = httptest.NewServer(h.Routes())
	return nil
}

func (a *apiFeature) iSendAGETRequestTo(path string) error {
	resp, err := http.Get(a.server.URL + path)
	if err != nil {
		return err
	}
	a.response = resp
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	a.body = string(body)
	return nil
}

func (a *apiFeature) iSendAPOSTRequestTo(path string) error {
	resp, err := http.Post(a.server.URL+path, "application/json", nil)
	if err != nil {
		return err
	}
	a.response = resp
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	a.body = string(body)
	return nil
}

func (a *apiFeature) theResponseStatusCodeShouldBe(code int) error {
	if a.response.StatusCode != code {
		return fmt.Errorf("expected status %d, got %d", code, a.response.StatusCode)
	}
	return nil
}

func (a *apiFeature) theResponseContentTypeShouldBe(contentType string) error {
	ct := a.response.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, contentType) {
		return fmt.Errorf("expected content type %s, got %s", contentType, ct)
	}
	return nil
}

func (a *apiFeature) theResponseBodyShouldContain(text string) error {
	if !strings.Contains(a.body, text) {
		return fmt.Errorf("expected body to contain %q, got %q", text, a.body)
	}
	return nil
}

func (a *apiFeature) cleanup(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
	if a.server != nil {
		a.server.Close()
	}
	if a.response != nil && a.response.Body != nil {
		a.response.Body.Close()
	}
	return ctx, nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	api := &apiFeature{}

	ctx.After(api.cleanup)

	ctx.Step(`^the server is running$`, api.theServerIsRunning)
	ctx.Step(`^I send a GET request to "([^"]*)"$`, api.iSendAGETRequestTo)
	ctx.Step(`^I send a POST request to "([^"]*)"$`, api.iSendAPOSTRequestTo)
	ctx.Step(`^the response status code should be (\d+)$`, api.theResponseStatusCodeShouldBe)
	ctx.Step(`^the response content type should be "([^"]*)"$`, api.theResponseContentTypeShouldBe)
	ctx.Step(`^the response body should contain "([^"]*)"$`, api.theResponseBodyShouldContain)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
