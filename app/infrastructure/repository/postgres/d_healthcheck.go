package postgres

// HealthCheck returns database health check.
func (d DB) HealthCheck() error {
	err := d.db.Exec(`SELECT 1`).Error
	if err != nil {
		return err
	}

	return nil
}
