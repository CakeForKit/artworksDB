package integration

import (
	"context"
	"errors"
	"testing"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/adminrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/artworkrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/authorrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/collectionrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/employeerep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/eventrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/fixtures"
	fixturesrep "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/fixtures/fixtures_rep"
	testobj "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/testObj"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type EventRepSuite struct {
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

func TestEventRepSuite(t *testing.T) {
	suite.RunSuite(t, new(EventRepSuite))
}

func (s *EventRepSuite) BeforeAll(t provider.T) {
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

func (s *EventRepSuite) BeforeEach(t provider.T) {
	t.Tags("integration", "event")
}

func (s *EventRepSuite) AfterAll(t provider.T) {
	if s.eventRep != nil {
		s.eventRep.Close()
	}
}

func (s *EventRepSuite) TestEventRep_GetAllEvents(t provider.T) {
	t.Parallel()

	t.Run("Success with empty filter", func(t provider.T) {
		events := []*models.Event{
			s.eventCreator.EventP(uuid.New()),
			s.eventCreator.EventP(uuid.New()),
		}
		fixturesrep.AddTestEvents(t, s.ctx, events, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)
		defer fixturesrep.DelTestEvents(t, s.ctx, events, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)

		filterOps := s.eventCreator.EventFilterEmpty()

		// ACT
		resEvents, err := s.eventRep.GetAll(s.ctx, filterOps)

		t.Require().NoError(err)
		fixturesrep.AssertEventsAreInRes(t, events, resEvents)
	})
	t.Run("Success with title filter", func(t provider.T) {
		event1 := s.eventCreator.EventP(uuid.New())
		event2 := s.eventCreator.EventP(uuid.New())
		events := []*models.Event{event1, event2}
		fixturesrep.AddTestEvents(t, s.ctx, events, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)
		defer fixturesrep.DelTestEvents(t, s.ctx, events, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)

		filterOps := &jsonreqresp.EventFilter{
			Title:     event1.GetTitle(),
			DateBegin: time.Time{},
			DateEnd:   time.Time{},
		}

		// ACT
		resEvents, err := s.eventRep.GetAll(s.ctx, filterOps)

		t.Require().NoError(err)
		fixturesrep.AssertEventsAreInRes(t, []*models.Event{event1}, resEvents)
		for _, v := range resEvents {
			t.Assert().True(v.GetTitle() == filterOps.Title)
		}
	})

	t.Run("Success with date range filter", func(t provider.T) {
		baseTime := time.Now().UTC()
		event1 := s.eventCreator.EventWithDatesP(uuid.New(), baseTime.AddDate(0, -1, 0), baseTime.AddDate(0, 1, 0))
		event2 := s.eventCreator.EventWithDatesP(uuid.New(), baseTime.AddDate(0, 2, 0), baseTime.AddDate(0, 3, 0))
		events := []*models.Event{event1, event2}
		fixturesrep.AddTestEvents(t, s.ctx, events, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)
		defer fixturesrep.DelTestEvents(t, s.ctx, events, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)

		filterOps := &jsonreqresp.EventFilter{
			Title:     "",
			DateBegin: baseTime.AddDate(0, -2, 0),
			DateEnd:   baseTime.AddDate(0, 1, 0),
		}

		// ACT
		resEvents, err := s.eventRep.GetAll(s.ctx, filterOps)

		t.Require().NoError(err)
		fixturesrep.AssertEventsAreInRes(t, []*models.Event{event1}, resEvents)
	})
}

func (s *EventRepSuite) TestEventRep_GetArtworkIDs(t provider.T) {
	t.Parallel()

	t.Run("Success", func(t provider.T) {
		artwork1 := s.artworkCreator.ArtworkP(uuid.New())
		artwork2 := s.artworkCreator.ArtworkP(uuid.New())
		artworkIDs := uuid.UUIDs{artwork1.GetID(), artwork2.GetID()}
		event := s.eventCreator.EventWithArtworksP(uuid.New(), artworkIDs)

		fixturesrep.AddTestEvents(t, s.ctx, []*models.Event{event}, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)
		defer fixturesrep.DelTestEvents(t, s.ctx, []*models.Event{event}, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)

		// ACT
		resArtworkIDs, err := s.eventRep.GetArtworkIDs(s.ctx, event.GetID())

		t.Require().NoError(err)
		t.Require().Equal(len(artworkIDs), len(resArtworkIDs))
		for _, expectedID := range artworkIDs {
			t.Assert().True(containsUUID(resArtworkIDs, expectedID))
		}
	})

	t.Run("Empty result for event without artworks", func(t provider.T) {
		event := s.eventCreator.EventWithArtworksP(uuid.New(), uuid.UUIDs{})
		fixturesrep.AddTestEvents(t, s.ctx, []*models.Event{event}, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)
		defer fixturesrep.DelTestEvents(t, s.ctx, []*models.Event{event}, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)

		// ACT
		resArtworkIDs, err := s.eventRep.GetArtworkIDs(s.ctx, event.GetID())

		t.Require().NoError(err)
		t.Require().Empty(resArtworkIDs)
	})

	t.Run("Not found", func(t provider.T) {
		// ACT
		_, err := s.eventRep.GetArtworkIDs(s.ctx, uuid.New())

		t.Require().ErrorIs(err, eventrep.ErrEventNotFound)
	})
}

func (s *EventRepSuite) TestEventRep_GetByID(t provider.T) {
	t.Parallel()

	t.Run("Success", func(t provider.T) {
		events := []*models.Event{
			s.eventCreator.EventP(uuid.New()),
		}
		fixturesrep.AddTestEvents(t, s.ctx, events, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)
		defer fixturesrep.DelTestEvents(t, s.ctx, events, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)

		// ACT
		resEvent, err := s.eventRep.GetByID(s.ctx, events[0].GetID())

		t.Require().NoError(err)
		t.Require().True(resEvent.Equal(events[0]))
	})

	t.Run("Not found", func(t provider.T) {
		// ACT
		_, err := s.eventRep.GetByID(s.ctx, uuid.New())

		t.Require().ErrorIs(err, eventrep.ErrEventNotFound)
	})
}

func (s *EventRepSuite) TestEventRep_GetEventsOfArtworkOnDate(t provider.T) {
	t.Parallel()

	t.Run("Success", func(t provider.T) {
		artwork := s.artworkCreator.ArtworkP(uuid.New())
		baseTime := time.Now().UTC()
		event1 := s.eventCreator.EventWithDateArtworksP(uuid.New(), baseTime.AddDate(0, -1, 0), baseTime.AddDate(0, 1, 0), uuid.UUIDs{artwork.GetID()})
		event2 := s.eventCreator.EventWithDateArtworksP(uuid.New(), baseTime.AddDate(0, 2, 0), baseTime.AddDate(0, 3, 0), uuid.UUIDs{artwork.GetID()})

		fixturesrep.AddTestEvents(t, s.ctx, []*models.Event{event1, event2}, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)
		defer fixturesrep.DelTestEvents(t, s.ctx, []*models.Event{event1, event2}, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)

		dateBeg := baseTime.AddDate(0, -2, 0)
		dateEnd := baseTime.AddDate(0, 2, 0)

		// ACT
		resEvents, err := s.eventRep.GetEventsOfArtworkOnDate(s.ctx, artwork.GetID(), dateBeg, dateEnd)

		t.Require().NoError(err)
		fixturesrep.AssertEventsAreInRes(t, []*models.Event{event1}, resEvents)
	})

	t.Run("Empty result for artwork without events in date range", func(t provider.T) {
		artwork := s.artworkCreator.ArtworkP(uuid.New())
		fixturesrep.AddTestArtworks(t, s.ctx, []*models.Artwork{artwork}, s.artworkRep, s.authorRep, s.collectionRep)
		defer fixturesrep.DelTestArtworks(t, s.ctx, []*models.Artwork{artwork}, s.artworkRep, s.authorRep, s.collectionRep)

		dateBeg := time.Now().AddDate(0, -1, 0)
		dateEnd := time.Now().AddDate(0, 1, 0)

		// ACT
		resEvents, err := s.eventRep.GetEventsOfArtworkOnDate(s.ctx, artwork.GetID(), dateBeg, dateEnd)

		t.Require().NoError(err)
		t.Require().Empty(resEvents)
	})
}

func (s *EventRepSuite) TestEventRep_GetCollectionsStat(t provider.T) {
	t.Parallel()

	t.Run("Success", func(t provider.T) {
		collection1 := s.collectionCreator.CollectionP(uuid.New())
		collection2 := s.collectionCreator.CollectionP(uuid.New())
		artworks := []*models.Artwork{
			s.artworkCreator.ArtworkWithCollectionP(uuid.New(), collection1),
			s.artworkCreator.ArtworkWithCollectionP(uuid.New(), collection1),
			s.artworkCreator.ArtworkWithCollectionP(uuid.New(), collection2),
		}
		artworkIDs := uuid.UUIDs{}
		for _, v := range artworks {
			artworkIDs = append(artworkIDs, v.GetID())
		}
		event := s.eventCreator.EventWithArtworksP(uuid.New(), artworkIDs)

		fixturesrep.AddTestEventsWithArtworks(t, s.ctx, []*models.Event{event}, artworks, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)
		defer fixturesrep.DelTestEvents(t, s.ctx, []*models.Event{event}, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)

		// ACT
		stats, err := s.eventRep.GetCollectionsStat(s.ctx, event.GetID())

		t.Require().NoError(err)
		t.Require().Len(stats, 2)

		checked := 0
		for _, stat := range stats {
			if stat.CollectionID() == collection1.GetID() {
				t.Assert().Equal(2, stat.ArtworksCount())
				checked += 1
			} else if stat.CollectionID() == collection2.GetID() {
				t.Assert().Equal(1, stat.ArtworksCount())
				checked += 1
			}
		}
		t.Assert().Equal(checked, 2)
	})

	t.Run("Empty result for event without artworks", func(t provider.T) {
		event := s.eventCreator.EventWithArtworksP(uuid.New(), uuid.UUIDs{})
		fixturesrep.AddTestEventsWithArtworks(t, s.ctx, []*models.Event{event}, []*models.Artwork{}, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)
		defer fixturesrep.DelTestEvents(t, s.ctx, []*models.Event{event}, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)

		// ACT
		stats, err := s.eventRep.GetCollectionsStat(s.ctx, event.GetID())

		t.Require().NoError(err)
		t.Require().Empty(stats)
	})
}

func (s *EventRepSuite) TestEventRep_CheckEmployeeByID(t provider.T) {
	t.Parallel()

	t.Run("Employee exists", func(t provider.T) {
		employee := s.employeeCreator.DefaultEmployeeP(uuid.New(), uuid.New())
		fixturesrep.AddTestEmployees(t, s.ctx, []*models.Employee{employee}, s.employeeRep, s.adminRep)
		defer fixturesrep.DelTestEmployees(t, s.ctx, []*models.Employee{employee}, s.employeeRep, s.adminRep)

		// ACT
		exists, err := s.eventRep.CheckEmployeeByID(s.ctx, employee.GetID())

		t.Require().NoError(err)
		t.Require().True(exists)
	})

	t.Run("Employee does not exist", func(t provider.T) {
		// ACT
		exists, err := s.eventRep.CheckEmployeeByID(s.ctx, uuid.New())

		t.Require().NoError(err)
		t.Require().False(exists)
	})
}

func (s *EventRepSuite) TestEventRep_Add(t provider.T) {
	t.Parallel()

	t.Run("Duplicate error", func(t provider.T) {
		event := s.eventCreator.EventP(uuid.New())
		fixturesrep.AddTestEvents(t, s.ctx, []*models.Event{event}, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)
		defer fixturesrep.DelTestEvents(t, s.ctx, []*models.Event{event}, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)

		// ACT - try to add event with same ID
		err := s.eventRep.Add(s.ctx, event)

		t.Require().Error(err)
	})

	t.Run("Employee not found error", func(t provider.T) {
		event := s.eventCreator.EventP(uuid.New())
		// Don't add employee to repository

		// ACT
		err := s.eventRep.Add(s.ctx, event)

		t.Require().Error(err)
	})
}

func (s *EventRepSuite) TestEventRep_SetNotValid(t provider.T) {
	t.Parallel()

	t.Run("Success", func(t provider.T) {
		event := s.eventCreator.EventP(uuid.New())
		fixturesrep.AddTestEvents(t, s.ctx, []*models.Event{event}, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)
		defer fixturesrep.DelTestEvents(t, s.ctx, []*models.Event{event}, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)

		// ACT
		err := s.eventRep.SetNotValid(s.ctx, event.GetID())

		t.Require().NoError(err)

		// Verify event was marked as not valid
		resEvent, err := s.eventRep.GetByID(s.ctx, event.GetID())
		t.Require().NoError(err)
		t.Assert().False(resEvent.IsValid())
	})

	t.Run("Not found", func(t provider.T) {
		// ACT
		err := s.eventRep.SetNotValid(s.ctx, uuid.New())

		t.Require().Error(err)
	})
}

func (s *EventRepSuite) TestEventRep_Delete(t provider.T) {
	t.Parallel()

	t.Run("Success", func(t provider.T) {
		event := s.eventCreator.EventP(uuid.New())
		fixturesrep.AddTestEvents(t, s.ctx, []*models.Event{event}, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)

		// ACT
		err := s.eventRep.Delete(s.ctx, event.GetID())

		t.Require().NoError(err)

		// Verify event was deleted
		_, err = s.eventRep.GetByID(s.ctx, event.GetID())
		t.Require().ErrorIs(err, eventrep.ErrEventNotFound)
	})

	t.Run("Not found", func(t provider.T) {
		// ACT
		err := s.eventRep.Delete(s.ctx, uuid.New())

		t.Require().Error(err)
	})
}

func (s *EventRepSuite) TestEventRep_Update(t provider.T) {
	t.Parallel()

	t.Run("Success", func(t provider.T) {
		event := s.eventCreator.EventP(uuid.New())
		fixturesrep.AddTestEvents(t, s.ctx, []*models.Event{event}, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)
		defer fixturesrep.DelTestEvents(t, s.ctx, []*models.Event{event}, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)

		newEvent, err := models.NewEvent(
			event.GetID(),
			"Updated Event Title",
			event.GetDateBegin().AddDate(0, 0, 1),
			event.GetDateEnd().AddDate(0, 0, -1),
			"Updated Event Adress",
			!event.GetAccess(),
			event.GetEmployeeID(),
			event.GetTicketCount()+10,
			event.IsValid(),
			event.GetArtworkIDs(),
		)
		t.Require().NoError(err)
		funcUpdate := func(e *models.Event) (*models.Event, error) {
			return &newEvent, nil
		}

		// ACT
		err = s.eventRep.Update(s.ctx, event.GetID(), funcUpdate)

		t.Require().NoError(err)

		updatedEvent, err := s.eventRep.GetByID(s.ctx, event.GetID())
		t.Require().NoError(err)
		t.Require().True(updatedEvent.Equal(&newEvent))
	})

	t.Run("Not found", func(t provider.T) {
		funcUpdate := func(e *models.Event) (*models.Event, error) {
			return e, nil
		}
		// ACT
		err := s.eventRep.Update(s.ctx, uuid.New(), funcUpdate)

		t.Require().ErrorIs(err, eventrep.ErrEventNotFound)
	})

	t.Run("Update function returns error", func(t provider.T) {
		event := s.eventCreator.EventP(uuid.New())
		fixturesrep.AddTestEvents(t, s.ctx, []*models.Event{event}, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)
		defer fixturesrep.DelTestEvents(t, s.ctx, []*models.Event{event}, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)

		expectedErr := errors.New("update function error")
		funcUpdate := func(e *models.Event) (*models.Event, error) {
			return nil, expectedErr
		}

		// ACT
		err := s.eventRep.Update(s.ctx, event.GetID(), funcUpdate)

		t.Require().ErrorIs(err, expectedErr)
	})
}

func (s *EventRepSuite) TestEventRep_AddArtworksToEvent(t provider.T) {
	t.Parallel()

	t.Run("Success", func(t provider.T) {
		event := s.eventCreator.EventWithArtworksP(uuid.New(), uuid.UUIDs{})
		fixturesrep.AddTestEventsWithArtworks(t, s.ctx, []*models.Event{event}, []*models.Artwork{}, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)
		defer fixturesrep.DelTestEvents(t, s.ctx, []*models.Event{event}, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)

		artwork1 := s.artworkCreator.ArtworkP(uuid.New())
		artwork2 := s.artworkCreator.ArtworkP(uuid.New())
		artworkIDs := uuid.UUIDs{artwork1.GetID(), artwork2.GetID()}
		fixturesrep.AddTestArtworks(t, s.ctx, []*models.Artwork{artwork1, artwork2}, s.artworkRep, s.authorRep, s.collectionRep)

		// ACT
		err := s.eventRep.AddArtworksToEvent(s.ctx, event.GetID(), artworkIDs)

		t.Require().NoError(err)

		// Verify artworks were added
		resArtworkIDs, err := s.eventRep.GetArtworkIDs(s.ctx, event.GetID())
		t.Require().NoError(err)
		t.Require().Equal(len(artworkIDs), len(resArtworkIDs))
		for _, expectedID := range artworkIDs {
			t.Assert().True(containsUUID(resArtworkIDs, expectedID))
		}
	})

	t.Run("Duplicate artwork error", func(t provider.T) {
		artwork := s.artworkCreator.ArtworkP(uuid.New())
		event := s.eventCreator.EventP(uuid.New())

		fixturesrep.AddTestArtworks(t, s.ctx, []*models.Artwork{artwork}, s.artworkRep, s.authorRep, s.collectionRep)
		fixturesrep.AddTestEvents(t, s.ctx, []*models.Event{event}, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)
		defer fixturesrep.DelTestEvents(t, s.ctx, []*models.Event{event}, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)
		defer fixturesrep.DelTestArtworks(t, s.ctx, []*models.Artwork{artwork}, s.artworkRep, s.authorRep, s.collectionRep)

		// Add artwork first time
		err := s.eventRep.AddArtworksToEvent(s.ctx, event.GetID(), uuid.UUIDs{artwork.GetID()})
		t.Require().NoError(err)

		// ACT - try to add same artwork again
		err = s.eventRep.AddArtworksToEvent(s.ctx, event.GetID(), uuid.UUIDs{artwork.GetID()})

		t.Require().Error(err)
	})

	t.Run("Event not found", func(t provider.T) {
		artwork := s.artworkCreator.ArtworkP(uuid.New())
		fixturesrep.AddTestArtworks(t, s.ctx, []*models.Artwork{artwork}, s.artworkRep, s.authorRep, s.collectionRep)
		defer fixturesrep.DelTestArtworks(t, s.ctx, []*models.Artwork{artwork}, s.artworkRep, s.authorRep, s.collectionRep)

		// ACT
		err := s.eventRep.AddArtworksToEvent(s.ctx, uuid.New(), uuid.UUIDs{artwork.GetID()})

		t.Require().Error(err)
	})

	t.Run("Artwork not found", func(t provider.T) {
		event := s.eventCreator.EventP(uuid.New())
		fixturesrep.AddTestEvents(t, s.ctx, []*models.Event{event}, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)
		defer fixturesrep.DelTestEvents(t, s.ctx, []*models.Event{event}, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)

		// ACT
		err := s.eventRep.AddArtworksToEvent(s.ctx, event.GetID(), uuid.UUIDs{uuid.New()})

		t.Require().Error(err)
	})
}

func (s *EventRepSuite) TestEventRep_DeleteArtworkFromEvent(t provider.T) {
	t.Parallel()

	t.Run("Success", func(t provider.T) {
		artworks := []*models.Artwork{
			s.artworkCreator.ArtworkP(uuid.New()),
			s.artworkCreator.ArtworkP(uuid.New()),
		}
		artworkIDs := uuid.UUIDs{}
		for _, v := range artworks {
			artworkIDs = append(artworkIDs, v.GetID())
		}
		event := s.eventCreator.EventWithArtworksP(uuid.New(), artworkIDs)
		fixturesrep.AddTestEventsWithArtworks(t, s.ctx, []*models.Event{event}, artworks, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)
		err := event.DeleteArtwork(artworkIDs[0])
		t.Require().NoError(err)
		defer fixturesrep.DelTestEvents(t, s.ctx, []*models.Event{event}, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)

		// ACT
		err = s.eventRep.DeleteArtworkFromEvent(s.ctx, event.GetID(), artworks[0].GetID())

		t.Require().NoError(err)

		// Verify artwork was removed
		resArtworkIDs, err := s.eventRep.GetArtworkIDs(s.ctx, event.GetID())
		t.Require().NoError(err)
		t.Require().Len(resArtworkIDs, 1)
		t.Require().True(resArtworkIDs[0] == artworks[1].GetID())
	})
	/*
		t.Run("Artwork not in event", func(t provider.T) {
			event := s.eventCreator.EventWithArtworksP(uuid.New(), uuid.UUIDs{})

			fixturesrep.AddTestEventsWithArtworks(t, s.ctx, []*models.Event{event}, []*models.Artwork{}, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)
			defer fixturesrep.DelTestEvents(t, s.ctx, []*models.Event{event}, s.eventRep, s.employeeRep, s.adminRep, s.artworkRep, s.authorRep, s.collectionRep)

			// ACT - try to delete artwork that was never added to event
			err := s.eventRep.DeleteArtworkFromEvent(s.ctx, event.GetID(), uuid.New())

			t.Require().Error(err)
		})

		t.Run("Event not found", func(t provider.T) {
			artwork := s.artworkCreator.ArtworkP(uuid.New())
			fixturesrep.AddTestArtworks(t, s.ctx, []*models.Artwork{artwork}, s.artworkRep, s.authorRep, s.collectionRep)
			defer fixturesrep.DelTestArtworks(t, s.ctx, []*models.Artwork{artwork}, s.artworkRep, s.authorRep, s.collectionRep)

			// ACT
			err := s.eventRep.DeleteArtworkFromEvent(s.ctx, uuid.New(), artwork.GetID())

			t.Require().Error(err)
		})
	*/
}

func containsUUID(uuids uuid.UUIDs, target uuid.UUID) bool {
	for _, u := range uuids {
		if u == target {
			return true
		}
	}
	return false
}
