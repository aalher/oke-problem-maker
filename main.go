/*
Copyright 2019 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"k8s.io/klog/v2"

	"github.com/aalher/oke-problem-maker/makers"
)

type options struct {
	// Command line options. See flag descriptions for the description
	Rate     float32
	Duration time.Duration
	Problem  string
	Random   bool
}

// AddFlags adds log counter command line options to pflag.
func (o *options) AddFlags(fs *pflag.FlagSet) {
	fs.Float32Var(&o.Rate, "rate", 1.0,
		"Number of times the problem should be generated per second")
	fs.DurationVar(&o.Duration, "duration", time.Duration(1)*time.Second,
		"Duration for problem maker to keep generating problems")

	problems := makers.GetProblemTypes()
	fs.StringVar(&o.Problem, "problem", "",
		fmt.Sprintf("The type of problem to be generated. Supported types: %q",
			strings.Join(problems, ", ")))

	fs.BoolVar(&o.Random, "random", false, "Make problems at random selected duration/problems at a rate of 1.0.")
}

func main() {
	klogFlags := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(klogFlags)
	klogFlags.VisitAll(func(f *flag.Flag) {
		switch f.Name {
		case "v", "vmodule", "logtostderr":
			flag.CommandLine.Var(f.Value, f.Name, f.Usage)
		}
	})
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.CommandLine.MarkHidden("vmodule")
	pflag.CommandLine.MarkHidden("logtostderr")

	o := options{}
	o.AddFlags(pflag.CommandLine)
	pflag.Parse()

	if o.Problem == "" && !o.Random {
		klog.Fatalf("Please specify the type of problem to make using the --problem argument or choose random chaos.")
	}

	if !o.Random {
		generate(o)
	} else {
		rand.Seed(time.Now().UnixNano())
		// Number of cases between 5-15

		o.Rate = 1.0
		o.Duration = time.Duration(rand.Intn(5)+1)
		cases := rand.Intn(6) + 10
		for i := 1; i <= cases; i++ {
			o.Problem = getRandomProblem(makers.ProblemGenerators)

			go generate(o)
		}
	}
}

func generate(o options) {
	problemGenerator, ok := makers.ProblemGenerators[o.Problem]
	if !ok {
		klog.Fatalf("Expected to see a problem type of one of %q, but got %q.",
			makers.GetProblemTypes(), o.Problem)
	}

	periodMilli := int(1000.0 / o.Rate)
	ticker := time.NewTicker(time.Duration(periodMilli) * time.Millisecond)
	defer ticker.Stop()

	done := make(chan bool)
	go func() {
		time.Sleep(o.Duration)
		done <- true
	}()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			klog.Infof("Generating problem: %q", o.Problem)
			problemGenerator()
		}
	}
}

func getRandomProblem(m map[string]func()) string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}

	rand.Seed(time.Now().UnixNano()) // Seed random generator
	randomIndex := rand.Intn(len(keys))

	return keys[randomIndex]
}
