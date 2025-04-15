CREATE TABLE IF NOT EXISTS refresh_tokens (
                                              id SERIAL PRIMARY KEY,
                                              user_id UUID NOT NULL,
                                              access_jti VARCHAR(128) NOT NULL, -- идентификатор access-токена
    refresh_token_hash VARCHAR(255) NOT NULL,
    client_ip VARCHAR(45) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
                             used BOOLEAN DEFAULT FALSE
                             );

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_access_jti ON refresh_tokens(access_jti);
