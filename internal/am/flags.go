package am

var Flags = map[string]interface{}{
	Key.ServerWebHost:    "localhost",
	Key.ServerWebPort:    "8080",
	Key.ServerWebEnabled: true,
	Key.ServerAPIHost:    "localhost",
	Key.ServerAPIPort:    "8081",
	Key.ServerAPIEnabled: true,
	Key.ServerPreviewHost:    "localhost",
	Key.ServerPreviewPort:    "8082",
	Key.ServerPreviewEnabled: true,
	Key.ServerResPath:    "/res",

	// SSG
	Key.SSGHeaderStyle: "stacked",

	// NOTE: These values should be overridden in production with secure values.
	Key.SecHashKey:  "0123456789abcdef0123456789abcdef",
	Key.SecBlockKey: "0123456789abcdef0123456789abcdef",
}
