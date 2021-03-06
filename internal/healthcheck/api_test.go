package healthcheck

import (
	"net/http"
	"testing"

	"github.com/pohsi/pktrade/internal/test"
	"github.com/pohsi/pktrade/pkg/log"
)

func TestAPI(t *testing.T) {
	logger, _ := log.NewForTest()
	router := test.MockRouter(logger)
	RegisterHandlers(router, "Test")
	test.Endpoint(t, router, test.APITestCase{
		Name:         "ok",
		Method:       "GET",
		URL:          "/healthcheck",
		Body:         "",
		Header:       nil,
		WantStatus:   http.StatusOK,
		WantResponse: `"OK: Test"`,
	})
}
