# Git Conventions

We use Git for our version control system. The `master` branch is the home of the current development candidate. Releases are tagged.

We accept changes to the code via GitHub Pull Requests (PRs). One workflow for doing this is as follows:

1. Go to your `$GOPATH/src/deis` directory and `git clone` the `github.com/deis/steward` repository.
2. Fork that repository into your GitHub account
3. Add your repository as a remote for `$GOPATH/src/github.com/deis/steward`
4. Create a new working branch (`git checkout -b my-feature`) and do your work on that branch.
5. When you are ready for us to review, push your branch to GitHub, and then open a new pull request with us.

All git commit messages should loosely follow [semantic commit messages](http://karma-runner.github.io/0.13/dev/git-commit-msg.html). We've relaxed the requirement that a full commit message body be present, however you must indicate if your commit closes any issue. See below for an example commit message:

```console
test(cf mode): Add more integration tests

Fixes #1234
```

Common commit types:

- `fix`: Fix a bug or error
- `feat`: Add a new feature
- `ref`: refactor some code
- `doc`: Change documentation
- `test`: Improve testing

Common scopes:

- `k8s`: general interaction with the Kubernetes API
- `k8s/claim`: CRUD actions on claims
- `k8s/claim/state`: the claim state machine
- `mode`: generic mode interfaces and data types
- `mode/{cf,helm,cmd}`: specific mode functionality
- `*`: two or more scopes

Read more:
- The [Deis Guidelines](https://github.com/deis/workflow/blob/master/src/contributing/submitting-a-pull-request.md)
  were the inspiration for this section.
- Karma Runner [defines](http://karma-runner.github.io/0.13/dev/git-commit-msg.html) the semantic commit message idea.

### Go Conventions

We follow the Go coding style standards very closely. Typically, running `go fmt` will make your code beautiful for you.

We also typically follow the conventions recommended by `go lint` and `govet`. We encourage you to install an extension to your IDE that automatically runs `go fmt` and `go vet` against your code as you develop.

Read more:

- Effective Go [introduces formatting](https://golang.org/doc/effective_go.html#formatting).
- The Go Wiki has a great article on [formatting](https://github.com/golang/go/wiki/CodeReviewComments).
