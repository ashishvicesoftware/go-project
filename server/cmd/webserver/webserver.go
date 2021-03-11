package main

import (
	"encoding/json"
	"fmt"
	"time"

	//	"io/ioutil"
	"net/http"
	"net/smtp"

	//	"os"
	"strconv"

	//  "errors"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/ashishvicesoftware/go-project/server/cmd/webserver/docs"
 
	"github.com/ashishvicesoftware/go-project/server/cmd/webserver/models"
	"github.com/ashishvicesoftware/go-project/server/pkg/database"
 	"go.uber.org/zap"
	// "github.com/auth0/go-jwt-middleware"
	// "github.com/dgrijalva/jwt-go"
)

https://github.com/ashishyadav21/golang.git/
// type Product struct {
// 	Id          int
// 	Name        string
// 	Slug        string
// 	Description string
// }

// var products = []Product{
// 	Product{Id: 1, Name: "Hover Shooters", Slug: "hover-shooters", Description: "Shoot your way to the top on 14 different hoverboards"},
// 	Product{Id: 2, Name: "Ocean Explorer", Slug: "ocean-explorer", Description: "Explore the depths of the sea in this one of a kind underwater experience"},
// 	Product{Id: 3, Name: "Dinosaur Park", Slug: "dinosaur-park", Description: "Go back 65 million years in the past and ride a T-Rex"},
// 	Product{Id: 4, Name: "Cars VR", Slug: "cars-vr", Description: "Get behind the wheel of the fastest cars in the world."},
// 	Product{Id: 5, Name: "Robin Hood", Slug: "robin-hood", Description: "Pick up the bow and arrow and master the art of archery"},
// 	Product{Id: 6, Name: "Real World VR", Slug: "real-world-vr", Description: "Explore the seven wonders of the world in VR"},
// }

// var mySigningKey = []byte("vGLS8dEvEmC9yYOuLtVHkXIpf2nmvprSw8YQAuza2TXTtCexuh8E4bittdmvTI2i")

// TODO look for a better way to handle userId
var userId = ""

type webserver struct {
	addr string
	db   database.DB
}

/**************************************************************************
	 StatusContinue           = 100
	 StatusSwitchingProtocols = 101

     StatusOK                   = 200
     StatusCreated              = 201
     StatusAccepted             = 202
     StatusNonAuthoritativeInfo = 203
     StatusNoContent            = 204
     StatusResetContent         = 205
     StatusPartialContent       = 206

     StatusMultipleChoices   = 300
     StatusMovedPermanently  = 301
     StatusFound             = 302
     StatusSeeOther          = 303
     StatusNotModified       = 304
     StatusUseProxy          = 305
     StatusTemporaryRedirect = 307

     StatusBadRequest                   = 400
     StatusUnauthorized                 = 401
     StatusPaymentRequired              = 402
     StatusForbidden                    = 403
     StatusNotFound                     = 404
     StatusMethodNotAllowed             = 405
     StatusNotAcceptable                = 406
     StatusProxyAuthRequired            = 407
     StatusRequestTimeout               = 408
     StatusConflict                     = 409
     StatusGone                         = 410
     StatusLengthRequired               = 411
     StatusPreconditionFailed           = 412
     StatusRequestEntityTooLarge        = 413
     StatusRequestURITooLong            = 414
     StatusUnsupportedMediaType         = 415
     StatusRequestedRangeNotSatisfiable = 416
     StatusExpectationFailed            = 417
     StatusTeapot                       = 418

     StatusInternalServerError     = 500
     StatusNotImplemented          = 501
     StatusBadGateway              = 502
     StatusServiceUnavailable      = 503
     StatusGatewayTimeout          = 504
	 StatusHTTPVersionNotSupported = 505
	 ************************************************************************/

func (ws *webserver) Start() {
	r := ws.router()

	logWriter.Info("starting http server", zap.String("addr", ws.addr))
	//log.Fatal(http.ListenAndServe(ws.addr, r))
	logWriter.Fatal(http.ListenAndServe(ws.addr, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r)))

}

