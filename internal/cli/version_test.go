package cli

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/10gen/realm-cli/internal/utils/api"
	"github.com/10gen/realm-cli/internal/utils/test/assert"
)

func TestVersionCheck(t *testing.T) {
	origOSArch := osArch
	osArch = "macos-amd64"
	defer func() { osArch = origOSArch }()

	for _, tc := range []struct {
		description string
		nextVersion string
		expectedMsg string
	}{
		{
			description: "should return nothing if the current version is not greater than the current version",
			nextVersion: "0.0.0",
		},
		{
			description: "should return a helpful message if there is a newer version",
			nextVersion: "0.1.0",
			expectedMsg: "New version (v0.1.0) of CLI available: http://whatever.com/test",
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			client := testClient{http.StatusOK, tc.nextVersion, osArch, "http://whatever.com/test"}

			versionMsg, err := checkVersion(client)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedMsg, versionMsg)
		})
	}

	t.Run("should return an error if the client request fails", func(t *testing.T) {
		var client failClient

		_, err := checkVersion(client)
		assert.Equal(t, errors.New("something bad happened"), err)
	})

	t.Run("should return an error if the response status code is not ok", func(t *testing.T) {
		client := testClient{statusCode: http.StatusInternalServerError}

		_, err := checkVersion(client)
		assert.Equal(t, api.ErrUnexpectedStatusCode{"get cli version manifest", http.StatusInternalServerError}, err)
	})

	t.Run("should return an error if the next cli version is not semantic", func(t *testing.T) {
		client := testClient{statusCode: http.StatusOK, version: "0.0"}

		_, err := checkVersion(client)
		assert.Equal(t, errors.New("failed to parse version v0.0"), err)
	})

	t.Run("should return an error if the current cli version is not semantic", func(t *testing.T) {
		origVersion := Version
		Version = "0.0"
		defer func() { Version = origVersion }()

		client := testClient{statusCode: http.StatusOK, version: "0.0.0"}

		_, err := checkVersion(client)
		assert.Equal(t, errors.New("failed to parse version v0.0"), err)
	})

	t.Run("should return an error if the cli os architecture is unrecognized", func(t *testing.T) {
		client := testClient{statusCode: http.StatusOK, version: "0.1.0", osArch: "some-other-arch"}

		_, err := checkVersion(client)
		assert.Equal(t, fmt.Errorf("unrecognized CLI OS build: %s", osArch), err)
	})
}

type testClient struct {
	statusCode int
	version    string
	osArch     string
	url        string
}

func (client testClient) Get(url string) (*http.Response, error) {
	return &http.Response{
		StatusCode: client.statusCode,
		Body: ioutil.NopCloser(strings.NewReader(fmt.Sprintf(`{
"version" :%q,
"info": {
	%q: {
			"url": %q
	}
}
}`, client.version, client.osArch, client.url))),
	}, nil
}

type failClient struct {
}

func (client failClient) Get(url string) (*http.Response, error) {
	return nil, errors.New("something bad happened")
}
