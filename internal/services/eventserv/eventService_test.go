package eventserv_test

import (
	"context"
	"errors"
	"testing"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/artworkrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/eventrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/eventserv"
	testobj "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/testObj"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/mock"
)

type EventServiceSuite struct {
	suite.Suite
}

func TestEventService(t *testing.T) {
	suite.RunSuite(t, new(EventServiceSuite))
}

func (s *EventServiceSuite) TestEventServ_GetAll(t provider.T) {
	eventCreator := testobj.NewEventMother()

	t.WithNewStep("success return 2 events", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		events := []*models.Event{
			eventCreator.EventP(uuid.New()),
			eventCreator.EventP(uuid.New()),
		}

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)
		mockEventRep.On(
			"GetAll", ctx,
			mock.AnythingOfType("*jsonreqresp.EventFilter"),
		).Return(events, nil)

		eventServ := eventserv.NewEventService(mockEventRep, mockArtworkRep)

		// ACT
		resEvents, err := eventServ.GetAll(ctx)

		sCtx.Require().NoError(err)
		sCtx.Require().True(len(events) == len(resEvents))
		for i := range len(resEvents) {
			sCtx.Assert().True(events[i].Equals(resEvents[i]))
		}
		mockEventRep.AssertCalled(t, "GetAll", ctx,
			mock.AnythingOfType("*jsonreqresp.EventFilter"))
	})

	t.WithNewStep("success return 0 events", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		events := []*models.Event{}

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)
		mockEventRep.On(
			"GetAll", ctx,
			mock.AnythingOfType("*jsonreqresp.EventFilter"),
		).Return(events, nil)

		eventServ := eventserv.NewEventService(mockEventRep, mockArtworkRep)

		// ACT
		resEvents, err := eventServ.GetAll(ctx)

		sCtx.Require().NoError(err)
		sCtx.Require().True(len(resEvents) == 0)
		mockEventRep.AssertCalled(t, "GetAll", ctx,
			mock.AnythingOfType("*jsonreqresp.EventFilter"))
	})

	t.WithNewStep("error in rep", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		events := []*models.Event{}
		expectedErr := errors.New("eventRep error")

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)
		mockEventRep.On(
			"GetAll", ctx,
			mock.AnythingOfType("*jsonreqresp.EventFilter"),
		).Return(events, expectedErr)

		eventServ := eventserv.NewEventService(mockEventRep, mockArtworkRep)

		// ACT
		_, err := eventServ.GetAll(ctx)

		sCtx.Require().ErrorIs(err, expectedErr)
		mockEventRep.AssertCalled(t, "GetAll", ctx,
			mock.AnythingOfType("*jsonreqresp.EventFilter"))
	})
}

func (s *EventServiceSuite) TestEventService_GetArtworksFromEvent(t provider.T) {
	artworkCreator := testobj.NewArtworkMother()
	// eventCreator := testobj.NewEventMother()

	t.WithNewStep("success get artworks from event", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		eventID := uuid.New()
		artworkIDs := uuid.UUIDs{uuid.New(), uuid.New()}
		artworks := []*models.Artwork{
			artworkCreator.ArtworkP(artworkIDs[0]),
			artworkCreator.ArtworkP(artworkIDs[1]),
		}

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)

		mockEventRep.On("GetArtworkIDs", ctx, eventID).Return(artworkIDs, nil)
		mockArtworkRep.On("GetByID", ctx, artworkIDs[0]).Return(artworks[0], nil)
		mockArtworkRep.On("GetByID", ctx, artworkIDs[1]).Return(artworks[1], nil)

		eventServ := eventserv.NewEventService(mockEventRep, mockArtworkRep)

		// ACT
		resArtworks, err := eventServ.GetArtworksFromEvent(ctx, eventID)

		sCtx.Require().NoError(err)
		sCtx.Require().True(len(artworks) == len(resArtworks))
		for i := range len(resArtworks) {
			sCtx.Require().True(artworks[i].Equals(resArtworks[i]))
		}
		mockEventRep.AssertCalled(t, "GetArtworkIDs", ctx, eventID)
		mockArtworkRep.AssertCalled(t, "GetByID", ctx, artworkIDs[0])
		mockArtworkRep.AssertCalled(t, "GetByID", ctx, artworkIDs[1])
	})

	t.WithNewStep("error getting artwork IDs", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		eventID := uuid.New()
		expectedErr := errors.New("get artwork IDs error")

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)

		mockEventRep.On("GetArtworkIDs", ctx, eventID).Return(nil, expectedErr)

		eventServ := eventserv.NewEventService(mockEventRep, mockArtworkRep)

		// ACT
		_, err := eventServ.GetArtworksFromEvent(ctx, eventID)

		sCtx.Require().Error(err)
		sCtx.Assert().ErrorIs(err, expectedErr)
		mockEventRep.AssertCalled(t, "GetArtworkIDs", ctx, eventID)
		mockArtworkRep.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
	})

	t.WithNewStep("error getting artwork by ID", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		eventID := uuid.New()
		artworkIDs := uuid.UUIDs{uuid.New()}
		expectedErr := errors.New("artwork not found")

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)

		mockEventRep.On("GetArtworkIDs", ctx, eventID).Return(artworkIDs, nil)
		mockArtworkRep.On("GetByID", ctx, artworkIDs[0]).Return(nil, expectedErr)

		eventServ := eventserv.NewEventService(mockEventRep, mockArtworkRep)

		// ACT
		_, err := eventServ.GetArtworksFromEvent(ctx, eventID)

		sCtx.Require().Error(err)
		sCtx.Assert().ErrorIs(err, expectedErr)
		mockEventRep.AssertCalled(t, "GetArtworkIDs", ctx, eventID)
		mockArtworkRep.AssertCalled(t, "GetByID", ctx, artworkIDs[0])
	})

}