func (ws *webserver) router() *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	r.NotFoundHandler = handler(notFoundHandler)

	// handle /ping for convenience. we'll also handle /api/v1/ping with the same function.
	r.HandleFunc("/ping", handler(ws.handlePing)).Methods("GET")

	apiv1 := r.PathPrefix("/api/v1").Subrouter()

	apiv1.HandleFunc("/ping", handler(ws.handlePing)).Methods("GET")

	// apiv1.HandleFunc("/contacts", handler(ws.handleGetContacts)).Methods("GET")
	// apiv1.HandleFunc("/contacts/{contactID}", handler(ws.handleGetContact)).Methods("GET")
	// apiv1.HandleFunc("/contacts", handler(ws.handlePostContact)).Methods("POST")
	// apiv1.HandleFunc("/contacts/{contactID}", handler(ws.handlePutContact)).Methods("PUT")
	// apiv1.HandleFunc("/contacts/{contactID}", handler(ws.handleDeleteContact)).Methods("DELETE")

	// apiv1.HandleFunc("/contacts/{contactID}/addresses", handler(ws.handleGetContactAddresses)).Methods("GET")
	// apiv1.HandleFunc("/contacts/{contactID}/addresses/{addressID}", handler(ws.handleGetContactAddress)).Methods("GET")
	// apiv1.HandleFunc("/contacts/{contactID}/addresses", handler(ws.handlePostContactAddresses)).Methods("POST")
	// apiv1.HandleFunc("/contacts/{contactID}/addresses/{addressID}", handler(ws.handlePutContactAddress)).Methods("PUT")
	// apiv1.HandleFunc("/contacts/{contactID}/addresses/{addressID}", handler(ws.handleDeleteContactAddress)).Methods("DELETE")
	// For Authentic API's:-
	apiv1.Handle("/projects", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handleGetProjects)))).Methods("GET", "OPTIONS")

	// Parcel / Property information
	// Property Search by Address or APN (County Assessor parcel number)
	apiv1.Handle("/property/{searchText}", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handleGetSearchProperty)))).Methods("GET", "OPTIONS")
	// Parcel Information by coordinates
	apiv1.Handle("/property/parcelByCoords/geo", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handleGetSearchPropertyByCoords)))).Queries("lat", "{lat}", "lon", "{lon}").Methods("GET", "OPTIONS")

	// Code Search URLs
	// Get the list of sections that match a search term /q?codeType={codeType}
	apiv1.Handle("/codeSections/{jurisdiction}/list", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handleGetCodeSectionList)))).Queries("codeType", "{codeType}", "searchText", "{searchText}").Methods("GET", "OPTIONS")
	// Get the content to display for a code section
	apiv1.Handle("/codeSections/{jurisdiction}/content", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handleGetCodeSectionContent)))).Queries("codeType", "{codeType}", "sectionNumber", "{sectionNumber}").Methods("GET", "OPTIONS")

	apiv1.Handle("/projects/{projectID}", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handleGetProject)))).Methods("GET", "OPTIONS")

	apiv1.Handle("/user/projects/recent", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handleGetRecentProject)))).Methods("GET", "OPTIONS")

	apiv1.Handle("/recent/activity", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handleGetRecentActivity)))).Methods("GET", "OPTIONS")

	// apiv1.Handle("/notes", negroni.New(
	// 	negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
	// 	negroni.Wrap(handler(ws.handleGetNote)))).Methods("GET","OPTIONS")

	apiv1.Handle("/users", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handleGetUserMiddleware)))).Methods("GET", "OPTIONS")
	apiv1.Handle("/users", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handlePutUser)))).Methods("PUT", "OPTIONS")

	apiv1.Handle("/invite/resetUserRoles", negroni.New(
			negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
			negroni.Wrap(handler(ws.handleResetUserRoles)))).Methods("PUT", "OPTIONS")
	// apiv1.Handle("/user/{userID}", negroni.New(
	// 	negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
	// 	negroni.Wrap(handler(ws.handleGetUser)))).Methods("GET", "OPTIONS")

	apiv1.Handle("/projects", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handlePostProject)))).Methods("POST", "OPTIONS")
	apiv1.Handle("/projects/{projectID}/properties/{propertyID}", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handleDeleteProjectProperty)))).Methods("DELETE", "OPTIONS")
	apiv1.Handle("/groups/{projectID}", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handlePostGroup)))).Methods("POST", "OPTIONS")
	// apiv1.Handle("/codeCollection/{codeID}", negroni.New(
	// 	negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
	// 	negroni.Wrap(handler(ws.handlePostCodeCollection)))).Methods("POST", "OPTIONS")
	apiv1.Handle("/files", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handlePostFiles)))).Methods("POST", "OPTIONS")

	apiv1.HandleFunc("/codeCollection", handler(ws.handlePostCodeCollection)).Methods("POST", "OPTIONS")

	apiv1.HandleFunc("/codeCollections/{codeCollectionID}/codeLibrary/{codeLibraryID}", handler(ws.handleCreateCodeCollectionsFromCodeLibrary)).Methods("POST", "OPTIONS")

	// apiv1.Handle("/projects/{projectID}", negroni.New(
	// 	negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
	// 	negroni.Wrap(handler(ws.handlePutProject)))).Methods("PUT","OPTIONS")

	apiv1.Handle("/projects/{projectID}", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handleDeleteProject)))).Methods("DELETE", "OPTIONS")
	apiv1.Handle("/project/duplicate/{projectID}", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handleDuplicateProject)))).Methods("POST", "OPTIONS")
	apiv1.Handle("/duplicate/{projectID}/property/{propertyID}", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handleDuplicateProperty)))).Methods("POST", "OPTIONS")
	apiv1.Handle("/users", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handlePostUser)))).Methods("POST", "OPTIONS")
	apiv1.Handle("/user_codes/{codeCollectionID}", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handlPostCodes)))).Methods("POST", "OPTIONS")
	apiv1.Handle("/library_codes/{codeLibraryID}", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handlePostCodesToCodeLibrary)))).Methods("POST", "OPTIONS")
	apiv1.Handle("/user_codes", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handleGetCodes)))).Methods("GET", "OPTIONS")
	apiv1.Handle("/invite/{inviteID}", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handleAcceptInvitation)))).Methods("PUT", "OPTIONS")
	apiv1.Handle("/invite/{inviteID}/group/{groupId}", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handleAcceptInvitationForGroup)))).Methods("PUT", "OPTIONS")

	apiv1.Handle("/codeLibrary", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handlePostCodeLibrary)))).Methods("POST", "OPTIONS")
	apiv1.Handle("/codeLibrary", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handleGetAllLibrary)))).Methods("GET", "OPTIONS")
	apiv1.Handle("/invite/project/{projectID}/group/{groupID}", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handleSendInvitationEmailToGroup)))).Methods("POST", "OPTIONS")

	apiv1.Handle("/company/detail", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handlePostCompany)))).Methods("POST", "OPTIONS")
	apiv1.Handle("/company", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handlePutCompanyDetails)))).Methods("PUT", "OPTIONS")
	apiv1.Handle("/company", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handleGetCompanyDetails)))).Methods("GET", "OPTIONS")
	apiv1.Handle("/company", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handlePostCompanyDetails)))).Methods("POST", "OPTIONS")

	apiv1.Handle("/resetPassword/{userID}", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handleResetPassword)))).Methods("GET", "OPTIONS")

	apiv1.Handle("/productTour", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handleUpdateIsIntro)))).Methods("PUT", "OPTIONS")

	apiv1.Handle("/invitation/{ProjectID}", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handleGetManageCollaborator)))).Methods("GET", "OPTIONS")
 

	// apiv1.Handle("/invite/project/{projectID}", negroni.New(
	// 	negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
	// 	negroni.Wrap(handler(ws.handleSendInvitationEmail)))).Methods("POST", "OPTIONS")
	// apiv1.Handle("/projects/{projectID}/property/{propertyID}/notes/{noteID}", negroni.New(
	// 	negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
	// 	negroni.Wrap(handler(ws.handlDeleteNote)))).Methods("DELETE","OPTIONS")
	// apiv1.Handle("/projects/{projectID}/property/{propertyID}/notes", negroni.New(
	// 	negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
	// 	negroni.Wrap(handler(ws.handlPostNote)))).Methods("POST","OPTIONS")
	// apiv1.Handle("/projects/{projectID}/properties", negroni.New(
	// 	negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
	// 	negroni.Wrap(handler(ws.handlePostProperty)))).Methods("POST","OPTIONS")
	apiv1.Handle("/projects/{projectID}/properties", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handleGetProperties)))).Methods("GET","OPTIONS")
	// apiv1.Handle("/property/{propertyId}/notes", negroni.New(
	// 		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
	// 		negroni.Wrap(handler(ws.handleGetNote)))).Methods("GET","OPTIONS")

	// NonAuthentic API's:-

	
	apiv1.HandleFunc("/remove/user/{inviteGroupID}", handler(ws.handlDeleteInviteGroupAssociations)).Methods("DELETE", "OPTIONS")

	apiv1.HandleFunc("/codeLibrary/{codeLibraryID}", handler(ws.handleGetCodeLibrary)).Methods("GET", "OPTIONS")
	apiv1.HandleFunc("/codeLibrary/{codeLibraryID}", handler(ws.handleDeleteCodeLibrary)).Methods("DELETE", "OPTIONS")

	apiv1.HandleFunc("/codeCollections", handler(ws.handleGetCollection)).Methods("GET", "OPTIONS")

	// apiv1.HandleFunc("/projects/{projectID}/properties", handler(ws.handleGetProperties)).Methods("GET", "OPTIONS")
	apiv1.HandleFunc("/property/{propertyID}/notes", handler(ws.handlePostNote)).Methods("POST", "OPTIONS")
	apiv1.HandleFunc("/projects/{projectID}/property/{propertyID}/notes/{noteID}", handler(ws.handleDeleteNote)).Methods("DELETE", "OPTIONS")
	apiv1.HandleFunc("/group/{groupID}/properties", handler(ws.handlePostPropertyToGroup)).Methods("POST", "OPTIONS")
	// apiv1.HandleFunc("/projects", handler(ws.handlePostProject)).Methods("POST", "OPTIONS")
	// apiv1.HandleFunc("/projects/{projectID}", handler(ws.handleGetProject)).Methods("GET", "OPTIONS")
	// apiv1.HandleFunc("/codeCollectionAssociation", handler(ws.handlePostCodeCollectionAssociation)).Methods("POST")

	apiv1.HandleFunc("/property/{propertyID}/notes", handler(ws.handleGetNote)).Methods("GET", "OPTIONS")
	apiv1.HandleFunc("/note/{noteID}", handler(ws.handlePutNote)).Methods("PUT", "OPTIONS")
	//apiv1.HandleFunc("/properties/{propertyID}", handler(ws.handleGetProperty)).Methods("GET", "OPTIONS")
	apiv1.Handle("/properties/{propertyID}", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handleGetProperty)))).Methods("GET", "OPTIONS")
	// apiv1.HandleFunc("/projects", handler(ws.handleGetProjects)).Methods("GET", "OPTIONS")
	//apiv1.HandleFunc("/projects/{projectID}", handler(ws.handlePutProject)).Methods("PUT", "OPTIONS")

	apiv1.Handle("/projects/{projectID}", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handlePutProject)))).Methods("PUT", "OPTIONS")

	apiv1.Handle("/inviter/details/{projectID}",negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handleGetInviterDetails)))).Methods("POST", "OPTIONS")



	apiv1.HandleFunc("/user/{userID}", handler(ws.handleGetUser)).Methods("GET", "OPTIONS")
	apiv1.HandleFunc("/user/{userID}", handler(ws.handleDeleteUser)).Methods("DELETE", "OPTIONS")

	// apiv1.HandleFunc("/users", handler(ws.handlePostUser)).Methods("POST", "OPTIONS")
	apiv1.HandleFunc("/property/{propertyID}", handler(ws.handlePutProperty)).Methods("PUT", "OPTIONS")
	// apiv1.HandleFunc("/project/duplicate/{projectID}", handler(ws.handleDuplicateProject)).Methods("POST", "OPTIONS")
	// apiv1.HandleFunc("/property/duplicate/{propertyID}", handler(ws.handleDuplicateProperty)).Methods("POST", "OPTIONS")
	apiv1.HandleFunc("/invite/project/{projectID}", handler(ws.handleSendInvitationEmail)).Methods("POST", "OPTIONS")
	apiv1.HandleFunc("/invite/resend/{inviteID}", handler(ws.handleReSendInvitationEmail)).Methods("POST", "OPTIONS")
	
	apiv1.HandleFunc("/invite/resend/{inviteID}/group/{groupId}", handler(ws.handleResendInvitationEmailToGroup)).Methods("POST", "OPTIONS")

	apiv1.HandleFunc("/collaborator/{collaboratorID}", handler(ws.handlDeleteInviteCollaborator)).Methods("DELETE", "OPTIONS")
	//---------------------
	apiv1.HandleFunc("/invitation/{inviteID}", handler(ws.handleAcceptInviteCollaborator)).Methods("PUT", "OPTIONS")
	// apiv1.HandleFunc("/invitation/{projectID}", handler(ws.handleGetManageCollaborator)).Methods("GET", "OPTIONS")
	apiv1.HandleFunc("/send/password-reset", handler(ws.handleSendResetPasswordToCollaborator)).Methods("POST", "OPTIONS")

	// apiv1.HandleFunc("/invite/{inviteID}\\", handler(ws.handleAcceptInvitation)).Methods("GET", "OPTIONS")
	// ---------------------
	// apiv1.HandleFunc("/invitation/{invitationID}", handler(ws.handleAcceptInviteCollaborator)).Methods("PUT", "OPTIONS")
	// apiv1.HandleFunc("/invitation/{ProjectID}", handler(ws.handleGetManageCollaborator)).Methods("GET", "OPTIONS")

	apiv1.Handle("/invite/email/project/{ProjectID}", negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler(ws.handleGetInviteCollaboratorByEmail)))).Methods("POST", "OPTIONS")
	// ------------------
	apiv1.HandleFunc("/submit/feedback", handler(ws.handlePostSubmitFeedback)).Methods("POST", "OPTIONS")
	apiv1.HandleFunc("/search_codes", handler(ws.handleGetSearchCodes)).Methods("GET", "OPTIONS")
	apiv1.HandleFunc("/codes/{codeID}/notes", handler(ws.handlPostCodesNote)).Methods("POST", "OPTIONS")
	apiv1.HandleFunc("/codes/{codeID}/notes", handler(ws.handleGetCodeNote)).Methods("GET", "OPTIONS")
	apiv1.HandleFunc("/codes/notes/{noteID}", handler(ws.handlePutCodeNote)).Methods("PUT", "OPTIONS")
	apiv1.HandleFunc("/code/{codeID}/note/{noteID}", handler(ws.handlDeleteCodeNote)).Methods("DELETE", "OPTIONS")
	apiv1.HandleFunc("/note/{noteID}/noteThreads", handler(ws.handleGetCodeNoteComment)).Methods("GET", "OPTIONS")
	apiv1.HandleFunc("/note/{noteID}/noteThreads", handler(ws.handlPostCodesNoteComment)).Methods("POST", "OPTIONS")
	apiv1.HandleFunc("/codes/notes/{noteID}", handler(ws.handleGetCodesNote)).Methods("GET", "OPTIONS")
	apiv1.HandleFunc("/groups/{groupID}", handler(ws.handleDeleteGroup)).Methods("DELETE", "OPTIONS")
	apiv1.HandleFunc("/groups/{groupID}", handler(ws.handleGetGroup)).Methods("GET", "OPTIONS")
	apiv1.HandleFunc("/groups/{groupID}", handler(ws.handlePutGroup)).Methods("PUT", "OPTIONS")
	apiv1.HandleFunc("/codeCollection/{codeCollectionID}", handler(ws.handleDeleteCodeCollection)).Methods("DELETE", "OPTIONS")
	apiv1.HandleFunc("/codeCollection/{codeCollectionID}", handler(ws.handleGetCodeCollection)).Methods("GET", "OPTIONS")
	apiv1.HandleFunc("/codeCollection/{codeCollectionID}", handler(ws.handlePutCodeCollection)).Methods("PUT", "OPTIONS")
	apiv1.HandleFunc("/projectGroup/{projectGroupID}", handler(ws.handleGetProjectGroupAssociations)).Methods("GET", "OPTIONS")
	apiv1.HandleFunc("/projectGroup/{projectGroupID}", handler(ws.handleDeleteProjectGroupAssociations)).Methods("DELETE", "OPTIONS")
	apiv1.HandleFunc("/codeCollectionAssociation/{codeCollectionAssociationID}", handler(ws.handleGetCodeCollectionAssociations)).Methods("GET", "OPTIONS")
	apiv1.HandleFunc("/codeCollectionAssociation/{codeCollectionAssociationID}", handler(ws.handleDeleteCodeCollectionAssociations)).Methods("DELETE", "OPTIONS")
	apiv1.HandleFunc("/codeGroupsAssociation", handler(ws.handlePostCodeGroupAssociations)).Methods("POST", "OPTIONS")
	apiv1.HandleFunc("/codeGroupsAssociation/{codeGroupsAssociationID}", handler(ws.handleGetCodeGroupAssociations)).Methods("GET", "OPTIONS")
	apiv1.HandleFunc("/codeGroupsAssociation/{codeGroupsAssociationID}", handler(ws.handleDeleteCodeGroupAssociations)).Methods("DELETE", "OPTIONS")
	apiv1.HandleFunc("/propertyGroupsAssociation", handler(ws.handlePostPropertyGroupAssociations)).Methods("POST", "OPTIONS")
	apiv1.HandleFunc("/propertyGroupsAssociation/{propertyGroupsAssociationID}", handler(ws.handleGetPropertyGroupAssociations)).Methods("GET", "OPTIONS")
	apiv1.HandleFunc("/propertyGroupsAssociation/{propertyGroupsAssociationID}", handler(ws.handleDeletePropertyGroupAssociations)).Methods("DELETE", "OPTIONS")
	apiv1.HandleFunc("/files/{filesID}", handler(ws.handleGetFiles)).Methods("GET", "OPTIONS")
	apiv1.HandleFunc("/files/{filesID}", handler(ws.handlePutFiles)).Methods("PUT", "OPTIONS")
	apiv1.HandleFunc("/files/{filesID}", handler(ws.handleDeleteFiles)).Methods("DELETE", "OPTIONS")
	apiv1.HandleFunc("/filesGroup", handler(ws.handlePostFileGroupAssociations)).Methods("POST", "OPTIONS")
	apiv1.HandleFunc("/filesGroup/{filesGroupID}", handler(ws.handleGetFileGroupAssociations)).Methods("GET", "OPTIONS")
	apiv1.HandleFunc("/filesGroup/{filesGroupID}", handler(ws.handleDeleteFileGroupAssociations)).Methods("DELETE", "OPTIONS")
	apiv1.HandleFunc("/groupCodeCollection", handler(ws.handlePostGroupCodeAssociations)).Methods("POST", "OPTIONS")
	apiv1.HandleFunc("/groupCodeCollection/{groupCodeCollectionID}", handler(ws.handleGetGroupCodeAssociations)).Methods("GET", "OPTIONS")
	apiv1.HandleFunc("/groupCodeCollection/{groupCodeCollectionID}", handler(ws.handleDeleteGroupCodeAssociations)).Methods("DELETE", "OPTIONS")
	apiv1.HandleFunc("/groupUser", handler(ws.handlePostGroupUserAssociations)).Methods("POST", "OPTIONS")
	apiv1.HandleFunc("/groupUser/{groupUserID}", handler(ws.handleGetGroupUserAssociations)).Methods("GET", "OPTIONS")
	apiv1.HandleFunc("/groupUser/{groupUserID}", handler(ws.handleDeleteGroupUserAssociations)).Methods("DELETE", "OPTIONS")
	apiv1.HandleFunc("/group/{groupID}/collection/{collectionID}", handler(ws.handlePostCodeCollectionsToGroup)).Methods("POST", "OPTIONS")
	apiv1.HandleFunc("/duplicate/project/{projectID}/group/{groupID}", handler(ws.handleDuplicateGroup)).Methods("POST", "OPTIONS")
	apiv1.HandleFunc("/groups/{groupID}/properties/{propertyID}", handler(ws.handleDeleteGroupProperty)).Methods("DELETE", "OPTIONS")
	apiv1.HandleFunc("/groups/{groupID}/codeCollections/{codeCollectionID}", handler(ws.handleDeleteGroupCodeCollection)).Methods("DELETE", "OPTIONS")
	apiv1.HandleFunc("/groups/{groupID}/files/{fileID}", handler(ws.handleDeleteGroupFiles)).Methods("DELETE", "OPTIONS")
	apiv1.HandleFunc("/user/details/{email}", handler(ws.handleGetUserDetailsByEmail)).Methods("GET", "OPTIONS")
	apiv1.HandleFunc("/invite/invitedUser/{projectID}", handler(ws.handleGetInvitiedUsers)).Methods("GET", "OPTIONS")

	apiv1.HandleFunc("/roles", handler(ws.handleGetRoles)).Methods("GET", "OPTIONS")
	apiv1.HandleFunc("/sendEmailToSupportTeam", handler(ws.handleSendEmailToSupportTeam)).Methods("POST", "OPTIONS")
	// apiv1.HandleFunc("/company", handler(ws.handlePutCompanyDetails)).Methods("PUT", "OPTIONS")
	apiv1.HandleFunc("/demoProject", handler(ws.handleGetDemoProjects)).Methods("GET", "OPTIONS")
	apiv1.HandleFunc("/demoProjects/{projectID}", handler(ws.handleGetDemoProject)).Methods("GET", "OPTIONS")
	// apiv1.HandleFunc("/invite/resetUserRoles", handler(ws.handleResetUserRoles)).Methods("PUT", "OPTIONS")

	// apiv1.HandleFunc("/save/selectedCode", handler(ws.handlPostCodes)).Methods("POST", "OPTIONS")
	// apiv1.HandleFunc("/company/detail", handler(ws.handlePostCompany)).Methods("POST", "OPTIONS")
	// apiv1.HandleFunc("/projects/{projectID}", handler(ws.handleDeleteProject)).Methods("DELETE", "OPTIONS")

	return r
}

