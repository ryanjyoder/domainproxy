package main

import (
	"fmt"
	"log"

	lxd "github.com/lxc/lxd/client"
)

func main() {
	fmt.Println("launchContainer:", launchContainer())
}

func launchContainer() error {
	// Connect to LXD over the Unix socket
	c, err := lxd.ConnectLXDUnix("/var/snap/lxd/common/lxd/unix.socket", nil)
	if err != nil {
		log.Println("couldnt connect to unix socket", err)
		return err
	}
	// Connect to the remote SimpleStreams server
	d, err := lxd.ConnectSimpleStreams("https://images.linuxcontainers.org", nil)
	if err != nil {
		return err
	}

	// Resolve the alias
	alias, _, err := d.GetImageAlias("centos/7")
	if err != nil {
		return err
	}

	// Get the image information
	image, _, err := d.GetImage(alias.Target)
	if err != nil {
		return err
	}

	// Ask LXD to copy the image from the remote server
	op, err := c.CopyImage(d, *image, nil)
	if err != nil {
		return err
	}

	// And wait for it to finish
	err = op.Wait()
	if err != nil {
		return err
	}

	return nil

}
