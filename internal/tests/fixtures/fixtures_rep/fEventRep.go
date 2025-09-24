package fixturesrep

import (
	"context"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models"
	jsonreqresp "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/models/json_req_resp"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/adminrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/artworkrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/authorrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/collectionrep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/employeerep"
	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/repository/eventrep"
	testobj "git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/testObj"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func getEmployeesOfEvents(events []*models.Event) (employees []*models.Employee) {
	employeeCreator := testobj.NewEmployeeMother()
	employeesMap := make(map[uuid.UUID]*models.Employee, 0)
	for _, v := range events {
		if _, ok := employeesMap[v.GetEmployeeID()]; !ok {
			empl := employeeCreator.DefaultEmployeeP(v.GetEmployeeID(), uuid.New())
			employees = append(employees, empl)
			employeesMap[v.GetEmployeeID()] = empl
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
	employees := getEmployeesOfEvents(events)
	AddTestEmployees(t, ctx, employees, employeeRep, adminRep)

	t.WithNewStep("Add Test Artworks for Events", func(sCtx provider.StepCtx) {
		artworksMap := make(map[uuid.UUID]*models.Artwork, 0)
		artworks := make([]*models.Artwork, 0)
		artworkCreator := testobj.NewArtworkMother()
		for _, e := range events {
			for _, ids := range e.GetArtworkIDs() {
				if _, ok := artworksMap[ids]; !ok {
					a := artworkCreator.ArtworkP(ids)
					artworksMap[ids] = a
					artworks = append(artworks, a)
				}
			}
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

func DelTestEvent(
	t provider.T, ctx context.Context, events []*models.Event,
	eventRep eventrep.EventRep,
	employeeRep employeerep.EmployeeRep,
	adminRep adminrep.AdminRep,
	artworkRep artworkrep.ArtworkRep,
	authorRep authorrep.AuthorRep,
	collectionRep collectionrep.CollectionRep,
) {
	t.WithNewStep("Delete Test Artworks for Events", func(sCtx provider.StepCtx) {
		artworksMap := make(map[uuid.UUID]*models.Artwork, 0)
		artworks := make([]*models.Artwork, 0)
		artworkCreator := testobj.NewArtworkMother()
		for _, e := range events {
			for _, ids := range e.GetArtworkIDs() {
				if _, ok := artworksMap[ids]; !ok {
					a := artworkCreator.ArtworkP(ids)
					artworksMap[ids] = a
					artworks = append(artworks, a)
				}
			}
		}
		DelTestArtworks(t, ctx, artworks, artworkRep, authorRep, collectionRep)
	})

	t.WithNewStep("Delete Test Events", func(sCtx provider.StepCtx) {
		for _, u := range events {
			err := eventRep.Delete(ctx, u.GetID())
			sCtx.Assert().NoError(err)
		}
	})
	employees := getEmployeesOfEvents(events)
	DelTestEmployees(t, ctx, employees, employeeRep, adminRep)
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
