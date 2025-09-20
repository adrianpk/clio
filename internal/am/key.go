package am

type Keys struct {
	AppEnv string

	ServerWebHost      string
	ServerWebPort      string
	ServerWebEnabled   string
	ServerAPIHost      string
	ServerAPIPort      string
	ServerAPIEnabled   string
	ServerResPath      string
	ServerIndexEnabled string

	DBSQLiteDSN string

	SecCSRFKey       string
	SecCSRFRedirect  string
	SecEncryptionKey string
	SecHashKey       string
	SecBlockKey      string
	SecBypassAuth    string

	ButtonStyleGray   string
	ButtonStyleBlue   string
	ButtonStyleRed    string
	ButtonStyleGreen  string
	ButtonStyleYellow string

	NotificationSuccessStyle string
	NotificationInfoStyle    string
	NotificationWarnStyle    string
	NotificationErrorStyle   string
	NotificationDebugStyle   string

	RenderWebErrors string
	RenderAPIErrors string

	SSGWorkspacePath string
	SSGDocsPath      string
	SSGMarkdownPath  string
	SSGHTMLPath      string
	SSGLayoutPath    string
	SSGAssetsPath    string
	SSGImagesPath    string
}

var Key = Keys{
	AppEnv: "app.env",

	ServerWebHost:      "server.web.host",
	ServerWebPort:      "server.web.port",
	ServerWebEnabled:   "server.web.enabled",
	ServerAPIHost:      "server.api.host",
	ServerAPIPort:      "server.api.port",
	ServerAPIEnabled:   "server.api.enabled",
	ServerResPath:      "server.res.path",
	ServerIndexEnabled: "server.index.enabled",

	DBSQLiteDSN: "db.sqlite.dsn",

	SecCSRFKey:       "sec.csrf.key",
	SecCSRFRedirect:  "sec.csrf.redirect",
	SecEncryptionKey: "sec.encryption.key",
	SecHashKey:       "sec.hash.key",
	SecBlockKey:      "sec.block.key",
	SecBypassAuth:    "sec.bypass.auth",

	ButtonStyleGray:   "button.style.gray",
	ButtonStyleBlue:   "button.style.blue",
	ButtonStyleRed:    "button.style.red",
	ButtonStyleGreen:  "button.style.green",
	ButtonStyleYellow: "button.style.yellow",

	NotificationSuccessStyle: "notification.success.style",
	NotificationInfoStyle:    "notification.info.style",
	NotificationWarnStyle:    "notification.warn.style",
	NotificationErrorStyle:   "notification.error.style",
	NotificationDebugStyle:   "notification.debug.style",

	RenderWebErrors: "render.web.errors",
	RenderAPIErrors: "render.api.errors",

	SSGWorkspacePath: "ssg.workspace.path",
	SSGDocsPath:      "ssg.docs.path",
	SSGMarkdownPath:  "ssg.markdown.path",
	SSGHTMLPath:      "ssg.html.path",
	SSGLayoutPath:    "ssg.layout.path",
	SSGAssetsPath:    "ssg.assets.path",
	SSGImagesPath:    "ssg.images.path",
}
