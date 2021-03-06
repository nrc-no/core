package seeder

import (
	"context"

	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/client"
)

// TODO include Spanish translations

func (s *Seed) seedGlobal(ctx context.Context, client client.Client) error {

	emptyDb := types.Database{}
	if s.globalDatabase == emptyDb {
		if err := client.CreateDatabase(ctx, &types.Database{
			Name: "Global",
		}, &s.globalDatabase); err != nil {
			return err
		}
	}

	emptyFolder := types.Folder{}
	if s.globalBioDataFolder == emptyFolder {
		if err := client.CreateFolder(ctx, &types.Folder{
			DatabaseID: s.globalDatabase.ID,
			Name:       GlobalBioDataFolderName,
		}, &s.globalBioDataFolder); err != nil {
			return err
		}
	}

	if err := s.seedGlobalForms(ctx, client); err != nil {
		return err
	}

	return nil
}

func (s *Seed) seedGlobalForms(ctx context.Context, client client.Client) error {
	// Root entities ---------------------------------

	if len(s.globalRootIndividualForm.ID) == 0 {
		if err := client.CreateForm(ctx, &types.FormDefinition{
			DatabaseID: s.globalDatabase.ID,
			FolderID:   s.globalBioDataFolder.ID,
			Name:       GlobalIndividualFormName,
			Type:       types.RecipientFormType,
			Fields: types.FieldDefinitions{
				{
					Name:        "Full Name",
					Description: "The full name of the individual",
					Required:    true,
					FieldType: types.FieldType{
						Text: &types.FieldTypeText{},
					},
				}, {
					Name:        "Preferred Name",
					Description: "The name which will be used to refer to the beneficiary within Core",
					Required:    true,
					FieldType: types.FieldType{
						Text: &types.FieldTypeText{},
					},
				},
			},
		}, &s.globalRootIndividualForm); err != nil {
			return err
		}
	}

	if len(s.globalRootHouseholdForm.ID) == 0 {
		if err := client.CreateForm(ctx, &types.FormDefinition{
			DatabaseID: s.globalDatabase.ID,
			FolderID:   s.globalBioDataFolder.ID,
			Name:       GlobalHouseholdFormName,
			Type:       types.RecipientFormType,
			Fields: types.FieldDefinitions{
				text("Household name", true),
			},
		}, &s.globalRootHouseholdForm); err != nil {
			return err
		}
	}
	if len(s.globalRootBeneficiaryForm.ID) == 0 {
		if err := client.CreateForm(ctx, &types.FormDefinition{
			DatabaseID: s.globalDatabase.ID,
			FolderID:   s.globalBioDataFolder.ID,
			Name:       GlobalIndividualBeneficiaryFormName,
			Type:       types.RecipientFormType,
			Fields: types.FieldDefinitions{
				&types.FieldDefinition{
					Name:     "Individual",
					Key:      true,
					Required: true,
					FieldType: types.FieldType{
						Reference: &types.FieldTypeReference{
							DatabaseID: s.globalDatabase.ID,
							FormID:     s.globalRootIndividualForm.ID,
						},
					},
				},
				yesNo("Has the beneficiary consented to NRC using their data?"),
				text("URL to proof of beneficiary consent", true),
				yesNo("Beneficiary prefers to remain anonymous?"),
				yesNo("Is the beneficiary a minor?"),
				yesNo("Beneficiary presents protection concerns?"),
				yesNo("Would you say you experience some form of physical challenges?"),
				dropdown("How would you define the intensity of such challenges?", wgShortSet, true),
				yesNo("Would you say you experience some form of sensory challenges?"),
				dropdown("How would you define the intensity of such challenges?", wgShortSet, true),
				yesNo("Would you say you experience some form of mental challenges?"),
				dropdown("How would you define the intensity of such challenges?", wgShortSet, true),
				dropdown("Displacement Status", globalDisplacementStatuses, true),
				dropdown("Gender", globalGenders, true),
				&types.FieldDefinition{
					Name:     "Affiliated Household",
					Required: true,
					FieldType: types.FieldType{
						Reference: &types.FieldTypeReference{
							DatabaseID: s.globalDatabase.ID,
							FormID:     s.globalRootHouseholdForm.ID,
						},
					}},
				yesNo("Are you a representative for the household?"),
			},
		}, &s.globalRootBeneficiaryForm); err != nil {
			return err
		}
	}

	return nil
}
