## Why make this pull request?

[Explain why you are making this pull request and what problem it solves.]

## What has changed

[Summarize what components of the repo is updated]

[Link to apis/sdk-go PRs if it's on top of any API changes]

- API change link: ...

## How to test this pull request?

[If write integration test in this SDK for server, merge this PR before the server PR. Then after server PR is merged, rerun the ci test]

## Checklist before merge
[ ] If applicable, merge the apis/apis PRs to main branch
[ ] Update `go.mod` to use the commitID of the main branches for apis

## Action After merge
[ ] If depending on server, rerun the [ci-integration test](https://github.com/xcherryio/sdk-go/actions) after the server has published the new image