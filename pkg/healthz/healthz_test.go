package healthz

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHealthzHandler(t *testing.T) {
	assert := require.New(t)

	backEndTs := httptest.NewServer(http.HandlerFunc(Handler()))

	res, err := http.Get(backEndTs.URL)

	body, _ := ioutil.ReadAll(res.Body)
	assert.NoError(err)
	assert.Equal("200 OK", res.Status)
	assert.NotEmpty(body)
}