func (s *EventServiceSuite) TestEventService_Add(t provider.T) {
	eventCreator := testobj.NewEventMother()

	t.WithNewStep("success add event", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		employeeID := uuid.New()
		eventReq := eventCreator.EventAdd(employeeID)

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)

		mockEventRep.On("CheckEmployeeByID", ctx, employeeID).Return(true, nil)
		mockEventRep.On("GetEventsOfArtworkOnDate", ctx, mock.Anything, eventReq.DateBegin, eventReq.DateEnd).
			Return(nil, eventrep.ErrEventNotFound)
		mockEventRep.On("Add", ctx, mock.AnythingOfType("*models.Event")).Return(nil)
		mockEventRep.On("AddArtworksToEvent", ctx, mock.Anything, mock.Anything).Return(nil)

		eventServ := eventserv.NewEventService(mockEventRep, mockArtworkRep)

		// ACT
		err := eventServ.Add(ctx, eventReq)

		sCtx.Require().NoError(err)
		mockEventRep.AssertCalled(t, "CheckEmployeeByID", ctx, employeeID)
		mockEventRep.AssertCalled(t, "Add", ctx, mock.AnythingOfType("*models.Event"))
		mockEventRep.AssertCalled(t, "AddArtworksToEvent", ctx, mock.Anything, mock.Anything)
	})

	t.WithNewStep("error employee not found", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		employeeID := uuid.New()
		eventReq := eventCreator.EventAdd(employeeID)

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)

		mockEventRep.On("CheckEmployeeByID", ctx, employeeID).Return(false, nil)

		eventServ := eventserv.NewEventService(mockEventRep, mockArtworkRep)

		// ACT
		err := eventServ.Add(ctx, eventReq)

		sCtx.Require().Error(err)
		sCtx.Require().ErrorIs(err, eventrep.ErrAddNoEmployee)
		mockEventRep.AssertCalled(t, "CheckEmployeeByID", ctx, employeeID)
	})

	t.WithNewStep("error artwork busy", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		employeeID := uuid.New()
		eventReq := eventCreator.EventAdd(employeeID)

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)

		mockEventRep.On("CheckEmployeeByID", ctx, employeeID).Return(true, nil)
		mockEventRep.On("GetEventsOfArtworkOnDate", ctx, mock.Anything, eventReq.DateBegin, eventReq.DateEnd).
			Return([]*models.Event{eventCreator.EventP(uuid.New())}, nil)

		eventServ := eventserv.NewEventService(mockEventRep, mockArtworkRep)

		// ACT
		err := eventServ.Add(ctx, eventReq)

		sCtx.Require().Error(err)
		sCtx.Require().ErrorIs(err, eventserv.ErrArtworkBusy)
		mockEventRep.AssertCalled(t, "CheckEmployeeByID", ctx, employeeID)
		mockEventRep.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	})

}

func (s *EventServiceSuite) TestEventService_Delete(t provider.T) {
	t.WithNewStep("success delete event", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		eventID := uuid.New()

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)

		mockEventRep.On("Delete", ctx, eventID).Return(nil)

		eventServ := eventserv.NewEventService(mockEventRep, mockArtworkRep)

		// ACT
		err := eventServ.Delete(ctx, eventID)

		sCtx.Require().NoError(err)
		mockEventRep.AssertCalled(t, "Delete", ctx, eventID)
	})

	t.WithNewStep("error in delete", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		eventID := uuid.New()
		expectedErr := errors.New("delete error")

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)

		mockEventRep.On("Delete", ctx, eventID).Return(expectedErr)

		eventServ := eventserv.NewEventService(mockEventRep, mockArtworkRep)

		// ACT
		err := eventServ.Delete(ctx, eventID)

		sCtx.Require().ErrorIs(err, expectedErr)
		mockEventRep.AssertCalled(t, "Delete", ctx, eventID)
	})
}

