package lookup

import (
	"reflect"
	"testing"
)

func Test_findIPInResultsMap(t *testing.T) {
	type args struct {
		data map[string]int
		ip   string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]int
		wantErr bool
	}{
		args {
			"Find IP in Map",
			map[string]int{
				"129.0.0.0/24" : 1299
				"143.20.19.0/16" : 1011
			},
			"129.0.0.1",
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findIPInResultsMap(tt.args.data, tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("findIPInResultsMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findIPInResultsMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
