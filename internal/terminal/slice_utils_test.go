package terminal

import "testing"

func TestLimitSlice(t *testing.T) {
	tests := []struct {
		name      string
		slice     []int
		limit     int
		wantMore  bool
		wantItems []int
	}{
		{
			name:      "less than limit",
			slice:     []int{1, 2},
			limit:     3,
			wantMore:  false,
			wantItems: []int{1, 2},
		},
		{
			name:      "equal to limit",
			slice:     []int{1, 2, 3},
			limit:     3,
			wantMore:  false,
			wantItems: []int{1, 2, 3},
		},
		{
			name:      "greater than limit",
			slice:     []int{1, 2, 3, 4},
			limit:     3,
			wantMore:  true,
			wantItems: []int{1, 2, 3},
		},
		{
			name:      "zero limit",
			slice:     []int{1, 2},
			limit:     0,
			wantMore:  true,
			wantItems: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMore, gotItems := limitSlice(tt.slice, tt.limit)
			if gotMore != tt.wantMore {
				t.Fatalf("want more=%v, got %v", tt.wantMore, gotMore)
			}
			if len(gotItems) != len(tt.wantItems) {
				t.Fatalf("want %d items, got %d", len(tt.wantItems), len(gotItems))
			}
			for i := range tt.wantItems {
				if gotItems[i] != tt.wantItems[i] {
					t.Fatalf("item %d mismatch: want %d got %d", i, tt.wantItems[i], gotItems[i])
				}
			}
		})
	}
}

func TestLimitedItemsAsString(t *testing.T) {
	tests := []struct {
		name  string
		items []string
		want  string
	}{
		{
			name:  "empty",
			items: []string{},
			want:  "",
		},
		{
			name:  "under limit",
			items: []string{"a", "b"},
			want:  "a\nb",
		},
		{
			name:  "at limit",
			items: []string{"a", "b", "c"},
			want:  "a\nb\nc",
		},
		{
			name:  "over limit",
			items: []string{"a", "b", "c", "d"},
			want:  "a\nb\nc\n...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := limitedItemsAsString(tt.items, func(s string) string { return s })
			if got != tt.want {
				t.Fatalf("want %q, got %q", tt.want, got)
			}
		})
	}
}
