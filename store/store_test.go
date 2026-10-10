package store

import (
	"testing"
)

func TestGetString(t *testing.T) {

	key := "existing"
	val := "Theo"

	SetString(key, val)
	// TODO set list here
	cases := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"Key Exists", "existing", "Theo", false},
		{"Non-existent Key", "non_existing_key", "", true},
		{"Wrong Type", "array_key", "", true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			got, err := GetString(c.input)

			if c.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got none.")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected an error: %v", err)
			}
			if got != c.want {
				t.Errorf("got %s: want %s", got, c.want)
			}

		})
	}
}

func TestSetList(t *testing.T) {
	cases := []struct {
		name    string
		key     string
		value   []string
		want    int
		wantErr bool
	}{
		{"Valid Input 2 values", "mylist", []string{"Brian", "Espina"}, 2, false},
		{"Valid Input 1 value", "mylist", []string{"Ramirez"}, 1, false},
		{"Valid Input 3 values", "mylist", []string{"Theo", "Pinion", "Espina"}, 3, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := SetList(c.key, c.value...)
			if c.wantErr {
				if err == nil {
					t.Fatalf("Expected an error, got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected an error: %v", err)
			}

			if got != c.want {
				t.Errorf("want %d, got %d", c.want, got)
			}

		})
	}
}
