package db_test

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/acaloiaro/frm"
	"github.com/acaloiaro/frm/internal"
	"github.com/acaloiaro/frm/types"
	"github.com/google/uuid"
)

var (
	workspaceID = uuid.New()
)

func TestCreateAndUpdateForm(t *testing.T) {
	const nameUpdate = "Tell us about you"
	fieldID := uuid.MustParse("1afd4bf9-42a4-4dfe-b359-d46a65ce5ba5")
	fieldID2 := uuid.MustParse("2afd4bf9-42a4-4dfe-b359-d46a65ce5ba5")
	fieldID3 := uuid.MustParse("3afd4bf9-42a4-4dfe-b359-d46a65ce5ba5")
	fields := types.FormFields{
		fieldID.String(): types.FormField{
			ID:       fieldID,
			Label:    "hello world",
			Type:     types.FormFieldTypeTextMultiple,
			Required: true,
		},
	}
	updatedFields := types.FormFields{
		fieldID.String(): types.FormField{
			ID:          fieldID,
			Order:       1,
			Label:       "What's your name?",
			Placeholder: "Tell us your name",
			Type:        types.FormFieldTypeTextSingle,
			Required:    true,
		},
		fieldID2.String(): types.FormField{
			ID:          fieldID2,
			Order:       2,
			Label:       "What are your view on bears?",
			Type:        types.FormFieldTypeTextMultiple,
			Placeholder: "Tell us how you generally feel about bears",
			Required:    false,
		},
		fieldID3.String(): types.FormField{
			ID:          fieldID3,
			Order:       3,
			Label:       "Which type of bear is best?",
			Placeholder: "Choose a bear",
			Type:        types.FormFieldTypeMultiSelect,
			Options: []types.Option{
				{ID: uuid.New(), Label: "Black bear", Value: "Black bear", Selected: false},
				{ID: uuid.New(), Label: "Brown bear", Value: "Brown bear", Selected: false},
				{ID: uuid.New(), Label: "Blue bear", Value: "Blue bear", Selected: false},
			},
			Required: true,
		},
	}
	ctx := context.Background()
	frms, err := frm.New(frm.Args{
		PostgresURL: os.Getenv("POSTGRES_URL"),
		WorkspaceID: workspaceID,
	})
	if err != nil {
		t.Error(err)
		return
	}
	err = frms.Init(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	f, err := internal.Q(ctx, frms.PostgresURL).SaveDraft(ctx, internal.SaveDraftParams{
		Name:        "hello world",
		Fields:      fields,
		WorkspaceID: frms.WorkspaceID,
	})
	if err != nil {
		t.Error(err)
		return
	}

	f, err = internal.Q(ctx, frms.PostgresURL).SaveDraft(ctx, internal.SaveDraftParams{
		ID:          f.ID,
		Name:        nameUpdate,
		Fields:      updatedFields,
		WorkspaceID: frms.WorkspaceID,
	})
	if err != nil {
		t.Error(err)
		return
	}

	if f.Name != nameUpdate {
		t.Error(fmt.Errorf("name did not update: %s!=%s", f.Name, nameUpdate))
	}

	if !reflect.DeepEqual(f.Fields, updatedFields) {
		t.Error(fmt.Errorf("fields did not update: %v!=%v", f.Fields, updatedFields))
	}
}
