module github.com/imerkle/rosetta-solana-go

require (
	contrib.go.opencensus.io/exporter/stackdriver v0.13.5 // indirect
	github.com/coinbase/rosetta-sdk-go v0.6.10
	github.com/dfuse-io/binary v0.0.0-20210216024852-4ae6830a495d
	github.com/dfuse-io/logging v0.0.0-20210109005628-b97a57253f70 // indirect
	github.com/dfuse-io/solana-go v0.2.0
	github.com/fatih/color v1.10.0
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/iancoleman/strcase v0.1.3
	github.com/mitchellh/copystructure v1.2.0
	github.com/mr-tron/base58 v1.2.0
	github.com/portto/solana-go-sdk v1.1.1
	github.com/spf13/cobra v1.1.1
	github.com/teris-io/shortid v0.0.0-20201117134242-e59966efd125 // indirect
	github.com/test-go/testify v1.1.4
	github.com/tidwall/gjson v1.7.5 // indirect
	go.opencensus.io v0.23.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b // indirect
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	golang.org/x/sys v0.0.0-20210426230700-d19ff857e887 // indirect
	golang.org/x/term v0.0.0-20210422114643-f5beecf764ed // indirect
	gotest.tools v2.2.0+incompatible
)

go 1.15

replace github.com/portto/solana-go-sdk => github.com/imerkle/solana-go-sdk v0.0.8

//replace github.com/portto/solana-go-sdk => ../solana-go-sdk
