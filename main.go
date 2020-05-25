package main

import (
	"flag"
	"fmt"
	"os"
)

type DefaultFlags struct {
	rabbitHost *string
	rabbitPort *int
	rabbitUser *string
	rabbitPass *string
	queue      *string
	durable    *bool
	verbose    *bool
}

func main() {

	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	message := sendCmd.String("message", "", "The message to send. e.g. message=\"Hello World!\"")
	sendCount := sendCmd.Int("count", 1, "The number of messages to send. e.g. count=100")

	drainCmd := flag.NewFlagSet("drain", flag.ExitOnError)
	drainConcurrency := drainCmd.Int("concurrency", 1, "The number of messages to drain at a time. e.g. concurrency=1")
	drainDwell := drainCmd.Int("dwell", 1, "The time in ms to dwell after processing each message. e.g. dwell=1")

	receiveCmd := flag.NewFlagSet("receive", flag.ExitOnError)
	receiveCount := receiveCmd.Int("count", 1, "The number of messages to receive. e.g. receive=1")

	if len(os.Args) < 2 || os.Args[1] == "-h" {
		PrintDefaults()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "send":
		defaultFlags := addDefaultFlags(sendCmd)
		sendCmd.Parse(os.Args[2:])
		rabbitConn := makeConnString(defaultFlags)
		send(rabbitConn, defaultFlags.queue, *message, *sendCount, defaultFlags.durable, defaultFlags.verbose)
	case "drain":
		defaultFlags := addDefaultFlags(drainCmd)
		drainCmd.Parse(os.Args[2:])
		rabbitConn := makeConnString(defaultFlags)
		drain(rabbitConn, defaultFlags.queue, drainConcurrency, drainDwell, defaultFlags.durable, defaultFlags.verbose)
	case "receive":
		defaultFlags := addDefaultFlags(receiveCmd)
		receiveCmd.Parse(os.Args[2:])
		rabbitConn := makeConnString(defaultFlags)
		receive(rabbitConn, defaultFlags.queue, receiveCount, defaultFlags.durable, defaultFlags.verbose)
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

}

func PrintDefaults() {
	fmt.Println("send, drain or receive command expected.")
	fmt.Println("usage: ")
	fmt.Println("  rabbitmq-drainer <send|receive|drain> -h")
	fmt.Println("")
	fmt.Println("examples: ")
	fmt.Println("  rabbitmq-drainer send -queue=hello -message=\"Hello World\" -count 10000")
	fmt.Println("  rabbitmq-drainer receive -queue=hello -count=1")
	fmt.Println("  rabbitmq-drainer drain -queue=hello -concurrency 10 -dwell 1")

}

func addDefaultFlags(flagSet *flag.FlagSet) DefaultFlags {
	rabbitHost := flagSet.String("host", "localhost", "The rabbit host. e.g. host=localhost")
	rabbitPort := flagSet.Int("port", 5672, "The rabbit port. e.g. port=5672")
	rabbitUser := flagSet.String("user", "guest", "The username. e.g. user=guest")
	rabbitPass := flagSet.String("pass", "guest", "The password. e.g. pass=guest")
	queue := flagSet.String("queue", "", "The queue name. e.g. queue=hello")
	durable := flagSet.Bool("durable", false, "Whether the queue is durable. e.g. -durable")
	verbose := flagSet.Bool("verbose", false, "Whether the output is verbose. e.g. -verbose")

	return DefaultFlags{
		rabbitHost,
		rabbitPort,
		rabbitUser,
		rabbitPass,
		queue,
		durable,
		verbose,
	}
}

func makeConnString(defaultFlags DefaultFlags) string {
	rabbitConn := fmt.Sprintf("amqp://%s:%s@%s:%d/", *defaultFlags.rabbitUser, *defaultFlags.rabbitPass, *defaultFlags.rabbitHost, *defaultFlags.rabbitPort)
	return rabbitConn
}
