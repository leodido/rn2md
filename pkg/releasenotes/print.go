package releasenotes

import (
	"bytes"
	"strings"
	"text/template"
	"time"
)

const templ = `## v{{ .Milestone }}


Released on {{ .Day }}


### Major Changes

{{ if .MajorNotes }}
{{ range .MajorNotes }}
* {{ .Description }} [[#{{ .Num }}]({{ .URI }})] - [{{ .Author }}]({{ .AuthorURL }})
{{ end }}
{{ end }}

{{ if .MinorNotes }}
### Minor Changes

{{ range .MinorNotes }}
* {{ .Description }} [[#{{ .Num }}]({{ .URI }})] - [{{ .Author }}]({{ .AuthorURL }})
{{ end }}
{{ end }}

{{ if .FixNotes }}
### Bug Fixes

{{ range .FixNotes }}
* {{ .Description }} [[#{{ .Num }}]({{ .URI }})] - [{{ .Author }}]({{ .AuthorURL }})
{{ end }}
{{ end }}

{{ if .RuleNotes }}
### Rule Changes

{{ range .RuleNotes }}
* {{ .Description }} [[#{{ .Num }}]({{ .URI }})] - [{{ .Author }}]({{ .AuthorURL }})
{{ end }}
{{ end }}

{{ if .NoneNotes }}
### Non user-facing changes

{{ range .NoneNotes }}
* {{ .Description }} [[#{{ .Num }}]({{ .URI }})] - [{{ .Author }}]({{ .AuthorURL }})
{{ end }}
{{ end }}`

type templateData struct {
	Milestone  string
	Day        string
	MajorNotes []ReleaseNote
	MinorNotes []ReleaseNote
	FixNotes   []ReleaseNote
	RuleNotes  []ReleaseNote
	NoneNotes  []ReleaseNote
}

// Print ...
func Print(milestone string, notes []ReleaseNote) (string, error) {
	majors := []ReleaseNote{}
	minors := []ReleaseNote{}
	fixes := []ReleaseNote{}
	rules := []ReleaseNote{}
	none := []ReleaseNote{}
	for _, n := range notes {
		switch n.Typology {
		case "fix":
			fixes = append(fixes, n)
			break
		case "rule":
			rules = append(rules, n)
			break
		case "BREAKING CHANGE":
			fallthrough
		case "new":
			majors = append(majors, n)
			break
		case "none":
			none = append(none, n)
			break
		default:
			minors = append(minors, n)
			break
		}
	}

	data := templateData{
		Milestone:  milestone,
		Day:        time.Now().Format("2006-01-02"),
		MinorNotes: minors,
		MajorNotes: majors,
		FixNotes:   fixes,
		RuleNotes:  rules,
		NoneNotes:  none,
	}

	t := template.New("changes")
	res, err := t.Parse(templ)
	if err != nil {
		return "", err
	}

	b := bytes.NewBuffer(nil)
	err = res.Execute(b, data)
	if err != nil {
		return "", err
	}

	result := strings.ReplaceAll(b.String(), "\n\n", "\n")

	return result, nil
}
