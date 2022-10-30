package merror

//func TestDiscover1(t *testing.T) {
//	ctx := context.Background()
//	cli, err := registry.NewClientWithProvider(ctx, "consul", "localhost:8500")
//	if err != nil {
//		panic(err)
//	}
//	cli.WatchService(ctx, "client_app_admin", nil, func(infos []registry.ServiceInfo) {
//		fmt.Printf("Got ServiceInfos %+v\n", infos)
//	})
//
//}
//
///*
//client_app_admin:
// app_name: client_app_admin
// service_name: client_app_admin
// provider: consul
// args:
//   address: localhost:8500
//*/
//func TestDiscover2(t *testing.T) {
//	ctx := context.Background()
//	r := zest.Region{}
//	fmt.Printf("region IsEmpty: %v\n", r.IsEmpty())
//	registry.WatchService(ctx, registry.ServiceWatchParams{
//		Provider:     "consul",
//		Address:      "localhost:8500",
//		ServiceParam: zest.ServiceParam{AppName: "client_app_admin", ServiceName: "client_app_admin", ZestRegion: r},
//	}, func(infos []registry.ServiceInfo) {
//		fmt.Printf("Got ServiceInfos %+v\n", infos)
//	})
//	time.Sleep(10 * time.Second)
//}
