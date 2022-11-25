package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/ipfs/go-cid"
	format "github.com/ipfs/go-ipld-format"
	iface "github.com/ipfs/interface-go-ipfs-core"
	cli "github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

var (
	cleanup      []func() error
	cleanupMutex sync.Mutex
)

func main() {
	// Do any cleanup on exit
	defer doCleanup()

	app := cli.NewApp()
	app.Name = "ipverify"
	app.Usage = "Verify IPFS objects."
	app.Version = "0.9.1"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "node",
			Aliases: []string{"n"},
			Usage:   "specify ipfs node strategy (\"local\", \"spawn\", \"temp\" or \"fallback\")",
			Value:   "fallback",
		},
		&cli.StringSliceFlag{
			Name:    "peers",
			Aliases: []string{"p"},
			Usage:   "specify a set of IPFS peers to connect to",
		},
		&cli.IntFlag{
			Name:    "offset",
			Aliases: []string{"o"},
			Usage:   "specify which line to start on in input file; note: 1-indexed",
			Value:   1,
		},
		&cli.BoolFlag{
			Name:    "show-stat",
			Aliases: []string{"ss"},
			Usage:   "show the node stat output",
			Value:   false,
		},
		&cli.IntFlag{
			Name:    "goroutines",
			Aliases: []string{"gs"},
			Usage:   "the number of goroutines used to verify cids",
			Value:   5,
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigExitCoder := make(chan cli.ExitCoder, 1)

	app.Action = func(c *cli.Context) error {
		if !c.Args().Present() {
			return fmt.Errorf("usage: ipverify <newline delimited file of cids>")
		}

		var ipfs iface.CoreAPI
		var err error
		switch c.String("node") {
		case "fallback":
			ipfs, err = http(ctx)
			if err == nil {
				break
			}
			fallthrough
		case "spawn":
			ipfs, err = spawn(ctx)
		case "local":
			ipfs, err = http(ctx)
		case "temp":
			ipfs, err = temp(ctx)
		default:
			return fmt.Errorf("no such 'node' strategy, %q", c.String("node"))
		}
		if err != nil {
			return err
		}

		connect(ctx, ipfs, c.StringSlice("peers"))

		cidFile := c.Args().First()
		data, err := os.ReadFile(cidFile)
		if err != nil {
			log.Fatal(err)
			return err
		}

		str := string(data)
		cidsstr := strings.Split(str, "\n")
		cidsstr = cidsstr[:len(cidsstr)-1]

		f, err := os.OpenFile("failed.cids", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
			return err
		}
		defer f.Close()

		offset := c.Int("offset") - 1
		if offset < 0 {
			offset = 0
		}

		g, ctx := errgroup.WithContext(ctx)
		g.SetLimit(c.Int("goroutines"))

		for i := offset; i < len(cidsstr); i++ {
			if (i+1)%10 == 0 {
				percent := float64((i - offset)) / float64(len(cidsstr)-offset-1) * float64(100)
				fmt.Printf("%d/%d\t%f%%\t%s\n", i-offset, len(cidsstr)-offset-1, percent, cidsstr[i])
			}
			cs := cidsstr[i]

			g.Go(func() error {
				cs := cs
				stat, err := getNodeStat(ctx, ipfs, cs)
				if err != nil {
					if err == context.DeadlineExceeded {
						if _, err := f.WriteString(fmt.Sprintf("%s\n", cs)); err != nil {
							fmt.Println("failed to write failed cid to file")
							fmt.Println(cs)
						}
						return nil
					}
					return err
				}
				if c.Bool("show-stat") {
					fmt.Println(stat.String())
				}
				return nil
			})
		}
		if err := g.Wait(); err != nil {
			return fmt.Errorf("errgroup wait failed: %w", err)
		}
		// 		ipr, err := ipfs.Block().Get(ctx, iPath)
		// 		if err != nil {
		// 			if err == context.Canceled {
		// 				return <-sigExitCoder
		// 			}
		// 			return cli.Exit(err, 2)
		// 		}
		// 		fmt.Println("block request made")
		// 		lipr := io.LimitReader(ipr, 4)

		// 		blk := []byte{}
		// 		fmt.Println("attempting to read from block")
		// 		n, err := lipr.Read(blk)
		// 		if err != nil {
		// 			if err == context.Canceled {
		// 				return <-sigExitCoder
		// 			}
		// 			return cli.Exit(err, 2)
		// 		}
		// 		fmt.Println("value of n is: ", n)
		// 		if n == 4 {
		// 			return cli.Exit(string(blk), 0)
		// 		}
		// 		if n != 4 {
		// 			if err == context.Canceled {
		// 				return <-sigExitCoder
		// 			}
		// 			return cli.Exit(err, 2)
		// 		}
		return nil
	}

	// Catch interrupt signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		sigExitCoder <- cli.Exit("", 128+int(sig.(syscall.Signal)))
		cancel()
	}()

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		doCleanup()
		os.Exit(1)
	}
}

func getNodeStat(ctx context.Context, ipfs iface.CoreAPI, cs string) (*format.NodeStat, error) {
	//TODO: make the timeout value a cli flag
	ctx, cancel := context.WithTimeout(ctx, 40*time.Second)
	defer cancel()

	cidarg, err := cid.Parse(cs)
	if err != nil {
		return nil, err
	}

	nod, err := ipfs.Dag().Get(ctx, cidarg)
	if err != nil {
		return nil, err
	}

	stat, err := nod.Stat()
	if err != nil {
		return nil, err
	}

	return stat, nil
}

func addCleanup(f func() error) {
	cleanupMutex.Lock()
	defer cleanupMutex.Unlock()
	cleanup = append(cleanup, f)
}

func doCleanup() {
	cleanupMutex.Lock()
	defer cleanupMutex.Unlock()

	for _, f := range cleanup {
		if err := f(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
