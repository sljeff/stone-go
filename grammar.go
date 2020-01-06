package stone

type Stone struct {
	Sections []*Section `( @@ )*`
}

type Section struct {
	Name       string  `"[" @Ident "]"`
	ParentName *string `("<" "[" @Ident "]")?`
	Stmts      []*Stmt `@@*`
}

type Stmt struct {
	Assignment *Assignment `@@`
	Deletion   *LeftValue  `| "DELETE" @@`
}

type Assignment struct {
	Left  LeftValue  `@@ "="`
	Right RightValue `@@`
}

type LeftValue struct {
	Identifier string  `@Ident`
	KeyString  *string `("[" (@String `
	KeyInt     *int    `| @Int) "]")?`
}

type RightValue struct {
	String     *string     `@String`
	Int        *int        `| @Int `
	Float      *float64    `| @Float`
	Bool       *Boolean    `| @("TRUE"|"FALSE")`
	LeftValue  *LeftValue  `| @@`
	ArrayValue *ArrayValue `| @@`
	MapValue   *MapValue   `| @@`
	EnvValue   *EnvValue   `| @@`
}

type ArrayValue struct {
	Array []*RightValue `("[" "]") | ("[" @@ ("," @@)* ","? "]")`
}

type MapValue struct {
	KVs []*KV `("{" "}") | ("{" @@ ("," @@)* ","? "}")`
}

type KV struct {
	K string     `@String ":"`
	V RightValue `@@`
}

type EnvValue struct {
	Env string `"$" "{" @Ident "}"`
}

type Boolean bool

func (b *Boolean) Capture(values []string) error {
	*b = values[0] == "TRUE"
	return nil
}
