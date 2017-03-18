package sisito

import (
	. "."
	"github.com/bouk/monkey"
	"github.com/stretchr/testify/assert"
	"gopkg.in/gorp.v1"
	"reflect"
	"testing"
)

func TestRecentlyBounced(t *testing.T) {
	assert := assert.New(t)
	driver := &Driver{DbMap: &gorp.DbMap{}}

	var guard *monkey.PatchGuard
	guard = monkey.PatchInstanceMethod(
		reflect.TypeOf(driver.DbMap), "Select",
		func(_ *gorp.DbMap, i interface{}, query string, args ...interface{}) ([]interface{}, error) {
			defer guard.Unpatch()
			guard.Restore()

			assert.Equal(`
    SELECT bm.*, IF(wm.id IS NULL, 0, 1) AS whitelisted
      FROM bounce_mails bm LEFT JOIN whitelist_mails wm
        ON bm.recipient = wm.recipient AND bm.senderdomain = wm.senderdomain
     WHERE bm.recipient = ?
       AND bm.senderdomain = ?
  ORDER BY bm.id DESC
     LIMIT 1`, query)

			assert.Equal([]interface{}{"foo@example.com", "example.net"}, args)

			rows := i.(*[]BounceMail)
			*rows = append(*rows, BounceMail{Id: 1})

			return nil, nil
		})

	rows, _ := driver.RecentlyBounced("recipient", "foo@example.com", "example.net")

	assert.Equal([]BounceMail{BounceMail{Id: 1}}, rows)
}

func TestRecentlyBouncedWithoutSenderdomain(t *testing.T) {
	assert := assert.New(t)
	driver := &Driver{DbMap: &gorp.DbMap{}}

	var guard *monkey.PatchGuard
	guard = monkey.PatchInstanceMethod(
		reflect.TypeOf(driver.DbMap), "Select",
		func(_ *gorp.DbMap, i interface{}, query string, args ...interface{}) ([]interface{}, error) {
			defer guard.Unpatch()
			guard.Restore()

			assert.Equal(`
    SELECT bm.*, IF(wm.id IS NULL, 0, 1) AS whitelisted
      FROM bounce_mails bm LEFT JOIN whitelist_mails wm
        ON bm.recipient = wm.recipient AND bm.senderdomain = wm.senderdomain
     WHERE bm.recipient = ?
  ORDER BY bm.id DESC
     LIMIT 1`, query)

			assert.Equal([]interface{}{"foo@example.com"}, args)

			rows := i.(*[]BounceMail)
			*rows = append(*rows, BounceMail{Id: 1})

			return nil, nil
		})

	rows, _ := driver.RecentlyBounced("recipient", "foo@example.com", "")

	assert.Equal([]BounceMail{BounceMail{Id: 1}}, rows)
}

func TestIsBounced(t *testing.T) {
	assert := assert.New(t)
	driver := &Driver{DbMap: &gorp.DbMap{}}

	var guard *monkey.PatchGuard
	guard = monkey.PatchInstanceMethod(
		reflect.TypeOf(driver.DbMap), "SelectInt",
		func(_ *gorp.DbMap, query string, args ...interface{}) (int64, error) {
			defer guard.Unpatch()
			guard.Restore()

			assert.Equal(`
    SELECT 1
      FROM bounce_mails bm LEFT JOIN whitelist_mails wm
        ON bm.recipient = wm.recipient AND bm.senderdomain = wm.senderdomain
     WHERE bm.recipient = ?
       AND bm.senderdomain = ?
       AND wm.id IS NULL
     LIMIT 1`, query)

			assert.Equal([]interface{}{"foo@example.com", "example.net"}, args)

			return 1, nil
		})

	count, _ := driver.IsBounced("recipient", "foo@example.com", "example.net")

	assert.Equal(count, true)
}

func TestIsBouncedWithoutSenderdomain(t *testing.T) {
	assert := assert.New(t)
	driver := &Driver{DbMap: &gorp.DbMap{}}

	var guard *monkey.PatchGuard
	guard = monkey.PatchInstanceMethod(
		reflect.TypeOf(driver.DbMap), "SelectInt",
		func(_ *gorp.DbMap, query string, args ...interface{}) (int64, error) {
			defer guard.Unpatch()
			guard.Restore()

			assert.Equal(`
    SELECT 1
      FROM bounce_mails bm LEFT JOIN whitelist_mails wm
        ON bm.recipient = wm.recipient AND bm.senderdomain = wm.senderdomain
     WHERE bm.recipient = ?
       AND wm.id IS NULL
     LIMIT 1`, query)

			assert.Equal([]interface{}{"foo@example.com"}, args)

			return 1, nil
		})

	count, _ := driver.IsBounced("recipient", "foo@example.com", "")

	assert.Equal(count, true)
}
