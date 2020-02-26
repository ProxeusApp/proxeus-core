package handlers

import (
	"strings"

	"github.com/ProxeusApp/proxeus-core/main/handlers/api"
	"github.com/ProxeusApp/proxeus-core/main/handlers/formbuilder"
	"github.com/ProxeusApp/proxeus-core/main/handlers/i18n"
	"github.com/ProxeusApp/proxeus-core/main/handlers/payment"
	"github.com/ProxeusApp/proxeus-core/main/handlers/template_ide"
	"github.com/ProxeusApp/proxeus-core/main/handlers/workflow"
	"github.com/ProxeusApp/proxeus-core/main/www"
	"github.com/ProxeusApp/proxeus-core/sys/model"

	"github.com/labstack/echo"
)

func MainHostedAPI(e *echo.Echo, s *www.Security, version string) {
	const (
		GET    = echo.GET
		POST   = echo.POST
		PUT    = echo.PUT
		DELETE = echo.DELETE
	)

	const (
		PUBLIC     = model.PUBLIC
		GUEST      = model.GUEST
		USER       = model.USER
		CREATOR    = model.CREATOR
		ADMIN      = model.ADMIN
		SUPERADMIN = model.SUPERADMIN
		ROOT       = model.ROOT
	)

	type r struct {
		m string     // http method
		a model.Role // access level
		p string     // path
		h func(echo.Context) error
	}

	routes := []r{
		{GET, PUBLIC, "/", api.PublicIndexHTMLHandler},

		//public access for shared by link
		{GET, PUBLIC, "/p/:type/:ID", api.SharedByLinkHTMLHandler},

		{GET, PUBLIC, "/validation", api.PublicIndexHTMLHandler},
		{GET, PUBLIC, "/login", api.PublicIndexHTMLHandler},
		{GET, PUBLIC, "/reset/password", api.PublicIndexHTMLHandler},
		{GET, PUBLIC, "/guest", api.PublicIndexHTMLHandler},

		{GET, CREATOR, "/admin", api.AdminIndexHandler},
		{GET, CREATOR, "/admin/*", api.AdminIndexHandler},
		{GET, USER, "/user", api.UserBackendHTMLHandler},
		{GET, USER, "/user/*", api.UserBackendHTMLHandler},

		{GET, PUBLIC, "/document", api.UserBackendHTMLHandler},
		{GET, PUBLIC, "/document/:ID/payment", api.UserBackendHTMLHandler},
		{GET, PUBLIC, "/document/:ID", api.UserBackendHTMLHandler},
	}

	routesNoCache := []r{
		// config
		{GET, PUBLIC, "/api/config", api.ConfigHandler(version)},

		// authentication
		{GET, PUBLIC, "/api/session/token", api.GetSessionTokenHandler},
		{DELETE, USER, "/api/session/token", api.DeleteSessionTokenHandler},
		{GET, PUBLIC, "/api/challenge", api.ChallengeHandler},       // Need session
		{POST, PUBLIC, "/api/change/bcaddress", api.UpdateAddress},  // Need session
		{POST, PUBLIC, "/api/change/email", api.ChangeEmailRequest}, // Need session
		{POST, PUBLIC, "/api/change/email/:token", api.ChangeEmail},
		{POST, PUBLIC, "/api/reset/password", api.ResetPasswordRequest},
		{POST, PUBLIC, "/api/reset/password/:token", api.ResetPassword},
		{POST, PUBLIC, "/api/register", api.RegisterRequest},
		{POST, PUBLIC, "/api/register/:token", api.Register},
		{POST, PUBLIC, "/api/login", api.LoginHandler},
		{POST, PUBLIC, "/api/logout", api.LogoutHandler},
		{GET, USER, "/api/user/create/api/key/:ID", api.CreateApiKeyHandler},
		{DELETE, USER, "/api/user/create/api/key/:ID", api.DeleteApiKeyHandler},

		// user
		{GET, PUBLIC, "/api/me", api.MeHandler}, // Need session
		{POST, USER, "/api/me", api.MeUpdateHandler},
		{POST, USER, "/api/my/profile/photo", api.PutProfilePhotoHandler},
		{GET, USER, "/api/my/profile/photo", api.GetProfilePhotoHandler},
		{GET, USER, "/api/profile/photo", api.GetProfilePhotoHandler},


		// document
		{GET, PUBLIC, "/api/document/:ID", api.DocumentHandler}, // Need session
		{GET, PUBLIC, "/api/document/list", workflow.ListPublishedHandler},
		{GET, PUBLIC, "/api/document/:ID/allAtOnce/schema", api.WorkflowSchema},
		{POST, PUBLIC, "/api/document/:ID/allAtOnce", api.WorkflowExecuteAtOnce},
		{POST, GUEST, "/api/document/:ID/edit/name", api.DocumentEditHandler},
		{POST, GUEST, "/api/document/:ID/name", api.DocumentEditHandler},
		{POST, GUEST, "/api/document/:ID/data", api.DocumentDataHandler},
		{POST, GUEST, "/api/document/:ID/next", api.DocumentNextHandler},
		{GET, PUBLIC, "/api/document/:ID/prev", api.DocumentPrevHandler}, // Why PUBLIC access for prev when next is GUEST
		{GET, GUEST, "/api/document/:ID/file/:inputName", api.DocumentFileGetHandler},
		{POST, PUBLIC, "/api/document/:ID/file/:inputName", api.DocumentFilePostHandler},                 // Should be GUEST
		{GET, PUBLIC, "/api/document/:ID/preview/:templateID/:lang/:format", api.DocumentPreviewHandler}, // Should be GUEST
		{GET, GUEST, "/api/document/:ID/delete", api.DocumentDeleteHandler},

		// user document
		{GET, GUEST, "/api/user/document", api.UserDocumentListHandler},
		{GET, USER, "/api/user/document/:ID", api.UserDocumentGetHandler},
		{GET, GUEST, "/api/user/document/file/:ID/:dataPath", api.UserDocumentFileHandler},
		{GET, USER, "/api/user/document/signingRequests/:ID/:docID", api.UserDocumentSignatureRequestGetByDocumentIDHandler},
		{POST, USER, "/api/user/document/signingRequests/:ID/:docID/add", api.UserDocumentSignatureRequestAddHandler},
		{POST, USER, "/api/user/document/signingRequests/:ID/:docID/revoke", api.UserDocumentSignatureRequestRevokeHandler},
		{POST, USER, "/api/user/document/signingRequests/:ID/:docID/reject", api.UserDocumentSignatureRequestRejectHandler},
		{GET, USER, "/api/user/document/signingRequests", api.UserDocumentSignatureRequestGetCurrentUserHandler},
		{POST, USER, "/api/user/delete", api.UserDeleteHandler},

		// import export
		{GET, USER, "/api/user/export", api.ExportUser},
		{GET, ROOT, "/api/settings/export", api.ExportSettings},
		{GET, USER, "/api/userdata/export", api.ExportUserData},
		{GET, USER, "/api/export/results", api.GetExportResults},
		{GET, USER, "/api/import/results", api.GetImportResults},
		{GET, USER, "/api/export", api.GetExport},
		{POST, USER, "/api/export", api.GetExport},
		{POST, USER, "/api/import", api.PostImport},

		// admin
		{GET, ROOT, "/api/init", api.GetInit},
		{POST, ROOT, "/api/init", api.PostInit},
		{GET, ROOT, "/api/switch/user/:address", api.SwitchUserHandler},
		{POST, ADMIN, "/api/admin/invite", api.InviteRequest},
		{GET, SUPERADMIN, " /api/admin/user/:ID", api.AdminUserGetHandler},
		{GET, USER, " /api/admin/user/list", api.AdminUserListHandler},
		{POST, USER, " /api/admin/user/list", api.AdminUserListHandler},
		{POST, SUPERADMIN, "/api/admin/user/update", api.AdminUserUpdateHandler},

		// i18n
		{GET, PUBLIC, "/api/admin/i18n/", i18n.IndexHandler},
		{GET, PUBLIC, "/api/admin/i18n/meta", i18n.MetaHandler},
		{GET, PUBLIC, "/api/admin/i18n/all", i18n.AllHandler},
		{GET, SUPERADMIN, "/api/i18n/export", i18n.ExportI18n},
		{GET, PUBLIC, "/api/i18n/meta", i18n.MetaHandler},
		{GET, PUBLIC, "/api/i18n/all", i18n.AllHandler},
		{GET, SUPERADMIN, "/api/admin/i18n/find", i18n.FindHandler},
		{GET, SUPERADMIN, "/api/admin/i18n/search", i18n.FormBuilderI18nSearchHandler},
		{GET, PUBLIC, "/api/i18n/search", i18n.FormBuilderI18nSearchHandler},
		{POST, SUPERADMIN, "/api/admin/i18n/update", i18n.UpdateHandler},
		{POST, SUPERADMIN, "/api/admin/i18n/fallback", i18n.SetFallbackHandler},
		{POST, SUPERADMIN, "/api/admin/i18n/lang", i18n.LangSwitchHandler},
		{POST, PUBLIC, "/api/admin/i18n/translate", i18n.TranslateHandler},

		// workflow
		{GET, CREATOR, "/api/admin/workflow/:ID/delete", workflow.DeleteHandler},
		{GET, USER, "/api/workflow/export", workflow.ExportWorkflow},
		{GET, USER, "/api/user/workflow/list", workflow.ListPublishedHandler},
		{GET, PUBLIC, "/api/admin/workflow/list", workflow.ListHandler},      // Need session
		{GET, PUBLIC, "/api/admin/workflow/:ID", workflow.GetHandler},        // Need session
		{POST, PUBLIC, "/api/admin/workflow/update", workflow.UpdateHandler}, // Need session

		// payment
		{GET, USER, "/api/admin/payments/check", api.CheckForWorkflowPayment},
		{POST, USER, "/api/admin/payments", payment.CreateWorkflowPayment},
		{GET, USER, "/api/admin/payments/:paymentId", payment.GetWorkflowPaymentById},
		{GET, USER, "/api/admin/payments", payment.GetWorkflowPayment},
		{PUT, USER, "/api/admin/payments/:paymentId", payment.UpdateWorkflowPaymentPending},
		{POST, USER, "/api/admin/payments/:paymentId/cancel", payment.CancelWorkflowPayment},
		{GET, SUPERADMIN, "/api/admin/payments/list", payment.ListPayments},
		{DELETE, SUPERADMIN, "/api/admin/payments/:paymentId", payment.DeleteWorkflowPayment},

		// form
		{GET, PUBLIC, "/api/form/component", formbuilder.GetComponentsHandler}, // `Need session`
		{GET, USER, "/api/form/export", formbuilder.ExportForms},
		{GET, CREATOR, "/api/admin/form/:ID/delete", formbuilder.DeleteHandler},
		{GET, PUBLIC, "/api/admin/form/list", formbuilder.ListHandler}, // Need session
		{GET, USER, "/api/admin/:type/list", workflow.ListCustomNodeHandler},
		{GET, PUBLIC, "/api/admin/form/:formID", formbuilder.GetOneFormHandler}, // Need session
		{POST, PUBLIC, "/api/admin/form/update", formbuilder.UpdateFormHandler}, // Need session

		{GET, PUBLIC, "/api/admin/form/component", formbuilder.GetComponentsHandler}, // Need session
		{POST, SUPERADMIN, "/api/admin/form/component", formbuilder.SetComponentHandler},
		{DELETE, SUPERADMIN, "/api/admin/form/component/:id", formbuilder.DeleteComponentHandler},
		{GET, PUBLIC, "/api/admin/form/vars", formbuilder.VarsHandler}, // Need session

		{POST, PUBLIC, "/api/admin/form/test/setFormSrc/:id", formbuilder.SetFormSrcHandler},      // Need session
		{GET, PUBLIC, "/api/admin/form/test/data/:id", formbuilder.GetDataId},                     // Need session
		{POST, PUBLIC, "/api/admin/form/test/data/:id", formbuilder.TestFormDataHandler},          // Need session
		{GET, PUBLIC, "/api/admin/form/test/file/:id/:fieldname", formbuilder.GetFileIdFieldName}, // Need session
		{GET, PUBLIC, "/api/admin/form/file/types", formbuilder.GetFileTypes},
		{POST, PUBLIC, "/api/admin/form/test/file/:id/:fieldname", formbuilder.PostFileIdFieldName}, // Need session

		// template
		{GET, CREATOR, "/api/admin/template/:ID/delete", template_ide.DeleteHandler},
		{GET, USER, "/api/template/export", template_ide.ExportTemplate},
		{GET, PUBLIC, "/api/admin/template/vars", template_ide.VarsTemplateHandler},                   // Need session
		{GET, PUBLIC, "/api/admin/template/list", template_ide.ListHandler},                           // Need session
		{POST, PUBLIC, "/api/admin/template/update", template_ide.UpdateHandler},                      // Need session
		{GET, PUBLIC, "/api/admin/template/:id", template_ide.OneTmplHandler},                         // Need session
		{GET, PUBLIC, "/api/admin/template/download/:id/:lang", template_ide.DownloadTemplateHandler}, // Need session
		{POST, PUBLIC, "/api/admin/template/upload/:id/:lang", template_ide.UploadTemplateHandler},    // Need session
		{GET, PUBLIC, "/api/admin/template/delete/:id/:lang", template_ide.DeleteTemplateHandler},     // Need session

		{GET, PUBLIC, "/api/admin/template/ide/active/:id/:lang", template_ide.IdeSetActiveHandler},   // Need session
		{POST, PUBLIC, "/api/admin/template/ide/upload/:id/:lang", template_ide.IdePostUploadHandler}, // Need session
		{GET, PUBLIC, "/api/admin/template/ide/delete/:id/:lang", template_ide.IdeGetDeleteHandler},   // Need session
		{GET, PUBLIC, "/api/admin/template/ide/download/:id", template_ide.IdeGetDownloadHandler},     // Need session
		{GET, CREATOR, "/api/admin/template/ide/tmplAssistanceDownload", template_ide.IdeGetTmpAssDownload},
		{GET, PUBLIC, "/api/admin/template/ide/form", template_ide.IdeFormHandler}, // Need session

		// external node
		{POST, PUBLIC, "/api/admin/external/register", api.ExternalRegister},
		{POST, PUBLIC, "/api/admin/external/config/:id", api.ExternalConfigStore},
		{GET, PUBLIC, "/api/admin/external/config/:id", api.ExternalConfigRetrieve},
		{GET, CREATOR, "/api/admin/external/:name/:id", api.ExternalConfigurationPage}, // Need session

		// test
		{PUT, PUBLIC, "/api/test/payments", payment.PutTestPayment},
		{PUT, PUBLIC, "/api/test/signatures", api.PutTestSignature},
	}

	addEndpoint := func(r r, ms ...echo.MiddlewareFunc) {
		if r.a > model.PUBLIC {
			ms = append(ms, s.With(r.a))
		}
		e.Add(r.m, strings.TrimSpace(r.p), r.h, ms...)
	}

	for _, r := range routes {
		addEndpoint(r)
	}
	for _, r := range routesNoCache {
		addEndpoint(r, noCache)
	}
}

func noCache(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Response().Header()
		header.Set("Cache-Control", "no-store")
		return next(c)
	}
}
