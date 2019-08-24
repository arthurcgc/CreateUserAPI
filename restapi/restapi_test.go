package restapi

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/arthurcgc/CreateUserAPI/myuser"
)

func TestGetUser(t *testing.T) {
	type args struct {
		w  http.ResponseWriter
		r  *http.Request
		db *data.data
	}
	tests := []struct {
		name    string
		args    args
		want    *myuser.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUser(tt.args.w, tt.args.r, tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
