package permissions

var (
	AllUsersPermission    = Permission("dashboard.users.*")
	ViewUsersPermission   = Permission("dashboard.users.view")
	CreateUsersPermission = Permission("dashboard.users.create")
	UpdateUsersPermission = Permission("dashboard.users.update")
	DeleteUsersPermission = Permission("dashboard.users.delete")

	TerminateUserSessionsPermission = Permission("dashboard.users.sessions.terminate")
	ViewUserSessionsPermission      = Permission("dashboard.users.sessions.view")
)
