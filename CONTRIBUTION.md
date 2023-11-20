# Develop xCherry Golang SDK

## Repo layout

Any contribution is welcome. Even just a fix for typo in a code comment, or README/wiki.

Here is the repository layout if you are interested to learn about it:

* `integ/` the end to end integration tests.
    * `init.go` the initiation & registration of testing processes. It's using global variables just for convenience
    * `main_test` the setup + tear down for running local in-memory xCherry worker with GoSDK
    * `xyz_test` the test for a test case xyz
    * `xyz_process.go` the test process for a test xyz
    * `xyz_process_state_*` the test process states for a test xyz
* `xc` the main directory
  * `*_impl.go` these are implementation for SDK. Ideally we should put them in separate folder, but Golang doesn't allow circular dependency, and we hate to use alias across packages
  * `internal_*.go` these are implementation for SDK
  * `_test.go` the unit test
  * other `.go` the interfaces defined in this SDK for user to use


### Coding convention 
There are some conventions here:
* The private struct shouldn't let other structs to access their private fields even they are in the same package -- All the SDK implementation are in the same package because of circular dependency issues. It's possible to write the code with the random access, but it would be a nightmare to maintain. So it's recommended to always expose methods (like `getter`) instead.