package main

import (
	"fmt"

	"github.com/evrenios/letmein/misc"
	"github.com/kelseyhightower/envconfig"
)

var conf = generate()

type config struct {
	Banlist       []string         `envconfig:"BANLIST"`
	SGPorts       map[string]int64 `envconfig:"SG_PORTS"`
	AWSRegion     string           `envconfig:"AWS_REGION"`
	errChan       chan error
	restrictedIPs map[string]struct{}
	currentSet    map[int64]*misc.AuthReq
	slacker       *slackNotifier
}

func generate() *config {
	c := &config{
		currentSet:    make(map[int64]*misc.AuthReq),
		slacker:       newSlackNotifier(),
		errChan:       make(chan error, 0),
		restrictedIPs: make(map[string]struct{}),
	}
	c.loadENV()
	go c.readErrors()
	return c
}

func (c *config) loadENV() {
	if err := envconfig.Process("", c); err != nil {
		panic(err)
	}
	for _, ip := range c.Banlist {
		c.restrictedIPs[ip] = struct{}{}
	}

	fmt.Println("Restricted IPs -> ", c.Banlist)
	fmt.Println("Security Groupts with Ports to add -> ", c.SGPorts)

	if len(c.SGPorts) == 0 {
		panic("there is no given security group for server to add, exiting...")
	}
}

func (c *config) readErrors() {
	select {
	case err := <-c.errChan:
		fmt.Println("internal err: ", err.Error())
		c.slacker.notify(err.Error())
	}
}