// @Summary Ping server
// @Produce json
// @Success 200 {object} models.PingResponse
// @Router /ping [get]
func (ws *webserver) handlePing(w http.ResponseWriter, r *http.Request) error {
	return Ok(w, models.PingResponse{Message: "pong"})
}

// summary get recent project
func (ws *webserver) handleGetRecentProject(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", "https://app.permitdocs.com")

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	projects, err := ws.db.Projects.GetTop(middleware.UserId)
	fmt.Println(projects, "projects")

	if err != nil {
		return err
	}

	response := make([]models.ProjectResponse, 0)
	for _, project := range projects {
		if err != nil {
			return err
		}
		response = append(response, models.MapProjectResponse(project, nil, nil, 3, ""))
	}
	return Ok(w, response)
}

// get search property - depricated - function is in addressSearch.go
/*
func (ws *webserver) handleGetSearchProperty(w http.ResponseWriter, r *http.Request) error {
	// get url params

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["search"])
	if err != nil {
		return &invalidRequest{}
	}

	response, err := http.Get("https://scoutred.com/api/search/addresses?q=" + strconv.Itoa(id))
	fmt.Println("https://scoutred.com/api/search/addresses?q=" + strconv.Itoa(id))
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject models.ScoutRedResponse
	json.Unmarshal(responseData, &responseObject)

	bodyString := string(responseData)

	fmt.Println(bodyString)

	return Ok(w, responseObject)
}
*/

// summary get all detail
func (ws *webserver) handleGetDetail(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get Details
	details, err := ws.db.Details.GetAll()

	if err != nil {
		return err
	}

	response := make([]models.DetailResponse, 0)
	for _, detail := range details {
		if err != nil {
			return err
		}
		response = append(response, models.MapDetail(detail))
	}

	return Ok(w, response)

}

// get Recent Activity
func (ws *webserver) handleGetRecentActivity(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get Recent Activity
	recentactivity, err := ws.db.RecentsActivity.GetAll()

	if err != nil {
		return err
	}

	response := make([]models.RecentActivityResponse, 0)
	for _, activity := range recentactivity {
		if err != nil {
			return err
		}
		response = append(response, models.MapRecentActivity(activity))
	}

	return Ok(w, response)
}

// summary get Notes
// func (ws *webserver) handleGetNote(w http.ResponseWriter, r *http.Request) error {

// 	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

// 	(w).Header().Set("Access-Control-Allow-Origin", "*")
//     (w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
// 	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
// 	// get all Notes
// 	notes, err := ws.db.Notes.GetAll()

// 	if err != nil {
// 		return err
// 	}

// 	response := make([]models.NoteResponse, 0)
// 	for _, note := range notes {
// 		if err != nil {
// 			return err
// 		}
// 		response = append(response, models.MapNoteResponse(note))
// 	}

// 	return Ok(w, response)

// }

func (ws *webserver) handleGetUserMiddleware(w http.ResponseWriter, r *http.Request) error {
	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	users, err := ws.db.Users.Get(middleware.UserId)
	if err != nil {
		return err
	}

	response := models.MapUserResponse(users)

	return Ok(w, response)
}

// create a new project
func (ws *webserver) handlePostProject(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// create var ready to hold decoded json from body
	var request models.ProjectCreateRequest

	// decode body
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&request); err != nil {
		return &invalidRequest{}
	}

	// create contact
	create := models.MapCreateProjectRequest(request, middleware.UserId)
	project, err := ws.db.Projects.Create(create)
	if err != nil {
		return &notFound{"Not able to Create"}
	}

	userRoleAssociate := models.MapProjectUserRoleAssociations(project.CreatedBy, project.ID, 1, true)
	projectUserRole, err := ws.db.ProjectUserRoleAssociations.Create(userRoleAssociate)

	fmt.Println(project, "project")
	fmt.Println(projectUserRole, "projectUserRole")

	// create response
	response := models.MapProjectResponse(project, nil, nil, 0, "")

	return Ok(w, response)

}

// delete project's property
func (ws *webserver) handleDeleteProjectProperty(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get url params
	vars := mux.Vars(r)
	projectID, err := models.MapDecodeRequest(vars["projectID"])

	if err != nil {
		return &invalidRequest{}
	}

	propertyID, err := strconv.Atoi(vars["propertyID"])
	if err != nil {
		return &invalidRequest{}
	}

	results, err := ws.db.ProjectPropertyAssociations.GetAllByProjectID(projectID)
	if err != nil {
		return err
	}

	fmt.Println(results, "results")

	for _, result := range results {

		if result.PropertyID == propertyID {
			properties, err := ws.db.Properties.Get(result.PropertyID)
			if err != nil {
				return err
			} else {
				ws.db.Properties.Delete(properties.ID)
				ws.db.ProjectPropertyAssociations.Delete(result.ID)
			}

		}
	}
	// struct{}{} is an empty object, returns "{}" to the client
	return Ok(w, "Property has been deleted")
}

// get user details using user ID

func (ws *webserver) handleGetUser(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get url params
	vars := mux.Vars(r)
	fmt.Println(vars, "vars")

	q := strconv.Quote(vars["userID"])
	if q == "nil" {
		// fmt.Println(err,"checkerror")
		return &invalidRequest{}
	}

	userID, err := strconv.Unquote(q)

	// get user
	user, err := ws.db.Users.Get(userID)
 
	if err != nil {
		return err
	}

	// response := make([]models.UserResponse, 0)

	results, err := ws.db.ProjectUserRoleAssociations.GetAllByUserID(user.ID)
	if err != nil {
		return err
	}

	projects := make([]database.Project, 0)

	for _, result := range results {
		project, err := ws.db.Projects.Get(result.ProjectID)

		if err != nil {
			return err
		}

		fmt.Println(projects, "projects")

		projects = append(projects, project)
	}
	response :=  models.MapUserResponse(user)

	return Ok(w, response)
}

