package realm_test

import (
	"testing"

	"github.com/10gen/realm-cli/internal/cloud/realm"
	u "github.com/10gen/realm-cli/internal/utils/test"
	"github.com/10gen/realm-cli/internal/utils/test/assert"
)

func TestRealmDrafts(t *testing.T) {
	u.SkipUnlessRealmServerRunning(t)

	t.Run("Should fail without an auth client", func(t *testing.T) {
		for _, tc := range []struct {
			description string
			call        func(client realm.Client) error
		}{
			{
				"With a call to get draft",
				func(client realm.Client) error {
					_, err := client.Draft("groupID", "appID")
					return err
				},
			},
		} {
			t.Run(tc.description, func(t *testing.T) {
				client := realm.NewClient(u.RealmServerURL())

				assert.Equal(t, realm.ErrInvalidSession{}, tc.call(client))
			})
		}
	})

	t.Run("With an active session", func(t *testing.T) {
		client := newAuthClient(t)
		groupID := u.CloudGroupID()

		app, appErr := client.CreateApp(groupID, "users-test", realm.AppMeta{})
		assert.Nil(t, appErr)

		t.Run("Getting drafts should fail if there are none", func(t *testing.T) {
			_, err := client.Draft(groupID, app.ID)
			assert.Equal(t, realm.ErrDraftNotFound, err)
		})

		t.Run("Should create a draft", func(t *testing.T) {
			draft, draftErr := client.CreateDraft(groupID, app.ID)
			assert.Nil(t, draftErr)
			assert.True(t, draft.ID != "", "created draft id should not be empty")

			t.Run("And be able to retrieve the draft", func(t *testing.T) {
				found, err := client.Draft(groupID, app.ID)
				assert.Nil(t, err)
				assert.Equal(t, draft, found)
			})

			t.Run("And be able to get diffs of the draft", func(t *testing.T) {
				diffs, err := client.DiffDraft(groupID, app.ID, draft.ID)
				assert.Nil(t, err)
				assert.Equal(t, 0, diffs.Len())
			})

			t.Run("And be able to discard the draft", func(t *testing.T) {
				assert.Nil(t, client.DiscardDraft(groupID, app.ID, draft.ID))

				_, err := client.Draft(groupID, app.ID)
				assert.Equal(t, realm.ErrDraftNotFound, err)
			})
		})
	})
}
