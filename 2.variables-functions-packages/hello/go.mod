module hello

go 1.21.3

require (
	github.com/zenkiebear/learn-go/calculator v0.1.2
	rsc.io/quote v1.5.2
)

require (
	golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c // indirect
	rsc.io/sampler v1.3.0 // indirect
)

replace zenkie.cn/calculator v0.1.2 => ../calculator