// update project name
// apiv1.HandleFunc("/projects/{projectID}", handler(ws.handlePutProject)).Methods("PUT", "OPTIONS")
func (ws *webserver) handlePutProject(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// get url params
	vars := mux.Vars(r)
	projectID, err := models.MapDecodeRequest(vars["projectID"])

	// create var ready to hold decoded json from body
	var request models.ProjectNameUpdateRequest

	// decode body
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&request)

	if err != nil {
		return &invalidRequest{}
	}

	// update contact
	update := models.MapUpdateProjectRequest(projectID, request)

	project, err := ws.db.Projects.Update(update)

	if err != nil {
		return err
	}

	// create response
	response := models.MapProjectResponse(project, nil, nil, 0, "")

	return Ok(w, response)
}

// Delete Project
func (ws *webserver) handleDeleteProject(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get url params
	vars := mux.Vars(r)
	projectID, err := models.MapDecodeRequest(vars["projectID"])

	if err != nil {
		return &invalidRequest{}
	}

	fmt.Println(projectID, "projectID")

	allproject, err := ws.db.Projects.Get(projectID)

	fmt.Println("allproject", allproject)
	if err != nil {
		return &notFound{"projectID doesn't exist"}
	}
	fmt.Println(allproject.CreatedBy, "allproject.CreatedBy")
	if allproject.CreatedBy == middleware.UserId {
		fmt.Println(allproject.CreatedBy, "start")

		// ws.db.Projects.Delete(projectID, middleware.UserId)

		results, err := ws.db.ProjectPropertyAssociations.GetAllByProjectID(projectID)
		fmt.Println("results", results)

		if len(results) == 0 {
			fmt.Println(projectID, "projectID")
			ws.db.Projects.Delete(projectID, middleware.UserId)

		} else {
			for _, result := range results {
				allproperty, err := ws.db.Properties.Get(result.PropertyID)
				fmt.Println(allproperty, "allproperty")
				if err != nil {
					ws.db.Projects.Delete(result.ProjectID, middleware.UserId)
					ws.db.ProjectPropertyAssociations.Delete(result.ID)
					if err != nil {
						return &notFound{"projectID doesn't Found in PropertyAssociaet"}
					}
				} else {
					propertyNotes, err := ws.db.PropertyNoteAssociations.Get(allproperty.ID)
					fmt.Println(propertyNotes, "propertyNotes")
					if err != nil {
						ws.db.Projects.Delete(result.ProjectID, middleware.UserId)
						ws.db.Properties.Delete(result.PropertyID)
						ws.db.ProjectPropertyAssociations.Delete(result.ID)

					} else {
						ws.db.Notes.Delete(propertyNotes.ID)
						ws.db.PropertyNoteAssociations.Delete(allproperty.ID)
						ws.db.Properties.Delete(allproperty.ID)
						ws.db.ProjectPropertyAssociations.Delete(result.ID)
						ws.db.Projects.Delete(projectID, middleware.UserId)
					}

				}

			}
		}

		allProjects, err := ws.db.ProjectUserRoleAssociations.GetAllByUserProjectID(projectID)
		if err != nil {
			return err
		}
		fmt.Println(allProjects, "allProjects")

		for _, result := range allProjects {
			ws.db.ProjectUserRoleAssociations.Delete(result.ID)
		}

	} else {
		return &notFound{"only Admin can delete projects"}
	}
	return Ok(w, "Project has been deleted")
}

// Add note to property
func (ws *webserver) handlePostNote(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get url params
	vars := mux.Vars(r)

	PropertyID, err := models.MapDecodeRequest(vars["propertyID"])
	if err != nil {
		return &invalidRequest{}
	}

	// create var ready to hold decoded json from body
	var request models.NoteRequest

	// decode body
	dec := json.NewDecoder(r.Body)
	fmt.Println(dec, "dec")
	fmt.Println(json.NewDecoder(r.Body), "body")

	if err := dec.Decode(&request); err != nil {
		return &invalidRequest{}
	}

	// update contact
	create := models.MapNoteRequest(PropertyID, middleware.UserId, request)

	// check is project exist.
	result, err := ws.db.Properties.Get(PropertyID)
	if err != nil {
		return &notFound{"propertyID doesn't match"}
	}

	fmt.Println(result, "resultsresultsresults")
	response := make([]models.NoteResponse, 0)

	// if len(results) > 0 {
	// 	for _, result := range results {

	if result.ID == PropertyID {
		note, err := ws.db.Notes.Create(create)
		if err != nil {
			return err
		} else {
			propertyNotes := models.MapPropertyNoteRequest(note.ID, PropertyID)
			propertyNoteAssociate, err := ws.db.PropertyNoteAssociations.Create(propertyNotes)
			if err != nil {
				return err
			}

			fmt.Println("ashish", propertyNoteAssociate)
		}

		response = append(response, models.MapNoteResponse(note))
		// 	}
		// }
	}
	return Ok(w, response)
}

// Get Property WRT to ProjectID
func (ws *webserver) handleGetProperties(w http.ResponseWriter, r *http.Request) error {
	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get url params
	vars := mux.Vars(r)
	projectID, err := models.MapDecodeRequest(vars["projectID"])

	//This should handle invalid Decoding
	if err != nil {
		return &invalidRequest{}
	}

	// get project
	project, err := ws.db.Projects.Get(projectID)
	if err != nil {
		return &notFound{"ProjectID doesn't exist."}
	}

	projects, err := ws.db.ProjectPropertyAssociations.GetAllByProjectID(project.ID)
	if err != nil {
		return err
	}

	// var c database.Note
	response := make([]models.PropertyResponse, 0)
	for _, result := range projects {
		properties, err := ws.db.Properties.Get(result.PropertyID)
		fmt.Println(properties, "<<< properties")
		if err != nil {
			return &notFound{"PropertyID doesn't exist."}
		}

		propertynotes, err := ws.db.PropertyNoteAssociations.GetAllByPropertyID(properties.ID)
		if err != nil {
			return &notFound{"PropertyNoteAssociations error"}
		}
		note := make([]database.Note, 0)

		for _, propertynote := range propertynotes {
			fmt.Println(propertynote.NoteID, "propertynote.NoteID")
			getnote, err := ws.db.Notes.Get(propertynote.NoteID)

			if err != nil {

				return &notFound{"property error"}
			}
			note = append(note, getnote)
		}
		response = append(response, models.MapPropertyResponse(properties, note, project.ID, project.Name))
	}
	return Ok(w, response)
}

// Delete note
func (ws *webserver) handleDeleteNote(w http.ResponseWriter, r *http.Request) error {
	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get url params
	vars := mux.Vars(r)
	projectID, err := models.MapDecodeRequest(vars["projectID"])

	if err != nil {
		return &invalidRequest{}
	}
	propertyID, err := models.MapDecodeRequest(vars["propertyID"])
	if err != nil {
		return &invalidRequest{}
	}

	NoteID, err := models.MapDecodeRequest(vars["noteID"])

	if err != nil {
		return &invalidRequest{}
	}

	fmt.Println(projectID, propertyID, NoteID, "IDS")

	propertynotes, err := ws.db.PropertyNoteAssociations.GetAllByPropertyID(propertyID)
	if err != nil {
		return &notFound{"NoteID doesn't match"}
	}

	fmt.Println(propertynotes, "propertynotes")
	for _, propertyNote := range propertynotes {
		fmt.Println(propertyNote, "propertyNote")
		fmt.Println(propertyNote.PropertyID, "propertyNote.PropertyID")
		fmt.Println(propertyID, "propertyID")

		if propertyNote.NoteID == NoteID {
			fmt.Println(propertyNote.PropertyID, "propertyNote.AShish")

			err = ws.db.PropertyNoteAssociations.Delete(propertyNote.ID)
			if err != nil {
				return err
			} else {
				if err = ws.db.Notes.Delete(NoteID); err != nil {
					return err
				}
			}
		}

	}

	return Ok(w, "success")
}

// Get Note from PropertyID
func (ws *webserver) handleGetNote(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get url params

	vars := mux.Vars(r)
	fmt.Println(vars, "ashish")
	id, err := models.MapDecodeRequest(vars["propertyID"])
	if err != nil {
		return &invalidRequest{}
	}

	// get project
	properties, err := ws.db.PropertyNoteAssociations.GetAllByPropertyID(id)
	if err != nil {
		return &notFound{"Property doesn't exist."}
	}

	// var c database.Note
	response := make([]models.NoteResponse, 0)
	for _, result := range properties {
		notes, err := ws.db.Notes.Get(result.NoteID)

		if err != nil {
			return &notFound{"Note doesn't exist."}
		}

		response = append(response, models.MapNoteResponse(notes))
	}

	return Ok(w, response)
}

// Update Note
func (ws *webserver) handlePutNote(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get url params
	vars := mux.Vars(r)

	noteID, err := models.MapDecodeRequest(vars["noteID"])
	if err != nil {
		return &invalidRequest{}
	}

	fmt.Println(noteID, "check Id")

	// create var ready to hold decoded json from body
	var request models.NoteUpdateRequest

	// decode body
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&request); err != nil {
		return &invalidRequest{}
	}

	// update contact
	update := models.MapNoteUpdateRequest(noteID, request)

	fmt.Println(update, "updateNote")

	note, err := ws.db.Notes.Update(update)
	if err != nil {
		return err
	}

	fmt.Println(note, "note")
	// create response
	response := models.MapNoteResponse(note)

	return Ok(w, response)
}

// update user details
func (ws *webserver) handlePutUser(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// create var ready to hold decoded json from body
	var request models.UserUpdateRequest

	// decode body
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&request); err != nil {
		return &invalidRequest{}
	}

	update := models.MapUserUpdateRequest(middleware.UserId, request)

	user, err := ws.db.Users.UpdateUserDetails(update)
	if err != nil {
		return err
	}

	// create response
	response := models.MapUserResponse(user)
	fmt.Println(response, "response")

	return Ok(w, response)
}

func (ws *webserver) handleDeleteUser(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get url params
	vars := mux.Vars(r)
	Id := strconv.Quote(vars["userID"])
	if Id == "" {
		return &invalidRequest{}
	}

	fmt.Println(Id, "ID")

	err := ws.db.Users.Delete(Id)
	if err != nil {
		return &notFound{"UserId doesn't match"}
	}
	return Ok(w, "success")
}

// create a new user
func (ws *webserver) handlePostUser(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// create var ready to hold decoded json from body
	var request models.UserRequest

	// decode body
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&request); err != nil {
		return &invalidRequest{}
	}

	var response models.UserResponse

	// create user
	create := models.MapUserRequest(middleware.UserId, request)
	fmt.Println(create, "create")

	user, err := ws.db.Users.Get(middleware.UserId)
	if err != nil {
		user, err := ws.db.Users.Create(create)
		if err != nil {
			return err
		}
		response = models.MapUserResponse(user)
		return Ok(w, response)
	}

	response = models.MapUserResponse(user)

	return Ok(w, response)

}

