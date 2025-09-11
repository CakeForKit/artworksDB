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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestEventServ_GetAll(t *testing.T) {
	eventCreator := testobj.NewEventMother()

	t.Run("success return 2 events", func(t *testing.T) {
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

		require.Nil(t, err)
		require.True(t, len(events) == len(resEvents))
		for i := range len(resEvents) {
			require.True(t, events[i].Equals(resEvents[i]))
		}
		mockEventRep.AssertCalled(t, "GetAll", ctx,
			mock.AnythingOfType("*jsonreqresp.EventFilter"))
	})

	t.Run("success return 0 events", func(t *testing.T) {
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

		require.Nil(t, err)
		require.True(t, len(resEvents) == 0)
		mockEventRep.AssertCalled(t, "GetAll", ctx,
			mock.AnythingOfType("*jsonreqresp.EventFilter"))
	})

	t.Run("error in rep", func(t *testing.T) {
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

		require.ErrorIs(t, err, expectedErr)
		mockEventRep.AssertCalled(t, "GetAll", ctx,
			mock.AnythingOfType("*jsonreqresp.EventFilter"))
	})
}

func TestEventService_GetArtworksFromEvent(t *testing.T) {
	artworkCreator := testobj.NewArtworkMother()
	// eventCreator := testobj.NewEventMother()

	t.Run("success get artworks from event", func(t *testing.T) {
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

		require.Nil(t, err)
		require.True(t, len(artworks) == len(resArtworks))
		for i := range len(resArtworks) {
			require.True(t, artworks[i].Equals(resArtworks[i]))
		}
		mockEventRep.AssertCalled(t, "GetArtworkIDs", ctx, eventID)
		mockArtworkRep.AssertCalled(t, "GetByID", ctx, artworkIDs[0])
		mockArtworkRep.AssertCalled(t, "GetByID", ctx, artworkIDs[1])
	})

	t.Run("error getting artwork IDs", func(t *testing.T) {
		ctx := context.Background()
		eventID := uuid.New()
		expectedErr := errors.New("get artwork IDs error")

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)

		mockEventRep.On("GetArtworkIDs", ctx, eventID).Return(nil, expectedErr)

		eventServ := eventserv.NewEventService(mockEventRep, mockArtworkRep)

		// ACT
		_, err := eventServ.GetArtworksFromEvent(ctx, eventID)

		require.Error(t, err)
		require.Contains(t, err.Error(), "eventService.GetArtworkFromEvent")
		mockEventRep.AssertCalled(t, "GetArtworkIDs", ctx, eventID)
		mockArtworkRep.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
	})

	t.Run("error getting artwork by ID", func(t *testing.T) {
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

		require.Error(t, err)
		require.Contains(t, err.Error(), "eventService.GetArtworkFromEvent")
		mockEventRep.AssertCalled(t, "GetArtworkIDs", ctx, eventID)
		mockArtworkRep.AssertCalled(t, "GetByID", ctx, artworkIDs[0])
	})

}

func TestEventService_Add(t *testing.T) {
	eventCreator := testobj.NewEventMother()

	t.Run("success add event", func(t *testing.T) {
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

		require.Nil(t, err)
		mockEventRep.AssertCalled(t, "CheckEmployeeByID", ctx, employeeID)
		mockEventRep.AssertCalled(t, "Add", ctx, mock.AnythingOfType("*models.Event"))
		mockEventRep.AssertCalled(t, "AddArtworksToEvent", ctx, mock.Anything, mock.Anything)
	})

	t.Run("error employee not found", func(t *testing.T) {
		ctx := context.Background()
		employeeID := uuid.New()
		eventReq := eventCreator.EventAdd(employeeID)

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)

		mockEventRep.On("CheckEmployeeByID", ctx, employeeID).Return(false, nil)

		eventServ := eventserv.NewEventService(mockEventRep, mockArtworkRep)

		// ACT
		err := eventServ.Add(ctx, eventReq)

		require.Error(t, err)
		require.ErrorIs(t, err, eventrep.ErrAddNoEmployee)
		mockEventRep.AssertCalled(t, "CheckEmployeeByID", ctx, employeeID)
	})

	t.Run("error artwork busy", func(t *testing.T) {
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

		require.Error(t, err)
		require.ErrorIs(t, err, eventserv.ErrArtworkBusy)
		mockEventRep.AssertCalled(t, "CheckEmployeeByID", ctx, employeeID)
		mockEventRep.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	})

}

