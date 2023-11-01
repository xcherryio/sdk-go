# XDB Golang SDK
[![Go Reference](https://pkg.go.dev/badge/github.com/xdblab/xdb-golang-sdk.svg)](https://pkg.go.dev/github.com/xdblab/xdb-golang-sdk)
[![Go Report Card](https://goreportcard.com/badge/github.com/xdblab/xdb-golang-sdk)](https://goreportcard.com/report/github.com/xdblab/xdb-golang-sdk)

<!--- using build branch for coverage and build status. Because server depends on GoSDK, so main branch code not be testable until server has implemented the features --->
[![Coverage Status](https://codecov.io/github/xdblab/xdb-golang-sdk/coverage.svg?branch=build)](https://app.codecov.io/gh/xdblab/xdb-golang-sdk/branch/main)

[![Build status](https://github.com/xdblab/xdb-golang-sdk/actions/workflows/ci-tests.yml/badge.svg?branch=build)](https://github.com/xdblab/xdb-golang-sdk/actions/workflows/ci-tests.yml)


Golang SDK for [XDB](https://github.com/xdblab/xdb)

See [samples](https://github.com/xdblab/xdb-golang-samples) for how to use this SDK.
# Contribution
See [contribution guide](CONTRIBUTION.md)

# Development Plan

## 1.0
- [ ] StartProcessExecution API
  - [x] Basic
  - [x] ProcessIdReusePolicy
  - [ ] Process timeout
  - [ ] Retention policy after closed
- [x] Executing `wait_until`/`execute` APIs
  - [x] Basic
  - [x] Parallel execution of multiple states
  - [x] StateOption: WaitUntil/Execute API timeout and retry policy
  - [x] AsyncState failure policy for recovery
- [ ] StateDecision
  - [x] Single next State
  - [x] Multiple next states
  - [x] Force completing process
  - [x] Graceful completing process
  - [x] Force fail process
  - [x] Dead end
  - [ ] Conditional complete process with checking queue emptiness
- [ ] Commands
  - [ ] AnyOfCompletion and AllOfCompletion waitingType
  - [ ] TimerCommand
- [ ] LocalQueue
  - [ ] LocalQueueCommand
  - [ ] MessageId for deduplication
  - [ ] SendMessage API without RPC
- [ ] LocalAttribute persistence
  - [ ] LoadingPolicy (attribute selection + locking)
  - [ ] InitialUpsert
- [x] GlobalAttribute  persistence
  - [x] LoadingPolicy (attribute selection + locking)
  - [x] InitialUpsert
  - [x] Multi-tables 
- [ ] RPC
- [ ] API error handling for canceled, failed, timeout, terminated
- [ ] StopProcessExecution API
- [ ] WaitForStateCompletion API
- [ ] ResetStateExecution for operation
- [x] DescribeProcessExecution API
- [ ] WaitForProcessCompletion API
- [ ] History events for operation/debugging

## Future

- [ ] Skip timer API for testing/operation
- [ ] Dynamic attributes and queue definition
- [ ] State options overridden dynamically
- [ ] Consume more than one messages in a single command with FIFO/BestMatch policies
- [ ] WaitingType: AnyCombinationsOf
- [ ] GlobalQueue
- [ ] CronSchedule
- [ ] Batch operation
- [ ] DelayStart
- [ ] Caching (with Redis, etc)
- [ ] Custom Database Query
- [ ] SearchAttribute (with ElasticSearch, etc)
- [ ] ExternalAttribute (with S3, Snowflake, etc)