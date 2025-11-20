package creators

import (
	"context"
	"fmt"

	"github.com/CakeForKit/artworksDB.git/internal/cnfg"
	"github.com/CakeForKit/artworksDB.git/internal/repository/adminrep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/artworkrep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/authorrep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/buyticketstxrep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/collectionrep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/employeerep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/eventrep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/ticketpurchasesrep"
	"github.com/CakeForKit/artworksDB.git/internal/repository/userrep"
	"github.com/CakeForKit/artworksDB.git/internal/services/adminserv"
	"github.com/CakeForKit/artworksDB.git/internal/services/artworkserv"
	attemptsrep "github.com/CakeForKit/artworksDB.git/internal/services/auth/attempts_rep"
	authadmin "github.com/CakeForKit/artworksDB.git/internal/services/auth/auth_admin"
	authemployee "github.com/CakeForKit/artworksDB.git/internal/services/auth/auth_employee"
	authsessionrep "github.com/CakeForKit/artworksDB.git/internal/services/auth/auth_session_repository"
	authuser "github.com/CakeForKit/artworksDB.git/internal/services/auth/auth_user"
	"github.com/CakeForKit/artworksDB.git/internal/services/auth/authz"
	"github.com/CakeForKit/artworksDB.git/internal/services/auth/hasher"
	"github.com/CakeForKit/artworksDB.git/internal/services/auth/otp"
	"github.com/CakeForKit/artworksDB.git/internal/services/auth/token"
	"github.com/CakeForKit/artworksDB.git/internal/services/authorserv"
	"github.com/CakeForKit/artworksDB.git/internal/services/buyticketserv"
	"github.com/CakeForKit/artworksDB.git/internal/services/collectionserv"
	"github.com/CakeForKit/artworksDB.git/internal/services/emailserv"
	"github.com/CakeForKit/artworksDB.git/internal/services/eventserv"
	"github.com/CakeForKit/artworksDB.git/internal/services/mailing"
	"github.com/CakeForKit/artworksDB.git/internal/services/searcher"
	"github.com/CakeForKit/artworksDB.git/internal/services/userservice"
	"github.com/CakeForKit/artworksDB.git/internal/tracing"
)

func GetConfig() (
	appCnfg *cnfg.AppConfig,
	dbCreds *cnfg.DatebaseCredentials,
	redisCreds *cnfg.RedisCredentials,
	dbCnfg *cnfg.DatebaseConfig,
	err error,
) {
	appCnfg, err = cnfg.LoadAppConfig("./configs/", "config", "yaml")
	if err != nil {
		return
	}
	if appCnfg.Datebase == cnfg.PostgresDB {
		dbCreds, err = cnfg.LoadPgCredentials("./configs/", "db", "env")
		if err != nil {
			err = fmt.Errorf("cannot load PgCredentials: %v", err)
			return
		}
	} else if appCnfg.Datebase == cnfg.ClickHouseDB {
		dbCreds, err = cnfg.LoadClickHouseCredentials("./configs/", "clickhouse", "env")
		if err != nil {
			err = fmt.Errorf("cannot load ClickHouseCredentials: %v", err)
			return
		}
	}

	redisCreds, err = cnfg.LoadRedisCredentials("./configs/", "redis", "env")
	if err != nil {
		panic(fmt.Errorf("cannot load RedisCredentials: %v", err))
	}
	dbCnfg, err = cnfg.LoadDatebaseConfig("./configs/", "config", "yaml")
	if err != nil {
		panic(fmt.Errorf("cannot load DatebaseConfig: %v", err))
	}
	return // appCnfg, dbCreds, redisCreds, dbCnfg
}

func GetRepositories(
	ctx context.Context,
	appCnfg *cnfg.AppConfig,
	dbCreds *cnfg.DatebaseCredentials,
	redisCreds *cnfg.RedisCredentials,
	dbCnfg *cnfg.DatebaseConfig,
) (
	userRep userrep.UserRep, employeeRep employeerep.EmployeeRep,
	adminRep adminrep.AdminRep, collectionRep collectionrep.CollectionRep,
	authorRep authorrep.AuthorRep, artworkRep artworkrep.ArtworkRep,
	eventRep eventrep.EventRep, txRep buyticketstxrep.BuyTicketsTxRep,
	tPurchasesRep ticketpurchasesrep.TicketPurchasesRep,
	err error,
) {
	userRep, err = userrep.NewUserRep(ctx, appCnfg.Datebase, dbCreds, dbCnfg)
	if err != nil {
		return
	}
	employeeRep, err = employeerep.NewEmployeeRep(ctx, appCnfg.Datebase, dbCreds, dbCnfg)
	if err != nil {
		return
	}
	adminRep, err = adminrep.NewAdminRep(ctx, appCnfg.Datebase, dbCreds, dbCnfg)
	if err != nil {
		return
	}
	collectionRep, err = collectionrep.NewCollectionRep(ctx, appCnfg.Datebase, dbCreds, dbCnfg)
	if err != nil {
		return
	}
	authorRep, err = authorrep.NewAuthorRep(ctx, appCnfg.Datebase, dbCreds, dbCnfg)
	if err != nil {
		return
	}
	artworkRep, err = artworkrep.NewArtworkRep(ctx, appCnfg.Datebase, dbCreds, dbCnfg)
	if err != nil {
		return
	}
	eventRep, err = eventrep.NewEventRep(ctx, appCnfg.Datebase, dbCreds, dbCnfg)
	if err != nil {
		return
	}
	txRep, err = buyticketstxrep.NewBuyTicketsTxRep(ctx, redisCreds)
	if err != nil {
		return
	}
	tPurchasesRep, err = ticketpurchasesrep.NewTicketPurchasesRep(ctx, appCnfg.Datebase, dbCreds, dbCnfg)
	if err != nil {
		return
	}
	return
}

