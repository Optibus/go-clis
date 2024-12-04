package main

import (
	"check-js-deps/sets"
	"check-js-deps/workspace"
	"context"
	"errors"
	"fmt"
	"os"

	"golang.org/x/sync/errgroup"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
)

func check(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%serror: %v\n%s", Red, err, Reset)
		os.Exit(1)
		// panic(e)
	}
}

func debugPrint(args ...interface{}) {
	_, isOk := os.LookupEnv("DEBUG_GO")
	if isOk {
		fmt.Println(args...)
	}
}

func main() {
	// notPackes := []string{"!**calendar**"}
	// path1 := "/Users/orengriffin/dev/armada/mithra/apps/calendar-server/package.json"
	// path2 := "/Users/orengriffin/dev/armada/mithra/apps/server/"
	// bool1 := workspace.ExcludePackage(notPackes, path1)
	// bool2 := workspace.ExcludePackage(notPackes, path2)
	// fmt.Println(bool1, bool2)
	//
	//
	links := sets.NewDoubleSet()

	root, rootExists := os.LookupEnv("ROOT")
	if !rootExists {
		check(errors.New("Missing env var ROOT"))
	}
	if len(os.Args) < 2 {
		check(errors.New("missing first parm, e.g. 'apps/mithra'"))
	}
	pathFromArg := os.Args[1]
	packagePath := root + "/mithra/" + pathFromArg
	// packagePath := "/Users/orengriffin/dev/armada/mithra/apps/community-engagement-fe"
	links.Add(packagePath)
	// var ops atomic.Uint64
	g, ctx := errgroup.WithContext(context.Background())
	// var g errgroup.Group
	// var wg sync.WaitGroup
	// isOk, err := workspace.CheckProject(packagePath, links, &wg)
	_, err := workspace.CheckProject(ctx, packagePath, links, g)
	if err != nil {
		check(err)
	}
	// wg.Wait()
	err = g.Wait()
	debugPrint(links.GetNoneChecked())
	check(err)
	// fmt.Println("ops:", ops.Load())
	// fmt.Println(isOk)
	// mithraPath := "/Users/orengriffin/dev/armada/mithra/"
	// pnpmWorkspacePath := mithraPath + "pnpm-workspace.yaml"
	// packages, err := workspace.Read(pnpmWorkspacePath)
	// if err != nil {
	// 	check(err)
	// }
	// fmt.Println(packages)
	fmt.Printf("%s%s Deps Installed (oren's tool says so)%s\n", Green, pathFromArg, Reset)
}
