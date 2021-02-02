package main

func ExampleStorageClient() {
  addrs := []string{
    "127.0.0.1:8080",
    "127.0.0.1:8081",
    "127.0.0.1:8082",
  }

  mgr, err := NewManager(addrs, WithGrpcDialOptions(
    grpc.WithBlock(),
    grpc.WithInsecure(),
  ),
    WithDialTimeout(500*time.Millisecond),
  )
  if err != nil {
    log.Fatal(err)
  }