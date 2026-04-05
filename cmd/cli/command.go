package main

type CommandFunc = func(*Command, []string)

type Command struct {
	Name  string
	Descr string
	Run   CommandFunc
}