func GetAuthServs(
	appCnfg *cnfg.AppConfig,
	tracer *tracing.Tracer,
	userRep userrep.UserRep, employeeRep employeerep.EmployeeRep, adminRep adminrep.AdminRep,
) (
	authZ authz.AuthZ, hash hasher.Hasher,
	authUserServ authuser.AuthUser, authEmployeeServ authemployee.AuthEmployee, authAdminServ authadmin.AdminAuth,
) {
	authZ = authz.NewAuthZ()
	if tracer.IsEnabled() {
		authZ = authz.NewTracedAuthZ(authZ, tracer)
	}
	tokenMaker, err := token.NewTokenMaker(appCnfg.TokenSymmetricKey)
	if err != nil {
		panic(fmt.Errorf("cannot create token maker: %w", err))
	}
	hash, err = hasher.NewHasher()
	if err != nil {
		panic(err)
	}
	authUserSessionRep := authsessionrep.NewAuthUserSessionRep(appCnfg.DurationLoginSession)
	loginAttemptRep := attemptsrep.NewLoginAttemptUserRep(appCnfg.MaxLoginAttemps, appCnfg.DurationLoginSession)
	otpAttemptRep := attemptsrep.NewOTPAttemptRep(appCnfg.MaxLoginAttemps, appCnfg.DurationLoginSession)
	emailCnfg := cnfg.LoadEmailCnfg()
	emailServ := emailserv.NewEmailService(
		emailCnfg.Host,
		emailCnfg.Port,
		emailCnfg.Username,
		emailCnfg.Password,
		emailCnfg.From,
	)
	otpServ := otp.NewOTPService(*emailServ)
	authUserServ = authuser.NewAuthUser(
		*appCnfg, userRep,
		tokenMaker, hash,
		otpServ, authUserSessionRep,
		loginAttemptRep, otpAttemptRep)
	if tracer.IsEnabled() {
		authUserServ = authuser.NewTracedAuthUser(authUserServ, tracer)
	}
	authEmployeeServ = authemployee.NewAuthEmployee(*appCnfg, employeeRep, tokenMaker, hash)
	if tracer.IsEnabled() {
		authEmployeeServ = authemployee.NewTracedAuthEmployee(authEmployeeServ, tracer)
	}
	authAdminServ = authadmin.NewAuthAdmin(*appCnfg, adminRep, tokenMaker, hash)
	if tracer.IsEnabled() {
		authAdminServ = authadmin.NewTracedAuthAdmin(authAdminServ, tracer)
	}

	return
}

func GetServs(
	appCnfg *cnfg.AppConfig,
	tracer *tracing.Tracer,
	userRep userrep.UserRep, employeeRep employeerep.EmployeeRep, adminRep adminrep.AdminRep,
	authZ authz.AuthZ, hasher hasher.Hasher,
	txRep buyticketstxrep.BuyTicketsTxRep, tPurchasesRep ticketpurchasesrep.TicketPurchasesRep,
	collectionRep collectionrep.CollectionRep, authorRep authorrep.AuthorRep, artworkRep artworkrep.ArtworkRep,
	eventRep eventrep.EventRep,

) (
	userServ userservice.UserService, adminServ adminserv.AdminService,
	buyTicketServ buyticketserv.BuyTicketsServ, collectionServ collectionserv.CollectionServ,
	authroServ authorserv.AuthorServ, artworkServ artworkserv.ArtworkService,
	eventServ eventserv.EventService, searcherServ searcher.Searcher,
	mailingServ mailing.MailingService,
) {
	userServ = userservice.NewUserService(userRep, authZ, hasher)
	if tracer.IsEnabled() {
		userServ = userservice.NewTracedUserService(userServ, tracer)
	}
	adminServ = adminserv.NewAdminService(employeeRep, userRep, authZ)
	buyTicketServ, _ = buyticketserv.NewBuyTicketsServ(txRep, tPurchasesRep, *appCnfg, authZ, userRep, eventRep)
	collectionServ = collectionserv.NewCollectionServ(collectionRep)
	authroServ = authorserv.NewAuthorServ(authorRep)
	artworkServ = artworkserv.NewArtworkService(artworkRep, authorRep, collectionRep)
	eventServ = eventserv.NewEventService(eventRep, artworkRep)
	searcherServ = searcher.NewSearcher(artworkRep, eventRep)
	if tracer.IsEnabled() {
		searcherServ = searcher.NewTracedSearcher(searcherServ, tracer)
	}

	mailingServ = mailing.NewGmailSender(userRep, "museum", "museum@test.ru", "1234")
	return
}
