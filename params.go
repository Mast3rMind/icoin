package main

type params struct {
	port string
}

var mainNetParams = params{
	port: "1986",
}

var testNetParams = params{
	port: "11986",
}