//update property
func (ws *webserver) handlePutProperty(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get url params
	vars := mux.Vars(r)

	fmt.Println(vars, "vars")

	propertyID, err := models.MapDecodeRequest(vars["propertyID"])
	if err != nil {
		return &invalidRequest{}
	}

	fmt.Println(propertyID, " propertyID")

	// create var ready to hold decoded json from body
	var request models.UpdatePropertyRequest

	// decode body
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&request); err != nil {
		return &invalidRequest{}
	}
	projectProperty, err := ws.db.ProjectPropertyAssociations.GetAllByPropertyID(propertyID)
	if err != nil {
		return &notFound{"propertyID is  not associate with any project ID"}
	}

	var response models.PropertyResponse
	update := models.MapUpdatePropertyRequest(propertyID, middleware.UserId, request)

	property, err := ws.db.Properties.Update(update)
	if err != nil {
		return err
	}

	for _, projectPropertyList := range projectProperty {
		projectDetail, err := ws.db.Projects.Get(projectPropertyList.ProjectID)
		if err != nil {
			return &notFound{"PropertyNoteAssociations error"}
		}

		response = models.MapPropertyResponse(property, nil, projectDetail.ID, projectDetail.Name)
	}
	return Ok(w, response)
}

// Duplicate a Property
func (ws *webserver) handleDuplicateProperty(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	vars := mux.Vars(r)

	fmt.Println(vars, "vars")

	propertyID, err := models.MapDecodeRequest(vars["propertyID"])
	if err != nil {
		return &invalidRequest{}
	}

	projectID, err := models.MapDecodeRequest(vars["projectID"])

	if err != nil {
		return &invalidRequest{}
	}

	create := models.MapDuplicaetPropertyRequest(propertyID, middleware.UserId)

	property, err := ws.db.Properties.Duplicate(create)
	if err != nil {
		return &notFound{"Property can't be duplicate "}
	}

	projectProperty := models.MapPropertyProjectRequest(projectID, property.ID)

	projectPropertyAssociate, err := ws.db.ProjectPropertyAssociations.Create(projectProperty)
	if err != nil {
		return err
	}

	fmt.Println(projectPropertyAssociate, "projectPropertyAssociate")
	// create response
	response := models.MapPropertyResponse(property, nil, 0, "")

	return Ok(w, response)

}

// Accept Invitation
// func (ws *webserver) handleAcceptInviteCollaborator(w http.ResponseWriter, r *http.Request) error {

// 	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

// 	(w).Header().Set("Access-Control-Allow-Origin", "*")
// 	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
// 	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
// 	// get url params
// 	vars := mux.Vars(r)
// 	InvitationId, err := strconv.Atoi(vars["invitationID"])
// 	if err != nil {
// 		return &invalidRequest{}
// 	}

// 	// create var ready to hold decoded json from body
// 	var request models.InviteUpdateRequest

// 	// decode body
// 	dec := json.NewDecoder(r.Body)

// 	if err := dec.Decode(&request); err != nil {
// 		return &invalidRequest{}
// 	}

// 	update := models.MapInviteAcceptRequest(InvitationId, request)

// 	invite, err := ws.db.Invites.Update(update)
// 	if err != nil {
// 		return err
// 	}
// 	// create response
// 	response := models.MapInviteCollaboratorResponse(invite, nil)

// 	return Ok(w, response)
// }

// get ManageCollaborator detail


// func (ws *webserver) handleGetInviteCollaboratorByEmail(w http.ResponseWriter, r *http.Request) error {

// 	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

// 	(w).Header().Set("Access-Control-Allow-Origin", "*")
// 	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
// 	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
// 	// get url params
// 	vars := mux.Vars(r)
// 	fmt.Println(vars,"vars")
// 	projectID, err := strconv.Atoi(vars["ProjectID"])

// 	// create var ready to hold decoded json from body
// 	var request models.InviteCollaboratorRequest

// 	// // decode body
// 	dec := json.NewDecoder(r.Body)
// 	if err := dec.Decode(&request); err != nil {
// 		return &invalidRequest{}
// 	}

// 	fmt.Println(request.Email, "requestedEmail")

// 	response := make([]models.InviteCollaboratorResponse, 0)
// 	invited, err := ws.db.Invites.GetCollaboratorbyEmailAndProjectID(request.Email, projectID)
// 	if err != nil {
// 		return err
// 	}
// 	for _, invite := range invited {
		
// 		inviteGroupsAssociations, err := ws.db.InviteGroupAssociations.GetAllByInviteId(invite.ID)
// 		fmt.Println("inviteGroupsAssociations", inviteGroupsAssociations)
// 		if err != nil {
// 			return &invalidRequest{}
// 		}

// 		// group := make([]database.Group,0)
// 		// var group []database.Group
// 		var group []database.InviteGroupAssociations

// 		// for _, inviteGroupsAssociations := range inviteGroupsAssociations {
// 		// 	groups := ws.db.Groups.GetGroupForActive(inviteGroupsAssociations.GroupID)
// 		// 	fmt.Println(groups, "groups")

// 		// 	for _, activeGroup := range groups {
// 		// 		group = append(group, activeGroup)
// 		// 	}

// 		// }
	
// 		for _, inviteGroupsAssociations := range inviteGroupsAssociations {
// 				group = append(group, inviteGroupsAssociations)
// 		}
// 		response = append(response, models.MapInviteCollaboratorResponse(invite,group))
// 		fmt.Println(invited, "invited")
// 	}
// 	return Ok(w, response)
// }


// get project for both admin and collaboratot
func (ws *webserver) handleGetProjects(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// var response models.ProjectResponse

	response := make([]models.ProjectResponse, 0)
	results, err := ws.db.ProjectUserRoleAssociations.GetAll(middleware.UserId)

	if err != nil {
		return err
	}
	for _, result := range results {
		projects, err := ws.db.Projects.Get(result.ProjectID)

		properties, err := ws.db.ProjectPropertyAssociations.GetAllByProjectID(result.ProjectID)
		
		projectUserGroupAssociations, err := ws.db.ProjectUserGroupAssociations.GetAllByUserProjectID(middleware.UserId,result.ProjectID)
		if err != nil {
			return err
		}

		// BENJAMIN - TAKE A LOOK AT THE LOGIC HERE...
		// Pre merge
		// if result.RoleID == 1 {
		// 	projects.Role = "Admin"
		// }
		// if result.RoleID == 2 {
		// 	projects.Role = "Collaborator"
		// }
		// if result.RoleID == 3 {
		// 	projects.Role = "External User"
		// }

		// if err != nil {
		// 	return err

		// Post merge
		fmt.Println(projectUserGroupAssociations,"projectUserGroupAssociations:===1")

		fmt.Println(len(projectUserGroupAssociations),"len(projectUserGroupAssociations)")
		roleValue := 0

		if(len(projectUserGroupAssociations) == 0) {

			projectRole, err := ws.db.Roles.Get(result.RoleID)
			if err != nil {
				return err
			}
 				projects.Role = projectRole.Role
		}
		// end post-merge conflict

		for i := 0; i < len(projectUserGroupAssociations); i++ {

			projectRole, err := ws.db.Roles.Get(projectUserGroupAssociations[i].Role)
			if err != nil {
				return err
			}
			if projectUserGroupAssociations[i].Role == roleValue {	
				  projects.Role = projectRole.Role
				} else {
					roleValue = projectUserGroupAssociations[i].Role
					if (len(projectUserGroupAssociations) == 1) {
					projects.Role = projectRole.Role
					}else{
						projects.Role = "Various"

					}
				}
			fmt.Println(projectUserGroupAssociations[i],"projectUserGroupAssociations:===222")
		}	
		// for key, projectUserGroupAssociate := range projectUserGroupAssociations{

		// 	fmt.Println(projectUserGroupAssociations[key-1],"projectUserGroupAssociations:===1")

		// 	projectRole, err := ws.db.Roles.Get(projectUserGroupAssociate.Role)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	if projectUserGroupAssociate.Role == projectUserGroupAssociations[key-1].Role {	
		// 		projects.Role = projectRole.Role
		// 	} else {
		// 		projects.Role = "Various"
		// 	}
		// 	// if result.RoleID == 1 {
		// 	// 	projects.Role = "Admin"
		// 	// }
		// 	// if result.RoleID == 2 {
		// 	// 	projects.Role = "Collaborator"
		// 	// }
		// 	// if result.RoleID == 3 {
		// 	// 	projects.Role = "External User"
		// 	// }
		// }
			 
		userDetail, err := ws.db.Users.Get(projects.CreatedBy)

		response = append(response, models.MapProjectResponse(projects, nil, nil, len(properties), userDetail.ImageUrl))
	}

	return Ok(w, response)

}

// Submit Feedback
func (ws *webserver) handlePostSubmitFeedback(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	var request models.SubmitFeedbackRequest

	// decode body
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&request); err != nil {
		return &invalidRequest{}
	}
	fmt.Println(request, "request")
	from := "ashish.yadav@vicesoftware.com"
	pass := "9671761787"
	to := "ashish.yadav@vicesoftware.com"
	body := request.Description
	Subject := request.FeedbackOverview

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + Subject + "\n" +
		body + "\n\n"

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		fmt.Println("smtp error: %s", err)
		return err
	}

	fmt.Println("sent, visit http://foobarbazz.mailinator.com")

	return Ok(w, "success")
}

// create a company profile
func (ws *webserver) handlePostCompany(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// create var ready to hold decoded json from body
	var request models.CompanyRequest

	// decode body
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&request); err != nil {
		return &invalidRequest{}
	}

	fmt.Println(middleware.UserId, "middleware.UserId")
	create := models.MapCompanyRequest(middleware.UserId, request)
	fmt.Println(create, "create")

	company, err := ws.db.Companies.Create(create)
	if err != nil {
		return err
	}

	// create response
	response := models.MapCompanyResponse(company)

	return Ok(w, response)

}

//Get Roles
func (ws *webserver) handleGetRoles(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get all Notes
	roles, err := ws.db.Roles.GetAll()

	if err != nil {
		return err
	}

	response := make([]models.RoleResponse, 0)

	for _, role := range roles {

		response = append(response, models.MapRoleResponse(role))

	}
	return Ok(w, response)
}

// Get codes
func (ws *webserver) handleGetSearchCodes(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get all Notes
	codes, err := ws.db.Codes.GetAll()

	if err != nil {
		return err
	}

	response := make([]models.SearchCodesResponse, 0)

	for _, code := range codes {

		response = append(response, models.MapSearchCodes(code))

	}
	return Ok(w, response)
}

