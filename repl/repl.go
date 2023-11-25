package repl

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/raiguard/luapls/lua/lexer"
	"github.com/raiguard/luapls/lua/parser"
	"github.com/raiguard/luapls/lua/token"

	"github.com/chzyer/readline"
)

func Run() {
	rl, err := readline.New("(luapls) ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}

		l := lexer.New(line)

		fmt.Println("TOKENS:")
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Println(tok.String())
		}

		fmt.Println("AST:")

		p := parser.New(line)
		block := p.ParseBlock()
		for _, err := range p.Errors() {
			fmt.Fprintln(os.Stderr, err)
		}
		if len(p.Errors()) > 0 {
			continue
		}
		bytes, _ := json.MarshalIndent(&block, "", "  ")
		fmt.Println(string(bytes))

		// fmt.Println("NODES:")
		// ast.Walk(&block, func(n ast.Node) bool {
		// 	fmt.Printf("%T: {%d, %d}\n", n, n.Pos(), n.End())
		// 	return true
		// })
	}
}
