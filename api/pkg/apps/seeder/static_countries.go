package seeder

import "github.com/nrc-no/core/pkg/apps/iam"

var (
	// Countries
	ugandaCountry   = country(iam.UgandaCountry.ID, iam.UgandaCountry.Name)
	colombiaCountry = country(iam.ColombiaCountry.ID, iam.ColombiaCountry.Name)
)