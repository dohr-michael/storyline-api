storyline-api
===========
 

[![Build Status](https://dev.azure.com/dohrmichael/storyline/_apis/build/status/dohr-michael.storyline-api?branchName=master)](https://dev.azure.com/dohrmichael/storyline/_build/latest?definitionId=1&branchName=master)

- Install and run NATS server :
    - `go get github.com/nats-io/gnatsd`
    - run `gnatsd`

- Using [go mod](https://github.com/golang/go/wiki/Modules)
- Build in with [magefile](https://magefile.org/)

- Commands (from Magefile)
    - Run unit test : `mage test` 
    - Build locally : `mage build`
