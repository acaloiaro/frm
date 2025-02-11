package frm_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/acaloiaro/frm"
	"github.com/acaloiaro/frm/internal"
	"github.com/acaloiaro/frm/types"
	"github.com/google/uuid"
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
	}).SaveForm(ctx, internal.SaveFormParams{
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

	sc2, err := f.CreateShortCode(ctx, frm.CreateShortCodeArgs{
		FormID:    draft.ID,
		SubjectID: "foobar_idx",
	})
	if err != nil {
		t.Error(err)
	}

	if sc2.ShortCode != sc.ShortCode {
		t.Error("successive SaveShortCode calls should create same short code")
	}
}

func TestCopyForm(t *testing.T) {
	copiedFormNameSuffix := "(COPY)"
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
	}).SaveForm(ctx, internal.SaveFormParams{
		Name: "hello world",
		Fields: types.FormFields{
			uuid.New().String(): {
				Label:    "What's your name?",
				Order:    0,
				Required: true,
				Type:     types.FormFieldTypeTextSingle,
			},
		},
		WorkspaceID: "1",
	})
	if err != nil {
		t.Error(err)
		return
	}

	clonedForm, err := f.CopyForm(ctx, frm.CopyFormArgs{
		ID:         draft.ID,
		NameSuffix: copiedFormNameSuffix,
	})
	if err != nil {
		t.Error(err)
		return
	}
	expectedCopiedFormName := fmt.Sprintf("%s %s", draft.Name, copiedFormNameSuffix)
	if clonedForm.Name != expectedCopiedFormName {
		t.Error(fmt.Errorf("form name should be '%s'", draft.Name), "got:", clonedForm.Name)
		return
	}
	if len(clonedForm.Fields) != len(draft.Fields) {
		t.Error("cloned form should have the same number of fields as the original, but does not.", "actual number of fields:", len(clonedForm.Fields), "expected number of fields:", len(draft.Fields))
		return
	}
}