func (s *EventServiceSuite) TestEventService_Update(t provider.T) {
	eventCreator := testobj.NewEventMother()

	t.WithNewStep("success update event", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		eventID := uuid.New()
		updateFields := &jsonreqresp.EventUpdate{
			Title: "Updated Title",
		}

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)

		event := eventCreator.EventP(eventID)
		mockEventRep.On("GetByID", ctx, eventID).Return(event, nil)
		mockEventRep.On("Update", ctx, eventID, mock.AnythingOfType("func(*models.Event) (*models.Event, error)")).Return(nil)

		eventServ := eventserv.NewEventService(mockEventRep, mockArtworkRep)

		// ACT
		err := eventServ.Update(ctx, eventID, updateFields)

		sCtx.Require().NoError(err)
		mockEventRep.AssertCalled(t, "GetByID", ctx, eventID)
		mockEventRep.AssertCalled(t, "Update", ctx, eventID, mock.AnythingOfType("func(*models.Event) (*models.Event, error)"))
	})

	t.WithNewStep("error event not found", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		eventID := uuid.New()
		updateFields := &jsonreqresp.EventUpdate{}
		expectedErr := errors.New("event not found")

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)

		mockEventRep.On("GetByID", ctx, eventID).Return(nil, expectedErr)

		eventServ := eventserv.NewEventService(mockEventRep, mockArtworkRep)

		// ACT
		err := eventServ.Update(ctx, eventID, updateFields)

		sCtx.Require().Error(err)
		sCtx.Assert().ErrorIs(err, expectedErr)
		mockEventRep.AssertCalled(t, "GetByID", ctx, eventID)
		mockEventRep.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything)
	})

}

func (s *EventServiceSuite) TestEventService_AddArtworksToEvent(t provider.T) {
	t.WithNewStep("success add artworks to event", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		eventID := uuid.New()
		artworkIDs := uuid.UUIDs{uuid.New(), uuid.New()}

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)

		event := testobj.NewEventMother().EventP(eventID)
		oldArtworkIDs := uuid.UUIDs{uuid.New()}
		mockEventRep.On("GetByID", ctx, eventID).Return(event, nil)
		mockEventRep.On("GetArtworkIDs", ctx, eventID).Return(oldArtworkIDs, nil)
		mockEventRep.On("AddArtworksToEvent", ctx, eventID, artworkIDs).Return(nil)

		eventServ := eventserv.NewEventService(mockEventRep, mockArtworkRep)

		// ACT
		err := eventServ.AddArtworksToEvent(ctx, eventID, artworkIDs)

		sCtx.Require().NoError(err)
		mockEventRep.AssertCalled(t, "GetByID", ctx, eventID)
		mockEventRep.AssertCalled(t, "GetArtworkIDs", ctx, eventID)
		mockEventRep.AssertCalled(t, "AddArtworksToEvent", ctx, eventID, artworkIDs)
	})

	t.WithNewStep("error duplicate artwork IDs", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		eventID := uuid.New()
		artworkIDs := uuid.UUIDs{uuid.New(), uuid.New()}

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)

		event := testobj.NewEventMother().EventP(eventID)
		oldArtworkIDs := artworkIDs // Duplicate IDs
		mockEventRep.On("GetByID", ctx, eventID).Return(event, nil)
		mockEventRep.On("GetArtworkIDs", ctx, eventID).Return(oldArtworkIDs, nil)

		eventServ := eventserv.NewEventService(mockEventRep, mockArtworkRep)

		// ACT
		err := eventServ.AddArtworksToEvent(ctx, eventID, artworkIDs)

		sCtx.Require().Error(err)
		sCtx.Require().ErrorIs(err, models.ErrAddArtwork)
		mockEventRep.AssertCalled(t, "GetByID", ctx, eventID)
		mockEventRep.AssertCalled(t, "GetArtworkIDs", ctx, eventID)
		mockEventRep.AssertNotCalled(t, "AddArtworksToEvent", mock.Anything, mock.Anything, mock.Anything)
	})

}

func (s *EventServiceSuite) TestEventService_DeleteArtworkFromEvent(t provider.T) {
	t.WithNewStep("success delete artwork from event", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		eventID := uuid.New()
		artworkID := uuid.New()

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)

		mockEventRep.On("DeleteArtworkFromEvent", ctx, eventID, artworkID).Return(nil)

		eventServ := eventserv.NewEventService(mockEventRep, mockArtworkRep)

		// ACT
		err := eventServ.DeleteArtworkFromEvent(ctx, eventID, artworkID)

		sCtx.Require().NoError(err)
		mockEventRep.AssertCalled(t, "DeleteArtworkFromEvent", ctx, eventID, artworkID)
	})

	t.WithNewStep("error in delete artwork", func(sCtx provider.StepCtx) {
		ctx := context.Background()
		eventID := uuid.New()
		artworkID := uuid.New()
		expectedErr := errors.New("delete error")

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)

		mockEventRep.On("DeleteArtworkFromEvent", ctx, eventID, artworkID).Return(expectedErr)

		eventServ := eventserv.NewEventService(mockEventRep, mockArtworkRep)

		// ACT
		err := eventServ.DeleteArtworkFromEvent(ctx, eventID, artworkID)

		sCtx.Require().ErrorIs(err, expectedErr)
		mockEventRep.AssertCalled(t, "DeleteArtworkFromEvent", ctx, eventID, artworkID)
	})
}
