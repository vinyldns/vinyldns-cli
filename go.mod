module github.com/ssranjani06/vinyldns-cli

go 1.19

require (
	github.com/crackcomm/go-clitable v0.0.0-20151121230230-53bcff2fea36
	github.com/olekukonko/tablewriter v0.0.4
	github.com/urfave/cli v1.22.4
	github.com/vinyldns/go-vinyldns v0.9.16
)

replace github.com/vinyldns/go-vinyldns => github.com/ssranjani06/go-vinyldns v0.0.0-20240611144018-a38e17c2e3c3

require (
	github.com/aws/aws-sdk-go-v2 v1.26.1 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.11 // indirect
	github.com/aws/smithy-go v1.20.2 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.0-20190314233015-f79a8a8ca69d // indirect
	github.com/mattn/go-runewidth v0.0.7 // indirect
	github.com/russross/blackfriday/v2 v2.0.1 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
)
