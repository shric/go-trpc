package cmd

import (
	"fmt"
	"os"

	"github.com/hekmon/transmissionrpc"
	"github.com/shric/trpc/internal/filter"
)

// MoveOptions is all the command line options for the verify command
type MoveOptions struct {
	filter.Options `group:"filters"`
	ForceAll       bool `long:"force-all" description:"Really move all torrents"`
}

func getFnamesAndDest(args []string) (fnames []string, dest string) {
	dest = args[len(args)-1]
	fnames = args[:len(args)-1]
	return
}

// Move implements the verify command (hash check torrents)
func Move(c *Command) {
	opts, ok := c.CommandOptions.(MoveOptions)
	optionsCheck(ok)
	if len(c.PositionalArgs) == 0 {
		fmt.Fprintln(os.Stderr, "move: Destination required")
		return
	}
	if len(c.PositionalArgs) == 1 && !opts.ForceAll {
		fmt.Fprintln(os.Stderr, "Use --force-all if you really want to move all torrents")
		return
	}
	fnames, destination := getFnamesAndDest(c.PositionalArgs)
	ProcessTorrents(c.Client, opts.Options, fnames, []string{"name", "id"}, func(torrent *transmissionrpc.Torrent) {
		if !c.CommonOptions.DryRun {
			err := c.Client.TorrentSetLocation(*torrent.ID, destination, true)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		}
		c.status("Moving torrent", torrent)
	})
}
