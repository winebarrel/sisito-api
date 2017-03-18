package sisito

import (
	. "."
	"github.com/bouk/monkey"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestPing(t *testing.T) {
	assert := assert.New(t)

	server := NewServer(&Config{User: []UserConfig{}}, nil)

	ts := httptest.NewServer(server.Engine)
	res, _ := http.Get(ts.URL + "/ping")
	body, status := readResponse(res)

	assert.Equal(200, status)
	assert.Equal(body, `{"message":"pong"}`+"\n")
}

func TestRecentWithRecipient(t *testing.T) {
	assert := assert.New(t)

	driver := &Driver{}
	server := NewServer(&Config{User: []UserConfig{}}, driver)

	var guard *monkey.PatchGuard
	guard = monkey.PatchInstanceMethod(
		reflect.TypeOf(driver), "RecentlyBounced",
		func(_ *Driver, name string, value string, senderdomain string) (bounced []BounceMail, err error) {
			defer guard.Unpatch()
			guard.Restore()

			assert.Equal("recipient", name)
			assert.Equal("foo@example.com", value)
			assert.Equal("example.net", senderdomain)

			bounced = []BounceMail{BounceMail{Id: 1}}

			return
		})

	ts := httptest.NewServer(server.Engine)
	res, _ := http.Get(ts.URL + "/recent?recipient=foo@example.com&senderdomain=example.net")
	body, status := readResponse(res)

	assert.Equal(200, status)
	assert.Equal(body, `{"addresser":"",`+
		`"alias":"",`+
		`"created_at":"0001-01-01T00:00:00Z",`+
		`"deliverystatus":"",`+
		`"destination":"",`+
		`"diagnosticcode":"",`+
		`"digest":"",`+
		`"lhost":"",`+
		`"messageid":"",`+
		`"reason":"",`+
		`"recipient":"",`+
		`"rhost":"",`+
		`"senderdomain":"",`+
		`"smtpagent":"",`+
		`"smtpcommand":"",`+
		`"softbounce":false,"subject":"",`+
		`"timestamp":"0001-01-01T00:00:00Z",`+
		`"timezoneoffset":"",`+
		`"updated_at":"0001-01-01T00:00:00Z",`+
		`"whitelisted":false}`+"\n")
}