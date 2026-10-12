package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTgHttpClientUsesProxy(t *testing.T) {
	var requested string
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requested = r.URL.String()
	}))
	defer proxy.Close()

	client, err := newTgHttpClient(proxy.URL)
	require.NoError(t, err)

	rsp, err := client.Get("http://api.telegram.org/botTOKEN/getMe")
	require.NoError(t, err)
	rsp.Body.Close()

	assert.Equal(t, "http://api.telegram.org/botTOKEN/getMe", requested)
}

func TestTgHttpClientWithoutProxyIgnoresEnv(t *testing.T) {
	t.Setenv("HTTPS_PROXY", "http://127.0.0.1:1")

	client, err := newTgHttpClient("")
	require.NoError(t, err)

	assert.Nil(t, client.Transport.(*http.Transport).Proxy)
}

func TestTgHttpClientInvalidProxy(t *testing.T) {
	_, err := newTgHttpClient("://bad")
	assert.Error(t, err)
}
