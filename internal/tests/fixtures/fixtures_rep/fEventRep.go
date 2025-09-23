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

func AddTestEvents(
	t provider.T, ctx context.Context, events []*models.Event,
	eventRep eventrep.EventRep,
	employeeRep employeerep.EmployeeRep,
	adminRep adminrep.AdminRep,
	artworkRep artworkrep.ArtworkRep,
	authorRep authorrep.AuthorRep,
	collectionRep collectionrep.CollectionRep,
) {
	adminCreator := testobj.NewAdminMother()
	admin := adminCreator.DefaultAdmin(uuid.New())
	AddTestAdmin(t, ctx, adminRep, []*models.Admin{&admin})

	t.WithNewStep("Add Test Employees for Events", func(sCtx provider.StepCtx) {
		employees := make(map[uuid.UUID]*models.Employee, 0)
		employeeCreator := testobj.NewEmployeeMother()
		for _, e := range events {
			if _, ok := employees[e.GetEmployeeID()]; !ok {
				empl := employeeCreator.DefaultEmployeeP(e.GetEmployeeID(), admin.GetID())
				err := employeeRep.Add(ctx, empl)
				sCtx.Require().NoError(err)
			}
		}
	})

	t.WithNewStep("Add Test Artowrks for Events", func(sCtx provider.StepCtx) {
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
	admins := make(map[uuid.UUID]*models.Admin, 0)
	employees := make(map[uuid.UUID]*models.Employee, 0)

	t.WithNewStep("Search to delete Test Employees", func(sCtx provider.StepCtx) {
		for _, e := range events {
			if _, ok := employees[e.GetEmployeeID()]; !ok {
				empl, err := employeeRep.GetByID(ctx, e.GetEmployeeID())
				if err == nil {
					employees[empl.GetID()] = empl
				} else {
					sCtx.Assert().ErrorIs(err, employeerep.ErrEmployeeNotFound)
				}
			}
		}
	})
	t.WithNewStep("Search to delete Test Admins", func(sCtx provider.StepCtx) {
		for _, e := range employees {
			if _, ok := admins[e.GetAdminID()]; !ok {
				adm, err := adminRep.GetByID(ctx, e.GetAdminID())
				if err == nil {
					admins[adm.GetID()] = adm
				} else {
					sCtx.Assert().ErrorIs(err, adminrep.ErrAdminNotFound)
				}
			}
		}
	})

	t.WithNewStep("Delete Test Artowrks for Events", func(sCtx provider.StepCtx) {
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
	emplToDel := make([]*models.Employee, 0)
	for _, v := range employees {
		emplToDel = append(emplToDel, v)
	}
	DelTestEmployee(t, ctx, employeeRep, emplToDel)

	adminToDel := make([]*models.Admin, 0)
	for _, v := range admins {
		adminToDel = append(adminToDel, v)
	}
	DelTestAdmin(t, ctx, adminRep, adminToDel)
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