// Delete code note thread
func (ws *webserver) handlDeleteCodeNoteThread(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get url params
	vars := mux.Vars(r)
	codeID, err := models.MapDecodeRequest(vars["codeID"])
	if err != nil {
		return &invalidRequest{}
	}

	NoteID, err := models.MapDecodeRequest(vars["noteID"])
	if err != nil {
		return &invalidRequest{}
	}

	fmt.Println(codeID, NoteID, "IDS")

	codeNoteResults, err := ws.db.CodeNotesMetadataAssociations.GetAllByNoteID(NoteID)
	if err != nil {
		return &notFound{"NoteID doesn't match"}
	}

	fmt.Println(codeNoteResults, "codeNoteResults")
	for _, codeNoteResult := range codeNoteResults {
		fmt.Println(codeNoteResult, "codeNoteResult")
		fmt.Println(codeNoteResult.CodeID, "codeNoteResult.COdeID")
		fmt.Println(codeID, "codeID")

		if codeNoteResult.NoteID == NoteID {
			fmt.Println(codeNoteResult.CodeID, "codeNoteResult.CodeID")

			err = ws.db.CodeNotesMetadataAssociations.Delete(codeNoteResult.ID)
			if err != nil {
				return err
			} else {
				if err = ws.db.Codenotes.Delete(NoteID); err != nil {
					return err
				}
			}
		}
	}

	return Ok(w, "success")
}

func (ws *webserver) handleGetCodeCollectionAssociations(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get url params
	vars := mux.Vars(r)
	id := strconv.Quote(vars["codeCollectionAssociationID"])
	if id == "" {
		return &invalidRequest{}
	}

	fmt.Println(id, "id")
	check, err := strconv.Unquote(id)
	fmt.Println(check, err, "Checkid")

	// get code collection association detail
	codeCollectionDetail, err := ws.db.CodeCollectionAssociations.Get(check)
	if err != nil {
		return &notFound{"codeCollectionID doesn't found"}
	}
	return Ok(w, codeCollectionDetail)
}

// Delete record from code collection association
func (ws *webserver) handleDeleteCodeCollectionAssociations(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get url params
	vars := mux.Vars(r)
	id := strconv.Quote(vars["codeCollectionAssociationID"])
	if id == "nil" {
		return &invalidRequest{}
	}

	codeCollectionId, err := strconv.Unquote(id)
	fmt.Println(codeCollectionId, err, "error")

	if err = ws.db.CodeCollectionAssociations.Delete(codeCollectionId); err != nil {
		return err
	}

	// struct{}{} is an empty object, returns "{}" to the client
	return Ok(w, struct{}{})
}

func (ws *webserver) handleGetProjectGroupAssociations(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get url params
	vars := mux.Vars(r)
	id := strconv.Quote(vars["projectGroupID"])
	if id == "" {
		return &invalidRequest{}
	}

	fmt.Println(id, "id")
	checkId, err := strconv.Unquote(id)

	// get group
	projectGroupDetail, err := ws.db.ProjectGroupsAssociations.Get(checkId)
	if err != nil {
		return &notFound{"projectGroupID doesn't found"}
	}
	fmt.Println(projectGroupDetail, "projectGroupDetail")
	return Ok(w, projectGroupDetail)
}

// Delete record from project group association
func (ws *webserver) handleDeleteProjectGroupAssociations(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get url params
	vars := mux.Vars(r)
	id := strconv.Quote(vars["codeCollectionID"])
	if id == "nil" {
		return &invalidRequest{}
	}

	codecollectionId, err := strconv.Unquote(id)
	fmt.Println(codecollectionId, err, "error")

	if err = ws.db.ProjectGroupsAssociations.Delete(codecollectionId); err != nil {
		return err
	}
	return Ok(w, struct{}{})
}

// Update deatil of Groups
func (ws *webserver) handlePutGroup(w http.ResponseWriter, r *http.Request) error {
	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get url params
	vars := mux.Vars(r)
	groupsID := strconv.Quote(vars["groupID"])
	if groupsID == "" {
		return &invalidRequest{}
	}
	fmt.Println(groupsID, "groupsID")
	check, err := strconv.Unquote(groupsID)
	fmt.Println(check, err, "Checkid")

	// create var ready to hold decoded json from body
	var request models.GroupUpdateRequest
	fmt.Println(request, "request")
	// decode body
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&request); err != nil {
		return &invalidRequest{}
	}

	// update contact
	update := models.MapUpdateGroupRequest(check, request)
	groupUpdate, err := ws.db.Groups.Update(update)
	fmt.Println(groupUpdate, "groupUpdate")
	if err != nil {
		return err
	}

	// create response
	response := models.MapGroupResponse(groupUpdate, 0)
	fmt.Println(response, "response")

	return Ok(w, response)
}

func (ws *webserver) handlePostCodeCollection(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get url params

	// vars := mux.Vars(r)
	// CodeID, err := strconv.Atoi(vars["codeID"])
	// if err != nil {
	// 	return &invalidRequest{}
	// }

	// create var ready to hold decoded json from body
	var request models.CodeCollectionsRequest

	// decode body
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&request); err != nil {
		return &invalidRequest{}
	}
	fmt.Println(request, "request")

	// create code collection
	create := models.MapCreateCodeCollectionsRequest(middleware.UserId, request, true)
	codeCollection, err := ws.db.CodeCollections.Create(create)
	if err != nil {
		return err
	}

	// createCodeCollectionAssociations := models.MapCreateCodeCollectionAssociationRequest(CodeID, codeCollection.ID)
	// codeCollectionAssociate, err := ws.db.CodeCollectionAssociations.Create(createCodeCollectionAssociations)
	// if err != nil {
	// 	return err
	// }
	// fmt.Println(codeCollectionAssociate, "codeCollectionAssociate")
	// //response
	// response := models.MapCodeCollectionsResponse(codeCollection)
	return Ok(w, codeCollection)
}

// get codecollection details using code collection ID
func (ws *webserver) handleGetCodeCollection(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get url params
	fmt.Println("ashish")
	vars := mux.Vars(r)
	collectionId := strconv.Quote(vars["codeCollectionID"])
	if collectionId == "" {
		return &invalidRequest{}
	}

	collectionID, err := strconv.Unquote(collectionId)
	response := make([]models.CodesResponse, 0)

	fmt.Println(collectionID, "collectionID")
	// get collection detail
	codeCollectionsAssociate, err := ws.db.CodeCollectionAssociations.GetCodeByCollectionID(collectionID)
	if err != nil {
		return &notFound{"codeCollectionID doesn't found"}
	}

	fmt.Println(codeCollectionsAssociate, "codeCollectionsAssociate")

	//Get GroupId on the basis of Collection Id
	groupCollectionId, err := ws.db.GroupsCodeCollectionAssociations.GetGroupsByCollectionID(collectionID)
	if err != nil {
		return &notFound{"codeCollectionID doesn't found"}
	}

	fmt.Println(groupCollectionId, "groupCollectionId")
	//Get Project Id on the basis of Group Id
	projectGroupId, err := ws.db.ProjectGroupsAssociations.GetAllByGroupID(groupCollectionId.GroupID)
	if err != nil {
		return &notFound{"codeCollectionID doesn't found"}
	}

	fmt.Println(projectGroupId, "projectGroupId")
	for _, codeCollection := range codeCollectionsAssociate {

		fmt.Println(codeCollection, "codeCollection")

		code, err := ws.db.UserCodes.Get(codeCollection.CodeID)
		if err != nil {
			return &notFound{"codeCollectionID doesn't found"}
		}

		fmt.Println(code, "code")

		response = append(response, models.MapUserCodesResponse(code, projectGroupId.ProjectID))

	}

	return Ok(w, response)
}

// Update deatil of codeCollection
func (ws *webserver) handlePutCodeCollection(w http.ResponseWriter, r *http.Request) error {
	// get url params
	vars := mux.Vars(r)
	collectionID := strconv.Quote(vars["codeCollectionID"])
	if collectionID == "" {
		return &invalidRequest{}
	}
	fmt.Println(collectionID, "collectionID")
	id, err := strconv.Unquote(collectionID)
	fmt.Println(id, err, "Checkid")

	// create var ready to hold decoded json from body
	var request models.CodeCollectionsUpdateRequest
	fmt.Println(request, "request")
	// decode body
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&request); err != nil {
		return &invalidRequest{}
	}

	// update code collection
	update := models.MapUpdateCodeCollectionRequest(id, request)
	collectionUpdate, err := ws.db.CodeCollections.Update(update)
	if err != nil {
		return err
	}

	// create response
	response := models.MapCodeCollectionsResponse(collectionUpdate)
	fmt.Println(response, "response")

	return Ok(w, response)
}

// Soft Delete Code collection
func (ws *webserver) handleDeleteCodeCollection(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get url params
	vars := mux.Vars(r)
	codeId := strconv.Quote(vars["codeCollectionID"])
	if codeId == "" {
		return &invalidRequest{}
	}

	fmt.Println(codeId, "codeId ")
	checkID, err := strconv.Unquote(codeId)
	fmt.Println(checkID, err, "Checkid")

	codeCollections, err := ws.db.CodeCollections.Get(checkID)
	if err != nil {
		return &notFound{"codeCollectionID doesn't found"}
	}
	fmt.Println(codeCollections, "codeCollections")

	update := models.MapUpdateCodeCollectionStatus(checkID, false)
	fmt.Println(update, "update")

	updated, err := ws.db.CodeCollections.UpdateDelete(update)
	if err != nil {
		return err
	}
	fmt.Println(updated, "updated")

	return Ok(w, "Success")
}

//Post the detail in project group association
func (ws *webserver) handlePostCodeGroupAssociations(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get url params

	// create var ready to hold decoded json from body
	var request models.CodeGroupAssociationRequest

	// decode body
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&request); err != nil {
		return &invalidRequest{}
	}
	fmt.Println(request, "request")

	// create Project Group
	create := models.MapCreateCodeGroupsRequest(request)
	fmt.Println(create, "create")
	codeGroup, err := ws.db.CodeGroupsAssociations.Create(create)
	if err != nil {
		return err
	}

	//response
	response := models.MapCodeGroupsResponse(codeGroup)
	return Ok(w, response)
}

func (ws *webserver) handleGetCodeGroupAssociations(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get url params
	vars := mux.Vars(r)
	id := strconv.Quote(vars["codeGroupsAssociationID"])
	if id == "" {
		return &invalidRequest{}
	}

	fmt.Println(id, "id")
	checkId, err := strconv.Unquote(id)

	// get group
	codeGroupDetail, err := ws.db.CodeGroupsAssociations.Get(checkId)
	if err != nil {
		return &notFound{"projectGroupID doesn't found"}
	}
	return Ok(w, codeGroupDetail)
}

// Delete record from code group association
func (ws *webserver) handleDeleteCodeGroupAssociations(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get url params
	vars := mux.Vars(r)
	id := strconv.Quote(vars["codeGroupsAssociationID"])
	if id == "nil" {
		return &invalidRequest{}
	}

	codeGroupId, err := strconv.Unquote(id)
	fmt.Println(codeGroupId, err, "error")

	if err = ws.db.CodeGroupsAssociations.Delete(codeGroupId); err != nil {
		return err
	}

	// struct{}{} is an empty object, returns "{}" to the client
	return Ok(w, struct{}{})
}

