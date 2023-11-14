module radish

go 1.20

require radishcache v0.0.0

require (
	go.starlark.net v0.0.0-20231101134539-556fd59b42f6 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
)

replace radishcache => ./radishcache
