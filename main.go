/*
people who move bricks should not yield
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

var (
	parallels  = flag.Bool("parallels", false, "parallels make projects")
	configPath = flag.String("config", "", "json config used by goporter")
)

func main() {
	flag.Parse()

	if len(*configPath) == 0 {
		fmt.Fprintf(os.Stderr, "'--config' can not be empty")
		os.Exit(1)
	}

	config, err := os.Open(*configPath)
	if err != nil {
		panic(err)
	}

	defer config.Close()

	configBytes, err := ioutil.ReadAll(config)
	if err != nil {
		panic(err)
	}

	bricks := make([]brick, 0)

	err = json.Unmarshal(configBytes, &bricks)
	if err != nil {
		panic(err)
	}

	// todo(weilaaa): output to file
	if *parallels {
		parallelsExecuteAll(bricks)
	} else {
		serialExecuteAll(bricks)
	}
}

func parallelsExecuteAll(bricks []brick) {
	w := &sync.WaitGroup{}
	w.Add(len(bricks))

	for _, b := range bricks {
		go func(brick brick) {
			if err := brick.moving(); err != nil {
				colorPrint(red, "brick moving failed, brick: %v, err: %v", brick, err)
			}
			w.Done()
		}(b)
	}

	w.Wait()
}

func serialExecuteAll(bricks []brick) {
	for _, b := range bricks {
		if err := b.moving(); err != nil {
			colorPrint(red, "brick moving failed, brick: %v, err: %v", b, err)
			continue
		}
	}
}
