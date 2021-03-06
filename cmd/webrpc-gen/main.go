package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/webrpc/webrpc/gen"
	_ "github.com/webrpc/webrpc/gen/golang"
	_ "github.com/webrpc/webrpc/gen/javascript"
	_ "github.com/webrpc/webrpc/gen/typescript"
	"github.com/webrpc/webrpc/schema"
)

var flags = flag.NewFlagSet("webrpc-gen", flag.ExitOnError)

func main() {
	schemaFlag := flags.String("schema", "", "webrpc schema file (required)")
	pkgFlag := flags.String("pkg", "proto", "generated package name for target language, default: proto")
	outFlag := flags.String("out", "", "generated output file, default: stdout")
	clientFlag := flags.Bool("client", false, "enable webrpc client library generation, default: off")
	serverFlag := flags.Bool("server", false, "enable webrpc server library generation, default: off")
	websocketFlag := flags.Bool("websocket", false, "enable websocket transport generation, default: off")

	// registered/available target languages
	targets := []string{}
	for k, _ := range gen.Generators {
		targets = append(targets, k)
	}
	targetFlag := flags.String("target", "", fmt.Sprintf("target language for webrpc library generation, %s (required)", targets))
	targetExtra := flags.String("extra", "", "target language extra/custom options")

	flags.Parse(os.Args[1:])

	if *schemaFlag == "" {
		fmt.Println("oops, you must pass a -schema flag, see -h for help/usage")
		os.Exit(1)
	}
	if *targetFlag == "" {
		fmt.Println("oops, you must pass a -target flag, see -h for help/usage")
		os.Exit(1)
	}

	targetLang := *targetFlag
	if _, ok := gen.Generators[targetLang]; !ok {
		fmt.Printf("oops, you passed an invalid -target flag, try one of registered generators: %s\n", targets)
		os.Exit(1)
	}

	// Parse+validate the webrpc schema file
	schema, err := schema.ParseSchemaFile(*schemaFlag)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Call our target code-generator
	generator := gen.GetGenerator(*targetFlag)
	if generator == nil {
		fmt.Printf("error! unable to find generator for target '%s'\n", *targetFlag)
		os.Exit(1)
	}

	targetOpts := gen.TargetOptions{
		PkgName:   *pkgFlag,
		Client:    *clientFlag,
		Server:    *serverFlag,
		Extra:     *targetExtra,
		Websocket: *websocketFlag,
	}

	protoGen, err := generator.Gen(schema, targetOpts)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Write output
	if *outFlag != "" && *outFlag != "stdout" {
		outfile := *outFlag
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if outfile[0:1] != "/" {
			outfile = filepath.Join(cwd, outfile)
		}

		outdir := filepath.Dir(outfile)
		if _, err := os.Stat(outdir); os.IsNotExist(err) {
			err := os.MkdirAll(outdir, 0755)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}

		err = ioutil.WriteFile(outfile, []byte(protoGen), 0644)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	} else {
		fmt.Println(protoGen)
	}
}
