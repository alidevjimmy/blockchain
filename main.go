package main

func main() {
	bc := NewBlockchain("default_address", "BTC")

	cli := NewCLI(bc)

	if err := cli.Run(); err != nil {
		panic(err)
	}
}
