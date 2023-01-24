/*
 * @Author: xiaozuhui
 * @Date: 2023-01-18 16:22:45
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2023-01-18 16:46:51
 * @Description:
 */
package main

import (
	"fmt"
	"lipsoft/repl"
	"os"
	"os/user"
)

func main() {
	current, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Lipsoft programming language!\n", current.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}
