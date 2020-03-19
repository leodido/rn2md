# release note to markdown

> Generate markdown for your changelogs from release-note blocks

## Usage

```bash
rn2md -o falcosecurity -m 0.21.0 -r falco
```

## Help

```bash
./rn2md --help                                                                                             2.89 
Little configurable CLI to generate the markdown for your changelos from release-note blocks found into your project pull requests.

Usage:
  rn2md [flags]

Flags:
  -b, --branch string      the target branch you want to filter by the pull requests (default "master")
  -h, --help               help for rn2md
  -m, --milestone string   the milestone you want to filter by the pull requests
  -o, --org string         the github organization
  -r, --repo string        the github repository name
```