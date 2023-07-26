package main

func main() {
	bc := NewBlockchain()

	cli := NewCLI(bc)

	if err := cli.Run(); err != nil {
		panic(err)
	}
}
