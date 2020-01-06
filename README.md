# stone-go
Stone configuration language Golang implement.

# Example

```bash
cd example
go run main.go
# modify config.stone or add env, then retry it
```

# Grammar

## Reserved words

```
DELETE
TRUE
FALSE
```

## Basic Types

```
String
Int
Float
Bool  (TRUE FALSE)
Array (of any type)
Map   (key: String, value: any)
```

## Section

declare section:

```
[Identifier]
```

declare a section inherit another section:

```
[Identifier] < [Identifier]
```

## Stmt

Stmt include Deletion and Assignment.

## Deletion

```
DELETE LeftValue
```

## Assignment

```
LeftValue = RightValue
```

## LeftValue

LeftValue can be:

```
Identifier
Identifier[String] (item of Map)
Identifier[Int]    (item of Array)
```

## RightValue

RightValue can be:

```
Literal Basic Types
LeftValue
EnvValue
```

## EnvValue

```
${Identifier}
```