func (ws *webserver) handlePostPropertyGroupAssociations(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get url params

	// create var ready to hold decoded json from body
	var request models.PropertyGroupAssociationRequest

	// decode body
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&request); err != nil {
		return &invalidRequest{}
	}
	fmt.Println(request, "request")

	// create Project Group

	create := models.MapCreatePropertyGroupsRequest(request)
	fmt.Println(create, "create")
	propertyGroup, err := ws.db.PropertyGroupsAssociations.Create(create)
	if err != nil {
		return err
	}

	//response
	response := models.MapPropertyGroupsResponse(propertyGroup)
	return Ok(w, response)
}

func (ws *webserver) handleGetPropertyGroupAssociations(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get url params
	vars := mux.Vars(r)
	id := strconv.Quote(vars["propertyGroupsAssociationID"])
	if id == "" {
		return &invalidRequest{}
	}

	fmt.Println(id, "id")
	checkId, err := strconv.Unquote(id)

	// get group
	propertyGroupDetail, err := ws.db.PropertyGroupsAssociations.Get(checkId)
	if err != nil {
		return &notFound{"propertyGroupID doesn't found"}
	}
	return Ok(w, propertyGroupDetail)
}

// Delete record from property group association
func (ws *webserver) handleDeletePropertyGroupAssociations(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get url params
	vars := mux.Vars(r)
	id := strconv.Quote(vars["propertyGroupsAssociationID"])
	if id == "nil" {
		return &invalidRequest{}
	}

	propertyGroupId, err := strconv.Unquote(id)
	fmt.Println(propertyGroupId, err, "error")

	if err = ws.db.PropertyGroupsAssociations.Delete(propertyGroupId); err != nil {
		return err
	}

	// struct{}{} is an empty object, returns "{}" to the client
	return Ok(w, struct{}{})
}

func (ws *webserver) handlePostGroupCodeAssociations(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get url params

	// create var ready to hold decoded json from body
	var request models.GroupCodeCollectionAssociationRequest
	// decode body
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&request); err != nil {
		return &invalidRequest{}
	}
	fmt.Println(request, "request")

	// create Project Group

	create := models.MapCreateGroupCodeRequest(request)
	fmt.Println(create, "create")
	groupCollection, err := ws.db.GroupsCodeCollectionAssociations.Create(create)
	if err != nil {
		return err
	}

	//response
	response := models.MapGroupCodeResponse(groupCollection)
	return Ok(w, response)
}

func (ws *webserver) handleGetGroupCodeAssociations(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get url params
	vars := mux.Vars(r)
	id := strconv.Quote(vars["groupCodeCollectionID"])
	if id == "" {
		return &invalidRequest{}
	}

	fmt.Println(id, "id")
	checkId, err := strconv.Unquote(id)

	// get group
	codeGroupDetail, err := ws.db.GroupsCodeCollectionAssociations.Get(checkId)
	if err != nil {
		return &notFound{"code GroupID doesn't found"}
	}
	return Ok(w, codeGroupDetail)
}

func (ws *webserver) handleDeleteGroupCodeAssociations(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get url params
	vars := mux.Vars(r)
	id := strconv.Quote(vars["groupCodeCollectionID"])
	if id == "nil" {
		return &invalidRequest{}
	}

	collectionGroupId, err := strconv.Unquote(id)
	fmt.Println(collectionGroupId, err, "error")

	if err = ws.db.GroupsCodeCollectionAssociations.Delete(collectionGroupId); err != nil {
		return err
	}

	// struct{}{} is an empty object, returns "{}" to the client
	return Ok(w, struct{}{})
}

//Get Detail of Group User
func (ws *webserver) handlePostGroupUserAssociations(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get url params

	// create var ready to hold decoded json from body
	var request models.GroupUserAssociationRequest
	// decode body
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&request); err != nil {
		return &invalidRequest{}
	}

	// create Group User
	create := models.MapCreateGroupUserRequest(request)
	groupUser, err := ws.db.GroupUserAssociations.Create(create)
	if err != nil {
		return err
	}

	//response
	response := models.MapGroupUserResponse(groupUser)
	return Ok(w, response)
}

// Get the detail of group user association
func (ws *webserver) handleGetGroupUserAssociations(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get url params
	vars := mux.Vars(r)
	id := strconv.Quote(vars["groupUserID"])
	if id == "" {
		return &invalidRequest{}
	}

	fmt.Println(id, "id")
	checkId, err := strconv.Unquote(id)

	// get group
	groupUserDetail, err := ws.db.GroupUserAssociations.Get(checkId)
	if err != nil {
		return &notFound{"group UserID doesn't found"}
	}
	return Ok(w, groupUserDetail)
}

func (ws *webserver) handleDeleteGroupUserAssociations(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get url params
	vars := mux.Vars(r)
	id := strconv.Quote(vars["groupUserID"])
	if id == "nil" {
		return &invalidRequest{}
	}

	groupUserId, err := strconv.Unquote(id)
	fmt.Println(groupUserId, err, "error")

	if err = ws.db.GroupUserAssociations.Delete(groupUserId); err != nil {
		return err
	}

	// struct{}{} is an empty object, returns "{}" to the client
	return Ok(w, struct{}{})
}

// Get the detail of invited user group assoication
func (ws *webserver) handleGetInvitiedUsers(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// get url params
	vars := mux.Vars(r)

	projectID := strconv.Quote(vars["projectID"])
	if projectID == "" {
		return &invalidRequest{}
	}

	// checkId, err := strconv.Unquote(id)

	response := make([]models.InviteGroupUserResponse, 0)
	ProjectID, err := strconv.Atoi(projectID)

	invites, err := ws.db.Invites.GetCollaboratorbyProjectID(middleware.UserId, ProjectID)
	fmt.Println(invites, "invites")
	projectGroupsAssociations, err := ws.db.ProjectGroupsAssociations.GetGroupByProjectID(ProjectID)
	fmt.Println("projectGroupsAssociations", projectGroupsAssociations)
	if err != nil {
		return &invalidRequest{}
	}

	// group := make([]database.Group,0)
	var group []database.Group
	for _, invite := range invites {
		for _, projectGroupAsspociate := range projectGroupsAssociations {
			res := ws.db.Groups.GetGroupForActive(projectGroupAsspociate.GroupID)
			fmt.Println(res, "res")

			for _, activeGroup := range res {
				group = append(group, activeGroup)
			}

		}
		response = append(response, models.MapInviteGroupUserResponse(invite, group))

	}

	// db := c.db.Joins("JOIN projectGroupsAssociations ON projectGroupsAssociations.group_id = groups.id").Joins("JOIN invites ON invites.project_id = projectGroupsAssociations.project_id").Where("projectGroupsAssociations.project_id = ?", "14" , "AND groups.is_active = ?", "true").Find(&projectGroupsAssociations); db.Error != nil {
	// 	return db{}, db.Error
	// }

	return Ok(w, response)

	// // get InviteGroupUserDetail
	// //InviteGroupUserDetail, err := ws.db.Invite.GetInviteGroupAssociation(checkId)
	// if err != nil {
	// 	return &notFound{"InviteGroupUserDetail ProjectID doesn't found"}
	// }
	// 	//response
	// 	response := models.MapInviteGroupUserResponse(InviteGroupUserDetail)
	// 	return Ok(w, response)
}

// func (ws *webserver) handleDeleteGroupUserAssociations(w http.ResponseWriter, r *http.Request) error {

// 	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

// 	(w).Header().Set("Access-Control-Allow-Origin", "*")
// 	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
// 	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
// 	// get url params
// 	vars := mux.Vars(r)
// 	id := strconv.Quote(vars["groupUserID"])
// 	if id == "nil" {
// 		return &invalidRequest{}
// 	}

// 	groupUserId, err := strconv.Unquote(id)
// 	fmt.Println(groupUserId, err, "error")

// 	if err = ws.db.GroupUserAssociations.Delete(groupUserId); err != nil {
// 		return err
// 	}

// 	// struct{}{} is an empty object, returns "{}" to the client
// 	return Ok(w, struct{}{})
// }

// delete group's property
func (ws *webserver) handleDeleteGroupProperty(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get url params
	vars := mux.Vars(r)

	groupID := strconv.Quote(vars["groupID"])
	if groupID == "nil" {
		return &invalidRequest{}
	}
	groupId, err := strconv.Unquote(groupID)
	fmt.Println(groupId, err, "error")

	propertyID, err := models.MapDecodeRequest(vars["propertyID"])
	if err != nil {
		return &invalidRequest{}
	}

	results, err := ws.db.PropertyGroupsAssociations.GetPropertyByGroupID(groupId)
	if err != nil {
		return err
	}

	fmt.Println(results, "results")
	fmt.Println(propertyID, "propertyID")

	for _, groupProperty := range results {

		if groupProperty.PropertyID == propertyID && groupProperty.GroupID == groupId {
			ws.db.Properties.Delete(propertyID)

			ws.db.PropertyGroupsAssociations.Delete(groupProperty.ID)
		}

	}
	return Ok(w, "Property has been deleted")

}

// delete group's code
func (ws *webserver) handleDeleteGroupCodeCollection(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get url params
	vars := mux.Vars(r)

	groupID := strconv.Quote(vars["groupID"])
	if groupID == "nil" {
		return &invalidRequest{}
	}
	groupId, err := strconv.Unquote(groupID)
	fmt.Println(groupId, err, "error")

	codeCollectionID := strconv.Quote(vars["codeCollectionID"])
	if codeCollectionID == "nil" {
		return &invalidRequest{}
	}
	codeCollectionId, err := strconv.Unquote(codeCollectionID)
	fmt.Println(codeCollectionId, err, "error")

	results, err := ws.db.GroupsCodeCollectionAssociations.GetCodeCollectionsByGroupID(groupId)
	if err != nil {
		return err
	}

	fmt.Println(results, "results")

	for _, groupCode := range results {

		if groupCode.CodeCollectionID == codeCollectionId {

			updateCode := models.MapUpdateCodeCollectionStatus(codeCollectionId, false)
			fmt.Println(updateCode, "updateCode")

			updatedCode, err := ws.db.CodeCollections.UpdateDelete(updateCode)
			if err != nil {
				return err
			}
			fmt.Println(updatedCode, "updatedCode")
		}

	}
	return Ok(w, "Code has been deleted")

}

// Add property to project
// func (ws *webserver) handlePostProperty(w http.ResponseWriter, r *http.Request) error {

// 	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

