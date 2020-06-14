package cmd

import (
	"fmt"
	"os"

	"github.com/shric/trpc/internal/torrent"
)

type whichOptions struct {
	Missing bool `long:"missing" description:"Show only unassociated files/paths"`
}

// Which implements the which command (find which torrent a file/path is associated with.
func Which(c *Command) {
	_, ok := c.Options.(whichOptions)
	optionsCheck(ok)

	finder := torrent.NewFinder(c.Client)

	for _, f := range c.PositionalArgs {
		torrent, fileID := finder.Find(f)
		if torrent != nil {
			fmt.Printf("%s belongs to torrent %d: %s (File ID %d)\n",
				f, *torrent.ID, *torrent.Name, fileID)
		} else {
			fmt.Fprintln(os.Stderr, "Couldn't find a torrent for", f)
		}
	}
}