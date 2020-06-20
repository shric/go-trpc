package cmd

import (
	"fmt"
	"os"

	"github.com/hekmon/transmissionrpc"
	"github.com/shric/trpc/internal/filter"
	"github.com/shric/trpc/internal/torrent"
)

type startOptions struct {
	filter.Options `group:"filters"`
	Now            bool `long:"now" description:"Start torrent now, bypassing the queue"`
}

// Start starts torrents.
func Start(c *Command) {
	opts, ok := c.Options.(startOptions)
	optionsCheck(ok)

	startFunc := c.Client.TorrentStartIDs

	if opts.Now {
		startFunc = c.Client.TorrentStartNowIDs
	}

	torrent.ProcessTorrents(c.Client, opts.Options, c.PositionalArgs, []string{"name", "id", "status"},
		func(torrent *transmissionrpc.Torrent) {
			if *torrent.Status != transmissionrpc.TorrentStatusStopped {
				return
			}
			if !c.CommonOptions.DryRun {
				err := startFunc([]int64{*torrent.ID})
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}
			c.status("Started torrent", torrent)
		})
}

type stopOptions struct {
	filter.Options `group:"filters"`
}

// Stop stops torrents.
func Stop(c *Command) {
	opts, ok := c.Options.(stopOptions)
	optionsCheck(ok)
	torrent.ProcessTorrents(c.Client, opts.Options, c.PositionalArgs, []string{"name", "id", "status"},
		func(torrent *transmissionrpc.Torrent) {
			if *torrent.Status == transmissionrpc.TorrentStatusStopped {
				return
			}
			if !c.CommonOptions.DryRun {
				err := c.Client.TorrentStopIDs([]int64{*torrent.ID})
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}
			c.status("Stopped torrent", torrent)
		})
}
