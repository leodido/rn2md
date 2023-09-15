# release notes to markdown

> Generate markdown for your changelogs from release-note blocks

It expects release-note block rows to follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) format.

Then, it extracts the `type` and creates different sections in the resulting markdown for different `type`.

For example `new: ...` and `BREAKING CHANGE: ...` release-note rows populate the "Major Changes" section of the markdown.

## Usage

```bash
rn2md -r falcosecurity/falco -m 0.21.0
```

## Help

```
./rn2md --help
Little configurable CLI to generate the markdown for your changelogs from release-note blocks found into your project pull requests.

Usage:
  rn2md [flags]

Flags:
  -b, --branch string      the target branch you want to filter by the pull requests (default "master")
  -h, --help               help for rn2md
  -m, --milestone string   the milestone you want to filter by the pull requests
  -r, --repo string        the full github repository name (org/repo)
  -t, --token string       a GitHub personal API token to perform authenticated requests
```

## Using the github action in your repo

To automatically generate release notes markdown for your project milestone, you must just add a step to your workflow.

```yaml
  - name: rn2md
    uses: leodido/rn2md@latest
    with:
      # The milestone you want to filter by the pull requests. Required.
      milestone: 0.21.0
      # A github token needed for the github client API calls (listing repo PRs). Defaults to ${{ github.token }}
      token: mytoken
      # Full name for your repo. Defaults to ${{ github.repository }}
      repo: myorg/myrepo 
      # Target branch to filter by the pull requests. Defaults to ${{ github.event.repository.default_branch }}.
      branch: main
      # Target file to be generated. Defaults to ${{ github.workspace }}/release_body.md
      output: out.md
```

---

[![Analytics](https://ga-beacon.appspot.com/UA-49657176-1/rn2md?flat)](https://github.com/igrigorik/ga-beacon)
