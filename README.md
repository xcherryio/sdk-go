# XDB Golang SDK
[![Go Reference](https://pkg.go.dev/badge/github.com/xdblab/xdb-golang-sdk.svg)](https://pkg.go.dev/github.com/xdblab/xdb-golang-sdk)
[![Go Report Card](https://goreportcard.com/badge/github.com/xdblab/xdb-golang-sdk)](https://goreportcard.com/report/github.com/xdblab/xdb-golang-sdk)
[![Coverage Status](https://codecov.io/github/xdblab/xdb-golang-sdk/coverage.svg?branch=release)](https://app.codecov.io/gh/xdblab/xdb-golang-sdk/branch/main)

[![Build status](https://github.com/xdblab/xdb-golang-sdk/actions/workflows/ci-integ-test.yml/badge.svg?branch=release)](https://github.com/xdblab/xdb-golang-sdk/actions/workflows/ci-integ-test.yml)


Golang SDK for [XDB](https://github.com/xdblab/xdb)

See [samples](https://github.com/xdblab/xdb-golang-samples) for how to use this SDK.
# Contribution
See [contribution guide](CONTRIBUTION.md)

# Development Plan

## MVP
- [ ] Start ProcessExecution
  - [ ] ProcessIdReusePolicy
- [ ] Executing `wait_until`/`execute` APIs 
- [ ] StateDecision: force completing process
- [ ] Parallel execution of multiple states
- [ ] WaitForProcessCompletion API
- [ ] StateOption: WaitUntil/Execute API timeout and retry policy
- [ ] AnyOfCompletion and AllOfCompletion waitingType
- [ ] TimerCommand
- [ ] LocalQueue
  - [ ] LocalQueueCommand
  - [ ] Deduplication
- [ ] LocalAttribute + GlobalAttribute
  - [ ] LoadingPolicy
  - [ ] InitialUpsert
  - [ ] Multi-tables for GlobalAttribute
- [ ] Stop ProcessExecution
- [ ] Error handling for canceled, failed, timeout, terminated
- [ ] AsyncState failure policy for recovery 
- [ ] RPC
- [ ] Retention policy for execution
- [ ] WaitForStateCompletion API
- [ ] ResetProcessExecution
- [ ] Describe ProcessExecution API
- [ ] Atomic conditional complete workflow by checking queue emptiness

## Future

- [ ] Skip timer API for testing/operation
- [ ] Dynamic attributes and queue definition
- [ ] State options overridden dynamically
- [ ] Consume more than one messages in a single command
  - [ ] FIFO/BestMatch policies
- [ ] WaitingType: AnyCombinationsOf
- [ ] GlobalQueue
- [ ] CronSchedule
- [ ] DelayStart
- [ ] Caching
- [ ] Custom Database Query
- [ ] SearchAttribute
- [ ] ExternalAttribute