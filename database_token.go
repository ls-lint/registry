package main

func (database *database) tokens(userId uint) ([]*Token, error) {
	var tokens []*Token

	return tokens, database.connection.
		Where("tokens.user_id = ?", userId).
		Find(&tokens).
		Error
}

func (database *database) createToken(token *Token) (*Token, error) {
	return token, database.connection.
		FirstOrCreate(&token, &Token{
			Token:  token.getToken(),
			UserId: token.getUserId(),
		}).
		Error
}

func (database *database) deleteToken(token *Token) error {
	return database.connection.Delete(&token).Error
}
