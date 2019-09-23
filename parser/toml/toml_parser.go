package toml

import (
	"github.com/alecthomas/kong"
	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/lexer/ebnf"
	"github.com/alecthomas/repr"
	"os"
)

type TOML struct {
	Pos lexer.Position

	Entries []*Entry `{ @@ }`
}

type Entry struct {
	Field   *Field   ` @@`
	Section *Section `| @@`
}

type Field struct {
	Key   string `@Ident "="`
	Value *Value `@@`
}

type Value struct {
	String   *string  ` @String`
	DateTime *string  `| @DateTime`
	Date     *string  `| @Date`
	Time     *string  `| @Time`
	Bool     *bool    `| {@"true" | "false"}`
	Integer  *int64   `| @Int`
	Float    *float64 `| @Float`
	List     []*Value `| "[" [ @@ { "," @@ } ] "]"`
}

type Section struct {
	Name string `"[" @(Ident { "." Ident }) "]"`
	Fields []*Field `{ @@ }`
}

var (
	tomLexer = lexer.Must(ebnf.New(`
		Comment = "#" { "\u0000"…"\uffff"-"\n" } .
		DateTime = date "T" time [ "-" digit digit ":" digit digit ] .
		Date = date .
		Time = time .
		Ident = (alpha | "_") { "_" | alpha | digit } .
		String = "\"" { "\u0000"…"\uffff"-"\""-"\\" | "\\" any } "\"" .
		Int = [ "-" | "+" ] digit { digit } .
		Float = ("." | digit) {"." | digit }
		Punct = "!"…"/" | ":"…"@" | "["…` + "\"`\"" + ` | "{"…"~" .
		Whitespace = " " | "\t" | "\n" | "\r" .
		alpha = "a"…"z" | "A"…"Z" .
		digit = "0"…"9" .
		any = "\u0000"…"\uffff" .
		date = digit digit digit digit "-" digit digit "-" digit digit .
		time = digit digit ":" digit digit ":" digit digit [ "." { digit } ] .
	`))
	tomlParser = participle.MustBuild(&TOML{},
		participle.Lexer(tomLexer),
		participle.Unquote("String"),
		participle.Elide("Whitespace", "Comment"),
	)

	cli struct {
		File string `help:"TOML file to parse." arg:""`
	}
)

func main() {
	ctx := kong.Parse(&cli)
	toml := &TOML{}
	r, err := os.Open(cli.File)
	ctx.FatalIfErrorf(err)
	defer r.Close()
	err = tomlParser.Parse(r, toml)
	ctx.FatalIfErrorf(err)
	repr.Println(toml)
}