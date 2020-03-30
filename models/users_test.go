package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"testing"
	"time"
)

func getPostgresLoginString() string {
	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "future"
		dbname   = "wiki_test"
	)
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user,
		password, dbname)
}

func testingUserService() (*UserService, error) {
	us, err := NewUserService(getPostgresLoginString())
	if err != nil {
		return nil, err
	}
	us.db.LogMode(false)
	us.DestructiveReset()
	return us, nil
}

func TestCreateUser(t *testing.T) {
	us, err := testingUserService()
	if err != nil {
		t.Fatal()
	}
	user := User{
		Model: gorm.Model{},
		Name:  "Michael Scott",
		Email: "michael@test.com",
	}
	_, err = us.Create(&user)
	if err != nil {
		t.Fatal(err)
	}
	if user.ID == 0 {
		t.Errorf("Expected ID > 0. Received %d", user.ID)
	}
	if time.Since(user.CreatedAt) > 5*time.Second {
		t.Errorf("Expected CreatedAt to be recent. Received %s", user.CreatedAt)
	}
}

func TestUpdateUser(t *testing.T) {
	us, err := testingUserService()
	if err != nil {
		t.Fatal()
	}

	// Create user
	user := User{
		Model: gorm.Model{},
		Name:  "Nemod Marg",
		Email: "nemod@testupdateuser.com",
	}
	_, err = us.Create(&user)
	if err != nil {
		t.Fatal(err)
	}

	// Change user email address
	user.Email = "n@t.com"
	err = us.Update(&user)
	if err != nil {
		t.Fatal(err)
	}

	// Get updated user
	newUserData, err := us.ByID(user.ID)
	if err != nil {
		t.Fatal(err)
	}
	got := newUserData.Email
	want := "n@t.com"
	if got != want {
		t.Fatalf("want email %s, got email %s", want, got)
	}
}

func TestDeleteUser(t *testing.T) {
	us, err := testingUserService()
	if err != nil {
		t.Fatal(err)
	}

	// Create fake user
	userData := User{
		Model: gorm.Model{},
		Name:  "Delete",
		Email: "delete@delete.com",
	}
	_, err = us.Create(&userData)
	if err != nil {
		t.Fatal(err)
	}

	// Find created user and verify it was saved
	user, err := us.ByID(userData.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Delete new user
	err = us.Delete(user.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Try to find the deleted user again
	_, gotErr := us.ByID(userData.ID)
	wantErr := ErrNotFound
	if gotErr != wantErr {
		t.Fatalf("After deleting the user we got error '%v', but want error '%v'", gotErr, wantErr)
	}
}

func TestByEmail(t *testing.T) {
	us, err := testingUserService()
	if err != nil {
		t.Fatal(err)
	}

	// Create fake user
	userData := User{
		Model: gorm.Model{},
		Name:  "ByEmail",
		Email: "byemail@byemail.com",
	}
	_, err = us.Create(&userData)
	if err != nil {
		t.Fatal(err)
	}

	// Get user from DB by email
	user, err := us.ByEmail(userData.Email)
	if err != nil {
		t.Fatal(err)
	}

	// Check results
	want := User{Name: userData.Name, Email: userData.Email}
	got := User{Name: user.Name, Email: user.Email}
	if got != want {
		t.Fatalf("ByEmail() = %v, want %v", got, want)
	}
}

func TestByID(t *testing.T) {
	us, err := testingUserService()
	if err != nil {
		t.Fatal(err)
	}

	// Create fake user
	userData := User{
		Model: gorm.Model{},
		Name:  "ByID",
		Email: "byid@byid.com",
	}
	_, err = us.Create(&userData)
	if err != nil {
		t.Fatal(err)
	}

	// Get user from DB by id
	user, err := us.ByID(userData.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Check results
	want := User{Name: userData.Name, Email: userData.Email}
	got := User{Name: user.Name, Email: user.Email}
	if got != want {
		t.Fatalf("ByID() = %v, want %v", got, want)
	}
}
