package iam_test

import (
	. "github.com/nrc-no/core/pkg/iam"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestParty() {
	s.Run("API", func() { s.testPartyAPI() })
	s.SetupTest()
	s.Run("List filtering", func() { s.testPartyListFilter() })
	s.SetupTest()
	s.Run("Search", func() { s.testPartySearch() })
}

func (s *Suite) testPartyAPI() {
	party := mockParties(1)[0]
	party.PartyTypeIDs = []string{IndividualPartyType.ID}
	attribute := mockPartyAttributeDefinitions(1)[0]
	party.Attributes.Set(attribute.ID, "mock")

	// CREATE
	created, err := s.client.Parties().Create(s.Ctx, party)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	party.ID = created.ID
	assert.Equal(s.T(), party.PartyTypeIDs, created.PartyTypeIDs)
	assert.Equal(s.T(), party.Get(attribute.ID), created.Get(attribute.ID))

	// GET
	get, err := s.client.Parties().Get(s.Ctx, created.ID)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Equal(s.T(), created, get)

	// UPDATE
	party.PartyTypeIDs = []string{newUUID()}
	party.Attributes.Set(attribute.ID, "update")
	updated, err := s.client.Parties().Update(s.Ctx, party)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Equal(s.T(), party.ID, updated.ID)
	assert.Equal(s.T(), party.PartyTypeIDs, updated.PartyTypeIDs)
	assert.Equal(s.T(), party.Get(attribute.ID), updated.Get(attribute.ID))

	// GET
	get, err = s.client.Parties().Get(s.Ctx, updated.ID)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Equal(s.T(), updated, get)

	// LIST
	list, err := s.client.Parties().List(s.Ctx, PartyListOptions{})
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Contains(s.T(), list.Items, get)
}

func (s *Suite) createParties(nParties int, nPartyTypes int, nAttributes int) ([]*Party, []string, []string) {
	// Make mock structs
	parties := mockParties(nParties)
	var partyTypeIds []string
	for i := 0; i < nPartyTypes; i++ {
		partyTypeIds = append(partyTypeIds, newUUID())
	}
	var attributes []string
	for i := 0; i < nAttributes; i++ {
		attributes = append(attributes, newUUID())
	}

	// Create parties
	for i, party := range parties {
		// Set PartyTypes
		n := 1 + i%len(partyTypeIds)
		pts := partyTypeIds[0:n]
		party.PartyTypeIDs = pts

		// Set Attributes
		m := 1 + i%len(attributes)
		attrs := attributes[0:m]
		for _, attributeId := range attrs {
			party.Attributes.Set(attributeId, "mock")
		}

		// Add to map
		// Save the party to the DB
		created, err := s.client.Parties().Create(s.Ctx, party)
		if !assert.NoError(s.T(), err) {
			s.T().FailNow()
		}
		party.ID = created.ID
	}
	return parties, partyTypeIds, attributes
}

func (s *Suite) testPartyListFilter() {

	parties, partyTypeIds, attributes := s.createParties(4, 2, 2)

	s.Run("by type", func() { s.testPartyListFilterByType(parties, partyTypeIds) })
	s.Run("by attribute", func() { s.testPartyListFilterByAttribute(parties, attributes) })
	s.Run("by type and attribute", func() { s.testPartyListFilterByTypeAndAttribute(parties, partyTypeIds, attributes) })
}

func (s *Suite) testPartyListFilterByType(parties []*Party, partyTypeIds []string) {
	for _, partyTypeId := range partyTypeIds {
		list, err := s.client.Parties().List(s.Ctx, PartyListOptions{
			PartyTypeID: partyTypeId,
		})
		if !assert.NoError(s.T(), err) {
			s.T().FailNow()
		}
		// Get expected items
		var expected []string
		for _, party := range parties {
			if party.HasPartyType(partyTypeId) {
				expected = append(expected, party.ID)
			}
		}
		assert.Len(s.T(), list.Items, len(expected))
		for _, item := range list.Items {
			assert.Contains(s.T(), expected, item.ID)
		}
	}
}

func (s *Suite) testPartyListFilterByAttribute(parties []*Party, attributes []string) {
	for i := 1; i <= len(attributes); i++ {
		attrs := attributes[0:i]
		attributeOptions := make(map[string]string)
		for _, attributeId := range attrs {
			attributeOptions[attributeId] = "mock"
		}
		list, err := s.client.Parties().List(s.Ctx, PartyListOptions{
			Attributes: attributeOptions,
		})
		if !assert.NoError(s.T(), err) {
			s.T().FailNow()
		}

		// Get expected items
		var expected []string
		for _, party := range parties {
			var include = true
			for attributeId := range attributeOptions {
				if !party.HasAttribute(attributeId) {
					include = false
					break
				}
			}
			if include {
				expected = append(expected, party.ID)
			}
		}

		assert.Len(s.T(), list.Items, len(expected))
		for _, item := range list.Items {
			assert.Contains(s.T(), expected, item.ID)
		}
	}
}

func (s *Suite) testPartyListFilterByTypeAndAttribute(parties []*Party, partyTypeIds, attributes []string) {
	for i := 1; i <= len(attributes); i++ {
		for _, partyTypeId := range partyTypeIds {
			attributeOptions := make(map[string]string)
			for _, attributeId := range attributes[0:i] {
				attributeOptions[attributeId] = "mock"
			}
			list, err := s.client.Parties().List(s.Ctx, PartyListOptions{
				PartyTypeID: partyTypeId,
				Attributes:  attributeOptions,
			})
			if !assert.NoError(s.T(), err) {
				s.T().FailNow()
			}
			// Get expected items
			var expected []string
			for _, party := range parties {
				if !party.HasPartyType(partyTypeId) {
					continue
				}
				var include = true
				for id := range attributeOptions {
					if !party.HasAttribute(id) {
						include = false
						break
					}
				}
				if include {
					expected = append(expected, party.ID)
				}
			}

			// Compare list lengths
			assert.Len(s.T(), expected, len(list.Items))
			// ... and contents
			for _, item := range list.Items {
				assert.Contains(s.T(), expected, item.ID)
			}
		}
	}
}

func (s *Suite) testPartySearch() {
	parties, partyTypeIds, attributes := s.createParties(4, 2, 2)

	//Search by partyId
	for _, party := range parties {
		list, err := s.client.Parties().Search(s.Ctx, PartySearchOptions{
			PartyIDs: []string{party.ID},
		})
		if !assert.NoError(s.T(), err) {
			s.T().FailNow()
		}
		assert.Len(s.T(), list.Items, 1)
		assert.Equal(s.T(), party.ID, list.Items[0].ID)
	}

	//Search by multiple partyIds
	var partyIds []string
	partyIds = append(partyIds, parties[0].ID)
	partyIds = append(partyIds, parties[2].ID)
	list, err := s.client.Parties().Search(s.Ctx, PartySearchOptions{
		PartyIDs: partyIds,
	})
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	assert.Len(s.T(), list.Items, len(partyIds))
	for i, item := range list.Items {
		assert.Equal(s.T(), partyIds[i], item.ID)
	}

	//Search by partyTypeId
	for _, partyTypeId := range partyTypeIds {
		list, err := s.client.Parties().Search(s.Ctx, PartySearchOptions{
			PartyTypeIDs: []string{partyTypeId},
		})
		if !assert.NoError(s.T(), err) {
			s.T().FailNow()
		}
		// Get expected items
		var expected []string
		for _, party := range parties {
			if party.HasPartyType(partyTypeId) {
				expected = append(expected, party.ID)
			}
		}
		assert.Len(s.T(), list.Items, len(expected))
		for _, item := range list.Items {
			assert.Contains(s.T(), expected, item.ID)
		}
	}

	//Search by attributes
	for i := 1; i <= len(attributes); i++ {
		attrs := attributes[0:i]
		attributeOptions := make(map[string]string)
		for _, attributeId := range attrs {
			attributeOptions[attributeId] = "mock"
		}
		list, err := s.client.Parties().Search(s.Ctx, PartySearchOptions{
			Attributes: attributeOptions,
		})
		if !assert.NoError(s.T(), err) {
			s.T().FailNow()
		}

		// Get expected items
		var expected []string
		for _, party := range parties {
			var include = true
			for attributeId := range attributeOptions {
				if !party.HasAttribute(attributeId) {
					include = false
					break
				}
			}
			if include {
				expected = append(expected, party.ID)
			}
		}

		assert.Len(s.T(), list.Items, len(expected))
		for _, item := range list.Items {
			assert.Contains(s.T(), expected, item.ID)
		}
	}
}