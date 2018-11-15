package main

import (
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/docker/docker/client"
	"github.com/euforia/cchain/pkg/dockerimage"

	"gopkg.in/urfave/cli.v2"
)

func newCLIApp() *cli.App {
	return &cli.App{
		Commands: []*cli.Command{
			&cli.Command{
				Name: "image",
				Subcommands: []*cli.Command{
					cmdImageConf(),
				},
			},
		},
	}
}

func showImageConfSummary(img *dockerimage.DockerImage) {
	w := tabwriter.NewWriter(os.Stdout, 1, 2, 2, ' ', tabwriter.StripEscape)

	fmt.Fprintf(w, "Name\t%s\n", img.Name)
	fmt.Fprintf(w, "CWD\t%s\n", img.WorkingDir())
	w.Flush()
	fmt.Println()

	fmt.Println("Environment Variables")
	fmt.Println("---------------------")
	envVars := img.EnvVars()
	for k, v := range envVars {
		fmt.Fprintf(w, "%s\t%s\n", k, v)
	}
	w.Flush()
	fmt.Println()

	fmt.Println("Ports")
	fmt.Println("-----")
	ports := img.Ports()
	for _, v := range ports {
		fmt.Println(v)
	}
	fmt.Println()

	fmt.Println("Volumes")
	fmt.Println("-------")
	volumes := img.Volumes()
	for _, v := range volumes {
		fmt.Println(v)
	}
	fmt.Println()
}

func cmdImageConf() *cli.Command {
	return &cli.Command{
		Name: "conf",
		Action: func(ctx *cli.Context) error {
			client, err := client.NewEnvClient()
			if err != nil {
				return err
			}
			imgName := ctx.Args().First()
			if imgName == "" {
				return errors.New("image name required")
			}
			client.
			img, err := dockerimage.NewDockerImage(client, imgName)
			if err != nil {
				return err
			}

			showImageConfSummary(img)

			return nil
		},
	}
}

func main() {
	app := newCLIApp()
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
