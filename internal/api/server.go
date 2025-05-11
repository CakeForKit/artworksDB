package api

// type Server struct {
// 	router      *gin.Engine
// 	artworkServ artworkserv.ArtworkService
// }

// func NewServer() (*Server, error) {
// 	// appCnfg, err := cnfg.LoadAppConfig()
// 	// if err != nil {
// 	// 	fmt.Printf("cannot load config: %v", err)
// 	// 	return
// 	// }
// 	pgCreds, err := cnfg.LoadPgCredentials()
// 	if err != nil {
// 		log.Fatal("cannot load config:", err)
// 	}
// 	dbCnfg, err := cnfg.LoadDatebaseConfig("./configs/")
// 	if err != nil {
// 		log.Fatal("cannot load config:", err)
// 	}

// 	ctx := context.Background()
// 	artworkRep, err := artworkrep.NewArtworkRep(ctx, pgCreds, dbCnfg)
// 	if err != nil {
// 		return nil, fmt.Errorf("NewArtworkRep: %v", err)
// 	}
// 	artworkServ := artworkserv.NewArtworkService(artworkRep)

// 	s := &Server{
// 		artworkServ: artworkServ,
// 	}
// 	s.setupRouter()
// 	return s, nil
// }

// func (server *Server) setupRouter() {
// 	// Ручная настройка middleware для логирования и обработки ошибок panic.
// 	router := gin.New()
// 	router.Use(gin.Logger())
// 	router.Use(gin.Recovery())

// 	// route для Swagger - НЕ ТРОГАТЬ
// 	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
// 	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

// 	router.GET("/artworks", server.getAllArtworks)

// 	server.router = router
// }

// func (server *Server) Start(address string) error {
// 	return server.router.Run(address)
// }

// func (server *Server) errorResponse(err error) gin.H {
// 	return gin.H{"error": err.Error()}
// }
