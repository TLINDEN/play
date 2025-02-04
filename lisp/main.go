package main

/*
 - Provide a list per each hook to lisp env
 - the user defines functions and vars etc
 - the user adds his functions to the lists
 - we then iterate each list at hook execution point and eval them
*/

import (
	"errors"
	"fmt"
	"log"

	"github.com/alecthomas/repr"
	"github.com/glycerine/zygomys/zygo"
)

var Hooks map[string][]*zygo.SexpSymbol

func AddHook(env *zygo.Zlisp, name string, args []zygo.Sexp) (zygo.Sexp, error) {
	var hookname string

	if len(args) < 2 {
		return zygo.SexpNull, errors.New("argument of %add-hook should be: %hook-name %your-function")
	}

	switch t := args[0].(type) {
	case *zygo.SexpSymbol:
		hookname = t.Name()
	default:
		return zygo.SexpNull, errors.New("hook name must be a symbol!")
	}

	switch t := args[1].(type) {
	case *zygo.SexpSymbol:
		Hooks[hookname] = append(Hooks[hookname], t)
	default:
		return zygo.SexpNull, errors.New("hook function must be a symbol!")
	}

	return zygo.SexpNull, nil
}

func Runhook(env *zygo.Zlisp, hook *zygo.SexpSymbol, arg string) (bool, error) {
	env.Clear()
	var result bool

	// either use AddGlobal to put stuff into lisp directly or create
	// a lisp arg to the function call using printf, eg:
	//res, err := env.EvalString(fmt.Sprintf("(%s `%s`)", hook.Name(), arg))
	//env.AddGlobal("data", &zygo.SexpStr{S: arg})
	list := zygo.MakeList([]zygo.Sexp{&zygo.SexpInt{Val: 1}, &zygo.SexpInt{Val: 2}, &zygo.SexpInt{Val: 3}})
	env.AddGlobal("data", list)
	res, err := env.EvalString(fmt.Sprintf("(%s data)", hook.Name()))
	if err != nil {
		return false, err
	}

	switch t := res.(type) {
	case *zygo.SexpBool:
		result = t.Val
	default:
		return false, errors.New("filterhook shall return BOOL!")
	}

	return result, nil
}

type Table struct {
	Headers []string   `json:"headers" msg:"headers"`
	Rows    [][]string `json:"rows" msg:"rows"`
}

func main() {
	/*
		scanner := bufio.NewScanner(os.Stdin)
		Hooks = map[string][]*zygo.SexpSymbol{}

		env := zygo.NewZlisp()
		env.AddFunction("addhook", AddHook)

		var code []string
		for scanner.Scan() {
			code = append(code, scanner.Text())
		}
	*/

	env := zygo.NewZlisp()
	env.StandardSetup()

	zygo.GoStructRegistry.RegisterUserdef(
		&zygo.RegisteredType{GenDefMap: true, Factory: func(env *zygo.Zlisp, h *zygo.SexpHash) (interface{}, error) {
			return &Table{}, nil
		}}, true, "table")

	code := `
        // A defmap is needed to define the table struct inside env.
        // The registry doesn't know about env(s), so it 
        // can't do it for us automatically.
        (defmap table)

        // Create an instance of table, with some data in it.
        (def t 
           (table headers:  ["wood"  "metal"] 
                  rows:    [["oak"  "silver"]
                            ["pine" "tin"   ]]))`

	x, err := env.EvalString(code)
	if err != nil {
		log.Fatalf(env.GetStackTrace(err))
	}

	tmp, err := zygo.SexpToGoStructs(x, &Table{}, env, nil)
	if err != nil {
		log.Fatalf(err.Error())
	}

	switch f := tmp.(type) {
	case *Table:
		repr.Println(f)
	default:
		log.Fatalf("wrong type!")
	}

	/*
		err := env.LoadString(strings.TrimSpace(strings.Join(code, "\n")))
		if err != nil {
			log.Fatalf(env.GetStackTrace(err))
		}

		expr, err := env.Run()
		if err != nil {
			log.Fatalf(env.GetStackTrace(err))
		}

		env.DumpEnvironment()
		fmt.Println(expr)
		fmt.Println(Hooks)

		for _, hook := range Hooks["filterhook"] {
			result, err := Runhook(env, hook, "datentest")
			if err != nil {
				log.Fatalf("Failed: %s", err)
			}

			fmt.Printf("%s result: %t\n", hook.Name(), result)
		}
	*/
}
