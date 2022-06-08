package user_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gloompi/ultimate-service/business/core/user"
	"github.com/gloompi/ultimate-service/business/data/schema"
	"github.com/gloompi/ultimate-service/business/data/tests"
	"github.com/gloompi/ultimate-service/business/web/auth"
	"github.com/gloompi/ultimate-service/foundation/docker"
	"github.com/google/go-cmp/cmp"
)

var c *docker.Container

func TestMain(m *testing.M) {
	var err error
	c, err = tests.StartDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tests.StopDB(c)

	m.Run()
}

func Test_User(t *testing.T) {
	log, db, teardown := tests.NewUnit(t, c, "testuser")
	t.Cleanup(teardown)

	core := user.NewCore(log, db)

	t.Log("Given the need to work with User records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single User.", testID)
		{
			ctx := context.Background()
			now := time.Date(2018, time.October, 1, 0, 0, 0, 0, time.UTC)

			nu := user.NewUser{
				Name:            "Kubanychbek Esenzhanov",
				Email:           "gloompi@gmail.com",
				Roles:           []string{auth.RoleAdmin},
				Password:        "gophers",
				PasswordConfirm: "gophers",
			}

			usr, err := core.Create(ctx, nu, now)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create user : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create user.", tests.Success, testID)

			saved, err := core.QueryByID(ctx, usr.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve user by ID: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve user by ID.", tests.Success, testID)

			if diff := cmp.Diff(usr, saved); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same user. Diff:\n%s", tests.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same user.", tests.Success, testID)

			upd := user.UpdateUser{
				Name:  tests.StringPointer("Jacob Walker"),
				Email: tests.StringPointer("jacob@ardanlabs.com"),
			}

			if err := core.Update(ctx, usr.ID, upd, now); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to update user : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to update user.", tests.Success, testID)

			saved, err = core.QueryByEmail(ctx, *upd.Email)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve user by Email : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve user by Email.", tests.Success, testID)

			if saved.Name != *upd.Name {
				t.Errorf("\t%s\tTest %d:\tShould be able to see updates to Name.", tests.Failed, testID)
				t.Logf("\t\tTest %d:\tGot: %v", testID, saved.Name)
				t.Logf("\t\tTest %d:\tExp: %v", testID, *upd.Name)
			} else {
				t.Logf("\t%s\tTest %d:\tShould be able to see updates to Name.", tests.Success, testID)
			}

			if saved.Email != *upd.Email {
				t.Errorf("\t%s\tTest %d:\tShould be able to see updates to Email.", tests.Failed, testID)
				t.Logf("\t\tTest %d:\tGot: %v", testID, saved.Email)
				t.Logf("\t\tTest %d:\tExp: %v", testID, *upd.Email)
			} else {
				t.Logf("\t%s\tTest %d:\tShould be able to see updates to Email.", tests.Success, testID)
			}

			if err := core.Delete(ctx, usr.ID); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to delete user : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to delete user.", tests.Success, testID)

			_, err = core.QueryByID(ctx, usr.ID)
			if !errors.Is(err, user.ErrNotFound) {
				t.Fatalf("\t%s\tTest %d:\tShould NOT be able to retrieve user : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould NOT be able to retrieve user.", tests.Success, testID)
		}
	}
}

func Test_PagingUser(t *testing.T) {
	log, db, teardown := tests.NewUnit(t, c, "testpaging")
	t.Cleanup(teardown)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	schema.Seed(ctx, db)

	user := user.NewCore(log, db)

	t.Log("Given the need to page through User records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen paging through 2 users.", testID)
		{
			ctx := context.Background()

			users1, err := user.Query(ctx, 1, 1)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve users for page 1 : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve users for page 1.", tests.Success, testID)

			if len(users1) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single user : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single user.", tests.Success, testID)

			users2, err := user.Query(ctx, 2, 1)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve users for page 2 : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve users for page 2.", tests.Success, testID)

			if len(users2) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single user : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single user.", tests.Success, testID)

			if users1[0].ID == users2[0].ID {
				t.Logf("\t\tTest %d:\tUser1: %v", testID, users1[0].ID)
				t.Logf("\t\tTest %d:\tUser2: %v", testID, users2[0].ID)
				t.Fatalf("\t%s\tTest %d:\tShould have different users : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have different users.", tests.Success, testID)
		}
	}
}
