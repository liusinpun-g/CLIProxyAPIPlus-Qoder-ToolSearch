# Publishing this fork

## Repository identity

Suggested repository: `liusinpun-g/CLIProxyAPIPlus-Qoder-ToolSearch` (public).

Suggested description: `Community fork of CLIProxyAPI Plus with opt-in, in-process Qoder tool search for Codex / Responses. MIT licensed.`

Create a **GitHub fork of `kaitranntt/CLIProxyAPIPlus`**, optionally rename it, then push this integration branch to the fork. GitHub's fork banner records the platform relationship; `README.md` and `THIRD_PARTY_NOTICES.md` record the direct base and the additional imported source. GitHub may display a fork-network root rather than every intermediate fork, so the written ancestry remains authoritative for this integration.

The repository currently has a shallow local clone. Fetch missing history from the public base before pushing if Git reports a shallow-update error. Do not flatten the project into a new initial commit or replace the base's history.

## License choice

Continue using MIT for this fork and its new integration code. It allows commercial use, modification and redistribution while requiring preservation of copyright and permission notices. Preserve the upstream notices and imported core license; do not replace them with a public-domain dedication or imply ownership of upstream work. Dependency licenses remain separate. See `LICENSE` and `THIRD_PARTY_NOTICES.md`.

## Prepare the first release

1. Review the source diff and the public-facing author identity. Keep credentials, real configuration, private endpoint values, request logs and local probe output outside Git.
2. Include the new integration sources/tests, root READMEs, `LICENSE`, notices, release notes, publishing guide, `docs/qoder-tool-search.md`, imported license/provenance/checksums, scripts and workflow changes. Preserve unrelated working-tree edits outside the release commit.
3. Verify the imported files with `shasum -a 256 -c internal/responsestools/UPSTREAM.sha256`. Run the affected-package tests, focused race checks and a build described in `docs/qoder-tool-search.md`. Retain the known full-suite failure in the release notes unless a new run establishes otherwise.
4. Commit the reviewed change and push it to the fork. Make the reviewed integration branch the default branch (or merge it into the fork's default branch).
5. Create a new annotated tag `v7.3.12-qoder-tools.1` on that exact commit and push only that new tag. Do not force-push or replace upstream tags.
6. Enable Actions for the fork and run **Qoder tool-search prerelease**, selecting that existing tag. It builds the macOS arm64 archive, collects dependency notices and uploads a **draft prerelease** using `RELEASE_NOTES.md`. No upstream synchronization or Docker publishing is enabled.
7. Review the draft's commit, assets, SHA-256 checksums and notes, then publish the prerelease. Do not mark it as a stable/latest release.

The workflow uses the repository's `GITHUB_TOKEN`; it does not require a personal access token. Creating the fork and pushing from this Mac requires an authenticated GitHub session or Git transport. The release script refuses a dirty checkout or a tag that does not match HEAD.

## Local packaging

On a clean checkout of the tagged commit:

```sh
bash scripts/package-qoder-release.sh v7.3.12-qoder-tools.1
```

The script builds a CGO-disabled macOS arm64 executable, saves dependency notices with `github.com/google/go-licenses/v2@v2.0.1`, includes MIT/provenance documents and creates an archive plus `checksums.txt` under `dist/`. It does not publish anything or read local provider credentials. A license-collection error stops packaging.

The inherited `.goreleaser.yml` is retained as upstream reference and is not used by this fork's first prerelease workflow.
