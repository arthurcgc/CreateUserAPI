package user

import "testing"

func TestGetUserName(t *testing.T) {
	client := user{name: "Arthur", email: "arthur@gmail.com"}
	want := "Arthur"
	got := client.GetUserName()
	if want != got {
		t.Errorf("Wrong username")
	}
}

func TestGetUserEmail(t *testing.T) {
	client := user{name: "Arthur", email: "arthur@gmail.com"}
	want := "arthur@gmail.com"
	got := client.GetUserEmail()
	if want != got {
		t.Errorf("Wrong email")
	}
}
