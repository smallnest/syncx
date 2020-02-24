# syncx

[![License](https://img.shields.io/:license-apache%202-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![GoDoc](https://godoc.org/github.com/smallnest/syncx?status.png)](http://godoc.org/github.com/smallnest/syncx)  [![travis](https://travis-ci.org/smallnest/syncx.svg?branch=master)](https://travis-ci.org/smallnest/syncx) [![Go Report Card](https://goreportcard.com/badge/github.com/smallnest/syncx)](https://goreportcard.com/report/github.com/smallnest/syncx) [![coveralls](https://coveralls.io/repos/smallnest/syncx/badge.svg?branch=master&service=github)](https://coveralls.io/github/smallnest/syncx?branch=master) 


More Advanced sync primitives.

Implemented some Go sync primitives.


## Token

provides token implementation.

Only the one thats owns the `Token` can do stuff and then it can handoffs the token to others.


## Batch

provides batch implementation.

It is like `errgroup` and can return all errors results of each task.

## Any

provides partial batch implementation.

You can wait some tasks have finished and returns.

## Other advanced sync primitives

- [singleflight](https://github.com/golang/sync/tree/master/singleflight): provides a duplicate function call suppression
- [errgroup](https://github.com/golang/sync/blob/master/errgroup/errgroup.go): provides synchronization, error propagation, and Context cancelation for groups of goroutines working on subtasks of a common task
- [semaphore](https://github.com/golang/sync/tree/master/semaphore): provides a weighted semaphore implementation