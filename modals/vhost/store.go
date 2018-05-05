package vhost

import (
	"context"
	"database/sql"

	"golang.ysitd.cloud/db"
	"golang.ysitd.cloud/http/timing"
)

type Store struct {
	DB    *db.GeneralOpener `inject:""`
	Cache *Cache            `inject:""`
}

func (s *Store) GetVHost(ctx context.Context, hostname string) (host *VirtualHost, err error) {
	collector := ctx.Value("timing").(*timing.Collector)
	timer := collector.New("fetch_vhost", "Fetch Virtual Host")
	timer.Start()
	defer timer.Stop()

	if host := s.Cache.Get(hostname); host != nil {
		return host, nil
	}

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

	s.Cache.Set(hostname, &instance, 5)

	return &instance, nil
}
