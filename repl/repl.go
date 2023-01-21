/*
 * @Author: xiaozuhui
 * @Date: 2023-01-18 16:19:01
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2023-01-18 16:47:04
 * @Description:
 */
package repl

import (
	"bufio"
	"fmt"
	"io"
	"lipsoft/lexer"
	"lipsoft/token"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