// 	(w).Header().Set("Access-Control-Allow-Origin", "*")
// 	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
// 	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
// 	// get url params
// 	vars := mux.Vars(r)
// 	ProjectID, err := strconv.Atoi(vars["projectID"])
// 	if err != nil {
// 		return &invalidRequest{}
// 	}

// 	var request models.AddPropertyRequest
// 	dec := json.NewDecoder(r.Body)
// 	if err := dec.Decode(&request); err != nil {
// 		fmt.Println(err, "decodeErrr")
// 		return &invalidRequest{}
// 	}

// 	create := models.MapAddPropertyRequest(ProjectID, middleware.UserId, request)
// 	response := make([]models.PropertyResponse, 0)

// 	project, err := ws.db.Projects.Get(ProjectID)
// 	fmt.Println(project, "project")
// 	if err != nil {
// 		// return &notFound{"projectID doesn't match"}
// 		return err
// 	} else {
// 		property, err := ws.db.Properties.Create(create)
// 		if err != nil {

// 			return &notFound{"method doesn't match"}
// 		} else {

// 			projectProperty := models.MapPropertyProjectRequest(ProjectID, property.ID)

// 			projectPropertyAssociate, err := ws.db.ProjectPropertyAssociations.Create(projectProperty)
// 			if err != nil {
// 				return err
// 			}
// 			fmt.Println(projectPropertyAssociate, "1721")

// 			response = append(response, models.MapPropertyResponse(property, nil, project.ID, project.Name))

// 		}

// 		return Ok(w, response)
// 	}
// }

// func (ws *webserver) handlePostCodeCollectionAssociation(w http.ResponseWriter, r *http.Request) error {
// 	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

// 	(w).Header().Set("Access-Control-Allow-Origin", "*")
// 	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
// 	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
// 	// get url params

// 	// create var ready to hold decoded json from body
// 	var request models.CodeCollectionAssociationRequest

// 	// decode body
// 	dec := json.NewDecoder(r.Body)
// 	if err := dec.Decode(&request); err != nil {
// 		return &invalidRequest{}
// 	}
// 	fmt.Println(request, "request")

// 	// create codecollection

// 	create := models.MapCreateCodeCollectionAssociationRequest(request)
// 	fmt.Println(create, "create")
// 	codeCollection, err := ws.db.CodeCollectionAssociations.Create(create)
// 	fmt.Println(codeCollection, "codeCollection")
// 	if err != nil {
// 		return err
// 	}

// 	//response
// 	response := models.MapCodeCollectionAssociationResponse(codeCollection)
// 	fmt.Println(response, "response")
// 	return Ok(w, response)
// }

func (ws *webserver) handleGetUserDetailsByEmail(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get url params
	vars := mux.Vars(r)
	Email := strconv.Quote(vars["email"])
	if Email == "" {
		return &invalidRequest{}
	}

	email, err := strconv.Unquote(Email)
	if err != nil {
		return &invalidRequest{}
	}

	fmt.Println(email, "Email")

	user, err := ws.db.Users.GetUserByEmail(email)
	if err != nil {
		return &invalidRequest{}
	}
	fmt.Println(user, "user")

	return Ok(w, user)
}

//Save Company Details
func (ws *webserver) handlePostCompanyDetails(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	var request models.CompanyDetailsRequest

	// decode body
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&request); err != nil {
		return &invalidRequest{}
	}
	fmt.Println(request, "request")
	fmt.Println(time.Now(), " new timezone")
	// create code collection
	create := models.MapCreateCompanyDetailsRequest(middleware.UserId, request, true)
	companyDetails, err := ws.db.Accounts.Create(create)
	if err != nil {
		return err
	}

	return Ok(w, companyDetails)
}

//Get Company Details
func (ws *webserver) handleGetCompanyDetails(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get the company details
	fmt.Println("=== handleGetCompanyDetails ===")
	companyDetails, err := ws.db.Accounts.GetAll(middleware.UserId)

	if err != nil {
		fmt.Println(err, "=== handleGetCompanyDetails ===")
		return err
	}

	response := models.MapCompanyDetailsResponse(companyDetails)

	// for _, company := range companyDetails {

	// 	response = append(response, )

	// }
	return Ok(w, response)

}

// Save updated company details
func (ws *webserver) handlePutCompanyDetails(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	var request models.CompanyDetailsRequest
	fmt.Println("=== handlePUTCompanyDetails ===")

	// decode body
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&request); err != nil {
		return &invalidRequest{}
	}
	fmt.Println(request, "request")
	fmt.Println(time.Now(), " new timezone")
	response := make([]models.CompanyDetailsResponse, 0)
	company, err := ws.db.Accounts.GetAll(middleware.UserId)
	if err != nil {

		// create Accounts Details
		create := models.MapCreateCompanyDetailsRequest(middleware.UserId, request, true)
		companyDetails, err := ws.db.Accounts.Create(create)
		response = append(response, models.MapUpdateCompanyDetailsResponse(companyDetails))
		if err != nil {
			fmt.Println(err, "=== handlePUTCompanyDetails ===")
			return err
		}
	} else {
		update := models.MapUpdateCompanyDetailsRequest(request, company.ID, middleware.UserId)
		fmt.Println(company, "company")
		putCompanyDetails, err := ws.db.Accounts.Update(update, middleware.UserId)
		response = append(response, models.MapUpdateCompanyDetailsResponse(putCompanyDetails))
		if err != nil {
			return err
		}
	}

	return Ok(w, response)
}

func (ws *webserver) handleGetDemoProjects(w http.ResponseWriter, r *http.Request) error {

	// var response models.ProjectResponse

	response := make([]models.ProjectResponse, 0)

	projects, err := ws.db.DemoProjects.GetAll("permitDocs")

	// properties, err := ws.db.ProjectPropertyAssociations.GetAllByProjectID(result.ProjectID)
	if err != nil {
		return err
	}

	// userDetail, err := ws.db.Users.Get(projects.CreatedBy)

	for _, project := range projects {
		response = append(response, models.MapDemoProjectResponse(project))
	}

	fmt.Println(projects, "projects")

	return Ok(w, response)

}

// Update detail of User
func (ws *webserver) handleUpdateIsIntro(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("handleUpdateIsIntro-1 ")
	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// create var ready to hold decoded json from body
	var request models.DoIntroRequest
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&request); err != nil {
		return &invalidRequest{}
	}
	// update contact
	update := models.MapIsIntroRequest(middleware.UserId, request)
	fmt.Println("handleUpdateIsIntro-2 ")
	user, err := ws.db.Users.Update(update)
	if err != nil {
		return err
	}
	// create response
	response := models.MapUserResponse(user)

	return Ok(w, response)
}

//Reset User Roles
func (ws *webserver) handleResetUserRoles(w http.ResponseWriter, r *http.Request) error {
 	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// create var ready to hold decoded json from body
	var request models.UserRolesRequest
	dec := json.NewDecoder(r.Body)
 	if err := dec.Decode(&request); err != nil {
		return &invalidRequest{}
	}
	// update contact
	update := models.MapUserRolesRequest(request)
	 
	fmt.Println(update,"updateupdate")

	inviteDetail, err := ws.db.Invites.Get(update.InviteID)
	fmt.Println(inviteDetail,"inviteDetail")
	if err != nil {
		return err
	} 
 	
 	response := make([]models.InviteUserResponse, 0)

		  
	roleID, err := strconv.Atoi(update.ID)

	updateRole := models.MapUpdateProjectRoleRequest(update.InviteID,roleID)
	fmt.Println(updateRole,"ashish-update")
	ws.db.Invites.Update(updateRole)
	if (inviteDetail.IsInviteAccepted == 1) {
		fmt.Println("Project update")
		user, err := ws.db.Users.GetUserByEmail(inviteDetail.Email)
		fmt.Println(user,"users")
		projectUserRoleAssociate,err := ws.db.ProjectUserRoleAssociations.GetUserRoleByProjectID(inviteDetail.ProjectID,user.ID)
		fmt.Println(projectUserRoleAssociate,"projectUserRoleAssociate")
		if err != nil {
			return err
		}
		userRoleAssociate := models.MapUpdateRole(projectUserRoleAssociate.ID, updateRole.Role)
		fmt.Println(userRoleAssociate, "userRoleAssociatehas")
		ws.db.ProjectUserRoleAssociations.Update(userRoleAssociate)
	}	 
 	
	return Ok(w, response)
}


func (ws *webserver) handlDeleteInviteGroupAssociations(w http.ResponseWriter, r *http.Request) error {

	(w).Header().Set("Access-Control-Allow-Origin", CORS_URL)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// get url params
	vars := mux.Vars(r)
 	 
	InviteGroupID := strconv.Quote(vars["inviteGroupID"])

	inviteGroupID, err := strconv.Atoi(InviteGroupID)
	if err == nil {
		fmt.Println(InviteGroupID)
	}

 	// if InviteGroupID == "" {
	// 	return &invalidRequest{}
	// }
	// inviteGroupID, err := strconv.Unquote(InviteGroupID)

	inviteGroupAssociations, err := ws.db.InviteGroupAssociations.GetByInviteID(inviteGroupID)
	fmt.Println(inviteGroupAssociations,"inviteGroupAssociations--1")
	if err != nil {
		return err
	}

	inviteDetail,err := ws.db.Invites.Get(inviteGroupAssociations.InviteID)
   		if(inviteDetail.IsInviteAccepted == 1) {
			fmt.Println("step1")
			user,err := ws.db.Users.GetUserByEmail(inviteDetail.Email);
			projectUserGroupDetails,err := ws.db.ProjectUserGroupAssociations.GetAllByUserProjectID(user.ID,inviteDetail.ProjectID)
			projectUserGroupAssociations,err := ws.db.ProjectUserGroupAssociations.GetByUserIDProjectIDGroupID(user.ID,inviteGroupAssociations.GroupID,inviteDetail.ProjectID)
			fmt.Println(projectUserGroupAssociations,"projectUserGroupAssociations")
			if err != nil {
				return err
			}
			ws.db.ProjectUserGroupAssociations.Delete(projectUserGroupAssociations.ID)
			if(len(projectUserGroupDetails) == 1) {
				userInfo, err := ws.db.ProjectUserRoleAssociations.DeleteCollaboratorForProject(user.ID,inviteDetail.ProjectID)
				fmt.Println(userInfo, "userInfo")
				if err != nil {
					return err
					}	
			}
			
		 }
		
		inviteGroupAssociationsDetails, err := ws.db.InviteGroupAssociations.GetAll(inviteGroupAssociations.InviteID)
		fmt.Println(inviteGroupAssociationsDetails,"inviteGroupAssociationsDetails--2")
		if(len(inviteGroupAssociationsDetails) == 1){
		 ws.db.Invites.Delete(inviteGroupAssociationsDetails[0].InviteID)
		}
		ws.db.InviteGroupAssociations.Delete(inviteGroupAssociations.ID);

 	
	return Ok(w, "Successfully Deleted Collaborator")
}