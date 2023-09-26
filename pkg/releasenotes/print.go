package releasenotes

import (
	"bytes"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"
)

const templ = `## v{{ .Milestone }}


Released on {{ .Day }}


### Major Changes

{{ if .MajorNotes }}
{{ range .BreakingNotes }}
* {{ .Description }} [[#{{ .Num }}]({{ .URI }})] - [{{ .Author }}]({{ .AuthorURL }})
{{ end }}
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
	Milestone     string
	Day           string
	BreakingNotes []ReleaseNote
	MajorNotes    []ReleaseNote
	MinorNotes    []ReleaseNote
	FixNotes      []ReleaseNote
	RuleNotes     []ReleaseNote
	NoneNotes     []ReleaseNote
}

func (notes ReleaseNotes) Print(milestone string) error {
	var breaking []ReleaseNote
	var majors []ReleaseNote
	var minors []ReleaseNote
	var fixes []ReleaseNote
	var rules []ReleaseNote
	var none []ReleaseNote
	for _, n := range notes {
		switch n.Typology {
		case "BREAKING CHANGE":
			breaking = append(breaking, n)
		case "fix":
			fixes = append(fixes, n)
		case "rule":
			rules = append(rules, n)
		case "new", "feat":
			majors = append(majors, n)
		case "none":
			none = append(none, n)
		default:
			minors = append(minors, n)
		}
	}

	data := templateData{
		Milestone:     milestone,
		Day:           time.Now().Format("2006-01-02"),
		MinorNotes:    minors,
		BreakingNotes: breaking,
		MajorNotes:    majors,
		FixNotes:      fixes,
		RuleNotes:     rules,
		NoneNotes:     none,
	}

	t := template.New("changes")
	res, err := t.Parse(templ)
	if err != nil {
		return err
	}

	b := bytes.NewBuffer(nil)
	err = res.Execute(b, data)
	if err != nil {
		return err
	}

	result := strings.ReplaceAll(b.String(), "\n\n", "\n")
	fmt.Println(result)
	return nil
}

func (s *Statistics) Print() error {
	fmt.Println("### Statistics")
	fmt.Println("")

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Merged PRs", "Number"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	table.Append([]string{"Not user-facing", strconv.FormatInt(s.nonFacing, 10)})
	table.Append([]string{"Release note", strconv.FormatInt(s.total-s.nonFacing, 10)})
	table.Append([]string{"Total", strconv.FormatInt(s.total, 10)})

	table.Render() // Send output
	return nil
}
