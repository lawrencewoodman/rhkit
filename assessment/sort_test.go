package assessment

import (
	"github.com/vlifesystems/rhkit/aggregator"
	"testing"
)

func TestMakeSortOrders(t *testing.T) {
	cases := []struct {
		descs []SortDesc
		want  []SortOrder
	}{
		{descs: []SortDesc{
			SortDesc{"income", "descending"},
			SortDesc{"numMatches", "ascending"},
			SortDesc{"percentMatches", "descending"},
		},
			want: []SortOrder{
				SortOrder{"income", DESCENDING},
				SortOrder{"numMatches", ASCENDING},
				SortOrder{"percentMatches", DESCENDING},
			},
		},
		{descs: []SortDesc{},
			want: []SortOrder{},
		},
	}
	fields := []string{"in"}
	aggregatorDescs := []*aggregator.Desc{
		{Name: "income", Kind: "sum", Arg: "in"},
	}
	aggregatorSpecs, err := aggregator.MakeSpecs(fields, aggregatorDescs)
	if err != nil {
		t.Fatalf("MakeSpecs: %s", err)
	}
	for _, c := range cases {
		got, err := MakeSortOrders(aggregatorSpecs, c.descs)
		if err != nil {
			t.Fatalf("MakeSortOrders: %s", err)
		}
		if len(got) != len(c.want) {
			t.Errorf("MakeSortOrders got: %s, want: %s", got, c.want)
			continue
		}
		for i, s := range got {
			if s.Aggregator != c.want[i].Aggregator {
				t.Errorf("MakeSortOrders got: %s, want: %s", got, c.want)
				continue
			}
			if s.Direction != c.want[i].Direction {
				t.Errorf("MakeSortOrders got: %s, want: %s", got, c.want)
				continue
			}
		}
	}
}

func TestMakeSortOrders_errors(t *testing.T) {
	cases := []struct {
		descs   []SortDesc
		wantErr error
	}{
		{descs: []SortDesc{
			SortDesc{"income", "DESCENDING"},
			SortDesc{"numMatches", "ascending"},
			SortDesc{"percentMatches", "descending"},
		},
			wantErr: SortOrderError{
				Aggregator: "income",
				Direction:  "DESCENDING",
				Err:        ErrInvalidDirection,
			},
		},
		{descs: []SortDesc{
			SortDesc{"income", "Descending"},
			SortDesc{"numMatches", "ascending"},
			SortDesc{"percentMatches", "descending"},
		},
			wantErr: SortOrderError{
				Aggregator: "income",
				Direction:  "Descending",
				Err:        ErrInvalidDirection,
			},
		},
		{descs: []SortDesc{
			SortDesc{"income", "descending"},
			SortDesc{"numMatches", "ASCENDING"},
			SortDesc{"percentMatches", "descending"},
		},
			wantErr: SortOrderError{
				Aggregator: "numMatches",
				Direction:  "ASCENDING",
				Err:        ErrInvalidDirection,
			},
		},
		{descs: []SortDesc{
			SortDesc{"income", "descending"},
			SortDesc{"numMatches", "Ascending"},
			SortDesc{"percentMatches", "descending"},
		},
			wantErr: SortOrderError{
				Aggregator: "numMatches",
				Direction:  "Ascending",
				Err:        ErrInvalidDirection,
			},
		},
		{descs: []SortDesc{
			SortDesc{"income", "descending"},
			SortDesc{"numMatches", "ascending"},
			SortDesc{"percentMatches", "fred"},
		},
			wantErr: SortOrderError{
				Aggregator: "percentMatches",
				Direction:  "fred",
				Err:        ErrInvalidDirection,
			},
		},
		{descs: []SortDesc{
			SortDesc{"income", "descending"},
			SortDesc{"boris", "ascending"},
			SortDesc{"percentMatches", "descending"},
		},
			wantErr: SortOrderError{
				Aggregator: "boris",
				Direction:  "ascending",
				Err:        ErrUnrecognisedAggregator,
			},
		},
	}
	fields := []string{"in"}
	aggregatorDescs := []*aggregator.Desc{
		{Name: "income", Kind: "sum", Arg: "in"},
	}
	aggregatorSpecs, err := aggregator.MakeSpecs(fields, aggregatorDescs)
	if err != nil {
		t.Fatalf("MakeSpecs: %s", err)
	}
	for _, c := range cases {
		_, err := MakeSortOrders(aggregatorSpecs, c.descs)
		if err == nil || err.Error() != c.wantErr.Error() {
			t.Errorf("MakeSortOrders err: %s, wantErr: %s", err, c.wantErr)
		}
	}
}

func TestSortOrderError(t *testing.T) {
	e := SortOrderError{
		Aggregator: "numMatches",
		Direction:  "Ascending",
		Err:        ErrInvalidDirection,
	}
	want := "problem with sort order - aggregator: numMatches, direction: Ascending (invalid direction)"
	got := e.Error()
	if got != want {
		t.Errorf("Error() got: %s, want: %s", got, want)
	}
}
