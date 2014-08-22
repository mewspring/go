// lex is a tool which tokenizes the contents of the provided files.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/mewlang/go/lexer"
)

func main() {
	flag.Parse()
	for _, path := range flag.Args() {
		err := lex(path)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

// lex tokenizes the contents of the provided file.
func lex(path string) error {
	log.Println("Lexing:", path)
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	input := string(buf)

	tokens, err := lexer.Parse(input)
	if err != nil {
		log.Println(err)
	}
	for _, token := range tokens {
		fmt.Println("token:", token)
	}

	return nil
}
