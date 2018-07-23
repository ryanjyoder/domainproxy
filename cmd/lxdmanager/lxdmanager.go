package main

import (
	"fmt"
	"log"

	lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
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
	op, err := c.CopyImage(d, *image, &lxd.ImageCopyArgs{
		CopyAliases: true,
	})
	if err != nil {
		return err
	}

	// And wait for it to finish
	err = op.Wait()
	if err != nil {
		return err
	}

	// Container creation request
	req := api.ContainersPost{
		Name: "my-container",
		Source: api.ContainerSource{
			Type:  "image",
			Alias: "centls/7",
		},
	}

	// Get LXD to create the container (background operation)
	opCreate, err := c.CreateContainer(req)
	if err != nil {
		return err
	}

	// Wait for the operation to complete
	err = opCreate.Wait()
	if err != nil {
		return err
	}

	// Get LXD to start the container (background operation)
	reqState := api.ContainerStatePut{
		Action:  "start",
		Timeout: -1,
	}

	opStart, err := c.UpdateContainerState("my-container", reqState, "")
	if err != nil {
		return err
	}

	// Wait for the operation to complete
	err = opStart.Wait()
	if err != nil {
		return err
	}

	return nil

}
