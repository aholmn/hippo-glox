# hippo-glox
A stripped down version of the Lox programming language, see https://github.com/munificent/craftinginterpreters,  written in Go. 
This project was started for learning Go and the code is probaly the least idiomatic Go code out there. 

## Difference from original Lox
* no classes
* block statement is required after if, for, while statement
* uses function instead of fun for declaring functions
* uses let instead of var for declaring variables

## Run files
clone project in your go directory then run following
```sh
go run hippo-glox examples/fibonacci
```
