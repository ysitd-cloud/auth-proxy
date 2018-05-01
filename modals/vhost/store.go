package vhost

import (
	"context"
	"database/sql"

	"golang.ysitd.cloud/db"
)

type Store struct {
	DB *db.GeneralOpener `inject:""`
}

func (s *Store) GetVHost(ctx context.Context, hostname string) (host *VirtualHost, err error) {
	conn, err := s.DB.Open()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	query := "SELECT hostname, oauth_id, oauth_secret, backend_path FROM host WHERE hostname = $1"
	row := conn.QueryRowContext(ctx, query, hostname)

	var instance VirtualHost
	if err := row.Scan(
		&instance.Hostname,
		&instance.OauthID,
		&instance.OAuthSecret,
		&instance.BackendPath,
	); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &instance, nil
}
