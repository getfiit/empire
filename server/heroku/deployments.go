package heroku

import (
	"net/http"

	"github.com/getfiit/empire/pkg/image"
	streamhttp "github.com/getfiit/empire/pkg/stream/http"
	"github.com/remind101/empire/server/auth"

	"github.com/remind101/empire"
	"golang.org/x/net/context"
)

// PostDeployForm is the form object that represents the POST body.
type PostDeployForm struct {
	Image  image.Image
	Stream bool
}

// ServeHTTPContext implements the Handler interface.
func (h *Server) PostDeploys(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	opts, err := newDeployOpts(ctx, w, req)
	if err != nil {
		return err
	}

	_, err = h.Deploy(ctx, *opts)

	// We only return the MessageRequiredError since all other errors are
	// written to the stream.
	switch err := err.(type) {
	case *empire.MessageRequiredError:
		return err
	}

	return nil
}

func newDeployOpts(ctx context.Context, w http.ResponseWriter, req *http.Request) (*empire.DeployOpts, error) {
	var form PostDeployForm

	if err := Decode(req, &form); err != nil {
		return nil, err
	}

	m, err := findMessage(req)
	if err != nil {
		return nil, err
	}

	w.Header().Set("Content-Type", "application/json; boundary=NL")

	if form.Image.Tag == "" && form.Image.Digest == "" {
		form.Image.Tag = "latest"
	}

	opts := empire.DeployOpts{
		User:    auth.UserFromContext(ctx),
		Image:   form.Image,
		Output:  empire.NewDeploymentStream(streamhttp.StreamingResponseWriter(w)),
		Message: m,
		Stream:  form.Stream,
	}
	return &opts, nil
}
