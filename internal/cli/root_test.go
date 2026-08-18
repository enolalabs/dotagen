package cli

import "testing"

func TestNormalizeVersion(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name  string
		value string
		want  string
	}{
		{name: "release tag", value: "v2.8.0", want: "2.8.0"},
		{name: "release number", value: "2.8.0", want: "2.8.0"},
		{name: "development", value: "dev", want: "dev"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := normalizeVersion(tc.value); got != tc.want {
				t.Fatalf("normalizeVersion(%q) = %q, want %q", tc.value, got, tc.want)
			}
		})
	}
}
