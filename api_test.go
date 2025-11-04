package mockhttp

import (
	"net/http"
	"testing"

	"github.com/le-yams/gotestingmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_api_mock_should(t *testing.T) {
	t.Parallel()

	t.Run("return_underlying_server_url", func(t *testing.T) {
		t.Parallel()
		// Arrange & Act
		mockedAPI := API(testingmock.New(t))
		t.Cleanup(mockedAPI.Close)

		// Assert
		assert.Equal(t, mockedAPI.testServer.URL, mockedAPI.GetURL().String())
		assert.Equal(t, mockedAPI.GetURL().Host, mockedAPI.GetHost())
	})

	t.Run("register a cleanup to close underlying server", func(t *testing.T) {
		t.Parallel()
		// Arrange
		mockedT := testingmock.New(t)
		mockedAPI := API(mockedT)
		mockedAPI.Stub(http.MethodGet, "/endpoint").WithStatusCode(http.StatusOK)
		endpointURL := mockedAPI.GetURL().JoinPath("endpoint").String()

		response, err := http.Get(endpointURL)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)

		// Act
		cleanups := mockedT.GetCleanups()
		require.Len(t, cleanups, 1)
		cleanups[0]() // Execute cleanup like a testing.T would do

		// Assert
		_, err = http.Get(endpointURL) // Should fail because server is closed
		assert.Error(t, err)
	})
}
