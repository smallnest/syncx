# syncx

Advanced sync primitives.

Implemented some Go sync primitives.


## Token

provides token implementation.

Only the one thats owns the `Token` can do stuff and then it can handoffs the token to others.

- Accquire


## Batch


## Other advanced sync primitives

- [singleflight](https://github.com/golang/sync/tree/master/singleflight): provides a duplicate function call suppression
- [errgroup](https://github.com/golang/sync/blob/master/errgroup/errgroup.go): provides synchronization, error propagation, and Context cancelation for groups of goroutines working on subtasks of a common task
- [semaphore](https://github.com/golang/sync/tree/master/semaphore): provides a weighted semaphore implementation