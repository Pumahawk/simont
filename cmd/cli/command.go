package main

type CommandFunc = func(*Command, []string) int

type Command struct {
	Name  string
	Descr string
	Run   CommandFunc
}
