package heroku

import (
	"net/http"

	"github.com/remind101/empire"
	"github.com/getfiit/empire/pkg/heroku"
	"golang.org/x/net/context"
)

func (h *Server) PostCerts(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	a, err := findApp(ctx, h)
	if err != nil {
		return err
	}

	var form heroku.CertsAttachOpts

	if err := Decode(r, &form); err != nil {
		return err
	}

	opts := empire.CertsAttachOpts{
		App:  a,
		Cert: *form.Cert,
	}
	if form.Process != nil {
		opts.Process = *form.Process
	}

	if err := h.CertsAttach(ctx, opts); err != nil {
		return err
	}

	return Encode(w, nil)
}
