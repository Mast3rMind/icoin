package main

import (
	"github.com/zgreat/icoin/wire"
)

type params struct {
	port  string
	netID wire.NetID
}

var mainNetParams = params{
	port:  "1986",
	netID: wire.MainNetID,
}

var testNetParams = params{
	port:  "11986",
	netID: wire.TestNetID,
}
