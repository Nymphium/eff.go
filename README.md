Eff.go
=======

ONE-SHOT Algebiraic Effects for Golang!
It is based on channel, goroutine and [Russ Cox's coroutine using them](https://research.swtch.com/coro).

# One-shot algebraic effects
You can access the delimited continuation which can run only once. Even the limitation exists, you can write powerful control flow manipulation, like async/await, call/1cc.
We have an formal definition for the implementation, by showing a conversion from algebraic effects and handlers to asymmetric coroutines.
[See here (in Japanese)](https://nymphium.github.io/2018/12/09/asymmetric-coroutines%E3%81%AB%E3%82%88%E3%82%8Boneshot-algebraic-effects%E3%81%AE%E5%AE%9F%E8%A3%85.html).
