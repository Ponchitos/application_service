package api

import (
	"github.com/Ponchitos/application_service/server/api/middleware"
	applicationRoutes "github.com/Ponchitos/application_service/server/api/routes/application"
	"github.com/Ponchitos/application_service/server/internal/services/applications"
	"github.com/Ponchitos/application_service/server/tools/logger"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/pprof"
)

func NewHTTPHandler(appService applications.Service, logger logger.Logger) http.Handler {
	router := mux.NewRouter()

	attachProfiler(router)

	mainAPIRouter := router.PathPrefix("/api").Subrouter()
	firstVersionAPIRouter := mainAPIRouter.PathPrefix("/v1").Subrouter()

	firstVersionAPIRouter.Use(middleware.SetRequestID)
	firstVersionAPIRouter.Use(middleware.LoggerRequest(logger))

	firstVersionAPIRouter.Handle("/upload", applicationRoutes.UploadAPKFileRoute(appService, logger)).Methods("POST")
	firstVersionAPIRouter.Handle("/upload/complete", applicationRoutes.CompleteUploadAPKFileRoute(appService, logger)).Methods("POST")
	firstVersionAPIRouter.Handle("/upload/cancel", applicationRoutes.CancelApkUploadRoute(appService, logger)).Methods("POST")
	firstVersionAPIRouter.Handle("/applications/{applicationUId}/versions", applicationRoutes.GetApplicationVersions(appService, logger)).Methods("GET")
	firstVersionAPIRouter.Handle("/applications/{applicationUId}/versions/{versionUId}", applicationRoutes.GetApplicationInfoRoute(appService, logger)).Methods("GET")
	firstVersionAPIRouter.Handle("/applications", applicationRoutes.GetApplicationsRoute(appService, logger)).Methods("GET")
	firstVersionAPIRouter.Handle("/application/apk/delete", applicationRoutes.DeleteApkVersionFile(appService, logger)).Methods("POST")
	firstVersionAPIRouter.Handle("/applications/{applicationUId}/versions/{versionUId}/download", applicationRoutes.DownloadApkFileRoute(appService, logger)).Methods("GET")
	firstVersionAPIRouter.Handle("/applications/{applicationUId}/versions/{versionUId}/download/chunk", applicationRoutes.DownloadApkFileByChunkRoute(appService, logger)).Methods("GET")

	internalAPIRouter := firstVersionAPIRouter.PathPrefix("/internal").Subrouter()

	internalAPIRouter.Handle("/application/uninstall/complete", applicationRoutes.UninstallApplicationCompleteRoute(appService, logger)).Methods("POST")

	return firstVersionAPIRouter
}

func attachProfiler(router *mux.Router) {
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	// Manually add support for paths linked to by index page at /debug/pprof/
	router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	router.Handle("/debug/pprof/block", pprof.Handler("block"))
}
