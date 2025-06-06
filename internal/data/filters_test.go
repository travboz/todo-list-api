package data_test

import (
	"reflect"
	"testing"

	"github.com/travboz/backend-projects/todo-list-api/internal/data"
)

func TestFilters(t *testing.T) {
	t.Run("Calculate metadata gives correct metadata", func(t *testing.T) {
		totalRecords := 56
		CurrentPage := 2
		PageSize := 10

		want := data.Metadata{
			CurrentPage,
			PageSize,
			totalRecords,
		}

		got := data.CalculateMetadata(totalRecords, CurrentPage, PageSize)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %v; want: %v", got, want)
		}
	})

	t.Run("Gives empty on no records", func(t *testing.T) {
		totalRecords := 0
		CurrentPage := 2
		PageSize := 10

		want := data.Metadata{}

		got := data.CalculateMetadata(totalRecords, CurrentPage, PageSize)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %v; want: %v", got, want)
		}
	})

}
