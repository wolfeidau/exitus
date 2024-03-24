package metrics

import (
	"database/sql"
	"time"

	"github.com/rcrowley/go-metrics"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/exitus/pkg/conf"
)

// LogWriter write metrics to the logger.
type LogWriter struct{}

// New create a new log writer for metrics.
func New(cfg *conf.Config) *LogWriter {
	return &LogWriter{}
}

// Start start the writer.
func (mr *LogWriter) Start() {
	metrics.WriteJSON(metrics.DefaultRegistry, 30*time.Second, mr)
}

func (mr *LogWriter) Write(data []byte) (int, error) {
	log.Info().RawJSON("r", data).Msg("metrics")
	return 0, nil
}

// DBMonitor monitor a db pool.
type DBMonitor struct {
	db *sql.DB
}

// NewDBMonitor create a new db pool monitor.
func NewDBMonitor(db *sql.DB) *DBMonitor {
	return &DBMonitor{db: db}
}

// Start start the collector.
func (dbm *DBMonitor) Start() {
	s := dbm.db.Stats()
	ig := metrics.GetOrRegisterGauge("db.conn.idle", nil)
	ig.Update(int64(s.Idle))
	iu := metrics.GetOrRegisterGauge("db.conn.inuse", nil)
	iu.Update(int64(s.InUse))
	oc := metrics.GetOrRegisterGauge("db.conn.open", nil)
	oc.Update(int64(s.OpenConnections))
}
