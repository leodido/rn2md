# release note to markdown

> Generate markdown for your changelogs from release-note blocks

It expects release-note block rows to follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) format.

Then, it extracts the `type` and creates different sections in the resulting markdown for different `type`.

For example `new: ...` and `BREAKING CHANGE: ...` release-note rows populate the "Major Changes" section of the markdown.

## Usage

```bash
rn2md -o falcosecurity -m 0.21.0 -r falco
```

## Help

```
./rn2md --help
Little configurable CLI to generate the markdown for your changelos from release-note blocks found into your project pull requests.

Usage:
  rn2md [flags]

Flags:
  -b, --branch string      the target branch you want to filter by the pull requests (default "master")
  -h, --help               help for rn2md
  -m, --milestone string   the milestone you want to filter by the pull requests
  -o, --org string         the github organization
  -r, --repo string        the github repository name
  -t, --token string       a GitHub personal API token to perform authenticated requests
```

## TODO

- [ ] gen markdown table for statistics

---

[![Analytics](https://ga-beacon.appspot.com/UA-49657176-1/rn2md?flat)](https://github.com/igrigorik/ga-beacon)
