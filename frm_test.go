package frm_test

import (
	"context"
	"os"
	"testing"

	"github.com/acaloiaro/frm"
	"github.com/acaloiaro/frm/internal"
	"github.com/acaloiaro/frm/types"
)

func TestCreateShortCode(t *testing.T) {
	ctx := context.Background()
	f, err := frm.New(frm.Args{
		PostgresURL:         os.Getenv("POSTGRES_URL"),
		PostgresDisableSSL:  true,
		WorkspaceID:         "1",
		WorkspaceIDUrlParam: "client_id",
		PostgresSchema:      "frm_test",
	})
	if err != nil {
		t.Error(err)
	}

	draft, err := internal.Q(ctx, internal.DBArgs{
		URL:        os.Getenv("POSTGRES_URL"),
		DisableSSL: true,
		Schema:     "frm_test",
	}).SaveDraft(ctx, internal.SaveDraftParams{
		Name:        "hello world",
		Fields:      types.FormFields{},
		WorkspaceID: "1",
	})
	if err != nil {
		t.Error(err)
		return
	}

	sc, err := f.CreateShortCode(ctx, frm.CreateShortCodeArgs{
		FormID:    draft.ID,
		SubjectID: "foobar_idx",
	})
	if err != nil {
		t.Error(err)
	}

	if len(sc.ShortCode) != internal.DefaultShortcodeLen {
		t.Fatal("shortcode length is incorrect")
	}
}
