package seeder

import (
	"context"

	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/client"
)

func (seedUganda *Seed) seedUganda(ctx context.Context, client client.Client) error {
	var ugDB types.Database
	if err := client.CreateDatabase(ctx, &types.Database{
		Name: "Uganda",
	}, &ugDB); err != nil {
		return err
	}

	if err := seedUganda.seedUgandaGeneralForms(ctx, client, ugDB.ID); err != nil {
		return err
	}

	if err := seedUganda.seedUgandaIclaForms(ctx, client, ugDB.ID); err != nil {
		return err
	}

	if err := seedUganda.seedUgandaIntakeForms(ctx, client, ugDB.ID); err != nil {
		return err
	}

	if err := seedUganda.seedUgandaProtectionForms(ctx, client, ugDB.ID); err != nil {
		return err
	}

	return seedUganda.seedUgandaEiLfsForms(ctx, client, ugDB.ID)
}

func (s *Seed) seedUgandaGeneralForms(ctx context.Context, client client.Client, dbID string) error {
	var folder types.Folder
	if err := client.CreateFolder(ctx, &types.Folder{
		DatabaseID: dbID,
		Name:       "General",
	}, &folder); err != nil {
		return err
	}

	var externalReferral *types.FormDefinition

	var forms = []*types.FormDefinition{externalReferral}

	externalReferral = newFormDefinition(
		dbID,
		folder.ID,
		"External Referral Form",
		[]*types.FieldDefinition{
			newFieldDefinition("Individual Beneficiary", "The beneficiary this external referral form has been completed for", false, true, types.FieldType{
				Reference: &types.FieldTypeReference{
					DatabaseID: s.globalDatabase.ID,
					FormID:     s.globalRootBeneficiaryForm.ID,
				},
			}),
			yesNoQuestion("Was the referral accepted by the other provider"),
			newFieldDefinition("Provide any pertinent details on service needs / requests", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
			yesNoQuestion("This case is now closed"),
		},
	)

	for _, form := range forms {
		if err := client.CreateForm(ctx, form, nil); err != nil {
			return err
		}
	}
	return nil
}

func (s *Seed) seedUgandaIclaForms(ctx context.Context, client client.Client, dbID string) error {
	var folder types.Folder
	if err := client.CreateFolder(ctx, &types.Folder{
		DatabaseID: dbID,
		Name:       "ICLA",
	}, &folder); err != nil {
		return err
	}

	var (
		screening      *types.FormDefinition
    caseAssessment *types.FormDefinition
		followupReport *types.FormDefinition
		appointment    *types.FormDefinition
		referral       *types.FormDefinition
	)

	var forms = []*types.FormDefinition{screening, followupReport, appointment, referral}

	screening = newFormDefinition(
		dbID,
		folder.ID,
		"Screening",
		[]*types.FieldDefinition{
			{
				Name:     "Individual Beneficiary",
				Key:      true,
				Required: true,
				FieldType: types.FieldType{
					Reference: &types.FieldTypeReference{
						DatabaseID: s.globalDatabase.ID,
						FormID:     s.globalRootBeneficiaryForm.ID,
					},
				},
			},
			{
				Name:        "What/Describe legal issue/concern are you facing? ",
				Description: "Facts and details of the problem",
				Key:         false,
				Required:    true,
				FieldType:   types.FieldType{MultilineText: &types.FieldTypeMultilineText{}},
			},
			{
				Name:     "Select the legal issue of concern",
				Required: true,
				FieldType: types.FieldType{
					SingleSelect: &types.FieldTypeSingleSelect{Options: ugICLALegalIssues},
				},
			},
			ifOtherPleaseSpecify,
			{
				Name: "What action has been taken to solve the problem, if any?", Required: true, FieldType: types.FieldType{Text: &types.FieldTypeText{}}},
			{
				Name:      "Information of beneficiary's representative",
				FieldType: types.FieldType{}, // FIXME this is a header
			},
			yesNoQuestion("Is there a representative for this individual?"),
			newFieldDefinition("Full name of representative (lawyer or another person/Legal Guardian/Other)", "", false, false, types.FieldType{
				Text: &types.FieldTypeText{},
			}),
			newFieldDefinition("Reason for representative (instead of beneficiary):", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
			yesNoQuestion("Is the guardianship legal as per national legislation?"),
			newFieldDefinition("If yes, attach/upload the legal/court order", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
			{
				Name:      "RSD",
				FieldType: types.FieldType{}, // FIXME this is a header
			},
			newFieldDefinition("What is the individual's displacement status?", "", false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{Options: ugICLADisplacementStatus}}),
			yesNoQuestion("Are you at risk of being stateless?"),
			newFieldDefinition("Describe this in detail", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
			newFieldDefinition("What RSD documents do you have?", "", false, false, types.FieldType{MultiSelect: &types.FieldTypeMultiSelect{Options: ugICLARSDDocuments}}),
			newFieldDefinition("Comment on RSD documents", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
			newFieldDefinition("Specific RSD issues presented", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
			{
				Name:      "HLP",
				FieldType: types.FieldType{}, // FIXME this is a header
			},
			newFieldDefinition("What specific HLP concern is presented?", "", false, true, types.FieldType{MultiSelect: &types.FieldTypeMultiSelect{Options: ugICLASpecificHLPConcern}}),
			{
				Name:      "Housing",
				FieldType: types.FieldType{}, // FIXME this is a sub-section
			},
			newFieldDefinition("Does the individual stay in their own house or rent?", "", false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{Options: ugICLAHomeOwnership}}),
			yesNoQuestion("In case of rent, does the individual posses any agreement?"),
			newFieldDefinition("What kind of agreement or  proof does the individual possess?", "", false, false, types.FieldType{
				Text: &types.FieldTypeText{},
			}),
			yesNoQuestion("Have you been or are you at risk of eviction?"),
			newFieldDefinition("If yes, What eviction document or proof do you posses?", "", false, false, types.FieldType{MultiSelect: &types.FieldTypeMultiSelect{Options: ugICLAEvictionDocuments}}),
			newFieldDefinition("Comment on eviction document", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
			{
				Name:      "Land",
				FieldType: types.FieldType{}, // FIXME this is a sub-section
			},
			yesNoQuestion("Are you the legal owner of the land?"),
			newFieldDefinition("Nature of tenancy", "", false, false, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{Options: ugICLANatureLandTenure}}),
			newFieldDefinition("Nature of tenure", "", false, false, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{Options: ugICLANatureTenure}}),
			newFieldDefinition("Land supporting documents possessed", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
			yesNoQuestion("Have you been or are you at risk of eviction?"),
			newFieldDefinition("If yes, What eviction document or proof do you posses?", "", false, false, types.FieldType{MultiSelect: &types.FieldTypeMultiSelect{Options: ugICLAEvictionDocuments}}),
			newFieldDefinition("Specific land issues", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
			{
				Name:      "Property",
				FieldType: types.FieldType{}, // FIXME this is a sub-section
			},
			newFieldDefinition("Nature of property in contest", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
			yesNoQuestion("Do you have legal ownership of property?"),
			newFieldDefinition("Proof of property ownership", "(Supporting documents)", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
			newFieldDefinition("Inquiry on property acquisition", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
			{
				Name:      "LCD",
				FieldType: types.FieldType{}, // FIXME this is a header
			},
			newFieldDefinition("What documentation challenges do you have?", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
			newFieldDefinition("Type of document", "", false, false, types.FieldType{MultiSelect: &types.FieldTypeMultiSelect{Options: ugICLATypeOfDocumentation}}),
			newFieldDefinition("What action had been taken so far on this issue?", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
			{
				Name: "ELP",
				FieldType: types.FieldType{}, // FIXME this is a header
			},
							newFieldDefinition("Is it an employment or business challenge?", "", false, false, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{Options: ugICLATypeOfChallenge}}),
							newFieldDefinition("What employment challenges do you have?", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
							newFieldDefinition("What type of agreement do you have?", "", false, false, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{Options: ugICLATypeOfAgreement}}),
							newFieldDefinition("What actions have been taken?", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
							newFieldDefinition("What business related challenge do you have?", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
							newFieldDefinition("What business registration services do you need?", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
							newFieldDefinition("What actions have been taken?", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
			})

      caseAssessment = &types.FormDefinition{
        DatabaseID: dbID,
        FolderID: folder.ID,
        Name: "Case Assessment",
        Fields: types.FieldDefinitions{
						newFieldDefinition("Type of actions for case worker agreed upon with beneficiary ", "", false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						newFieldDefinition("Notes on types of actions agreed upon", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						newFieldDefinition("Are there any elements of risks for the safety or well-being of the beneficiary or that of a relative in relation to the suggested course of action", "", false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						newFieldDefinition("Narrative regarding elements of risks", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						newFieldDefinition("Any particular Protection Risks", "", false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						newFieldDefinition("Is the guardianship legal as per national legislation", "", false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						newFieldDefinition("If yes, please indicate what type", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						newFieldDefinition("Guardianship notes", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						newFieldDefinition("Are there any unintended negative consequences of the suggested course of actions for the beneficiary's family or larger community", "", false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						newFieldDefinition("Unintended consequences notes", "If any of the answers were 'yes', discuss with the beneficiary what might be done to avoid or minimise the risks or negative consequences", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						newFieldDefinition("Actions agreed upon with the beneficiary", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						newFieldDefinition("Does the beneficiary agree to continue with the case", "Discuss the pro's and con's of the suggested course of action, including the analysis of risks", false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						newFieldDefinition("Notes on pro's and con's of suggested course of action", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						newFieldDefinition("Is the guardianship legal as per national legislation", "", false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						newFieldDefinition("Is a Best Interest Determination needed for the case", "If 'yes', refer the case to social services or an appropriate child protection actor", false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						newFieldDefinition("Agreed follow up means", "", false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
        },
      }

			newFieldDefinition("Case Assessment - Case Plan", "ICLA case plan for an individual beneficiary", false, true, types.FieldType{
				SubForm: &types.FieldTypeSubForm{
					Fields: []*types.FieldDefinition{
					},
				},
			}),
			newFieldDefinition("Case Assessment - Action Plan", "ICLA action plan for an individual beneficiary", false, true, types.FieldType{
				SubForm: &types.FieldTypeSubForm{
					Fields: []*types.FieldDefinition{
						newFieldDefinition("Which service has the beneficiary together with staff agreed to take", "", false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						newFieldDefinition("If other specify", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						newFieldDefinition("Type of actions for case worker agreed upon with beneficiary ", "", false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						newFieldDefinition("Action comment", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
					},
				},
			}),
			newFieldDefinition("Case Assessment - Case Closure", "Used to close an ICLA case in Uganda", false, true, types.FieldType{
				SubForm: &types.FieldTypeSubForm{
					Fields: []*types.FieldDefinition{
						newFieldDefinition("Case closure date", "", false, false, types.FieldType{Date: &types.FieldTypeDate{}}),
						newFieldDefinition("What is the reason for the case closure?", "", false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						newFieldDefinition("Notes", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
					},
				},
			}),
		},
	)

	followupReport = newFormDefinition(
		dbID,
		folder.ID,
		"Followup Form",
		[]*types.FieldDefinition{
			newFieldDefinition("Individual Beneficiary", "The beneficiary this followup form has been completed for", false, true, types.FieldType{
				Reference: &types.FieldTypeReference{
					DatabaseID: s.globalDatabase.ID,
					FormID:     s.globalRootBeneficiaryForm.ID,
				},
			}),
			newFieldDefinition("Date of follow-up", "", false, true, types.FieldType{Date: &types.FieldTypeDate{}}),
			newFieldDefinition("Follow-up", "", false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
			newFieldDefinition("If \"Other\", specify", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
			newFieldDefinition("Notes from the follow-up undertaken", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
			newFieldDefinition("Copies of documents", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
		},
	)

	appointment = newFormDefinition(
		dbID,
		folder.ID,
		"Appointment Form",
		[]*types.FieldDefinition{
			newFieldDefinition("Individual Beneficiary", "The beneficiary this appointment form has been completed for", false, true, types.FieldType{
				Reference: &types.FieldTypeReference{
					DatabaseID: s.globalDatabase.ID,
					FormID:     s.globalRootBeneficiaryForm.ID,
				},
			}),
			newFieldDefinition("Name", "", false, true, types.FieldType{Text: &types.FieldTypeText{}}),
			newFieldDefinition("Place", "", false, true, types.FieldType{Text: &types.FieldTypeText{}}),
			newFieldDefinition("Date", "", false, true, types.FieldType{Date: &types.FieldTypeDate{}}),
			newFieldDefinition("Preferred contact method", "", false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
			newFieldDefinition("Appointment purpose", "", false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
			newFieldDefinition("Preferred date", "", false, true, types.FieldType{Date: &types.FieldTypeDate{}}),
		},
	)

	referral = newFormDefinition(
		dbID,
		folder.ID,
		"Referral Form",
		[]*types.FieldDefinition{
			newFieldDefinition("Individual Beneficiary", "The beneficiary this ICLA referral form has been completed for", false, true, types.FieldType{
				Reference: &types.FieldTypeReference{
					DatabaseID: s.globalDatabase.ID,
					FormID:     s.globalRootBeneficiaryForm.ID,
				},
			}),
			newFieldDefinition("Beneficiary's information", "", false, true, types.FieldType{
				SubForm: &types.FieldTypeSubForm{
					Fields: []*types.FieldDefinition{
						newFieldDefinition("Type of referral (internal)", "", false, false, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						newFieldDefinition("Type of referral (external)", "If yes, provide details below", false, false, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						newFieldDefinition("Organization", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
						newFieldDefinition("Contact person", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
						newFieldDefinition("Phone number", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
						newFieldDefinition("Email", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
						newFieldDefinition("Types of services/assistence requested", "", false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						newFieldDefinition("If \"Other\", specify", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
						newFieldDefinition("Reason for the referral", "", false, true, types.FieldType{Text: &types.FieldTypeText{}}),
					},
				},
			}),
			newFieldDefinition("Consent", "", false, true, types.FieldType{
				SubForm: &types.FieldTypeSubForm{
					Fields: []*types.FieldDefinition{
						newFieldDefinition("Has the beneficiary consented to the release of their information for the referral?", "", false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						newFieldDefinition("If 'yes', upload a signed consent form and proceed", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
						newFieldDefinition("If 'no', explain the reason why and do not refer the case", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
					},
				},
			}),
			newFieldDefinition("Means of referral", "", false, true, types.FieldType{
				SubForm: &types.FieldTypeSubForm{
					Fields: []*types.FieldDefinition{
						newFieldDefinition("Does the beneficiary have any restrictions to being referred?", "", false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						newFieldDefinition("Means of referral", "", false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						newFieldDefinition("Means and terms of receiving feedback from the client", "", false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
					},
				},
			}),
		},
	)

	for _, form := range forms {
		if err := client.CreateForm(ctx, form, nil); err != nil {
			return err
		}
	}
	return nil
}
func (s *Seed) seedUgandaIntakeForms(ctx context.Context, client client.Client, dbID string) error {
	var intakeFolder types.Folder
	if err := client.CreateFolder(ctx, &types.Folder{
		DatabaseID: dbID,
		Name:       "Intake",
	}, &intakeFolder); err != nil {
		return err
	}

	var (
		individualBeneficiaryIntake *types.FormDefinition
	)

	var forms = []*types.FormDefinition{individualBeneficiaryIntake}

	individualBeneficiaryIntake = newFormDefinition(
		dbID,
		intakeFolder.ID,
		"Individual Intake Form",
		[]*types.FieldDefinition{
			newFieldDefinition("Individual Beneficiary", "The beneficiary this intake form has been completed for", true, true, types.FieldType{
				Reference: &types.FieldTypeReference{
					DatabaseID: s.globalDatabase.ID,
					FormID:     s.globalRootBeneficiaryForm.ID,
				},
			}),
			newFieldDefinition("PII", "Personally Indentifiable Information", false, false, types.FieldType{
				SubForm: &types.FieldTypeSubForm{
					Fields: []*types.FieldDefinition{
						newFieldDefinition("Date of birth", "", false, true, types.FieldType{
							Date: &types.FieldTypeDate{},
						}),
						newFieldDefinition("Nationality", "", false, true, types.FieldType{
							SingleSelect: &types.FieldTypeSingleSelect{Options: ugNationality},
						}),
						newFieldDefinition("Relationship to HH representative", "", false, true, types.FieldType{
							SingleSelect: &types.FieldTypeSingleSelect{Options: ugRelationshipToHH},
						}),
						ifOtherPleaseSpecify,
					},
				},
			}),
			newFieldDefinition("Identification details", "", false, false, types.FieldType{
				SubForm: &types.FieldTypeSubForm{
					Fields: []*types.FieldDefinition{
						newFieldDefinition("Type of identification", "", false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{Options: ugTypeOfIdentification}}),
						ifOtherPleaseSpecify,
						newFieldDefinition("Identification number", "e.g.: ID Number", false, true, types.FieldType{Text: &types.FieldTypeText{}}),
						newFieldDefinition("Additional ID #1 (type & number)", "e.g.: UNHCR Number", false, true, types.FieldType{Text: &types.FieldTypeText{}}),
						newFieldDefinition("Additional ID #2 (type & number)", "e.g.: UNHCR Number", false, true, types.FieldType{Text: &types.FieldTypeText{}}),
					},
				},
			}),
			newFieldDefinition("NRC Details", "", false, false, types.FieldType{
				SubForm: &types.FieldTypeSubForm{
					Fields: []*types.FieldDefinition{
						newFieldDefinition("Enumerator name (if no Okta ID)", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
						newFieldDefinition("Date of identification", "", false, false, types.FieldType{Date: &types.FieldTypeDate{}}),
						newFieldDefinition("Location of identification", "", false, false, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{Options: ugLocationOfIdentification}}),
						ifOtherPleaseSpecify,
						newFieldDefinition("Source of identification", "", false, false, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{Options: ugSourceOfIdentification}}),
						ifOtherPleaseSpecify,
						newFieldDefinition("District", "", false, false, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{Options: ugAdmin2}}),
						newFieldDefinition("Subcounty", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
						newFieldDefinition("Parish", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
						newFieldDefinition("Village", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
					},
				},
			}),
			newFieldDefinition("Contact Info", "", false, false, types.FieldType{
				SubForm: &types.FieldTypeSubForm{
					Fields: []*types.FieldDefinition{
						newFieldDefinition("Preferred Language", "List the language preferred by the respondent ", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
						newFieldDefinition("Email address", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
						newFieldDefinition("Physical Address", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						newFieldDefinition("Phone number (1)", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
						newFieldDefinition("Phone number (2)", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
						newFieldDefinition("Preferred means of contact", "", false, false, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{Options: ugMeansOfContact}}),
						newFieldDefinition("Instructions for contact or other comments", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						yesNoQuestion("Can NRC initiate contact with the beneficiary?"),
						yesNoQuestion("Does the beneficiary require an interpreter?"),
						newFieldDefinition("Interpreter name", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
					},
				},
			}),
			newFieldDefinition("Consent", "", false, true, types.FieldType{
				SubForm: &types.FieldTypeSubForm{
					Fields: []*types.FieldDefinition{
						newFieldDefinition("Contact number of caregiver", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
						newFieldDefinition("Main obstacles you face in meeting accomodation needs", "", false, true, types.FieldType{
							SingleSelect: &types.FieldTypeSingleSelect{},
						}),
						newFieldDefinition("Can the HH meet WASH needs?", "", false, true, types.FieldType{
							SingleSelect: &types.FieldTypeSingleSelect{},
						}),
						newFieldDefinition("Main obstacles in meeting WASH needs", "", false, true, types.FieldType{
							SingleSelect: &types.FieldTypeSingleSelect{},
						}),
						newFieldDefinition("Summary narrative", "", false, false, types.FieldType{
							MultilineText: &types.FieldTypeMultilineText{},
						}),
						newFieldDefinition("NRC Staff Name", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
						newFieldDefinition("Please upload the consent form signed by the beneficiary", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
						yesNoQuestion("Has the client consented to NRC collecting and storing and using their data?"),
						yesNoQuestion("Has client consented to NRC sharing information about needs with relevant providers?"),
						yesNoQuestion("Can NRC staff initiate contact with the beneficiary?"),
					},
				},
			}),
			newFieldDefinition("Situation Analysis", "A form used to carry out a situation analysis for an individual in Uganda", false, true, types.FieldType{
				SubForm: &types.FieldTypeSubForm{
					Fields: []*types.FieldDefinition{
						newFieldDefinition("Do you think you are living a safe and dignified life? Are you achieving what you want? Are you able to live a good life? Probe for description.", "Multi-sectoral staff to enter the feedback from the individual.", false, true, types.FieldType{
							MultilineText: &types.FieldTypeMultilineText{},
						}),
						newFieldDefinition("How are you addressing these challenges and barriers? What is standing in your way? Can you give me some examples of how you are dealing with these challenges?", "Multi-sectoral staff to enter the feedback from the individual.", false, true, types.FieldType{
							MultilineText: &types.FieldTypeMultilineText{},
						}),
						newFieldDefinition("What are some solutions you see for this and how could we work together on these solutions? How could we work to reduce these challenges together?", "Multi-sectoral staff to enter the feedback from the individual.", false, true, types.FieldType{
							MultilineText: &types.FieldTypeMultilineText{},
						}),
						newFieldDefinition("If we were to work together on this, what could we do together? What would make the most difference for you?", "Multi-sectoral staff to enter the feedback from the individual.", false, true, types.FieldType{
							Text: &types.FieldTypeText{},
						}),
					},
				},
			}),
			newFieldDefinition("Response", "Response form for an individual in Uganda", false, true, types.FieldType{
				SubForm: &types.FieldTypeSubForm{
					Fields: []*types.FieldDefinition{
						newFieldDefinition("Which service has the individual/community requested as a starting point of support?", "Provide any pertinent details on service needs / requests.", false, true, types.FieldType{
							MultiSelect: &types.FieldTypeMultiSelect{Options: ugServicesRequested},
						}),
						newFieldDefinition("Provide any pertinent details on service needs / requests", "", false, true, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						newFieldDefinition("What other services has the individual /household requested/identified?", "", false, true, types.FieldType{
							MultiSelect: &types.FieldTypeMultiSelect{Options: ugServicesRequested},
						}),
						newFieldDefinition("Provide any pertinent details on service needs / requests for the other services requested.", "", false, true, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						newFieldDefinition("What is the perceived priority response level of the individual / household?", "", false, true, types.FieldType{
							SingleSelect: &types.FieldTypeSingleSelect{Options: ugPerceivedPriorityResponseLevel},
						}),
						newFieldDefinition("Provide any pertinent details on how priority was determined.", "", false, true, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						newFieldDefinition("Other information", "Comments or notes", false, true, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
					},
				},
			}),
			newFieldDefinition("Response - Vulnerability Assessment", "", false, true, types.FieldType{
				SubForm: &types.FieldTypeSubForm{
					Fields: []*types.FieldDefinition{
						newFieldDefinition("What needs does the client have?", "", false, true, types.FieldType{
							MultiSelect: &types.FieldTypeMultiSelect{Options: ugSpecificNeed}}),
						yesNoQuestion("Has the client suffered a protection incident?"),
						yesNoQuestion("Has a protection incident report been filled?"),
						yesNoQuestion("Does the client have specific protection needs?"),
						newFieldDefinition("If so, explain", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						yesNoQuestion("Does the client feel safe?"),
						newFieldDefinition("If not, explain", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						yesNoQuestion("Has any action already been taken on security?"),
						newFieldDefinition("Security actions taken to whom?", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
						newFieldDefinition("Security actions explanation", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						yesNoQuestion("Does the client want to take any step on security?"),
						newFieldDefinition("Security actions wanted with regards to whom?", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
						newFieldDefinition("Security actions wanted explanation", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						yesNoQuestion("Is any Housing Land Property (HLP) assistance needed?  "),
						newFieldDefinition("Explanation of HLP assistance needed", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						yesNoQuestion("Is assistance for Refugee Status Determination (RSD) needed? "),
						newFieldDefinition("Explanation of RSD assistance needed", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						yesNoQuestion("Is assistance for Legal and Civil Identity (LCD) needed? "),
						newFieldDefinition("Explanation of LCD assistance needed", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						yesNoQuestion("Is assistance for Employment Laws and Procedures needed? "),
						newFieldDefinition("Explanation of employment laws and procedures assistance needed", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						newFieldDefinition("Laws, rules, policy applicable", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						newFieldDefinition("Recommendation/Resolution", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						yesNoQuestion("Does the client require medical assistence?"),
						yesNoQuestion("Does the client require a referral?"),
						newFieldDefinition("Explanation of medical referral", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						yesNoQuestion("Has the client reported any psychosocial needs or have you noticed such needs?"),
						newFieldDefinition("Explanation of psychosocial needs", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						yesNoQuestion("Does the client have family support?"),
						newFieldDefinition("Explanation regarding client's familial support", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						newFieldDefinition("Actions to be taken", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						newFieldDefinition("Explanation regarding actions to be taken", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
					},
				},
			}),
		},
	)

	for _, form := range forms {
		if err := client.CreateForm(ctx, form, nil); err != nil {
			return err
		}
	}
	return nil
}

func (s *Seed) seedUgandaProtectionForms(ctx context.Context, client client.Client, dbID string) error {
	var protectionFolder types.Folder
	if err := client.CreateFolder(ctx, &types.Folder{
		DatabaseID: dbID,
		Name:       "Protection",
	}, &protectionFolder); err != nil {
		return err
	}
	var (
		protectionIntake         *types.FormDefinition
		protectionIncident       *types.FormDefinition
		protectionFollowupReport *types.FormDefinition
		protectionActionReport   *types.FormDefinition
		protectionReferral       *types.FormDefinition
	)
	var forms = []*types.FormDefinition{protectionIntake, protectionIncident, protectionFollowupReport, protectionActionReport, protectionReferral}

	protectionIntake = newFormDefinition(
		dbID,
		protectionFolder.ID,
		"Protection Case Opening Form",
		[]*types.FieldDefinition{
			newFieldDefinition("Individual Beneficiary", "The beneficiary this protection intake form has been completed for", false, true, types.FieldType{
				Reference: &types.FieldTypeReference{
					DatabaseID: s.globalDatabase.ID,
					FormID:     s.globalRootBeneficiaryForm.ID,
				},
			}),
			{
				Name:        "Intake Screening",
				Description: "A form used to collect intake details for the Uganda Protection team",
				FieldType: types.FieldType{
					SubForm: &types.FieldTypeSubForm{
						Fields: []*types.FieldDefinition{
							newFieldDefinition("Date of screening", "", false, true, types.FieldType{Date: &types.FieldTypeDate{}}),
							newFieldDefinition("Have you been exposed to any protection risks?", "", false, false, types.FieldType{MultiSelect: &types.FieldTypeMultiSelect{Options: ugProtectionRisk}}),
							ifOtherPleaseSpecify,
							newFieldDefinition("What type of protection concern experienced?", "", false, false, types.FieldType{MultiSelect: &types.FieldTypeMultiSelect{Options: ugProtectionConcern}}),
							newFieldDefinition("Provide details", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
							yesNoQuestion("Is this case confidential?"),
						},
					},
				}},
			{
				Name: "Protection Incident",
				FieldType: types.FieldType{
					SubForm: &types.FieldTypeSubForm{
						Fields: []*types.FieldDefinition{
							yesNoQuestion("Protection indident being reported?"),
							newFieldDefinition("If yes, provide protection incident details", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
							newFieldDefinition("Location of incident", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
							newFieldDefinition("Time of incident", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
							newFieldDefinition("Date of incident", "", false, false, types.FieldType{Date: &types.FieldTypeDate{}}),
							newFieldDefinition("Received by", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
							newFieldDefinition("Vulnerability", "", false, false, types.FieldType{MultiSelect: &types.FieldTypeMultiSelect{Options: ugProtectionVulnerability}}),
							newFieldDefinition("Description of the incident", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
							yesNoQuestion("Hase the incident been reported to the police?"),
							newFieldDefinition("Specify platform", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
						},
					},
				},
			},
			{
				Name:        "Social Status Assessment",
				Description: "Used to collect information regarding the social status of an individual beneficiary",
				FieldType: types.FieldType{
					SubForm: &types.FieldTypeSubForm{
						Fields: []*types.FieldDefinition{
							newFieldDefinition("Does the client have specific needs?", "", false, false, types.FieldType{MultiSelect: &types.FieldTypeMultiSelect{Options: ugProtectionSpecificNeeds}}),
							newFieldDefinition("Comment", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
							yesNoQuestion("Does any other member of the HH have specific needs?"),
							newFieldDefinition("What specific needs does the HH member have?", "", false, false, types.FieldType{MultiSelect: &types.FieldTypeMultiSelect{Options: ugProtectionSpecificNeeds}}),
							newFieldDefinition("Home situation", "", false, false, types.FieldType{MultiSelect: &types.FieldTypeMultiSelect{Options: ugProtectionHomeSituation}}),
							newFieldDefinition("Does the client have a disability", "", false, false, types.FieldType{MultiSelect: &types.FieldTypeMultiSelect{Options: ugProtectionDisability}}),
							yesNoQuestion("Does any other member of the HH live with disability?"),
							newFieldDefinition("Which other member of the HH lives with a disability?", "", false, false, types.FieldType{Reference: &types.FieldTypeReference{
								DatabaseID: s.globalDatabase.ID,
								FormID:     s.globalRootBeneficiaryForm.ID,
							}}),
							newFieldDefinition("If the HH member is not a registered beneficiary, provide details (Name, ID number, etc)", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
							newFieldDefinition("How many people are able to work in your HH?", "", false, false, types.FieldType{Quantity: &types.FieldTypeQuantity{}}),
							newFieldDefinition("How often do they work?", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
							newFieldDefinition("What do they do?", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
							yesNoQuestion("Do you receive humanitarian assistance?"),
							newFieldDefinition("Comment/Recent changes regarding humanitarian assistance", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
							newFieldDefinition("Comment on living situation", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
							{
								Name: "Household composition",
								FieldType: types.FieldType{
									SubForm: &types.FieldTypeSubForm{
										Fields: []*types.FieldDefinition{
											newFieldDefinition("Number of 0-5 year old girls", "", false, false, types.FieldType{Quantity: &types.FieldTypeQuantity{}}),
											newFieldDefinition("Number of 0-5 year old boys", "", false, false, types.FieldType{Quantity: &types.FieldTypeQuantity{}}),
											newFieldDefinition("Number of 6-12 year old girls", "", false, false, types.FieldType{Quantity: &types.FieldTypeQuantity{}}),
											newFieldDefinition("Number of 6-12 year old boys", "", false, false, types.FieldType{Quantity: &types.FieldTypeQuantity{}}),
											newFieldDefinition("Number of 13-17 year old girls", "", false, false, types.FieldType{Quantity: &types.FieldTypeQuantity{}}),
											newFieldDefinition("Number of 13-17 year old boys", "", false, false, types.FieldType{Quantity: &types.FieldTypeQuantity{}}),
											newFieldDefinition("Number of 18-59 year old females", "", false, false, types.FieldType{Quantity: &types.FieldTypeQuantity{}}),
											newFieldDefinition("Number of 18-59 year old males", "", false, false, types.FieldType{Quantity: &types.FieldTypeQuantity{}}),
											newFieldDefinition("Number of 59+ year old females", "", false, false, types.FieldType{Quantity: &types.FieldTypeQuantity{}}),
											newFieldDefinition("Number of 59+ year old males", "", false, false, types.FieldType{Quantity: &types.FieldTypeQuantity{}}),
										},
									},
								},
							},
							{
								Name: "Situation/Needs of Household",
								FieldType: types.FieldType{
									SubForm: &types.FieldTypeSubForm{
										Fields: []*types.FieldDefinition{
											newFieldDefinition("HH's ability to meet the food needs of all its members", "", false, false, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{Options: ugProtectionHHAbilityToMeetNeeds}}),
											newFieldDefinition("What are the main obstacles you face in meeting food needs?", "", false, false, types.FieldType{MultiSelect: &types.FieldTypeMultiSelect{Options: ugProtectionFoodNeedsObstacles}}),
											newFieldDefinition("What are the main obstacles you face in meeting accomodation needs?", "", false, false, types.FieldType{MultiSelect: &types.FieldTypeMultiSelect{Options: ugProtectionAccommodationNeedsObstacles}}),
											newFieldDefinition("Can the HH meet WASH needs?", "", false, false, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{Options: ugProtectionHHAbilityToMeetNeeds}}),
											newFieldDefinition("What are the main obstacles in meeting WASH needs?", "", false, false, types.FieldType{MultiSelect: &types.FieldTypeMultiSelect{Options: ugProtectionWASHNeedsObstacles}}),
											newFieldDefinition("Summary narrative", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
										},
									},
								},
							},
						},
					},
				},
			},
			newFieldDefinition("Main obstacles you face in meeting accomodation needs", "", false, true, types.FieldType{
				SingleSelect: &types.FieldTypeSingleSelect{},
			}),
			newFieldDefinition("Can the HH meet WASH needs?", "", false, true, types.FieldType{
				SingleSelect: &types.FieldTypeSingleSelect{},
			}),
			newFieldDefinition("Main obstacles in meeting WASH needs", "", false, true, types.FieldType{
				SingleSelect: &types.FieldTypeSingleSelect{},
			}),
			newFieldDefinition("Summary narrative", "", false, false, types.FieldType{
				MultilineText: &types.FieldTypeMultilineText{},
			}),
		},
	)

	protectionIncident = newFormDefinition(
		dbID,
		protectionFolder.ID,
		"Protection Incident Form",
		[]*types.FieldDefinition{
			newFieldDefinition("Individual Beneficiary", "The beneficiary this incident form has been completed for", false, true, types.FieldType{
				Reference: &types.FieldTypeReference{
					DatabaseID: s.globalDatabase.ID,
					FormID:     s.globalRootBeneficiaryForm.ID,
				},
			}),
			newFieldDefinition("Location of incident", "", false, true, types.FieldType{Text: &types.FieldTypeText{}}),
			newFieldDefinition("Time of incident", "", false, true, types.FieldType{Text: &types.FieldTypeText{}}),
			newFieldDefinition("Date of incident", "", false, true, types.FieldType{Date: &types.FieldTypeDate{}}),
			newFieldDefinition("Received by", "", false, true, types.FieldType{Text: &types.FieldTypeText{}}),
			newFieldDefinition("Vulnerability", "", false, true, types.FieldType{
				SingleSelect: &types.FieldTypeSingleSelect{},
			}),
			newFieldDefinition("Description of the incident", "i.e. Where, when, what, who involved", false, true, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
			newFieldDefinition("Has the incident been reported to the police?", "", false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
			newFieldDefinition("Comment", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
			newFieldDefinition("Has the incident been reported to", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
		},
	)

	protectionActionReport = newFormDefinition(
		dbID,
		protectionFolder.ID,
		"Protection Action Report Form",
		[]*types.FieldDefinition{
			newFieldDefinition("Individual Beneficiary", "The beneficiary this action report form has been completed for", false, true, types.FieldType{
				Reference: &types.FieldTypeReference{
					DatabaseID: s.globalDatabase.ID,
					FormID:     s.globalRootBeneficiaryForm.ID,
				},
			}),
			newFieldDefinition("Which service has the beneficiary together with staff agreed to take?", "", false, true, types.FieldType{
				SingleSelect: &types.FieldTypeSingleSelect{},
			}),
			newFieldDefinition("If \"Other\", specify", "", false, false, types.FieldType{
				Text: &types.FieldTypeText{},
			}),
			newFieldDefinition("Narrate/Describe the response action agreed upon", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
			newFieldDefinition("Agreed follow-up with beneficiary", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
		})

	protectionFollowupReport = newFormDefinition(
		dbID,
		protectionFolder.ID,
		"Protection Followup Form",
		[]*types.FieldDefinition{
			newFieldDefinition("Individual Beneficiary", "The beneficiary this followup form has been completed for", false, true, types.FieldType{
				Reference: &types.FieldTypeReference{
					DatabaseID: s.globalDatabase.ID,
					FormID:     s.globalRootBeneficiaryForm.ID,
				},
			}),
			newFieldDefinition("Follow up after", "", false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
			newFieldDefinition("Agreed follow-up with the beneficiary", "", false, true, types.FieldType{Text: &types.FieldTypeText{}}),
		},
	)

	protectionReferral = newFormDefinition(
		dbID,
		protectionFolder.ID,
		"Protection Referral Form",
		[]*types.FieldDefinition{
			newFieldDefinition("Individual Beneficiary", "The beneficiary this referral form has been completed for", false, true, types.FieldType{
				Reference: &types.FieldTypeReference{
					DatabaseID: s.globalDatabase.ID,
					FormID:     s.globalRootBeneficiaryForm.ID,
				},
			}),
			newFieldDefinition("Priority", "", false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
			newFieldDefinition("Referred via", "", false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
			newFieldDefinition("Referral date", "", false, true, types.FieldType{Date: &types.FieldTypeDate{}}),
			newFieldDefinition("Receiving agency", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
			newFieldDefinition("Name of partner case worker", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
			newFieldDefinition("Position of person receiving referral", "", false, true, types.FieldType{Text: &types.FieldTypeText{}}),
			newFieldDefinition("Contact of person receiving referral", "", false, true, types.FieldType{Text: &types.FieldTypeText{}}),
			newFieldDefinition("Consent to release information", "", false, false, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
			newFieldDefinition("Has the person expressed any restrictions on referrals?", "", false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
			newFieldDefinition("If yes, specify", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
			newFieldDefinition("Is the beneficiary a minor", "", false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
			newFieldDefinition("Name of the primary caregiver", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
			newFieldDefinition("Relationship to the child", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
			newFieldDefinition("Has the caregiver been informed of the referral?", "", false, false, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
			newFieldDefinition("If not, explain", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
			newFieldDefinition("Reason for referral", "", false, true, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
			newFieldDefinition("Type of referral", "", false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
		},
	)

	for _, form := range forms {
		if err := client.CreateForm(ctx, form, nil); err != nil {
			return err
		}
	}
	return nil
}

func (s *Seed) seedUgandaEiLfsForms(ctx context.Context, client client.Client, dbID string) error {
	var eiLfsFolder types.Folder
	if err := client.CreateFolder(ctx, &types.Folder{
		DatabaseID: dbID,
		Name:       "EI & LFS",
	}, &eiLfsFolder); err != nil {
		return err
	}

	var (
		eiLfsScreening     *types.FormDefinition
		educationScreening *types.FormDefinition
		actionDecision     *types.FormDefinition
	)

	var forms = []*types.FormDefinition{eiLfsScreening, educationScreening, actionDecision}
	eiLfsScreening = newFormDefinition(
		dbID,
		eiLfsFolder.ID,
		"EI & LFS Screening",
		[]*types.FieldDefinition{
			newFieldDefinition("Individual beneficiary", "", false, true, types.FieldType{Reference: &types.FieldTypeReference{DatabaseID: s.globalDatabase.ID, FormID: s.globalRootIndividualForm.ID}}),
			newFieldDefinition("What challenges are you currently facing? How are you coping currently?", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
			newFieldDefinition("What are you doing to improve your livelihood?", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
			newFieldDefinition("What solutions do you suggest to improve this situation?", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
			yesNoQuestion("Are you engaged in any form of livelihood generating activity?"),
			newFieldDefinition("If yes, what income generating activity are you engaged in?", "", false, false, types.FieldType{MultiSelect: &types.FieldTypeMultiSelect{Options: ugEILFSIncomeActivity}}),
			ifOtherPleaseSpecify,
			{
				Name: "Business",
				FieldType: types.FieldType{
					SubForm: &types.FieldTypeSubForm{
						Fields: []*types.FieldDefinition{
							newFieldDefinition("What was the source of your start up capital?", "", false, false, types.FieldType{MultiSelect: &types.FieldTypeMultiSelect{Options: ugEILFSStartupCapitalSource}}),
							ifOtherPleaseSpecify,
							newFieldDefinition("How much was the initial investment amount?", "", false, false, types.FieldType{Quantity: &types.FieldTypeQuantity{}}),
							newFieldDefinition("On average what is the monthly business profit?", "", false, false, types.FieldType{Quantity: &types.FieldTypeQuantity{}}),
							newFieldDefinition("How was the profit utilized?", "", false, false, types.FieldType{MultiSelect: &types.FieldTypeMultiSelect{Options: ugEILFSProfitUse}}),
							ifOtherPleaseSpecify,
							yesNoQuestion("Have you received any training in regards to what you are currently engaged in?"),
							newFieldDefinition("If yes, what kind of training have you received?", "", false, false, types.FieldType{MultiSelect: &types.FieldTypeMultiSelect{Options: ugEILFSTrainingKind}}),
							ifOtherPleaseSpecify,
						},
					},
				},
			},
			{
				Name: "Skills training",
				FieldType: types.FieldType{
					SubForm: &types.FieldTypeSubForm{
						Fields: types.FieldDefinitions{
							yesNoQuestion("Have you had any skills training?"),
							newFieldDefinition("What training course have you attended?", "", false, false, types.FieldType{MultiSelect: &types.FieldTypeMultiSelect{Options: ugEILFSSkillsTraining}}),
							ifOtherPleaseSpecify,
						},
					},
				},
			},
			{
				Name: "Farming/Agriculture",
				FieldType: types.FieldType{
					SubForm: &types.FieldTypeSubForm{
						Fields: types.FieldDefinitions{
							newFieldDefinition("What farming activities are you engaged in?", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
							newFieldDefinition("What challenges have you faced?", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
							yesNoQuestion("Do you have access to markets for you products?"),
							newFieldDefinition("What challenges have you faced in accessing markets?", "", false, false, types.FieldType{MultiSelect: &types.FieldTypeMultiSelect{Options: ugEILFSMarketChallenges}}),
						},
					},
				},
			},
			{
				Name: "Employment",
				FieldType: types.FieldType{
					SubForm: &types.FieldTypeSubForm{
						Fields: types.FieldDefinitions{
							newFieldDefinition("What kind of employment are you engaged in?", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
							newFieldDefinition("What challenges are you facing in you current employment?", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
							newFieldDefinition("What steps have you taken towards addressing the challenges?", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
							newFieldDefinition("How would you want NRC to help you address those challenges?", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						},
					},
				},
			},
			{
				Name: "Group Membership",
				FieldType: types.FieldType{
					SubForm: &types.FieldTypeSubForm{
						Fields: types.FieldDefinitions{
							yesNoQuestion("Do you belong to any business group?"),
							newFieldDefinition("If yes, which group do you belong to?", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
							newFieldDefinition("What is the core business area of the group?", "", false, false, types.FieldType{MultiSelect: &types.FieldTypeMultiSelect{Options: ugEILFSGroupBusinessArea}}),
						},
					},
				},
			},
			{
				Name: "Financial Access",
				FieldType: types.FieldType{
					SubForm: &types.FieldTypeSubForm{
						Fields: types.FieldDefinitions{
							yesNoQuestion("Do you access financial services?"),
							newFieldDefinition("If yes, what type of financial services do you have acces to?", "", false, false, types.FieldType{MultiSelect: &types.FieldTypeMultiSelect{Options: ugEILFSFinancialServices}}),
							newFieldDefinition("From which financial service providers?", "", false, false, types.FieldType{MultiSelect: &types.FieldTypeMultiSelect{Options: ugEILFSFinancialServiceProviders}}),
							ifOtherPleaseSpecify,
						},
					},
				},
			},
		})

	educationScreening = newFormDefinition(
		dbID,
		eiLfsFolder.ID, // FIXME double-check this
		"Education Screening",
		types.FieldDefinitions{
			newFieldDefinition("Individual beneficiary", "", false, true, types.FieldType{Reference: &types.FieldTypeReference{DatabaseID: s.globalDatabase.ID, FormID: s.globalRootBeneficiaryForm.ID}}),
			newFieldDefinition("What/Describe Education issue/concern are you facing?", "", false, true, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
			newFieldDefinition("What action have you taken to solve the problem?", "", false, true, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
			newFieldDefinition("What suggestions do you give for further action?", "", false, true, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
			newFieldDefinition("What is the vulnerability status of the beneficiary?", "", false, true, types.FieldType{MultiSelect: &types.FieldTypeMultiSelect{Options: ugEducationVulnerability}}),
			{
				Name: "Education background",
				FieldType: types.FieldType{
					SubForm: &types.FieldTypeSubForm{
						Fields: types.FieldDefinitions{
							yesNoQuestion("Are you in school?"),
							{
								Name: "If in school",
								FieldType: types.FieldType{
									SubForm: &types.FieldTypeSubForm{
										Fields: types.FieldDefinitions{
											newFieldDefinition("Specify school", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
											newFieldDefinition("Specify level/class", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
										},
									},
								},
							},
							{
								Name: "If not in school",
								FieldType: types.FieldType{
									SubForm: &types.FieldTypeSubForm{
										Fields: types.FieldDefinitions{
											newFieldDefinition("Specify class last attended", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
											newFieldDefinition("Reason why out of school", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
											newFieldDefinition("Last year in school", "", false, false, types.FieldType{Quantity: &types.FieldTypeQuantity{}}),
										},
									},
								},
							},
						},
					},
				},
			},
		})

	actionDecision = newFormDefinition(
		dbID,
		eiLfsFolder.ID, // FIXME stand-in, determine
		"Action decision",
		types.FieldDefinitions{
			newFieldDefinition("Individual beneficiary", "", false, true, types.FieldType{
				Reference: &types.FieldTypeReference{
					DatabaseID: s.globalDatabase.ID,
					FormID:     s.globalRootBeneficiaryForm.ID,
				}}),
			newFieldDefinition("Status", "", false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{Options: ugAction}}),
			{
				Name: "Individual risk assessment",
				FieldType: types.FieldType{
					SubForm: &types.FieldTypeSubForm{
						Fields: types.FieldDefinitions{
							yesNoQuestion("Are there any elements of risks for the safety or well-being of the beneficiary or that of a relative in relation to the suggested course of action?"),
							yesNoQuestion("If yes, are there Protection Risks?"),
							newFieldDefinition("If yes, please indicate what type", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
							newFieldDefinition("Notes/Narrative", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
							newFieldDefinition("If any of the answers were 'yes', discuss with the beneficiary what might be done to avoid or minimise the risks or negative consequences.", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
							newFieldDefinition("Actions agreed upon with the beneficiary", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
							yesNoQuestion("Discuss the pro's and con's of the suggested course of action, including the analysis of risks. Does the beneficiary agree to continue with the case?"),
							newFieldDefinition("Notes/Narrative", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
							yesNoQuestion("Is a Best Interest Determination needed for the case?"),
							newFieldDefinition("If 'yes', refer the case to social services or an appropriate child protection actor.", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
						},
					},
				},
			},
			{
				Name: "If action taken is referral",
				FieldType: types.FieldType{
					SubForm: &types.FieldTypeSubForm{
						Fields: types.FieldDefinitions{
							newFieldDefinition("Type of referral", "", false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{Options: ugReferralType}}),
							newFieldDefinition("If internal, choose CC", "", false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{Options: ugCC}}),
							{
								Name: "If external, specify",
								FieldType: types.FieldType{
									SubForm: &types.FieldTypeSubForm{
										Fields: types.FieldDefinitions{
											newFieldDefinition("Organisation", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
											newFieldDefinition("Contact person", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
											newFieldDefinition("Phone number", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
											newFieldDefinition("Email", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
										},
									},
								},
							},
							newFieldDefinition("Reason for referral", "", false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{Options: ugReferralReason}}),
							yesNoQuestion("Is the beneficiary a minor?"),
							newFieldDefinition("Name of the primary care giver", "", false, false, types.FieldType{
								Reference: &types.FieldTypeReference{
									DatabaseID: s.globalDatabase.ID,
									FormID:     s.globalRootBeneficiaryForm.ID,
								},
							}),
							newFieldDefinition("If not a registered beneficiary, specify name", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
							newFieldDefinition("Relationship to the child", "", false, false, types.FieldType{Text: &types.FieldTypeText{}}),
							newFieldDefinition("Means of referral", "", false, false, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{Options: ugReferralMeans}}),
							ifOtherPleaseSpecify,
						},
					},
				},
			},
			{
				Name: "Follow-up",
				FieldType: types.FieldType{
					SubForm: &types.FieldTypeSubForm{
						Fields: types.FieldDefinitions{
							newFieldDefinition("Agreed follow-up means", "", false, false, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{Options: ugAgreedFollowupMeans}}),
							ifOtherPleaseSpecify,
							newFieldDefinition("Notes from follow-up", "", false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						},
					},
				},
			},
			newFieldDefinition("Priority level", "", false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{Options: ugReferralPriority}}),
		},
	)

	for _, form := range forms {
		if err := client.CreateForm(ctx, form, nil); err != nil {
			return err
		}
	}
	return nil
}
