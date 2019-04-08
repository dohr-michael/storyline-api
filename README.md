storyline-api
===========
 

- Needs to have github.com/jteeuwen/go-bindata
    - `go get -u github.com/jteeuwen/go-bindata/...` in another terminal 


- Install and run NATS server :
    - `go get github.com/nats-io/gnatsd`
    - run `gnatsd`

- Using [go mod](https://github.com/golang/go/wiki/Modules)
- Build in with [magefile](https://magefile.org/)
- Release management with [GoReleaser](https://goreleaser.com/)

- Commands (from Magefile)
    - Run unit test : `mage test` 
    - Build locally : `mage build` 
    - Snapshot release (with container generation): `mage snapshot` 
    - Release (if tag is specify): `mage release`
- Instal dependencies
```
go mod download
```
- Run project
```
go run main.go start
```
- Hot Reload
```
go get github.com/codegangsta/gin
gin --appPort 8080 --buildArgs main.go -i run start
```

- Needs
    - Replacement in the project (in order):
        - `storyline-api` by your project name
        - `github.com/dohr-michael` by your github / git path / by your go package root.
        - `jdrbahamut` by your docker registry organization.
    - environment variables
        - GITHUB_TOKEN
    - Docker registry authentication
- Release Generate :
    - Binaries
        - darwin_amd64
        - darwin_386
        - linux_amd64
        - linux_arm64
        - linux_386
        - windows_amd64
        - windows_386
    - Docker Container
        - linux_amd64
    - Github release
    - Homebrew Tap
    - Scoop recipe
