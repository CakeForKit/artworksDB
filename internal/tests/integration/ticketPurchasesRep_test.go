package integration

/*
type TicketPurchasesRepSuite struct {
	fixtures.BaseIntegrationSuite
	ctx               context.Context
	eventCreator      testobj.EventMother
	artworkCreator    testobj.ArtworkMother
	collectionCreator testobj.CollectionMother
	employeeCreator   testobj.EmployeeMother
	eventRep          eventrep.EventRep
	employeeRep       employeerep.EmployeeRep
	adminRep          adminrep.AdminRep
	artworkRep        artworkrep.ArtworkRep
	authorRep         authorrep.AuthorRep
	collectionRep     collectionrep.CollectionRep
}

func TestTicketPurchasesRepSuite(t *testing.T) {
	suite.RunSuite(t, new(TicketPurchasesRepSuite))
}

func (s *TicketPurchasesRepSuite) BeforeAll(t provider.T) {
	s.BaseIntegrationSuite.BeforeAll(t)
	s.ctx = context.Background()
	s.eventCreator = testobj.NewEventMother()
	s.artworkCreator = testobj.NewArtworkMother()
	s.collectionCreator = testobj.NewCollectionMother()
	s.employeeCreator = testobj.NewEmployeeMother()

	t.WithNewStep("Create repositories", func(sCtx provider.StepCtx) {
		var err error = nil
		s.eventRep, err = eventrep.NewEventRep(s.ctx, cnfg.PostgresDB, s.DBCreds, s.DBCnfg)
		sCtx.Require().NoError(err)
		s.artworkRep, err = artworkrep.NewArtworkRep(s.ctx, cnfg.PostgresDB, s.DBCreds, s.DBCnfg)
		sCtx.Require().NoError(err)
		s.authorRep, err = authorrep.NewAuthorRep(s.ctx, cnfg.PostgresDB, s.DBCreds, s.DBCnfg)
		sCtx.Require().NoError(err)
		s.collectionRep, err = collectionrep.NewCollectionRep(s.ctx, cnfg.PostgresDB, s.DBCreds, s.DBCnfg)
		sCtx.Require().NoError(err)
		s.employeeRep, err = employeerep.NewEmployeeRep(s.ctx, cnfg.PostgresDB, s.DBCreds, s.DBCnfg)
		sCtx.Require().NoError(err)
		s.adminRep, err = adminrep.NewAdminRep(s.ctx, cnfg.PostgresDB, s.DBCreds, s.DBCnfg)
		sCtx.Require().NoError(err)
	})
}

func (s *TicketPurchasesRepSuite) BeforeEach(t provider.T) {
	t.Tags("integration", "event")
}

func (s *TicketPurchasesRepSuite) AfterAll(t provider.T) {
	if s.eventRep != nil {
		s.eventRep.Close()
	}
}

func (s *TicketPurchasesRepSuite) TestTicketPurchasesRep_GetAllTicketPurchasess(t provider.T) {
	t.Parallel()

	t.Run("Success with empty filter", func(t provider.T) {
		events := []*models.Event{
			s.eventCreator.EventP(uuid.New()),
			s.eventCreator.EventP(uuid.New()),
		}
		fixturesrep.AddTestEvents(t, s.ctx, events, s.eventRep, s.employeeRep, s.adminRep,
		s.artworkRep, s.authorRep, s.collectionRep)
		defer fixturesrep.DelTestEvents(t, s.ctx, events, s.eventRep, s.employeeRep, s.adminRep,
		s.artworkRep, s.authorRep, s.collectionRep)

		filterOps := s.eventCreator.EventFilterEmpty()

		// ACT
		resEvents, err := s.eventRep.GetAll(s.ctx, filterOps)

		t.Require().NoError(err)
		fixturesrep.AssertEventsAreInRes(t, events, resEvents)
	})

}
*/
