package main

import (
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
	"github.com/getfiit/empire/pkg/heroku"
)

var stream bool

var cmdDeploy = &Command{
	Run:             maybeMessage(runDeploy),
	Usage:           "deploy [<registry>]<image>:[<tag>] [-s]",
	OptionalApp:     true,
	OptionalMessage: true,
	Category:        "deploy",
	Short:           "deploy a docker image",
	Long: `
Deploy is used to deploy a docker image to an app.

Options:

    -s enable the status stream during the deployment. If this is enabled, the
    command will wait until the scheduler has finished deploying the new
    release.

Examples:

    $ emp deploy remind101/acme-inc:latest
    Pulling repository remind101/acme-inc
    345c7524bc96: Download complete
    a1dd7097a8e8: Download complete
    23debee88b99: Download complete
    31862d352883: Download complete
    c7388ff7ab91: Download complete
    78fb106ed050: Download complete
    133fcef559c4: Download complete
    Status: Image is up to date for remind101/acme-inc:latest
    Status: Created new release v1 for acme-inc
    $ emp releases
    v1    Jan 1 12:55  Deploy remind101/acme-inc:latest
`,
}

func init() {
	cmdDeploy.Flag.BoolVarP(&stream, "stream", "s", false, "boolean to enable the status stream")
}

type PostDeployForm struct {
	Image  string `json:"image"`
	Stream bool   `json:"stream"`
}

func runDeploy(cmd *Command, args []string) {
	r, w := io.Pipe()

	if len(args) < 1 {
		printFatal("You must specify an image to deploy")
	}

	image := args[0]
	message := getMessage()
	form := &PostDeployForm{Image: image, Stream: stream}

	var endpoint string
	appName, _ := app()
	if appName != "" {
		endpoint = fmt.Sprintf("/apps/%s/deploys", appName)
	} else {
		endpoint = "/deploys"
	}

	rh := heroku.RequestHeaders{CommitMessage: message}
	go func() {
		retry := func() {
			runDeploy(cmd, args)
		}
		cleanup := func() {
			must(w.Close())
		}
		defer retryMessageRequired(retry, cleanup)
		must(client.PostWithHeaders(w, endpoint, form, rh.Headers()))
	}()

	outFd, isTerminalOut := term.GetFdInfo(os.Stdout)
	must(jsonmessage.DisplayJSONMessagesStream(r, os.Stdout, outFd, isTerminalOut, nil))
}
