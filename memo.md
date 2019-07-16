### Newでオブジェクトを返し、と同時にそのオブジェクト用にbackgroundで何か動くとき

```go
obj := New()
obj.Close()
```

* Closeメソッドを生やす。Closeメソッドは以下1,2,3をやる。
    * (1) backgroundを止めようとする
    * (2) 本当に止まるまで同期的に待つ
    * (3) そしてreturnする

### ch := newChannel(stop)的な関数

```go
stop := make(chan struct{})
go func() {
    <-sig
    close(stop)
}()
ch := newChannel(stop)

for c := range ch {
    // ...
}
```

* stopをcloseしたら以下1,2,3が起こるようにする
    * (1) newChannelを呼んだ後、内部で動いているロジックがあるはずでそれを止めようとする
    * (2) 本当に止まるまで同期的に待つ
    * (3) そしてchをcloseする


### obj := New(); obj.Run(stop)的なメソッドがあるとき

```go
stop := make(chan struct{})
go func() {
    <-sig
    close(stop)
}()
obj := New()
obj.Run(stop)
```

* stopをcloseしたら以下1,2,3が起こるようにする
    * (1) Run内部のロジックを止めようとする
    * (2) 本当に止まるまで同期的に待つ
    * (3) Runがreturnする