func TestEventService_Delete(t *testing.T) {
	t.Run("success delete event", func(t *testing.T) {
		ctx := context.Background()
		eventID := uuid.New()

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)

		mockEventRep.On("Delete", ctx, eventID).Return(nil)

		eventServ := eventserv.NewEventService(mockEventRep, mockArtworkRep)

		// ACT
		err := eventServ.Delete(ctx, eventID)

		require.Nil(t, err)
		mockEventRep.AssertCalled(t, "Delete", ctx, eventID)
	})

	t.Run("error in delete", func(t *testing.T) {
		ctx := context.Background()
		eventID := uuid.New()
		expectedErr := errors.New("delete error")

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)

		mockEventRep.On("Delete", ctx, eventID).Return(expectedErr)

		eventServ := eventserv.NewEventService(mockEventRep, mockArtworkRep)

		// ACT
		err := eventServ.Delete(ctx, eventID)

		require.ErrorIs(t, err, expectedErr)
		mockEventRep.AssertCalled(t, "Delete", ctx, eventID)
	})
}

func TestEventService_Update(t *testing.T) {
	eventCreator := testobj.NewEventMother()

	t.Run("success update event", func(t *testing.T) {
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

		require.Nil(t, err)
		mockEventRep.AssertCalled(t, "GetByID", ctx, eventID)
		mockEventRep.AssertCalled(t, "Update", ctx, eventID, mock.AnythingOfType("func(*models.Event) (*models.Event, error)"))
	})

	t.Run("error event not found", func(t *testing.T) {
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

		require.Error(t, err)
		require.Contains(t, err.Error(), "eventService.Update")
		mockEventRep.AssertCalled(t, "GetByID", ctx, eventID)
		mockEventRep.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything)
	})

}

func TestEventService_AddArtworksToEvent(t *testing.T) {
	t.Run("success add artworks to event", func(t *testing.T) {
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

		require.Nil(t, err)
		mockEventRep.AssertCalled(t, "GetByID", ctx, eventID)
		mockEventRep.AssertCalled(t, "GetArtworkIDs", ctx, eventID)
		mockEventRep.AssertCalled(t, "AddArtworksToEvent", ctx, eventID, artworkIDs)
	})

	t.Run("error duplicate artwork IDs", func(t *testing.T) {
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

		require.Error(t, err)
		require.ErrorIs(t, err, models.ErrAddArtwork)
		mockEventRep.AssertCalled(t, "GetByID", ctx, eventID)
		mockEventRep.AssertCalled(t, "GetArtworkIDs", ctx, eventID)
		mockEventRep.AssertNotCalled(t, "AddArtworksToEvent", mock.Anything, mock.Anything, mock.Anything)
	})

}

func TestEventService_DeleteArtworkFromEvent(t *testing.T) {
	t.Run("success delete artwork from event", func(t *testing.T) {
		ctx := context.Background()
		eventID := uuid.New()
		artworkID := uuid.New()

		mockArtworkRep := new(artworkrep.MockArtworkRep)
		mockEventRep := new(eventrep.MockEventRep)

		mockEventRep.On("DeleteArtworkFromEvent", ctx, eventID, artworkID).Return(nil)

		eventServ := eventserv.NewEventService(mockEventRep, mockArtworkRep)

		// ACT
		err := eventServ.DeleteArtworkFromEvent(ctx, eventID, artworkID)

		require.Nil(t, err)
		mockEventRep.AssertCalled(t, "DeleteArtworkFromEvent", ctx, eventID, artworkID)
	})

	t.Run("error in delete artwork", func(t *testing.T) {
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

		require.ErrorIs(t, err, expectedErr)
		mockEventRep.AssertCalled(t, "DeleteArtworkFromEvent", ctx, eventID, artworkID)
	})
}
