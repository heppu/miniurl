# MiniURL
Skeleton for URL shortener project.

- [pr.yaml](.github/workflows/pr.yaml) defines GHA workflow for pull requests
- [Makefile](Makefile) has stubs for required steps and some helper code
- [openapi.yaml](openapi.yaml) defines expected REST API
- [index.html](ui/index.html) implement extremely cool UI for our app

Fork and fill in the blanks =)

## Requirements
- [Make](https://www.gnu.org/software/make/)
- [Go](https://go.dev/) >= 1.13
- Docker or Podman with compose for integration tests with [testcontainers](https://golang.testcontainers.org/)
- awk for `make help`
