package the_platinum_searcher_test

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"unicode/utf8"

	pt "github.com/monochromegane/the_platinum_searcher"
)

func TestSmallG(t *testing.T) {
	tests := []struct {
		args []string
		want []string
	}{
		{
			args: []string{"--nocolor", "-g", "context", "files"},
			want: []string{"files/context/context.txt"},
		},
		{
			args: []string{"--nocolor", "-g", "Context", "files"},
			want: []string{"files/depth/Context.txt"},
		},
		{
			args: []string{"--nocolor", "-g", "conteXt", "files"},
			want: []string{""},
		},
		{
			args: []string{"--nocolor", "-ig", "context", "files"},
			want: []string{"files/context/context.txt", "files/depth/Context.txt"},
		},
		{
			args: []string{"--nocolor", "-ig", "Context", "files"},
			want: []string{"files/context/context.txt", "files/depth/Context.txt"},
		},
		{
			args: []string{"--nocolor", "-ig", "conteXt", "files"},
			want: []string{"files/context/context.txt", "files/depth/Context.txt"},
		},
		{
			args: []string{"--nocolor", "-Sg", "context", "files"},
			want: []string{"files/context/context.txt", "files/depth/Context.txt"},
		},
		{
			args: []string{"--nocolor", "-Sg", "Context", "files"},
			want: []string{"files/depth/Context.txt"},
		},
		{
			args: []string{"--nocolor", "-Sg", "conteXt", "files"},
			want: []string{""},
		},
	}

	var buf bytes.Buffer
	for _, test := range tests {
		pt := pt.PlatinumSearcher{Out: &buf, Err: os.Stderr}
		pt.Run(test.args)
		got := []string{""}
		str := buf.String()
		if str != "" {
			got = strings.Fields(str)
		}
		if !compareIgnoringOrder(test.want, got) {
			t.Errorf("Args: %q, Want: %q, Got: %q\n", test.args, test.want, got)
		}
		buf.Reset()
	}
}

func TestBigG(t *testing.T) {
	tests := []struct {
		args []string
		want string
	}{
		{
			[]string{"--nocolor", "-G", "Context", "Capital", "files"},
			"files/depth/Context.txt:1:Capital test",
		},
		{
			[]string{"--nocolor", "-G", "Context", "capital", "files"},
			"",
		},
		{
			[]string{"--nocolor", "-G", "context", "Capital", "files"},
			"",
		},
		{
			[]string{"--nocolor", "-iG", "Context", "capital", "files"},
			"files/depth/Context.txt:1:Capital test",
		},
		{
			[]string{"--nocolor", "-iG", "Context", "Capital", "files"},
			"files/depth/Context.txt:1:Capital test",
		},
		{
			[]string{"--nocolor", "-iG", "Context", "capitaL", "files"},
			"files/depth/Context.txt:1:Capital test",
		},
		{
			[]string{"--nocolor", "-iG", "context", "Capital", "files"},
			"",
		},
		{
			[]string{"--nocolor", "-SG", "Context", "capital", "files"},
			"files/depth/Context.txt:1:Capital test",
		},
		{
			[]string{"--nocolor", "-SG", "Context", "Capital", "files"},
			"files/depth/Context.txt:1:Capital test",
		},
		{
			[]string{"--nocolor", "-SG", "Context", "capitaL", "files"},
			"",
		},
	}

	var buf bytes.Buffer
	for _, test := range tests {
		pt := pt.PlatinumSearcher{Out: &buf, Err: os.Stderr}
		pt.Run(test.args)
		got := buf.String()
		raw := strings.Split(got, "\n")
		if len(raw) > 0 {
			got = strings.Join(raw, ":")
			for r, n := utf8.DecodeLastRuneInString(got); r == ':'; r, n = utf8.DecodeLastRuneInString(got) {
				got = got[:len(got)-n]
			}
		}
		if test.want != got {
			t.Errorf("Args: %q, Want: %q, Got: %q\n", test.args, test.want, got)
		}
		buf.Reset()
	}
}

func compareIgnoringOrder(want, got []string) bool {
	if len(want) != len(got) {
		return false
	}
	seen := make(map[string]bool)
	for _, k := range want {
		seen[k] = true
	}
	for _, k := range got {
		if !seen[k] {
			return false
		}
	}
	return true
}
