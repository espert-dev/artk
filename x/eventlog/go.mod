module artk.dev/x/eventlog

go 1.22.0

require (
	artk.dev v0.4.0
	artk.dev/x/testlog v0.2.0
)

require (
	github.com/lmittmann/tint v1.0.5 // indirect
	github.com/neilotoole/slogt v1.1.0 // indirect
)

replace artk.dev => ../../

replace artk.dev/x/testlog => ../testlog
