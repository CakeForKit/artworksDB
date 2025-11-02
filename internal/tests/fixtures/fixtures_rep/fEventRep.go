package fixturesrep

import (
	"context"

	"github.com/CakeForKit/artworksDB.git/internal/models"
	jsonreqresp "github.com/CakeForKit/artworksDB.git/internal/models/json_req_resp"
	"github.com/CakeForKit/artworksDB.git/internal/repository/adminrep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/artworkrep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/authorrep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/collectionrep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/employeerep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/eventrep"
	testobj "github.com/CakeForKit/artworksDB.git/internal/tests/testObj"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func getEmployeeIDsOfEvents(events []*models.Event) (employeeIDs uuid.UUIDs) {
	employeesMap := make(map[uuid.UUID]bool, 0)
	for _, v := range events {
		id := v.GetEmployeeID()
		if _, ok := employeesMap[id]; !ok {
			employeesMap[id] = true
			employeeIDs = append(employeeIDs, id)
		}
	}
	return
}

func getArtworkIDsOfEvents(events []*models.Event) (artworkIDS uuid.UUIDs) {
	artworksMap := make(map[uuid.UUID]bool, 0)
	for _, e := range events {
		for _, ids := range e.GetArtworkIDs() {
			if _, ok := artworksMap[ids]; !ok {
				artworksMap[ids] = true
				artworkIDS = append(artworkIDS, ids)
			}
		}
	}
	return
}

func AddTestEvents(
	t provider.T, ctx context.Context, events []*models.Event,
	eventRep eventrep.EventRep,
	employeeRep employeerep.EmployeeRep,
	adminRep adminrep.AdminRep,
	artworkRep artworkrep.ArtworkRep,
	authorRep authorrep.AuthorRep,
	collectionRep collectionrep.CollectionRep,
) {
	t.WithNewStep("Add Test Employees", func(sCtx provider.StepCtx) {
		adminID := uuid.New()
		employeeIDs := getEmployeeIDsOfEvents(events)
		employeeCreator := testobj.NewEmployeeMother()
		employees := make([]*models.Employee, 0, len(employeeIDs))
		for _, id := range employeeIDs {
			employees = append(employees, employeeCreator.DefaultEmployeeP(id, adminID))
		}
		AddTestEmployees(t, ctx, employees, employeeRep, adminRep)
	})

	t.WithNewStep("Add Test Artworks", func(sCtx provider.StepCtx) {
		artworkIDs := getArtworkIDsOfEvents(events)
		artworkCreator := testobj.NewArtworkMother()
		artworks := make([]*models.Artwork, 0, len(artworkIDs))
		for _, id := range artworkIDs {
			artworks = append(artworks, artworkCreator.ArtworkP(id))
		}
		AddTestArtworks(t, ctx, artworks, artworkRep, authorRep, collectionRep)
	})

	t.WithNewStep("Add Test Events", func(sCtx provider.StepCtx) {
		for _, u := range events {
			err := eventRep.Add(ctx, u)
			sCtx.Require().NoError(err)
			err = eventRep.AddArtworksToEvent(ctx, u.GetID(), u.GetArtworkIDs())
			sCtx.Require().NoError(err)
		}
	})
}

func AddTestEventsWithArtworks(
	t provider.T, ctx context.Context, events []*models.Event, artworks []*models.Artwork,
	eventRep eventrep.EventRep,
	employeeRep employeerep.EmployeeRep,
	adminRep adminrep.AdminRep,
	artworkRep artworkrep.ArtworkRep,
	authorRep authorrep.AuthorRep,
	collectionRep collectionrep.CollectionRep,
) {
	t.WithNewStep("Add Test Employees", func(sCtx provider.StepCtx) {
		adminID := uuid.New()
		employeeIDs := getEmployeeIDsOfEvents(events)
		employeeCreator := testobj.NewEmployeeMother()
		employees := make([]*models.Employee, 0, len(employeeIDs))
		for _, id := range employeeIDs {
			employees = append(employees, employeeCreator.DefaultEmployeeP(id, adminID))
		}
		AddTestEmployees(t, ctx, employees, employeeRep, adminRep)
	})

	t.WithNewStep("Add Test Artworks", func(sCtx provider.StepCtx) {
		if len(artworks) > 0 {
			AddTestArtworks(t, ctx, artworks, artworkRep, authorRep, collectionRep)
		}
	})

	t.WithNewStep("Add Test Events", func(sCtx provider.StepCtx) {
		for _, u := range events {
			err := eventRep.Add(ctx, u)
			sCtx.Require().NoError(err)
			err = eventRep.AddArtworksToEvent(ctx, u.GetID(), u.GetArtworkIDs())
			sCtx.Require().NoError(err)
		}
	})
}

func DelTestEvents(
	t provider.T, ctx context.Context, events []*models.Event,
	eventRep eventrep.EventRep,
	employeeRep employeerep.EmployeeRep,
	adminRep adminrep.AdminRep,
	artworkRep artworkrep.ArtworkRep,
	authorRep authorrep.AuthorRep,
	collectionRep collectionrep.CollectionRep,
) {
	t.WithNewStep("Delete artwork-event relationships", func(sCtx provider.StepCtx) {
		for _, e := range events {
			artworkIDs := e.GetArtworkIDs()
			for _, v := range artworkIDs {
				err := eventRep.DeleteArtworkFromEvent(ctx, e.GetID(), v)
				sCtx.Assert().NoError(err)
			}
		}
	})
	t.WithNewStep("Delete Test Events", func(sCtx provider.StepCtx) {
		for _, u := range events {
			err := eventRep.Delete(ctx, u.GetID())
			sCtx.Assert().NoError(err)
		}
	})
	t.WithNewStep("Delete Test Artworks", func(sCtx provider.StepCtx) {
		artworkIDs := getArtworkIDsOfEvents(events)
		artworks := make([]*models.Artwork, 0, len(artworkIDs))
		for _, id := range artworkIDs {
			a, err := artworkRep.GetByID(ctx, id)
			sCtx.Assert().NoError(err)
			artworks = append(artworks, a)
		}
		DelTestArtworks(t, ctx, artworks, artworkRep, authorRep, collectionRep)
	})
	t.WithNewStep("Delete Test Employees", func(sCtx provider.StepCtx) {
		employeeIDs := getEmployeeIDsOfEvents(events)
		employees := make([]*models.Employee, 0, len(employeeIDs))
		for _, id := range employeeIDs {
			a, err := employeeRep.GetByID(ctx, id)
			sCtx.Assert().NoError(err)
			employees = append(employees, a)
		}
		DelTestEmployees(t, ctx, employees, employeeRep, adminRep)
	})
}

func AssertEventResponsesAreInRes(t provider.T, eventResp, expectedEventResp []jsonreqresp.EventResponse) {
	t.WithNewStep("Check eventResp sre in the result", func(sCtx provider.StepCtx) {
		foundAll := make([]bool, len(expectedEventResp))
		for i, ru := range expectedEventResp {
			for _, u := range eventResp {
				if ru.Equal(&u) {
					foundAll[i] = true
				}
			}
		}
		for _, v := range foundAll {
			sCtx.Assert().True(v)
		}
	})
}

func AssertEventsAreInRes(t provider.T, events, resEvents []*models.Event) {
	t.WithNewStep("Check events sre in the result", func(sCtx provider.StepCtx) {
		foundAll := make([]bool, len(events))
		for _, ru := range resEvents {
			for i, u := range events {
				if ru.Equal(u) {
					foundAll[i] = true
				}
			}
		}
		for _, v := range foundAll {
			sCtx.Assert().True(v)
		}
	})
}
