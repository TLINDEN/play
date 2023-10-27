package main

/*
 - Provide a list per each hook to lisp env
 - the user defines functions and vars etc
 - the user adds his functions to the lists
 - we then iterate each list at hook execution point and eval them
*/

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/glycerine/zygomys/zygo"
)

var Hooks map[string][]*zygo.SexpSymbol

func AddHook(env *zygo.Zlisp, name string, args []zygo.Sexp) (zygo.Sexp, error) {
	var hookname string

	if len(args) < 2 {
		return zygo.SexpNull, errors.New("argument of ^add-hook should be: ^hook-name ^your-function")
	}

	for _, number := range args {
		switch t := number.(type) {
		case *zygo.SexpStr:
			if hookname == "" {
				hookname = t.S
			} else {
				return zygo.SexpNull, errors.New("argument type of (addhook) should be string symbol")
			}
		case *zygo.SexpSymbol:
			if hookname == "" {
				return zygo.SexpNull, errors.New("argument type of (addhook) should be string symbol")
			} else {
				Hooks[hookname] = append(Hooks[hookname], t)
			}
		default:
			return zygo.SexpNull, errors.New("argument type of (addhook) should be string symbol")
		}
	}

	return zygo.SexpNull, nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	Hooks = map[string][]*zygo.SexpSymbol{}

	env := zygo.NewZlisp()
	env.AddFunction("addhook", AddHook)

	for scanner.Scan() {
		err := env.LoadString(strings.TrimSpace(scanner.Text()))
		if err != nil {
			fmt.Printf(env.GetStackTrace(err))
			env.Clear()
		}

		expr, err := env.Run()
		if err != nil {
			fmt.Printf(env.GetStackTrace(err))
			env.Clear()
		}

		env.DumpEnvironment()
		fmt.Println(expr)
		fmt.Println(Hooks)
	}

}
