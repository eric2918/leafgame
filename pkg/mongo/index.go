package mongo

func InitIndex() {
	// login db
	UniqueIndex(LOGIN_DB, ACCOUNT_COLLECTION, []string{"username"})

	// game db
	UniqueIndex(GAME_DB, PLAYER_COLLECTION, []string{"player_id"})
	UniqueIndex(GAME_DB, PLAYER_COLLECTION, []string{"account_id"})

	UniqueIndex(GAME_DB, ROLES_COLLECTION, []string{"role_id"})

	UniqueIndex(GAME_DB, SKILLS_COLLECTION, []string{"skill_id"})
}
