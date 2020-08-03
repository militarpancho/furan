package datalayer

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	pgtypeuuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/dollarshaveclub/furan/pkg/models"
)

// DataLayer describes an object that interacts with a data store
type DataLayer interface {
	CreateBuild(context.Context, models.Build) (uuid.UUID, error)
	GetBuildByID(context.Context, uuid.UUID) (models.Build, error)
	SetBuildCompletedTimestamp(context.Context, uuid.UUID, time.Time) error
	SetBuildStatus(context.Context, uuid.UUID, models.BuildStatus) error
	DeleteBuild(context.Context, uuid.UUID) error
	CancelBuild(context.Context, uuid.UUID) error
	ListenForCancellation(context.Context, uuid.UUID) error
	ListenForBuildEvents(ctx context.Context, id uuid.UUID, c chan<- string) error
	AddEvent(ctx context.Context, id uuid.UUID, event string) error
	SetBuildAsRunning(ctx context.Context, id uuid.UUID) error
	ListenForBuildRunning(ctx context.Context, id uuid.UUID) error
	SetBuildAsCompleted(ctx context.Context, id uuid.UUID, status models.BuildStatus) error
	ListenForBuildCompleted(ctx context.Context, id uuid.UUID) (models.BuildStatus, error)
}

// PostgresDBLayer is a DataLayer instance that utilizes a PostgreSQL database
type PostgresDBLayer struct {
	p *pgxpool.Pool
}

var _ DataLayer = &PostgresDBLayer{}

// NewPostgresDBLayer returns a data layer object backed by PostgreSQL
func NewPostgresDBLayer(pguri string) (*PostgresDBLayer, error) {
	dbcfg, err := pgxpool.ParseConfig(pguri)
	if err != nil {
		return nil, fmt.Errorf("error parsing pg db uri: %w", err)
	}
	dbcfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		conn.ConnInfo().RegisterDataType(pgtype.DataType{
			Value: &pgtypeuuid.UUID{},
			Name:  "uuid",
			OID:   pgtype.UUIDOID,
		})
		return nil
	}
	pool, err := pgxpool.ConnectConfig(context.Background(), dbcfg)
	if err != nil {
		return nil, fmt.Errorf("error creating pg connection pool: %w", err)
	}
	return &PostgresDBLayer{p: pool}, nil
}

// Close closes all database connections in the connection pool
func (dl *PostgresDBLayer) Close() {
	dl.p.Close()
}

// CreateBuild inserts a new build into the DB returning the ID
func (dl *PostgresDBLayer) CreateBuild(ctx context.Context, b models.Build) (uuid.UUID, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return id, fmt.Errorf("error generating build id: %w", err)
	}
	// Clear out GH credential if present
	b.Request.GetBuild().GithubCredential = ""
	_, err = dl.p.Exec(ctx,
		`INSERT INTO builds (id, github_repo, github_ref, encrypted_github_credential, image_repos, tags, commit_sha_tag, request, status) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9);`,
		id, b.GitHubRepo, b.GitHubRef, b.EncryptedGitHubCredential, b.ImageRepos, b.Tags, b.CommitSHATag, b.Request, models.BuildStatusNotStarted)
	if err != nil {
		return id, fmt.Errorf("error inserting build: %w", err)
	}
	return id, nil
}

var ErrNotFound = fmt.Errorf("not found")

// GetBuildByID fetches a build object from the DB
func (dl *PostgresDBLayer) GetBuildByID(ctx context.Context, id uuid.UUID) (models.Build, error) {
	out := models.Build{}
	var updated, completed pgtype.Timestamptz
	err := dl.p.QueryRow(ctx, `SELECT id, created, updated, completed, github_repo, github_ref, encrypted_github_credential, image_repos, tags, commit_sha_tag, request, status, events FROM builds WHERE id = $1;`, id).Scan(&out.ID, &out.Created, &updated, &completed, &out.GitHubRepo, &out.GitHubRef, &out.EncryptedGitHubCredential, &out.ImageRepos, &out.Tags, &out.CommitSHATag, &out.Request, &out.Status, &out.Events)
	if err != nil {
		if err == pgx.ErrNoRows {
			return out, ErrNotFound
		}
		return out, fmt.Errorf("error getting build by id: %w", err)
	}
	if updated.Status == pgtype.Present {
		out.Updated = updated.Time
	}
	if completed.Status == pgtype.Present {
		out.Completed = completed.Time
	}
	return out, nil
}

func (dl *PostgresDBLayer) SetBuildCompletedTimestamp(ctx context.Context, id uuid.UUID, completed time.Time) error {
	_, err := dl.p.Exec(ctx, `UPDATE builds SET completed = $1 WHERE id = $2;`, completed, id)
	return err
}

func (dl *PostgresDBLayer) SetBuildStatus(ctx context.Context, id uuid.UUID, s models.BuildStatus) error {
	_, err := dl.p.Exec(ctx, `UPDATE builds SET status = $1 WHERE id = $2;`, s, id)
	return err
}

// DeleteBuild removes a build from the DB.
func (dl *PostgresDBLayer) DeleteBuild(ctx context.Context, id uuid.UUID) (err error) {
	_, err = dl.p.Exec(ctx, `DELETE FROM builds WHERE id = $1;`, id)
	return err
}

// ListenForBuildEvents blocks and listens for the build events to occur for a build, writing any events that are received to c.
// If build is not currently listenable an error will be returned immediately.
// Always returns a non-nil error.
func (dl *PostgresDBLayer) ListenForBuildEvents(ctx context.Context, id uuid.UUID, c chan<- string) error {
	if c == nil {
		return fmt.Errorf("channel cannot be nil")
	}
	b, err := dl.GetBuildByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error getting build by id: %w", err)
	}
	if !b.EventListenable() {
		return fmt.Errorf("build status %v; no events are possible", b.Status.String())
	}
	conn, err := dl.p.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Release()
	if _, err := conn.Exec(ctx, fmt.Sprintf("LISTEN %s;", pgChanFromID(id))); err != nil {
		return fmt.Errorf("error listening on postgres channel: %w", err)
	}
	for {
		n, err := conn.Conn().WaitForNotification(ctx)
		if err != nil {
			return fmt.Errorf("error waiting for notification: %v: %w", id.String(), err)
		}
		c <- n.Payload
	}
}

// pgChanFromID returns a legal Postgres identifier from a build ID
func pgChanFromID(id uuid.UUID) string {
	return "build_" + strings.ReplaceAll(id.String(), "-", "_")
}

// AddEvent appends an event to a build and notifies any listeners to that channel
func (dl *PostgresDBLayer) AddEvent(ctx context.Context, id uuid.UUID, event string) error {
	b, err := dl.GetBuildByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error getting build by id: %w", err)
	}
	if !b.CanAddEvent() {
		return fmt.Errorf("build status %v; cannot add new event", b.Status.String())
	}
	txn, err := dl.p.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error opening txn: %w", err)
	}
	defer txn.Rollback(ctx)
	if _, err := txn.Exec(ctx, `UPDATE builds SET events = array_append(events, $1) WHERE id = $2;`, event, id); err != nil {
		return fmt.Errorf("error appending event: %w", err)
	}
	if _, err := txn.Exec(ctx, fmt.Sprintf("NOTIFY %s, '%s';", pgChanFromID(id), event)); err != nil {
		return fmt.Errorf("error notifying channel: %w", err)
	}
	if err := txn.Commit(ctx); err != nil {
		return fmt.Errorf("error committing txn: %w", err)
	}
	return nil
}

// pgCxlChanFromID returns a legal Postgres identifier for the cancellation notification from a build ID
func pgCxlChanFromID(id uuid.UUID) string {
	return "cxl_build_" + strings.ReplaceAll(id.String(), "-", "_")
}

// CancelBuild broadcasts a cancellation request for build id
func (dl *PostgresDBLayer) CancelBuild(ctx context.Context, id uuid.UUID) error {
	b, err := dl.GetBuildByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error getting build: %w", err)
	}
	if !b.Running() {
		return fmt.Errorf("build not cancellable: %v", b.Status.String())
	}
	txn, err := dl.p.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error opening txn: %w", err)
	}
	defer txn.Rollback(ctx)
	if _, err := txn.Exec(ctx, `UPDATE builds SET status = $1 WHERE id = $2;`, models.BuildStatusCancelRequested, id); err != nil {
		return fmt.Errorf("error cancelling build: %w", err)
	}
	q := fmt.Sprintf("NOTIFY %s, '%s';", pgCxlChanFromID(id), "cancel")
	if _, err := txn.Exec(ctx, q); err != nil {
		return fmt.Errorf("error notifying cancel channel: %w", err)
	}
	if err := txn.Commit(ctx); err != nil {
		return fmt.Errorf("error committing txn: %w", err)
	}
	return nil
}

// ListenForCancellation blocks and listens for cancellation requests for build id.
// If a cancellation request is received a nil error will be returned.
func (dl *PostgresDBLayer) ListenForCancellation(ctx context.Context, id uuid.UUID) error {
	b, err := dl.GetBuildByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error getting build by id: %w", err)
	}
	if !b.Running() {
		return fmt.Errorf("build not running: %v", b.Status.String())
	}
	conn, err := dl.p.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Release()
	q := fmt.Sprintf("LISTEN %s;", pgCxlChanFromID(id))
	if _, err := conn.Exec(ctx, q); err != nil {
		return fmt.Errorf("error listening on postgres cxl channel: %w", err)
	}
	_, err = conn.Conn().WaitForNotification(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for notification: %v: %w", id.String(), err)
	}
	return nil
}

// pgRunChanFromID returns a legal Postgres identifier for the running notification from a build ID
func pgRunChanFromID(id uuid.UUID) string {
	return "running_build_" + strings.ReplaceAll(id.String(), "-", "_")
}

// SetBuildAsRunning updates build status to Running and sends a notification for listeners
func (dl *PostgresDBLayer) SetBuildAsRunning(ctx context.Context, id uuid.UUID) error {
	txn, err := dl.p.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error opening txn: %w", err)
	}
	defer txn.Rollback(ctx)
	if _, err := txn.Exec(ctx, `UPDATE builds SET status = $1 WHERE id = $2;`, models.BuildStatusRunning, id); err != nil {
		return fmt.Errorf("error appending event: %w", err)
	}
	if _, err := txn.Exec(ctx, fmt.Sprintf("NOTIFY %s, 'running';", pgRunChanFromID(id))); err != nil {
		return fmt.Errorf("error notifying channel: %w", err)
	}
	if err := txn.Commit(ctx); err != nil {
		return fmt.Errorf("error committing txn: %w", err)
	}
	return nil
}

// ListenForBuildRunning blocks and listens for a build to be updated to Running.
// If a notification is received, a nil error will be returned
func (dl *PostgresDBLayer) ListenForBuildRunning(ctx context.Context, id uuid.UUID) error {
	b, err := dl.GetBuildByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error getting build by id: %w", err)
	}
	if b.Status != models.BuildStatusNotStarted {
		return fmt.Errorf("unexpected build status (wanted NotStarted): %v", b.Status.String())
	}
	conn, err := dl.p.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Release()
	q := fmt.Sprintf("LISTEN %s;", pgRunChanFromID(id))
	if _, err := conn.Exec(ctx, q); err != nil {
		return fmt.Errorf("error listening on postgres running channel: %w", err)
	}
	_, err = conn.Conn().WaitForNotification(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for notification: %v: %w", id.String(), err)
	}
	return nil
}

// pgCompletedChanFromID returns a legal Postgres identifier for the completed notification from a build ID
func pgCompletedChanFromID(id uuid.UUID) string {
	return "completed_build_" + strings.ReplaceAll(id.String(), "-", "_")
}

// SetBuildAsCompleted updates build status to completed (success, failure, skipped) and sends a notification for listeners
func (dl *PostgresDBLayer) SetBuildAsCompleted(ctx context.Context, id uuid.UUID, status models.BuildStatus) error {
	if !status.TerminalState() {
		return fmt.Errorf("invalid status for completed build: %v", status)
	}
	txn, err := dl.p.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error opening txn: %w", err)
	}
	defer txn.Rollback(ctx)
	if _, err := txn.Exec(ctx, `UPDATE builds SET status = $1 WHERE id = $2;`, status, id); err != nil {
		return fmt.Errorf("error setting status: %w", err)
	}
	if _, err := txn.Exec(ctx, fmt.Sprintf("NOTIFY %s, '%d';", pgCompletedChanFromID(id), status)); err != nil {
		return fmt.Errorf("error notifying channel: %w", err)
	}
	if err := txn.Commit(ctx); err != nil {
		return fmt.Errorf("error committing txn: %w", err)
	}
	return nil
}

// ListenForBuildCompleted blocks and listens for a build to be updated to completed.
// If a notification is received, the completed build status and a nil error are returned
func (dl *PostgresDBLayer) ListenForBuildCompleted(ctx context.Context, id uuid.UUID) (models.BuildStatus, error) {
	b, err := dl.GetBuildByID(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("error getting build by id: %w", err)

	}
	// if build is already finished, return status
	if b.Status.TerminalState() {
		return b.Status, nil
	}

	conn, err := dl.p.Acquire(ctx)
	if err != nil {
		return 0, fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Release()
	q := fmt.Sprintf("LISTEN %s;", pgCompletedChanFromID(id))
	if _, err := conn.Exec(ctx, q); err != nil {
		return 0, fmt.Errorf("error listening on postgres running channel: %w", err)
	}
	sn, err := conn.Conn().WaitForNotification(ctx)
	if err != nil {
		return 0, fmt.Errorf("error waiting for notification: %v: %w", id.String(), err)
	}
	if sn == nil {
		return 0, fmt.Errorf("nil notification")
	}
	si, err := strconv.Atoi(sn.Payload)
	if err != nil {
		return 0, fmt.Errorf("error parsing status received via notification: %w", err)
	}
	bs := models.BuildStatus(si)
	if !bs.TerminalState() {
		return 0, fmt.Errorf("invalid status for completed build: %v", bs)
	}
	return bs, nil
}

func (dl *PostgresDBLayer) spewerr(err error) {
	var pgerr *pgconn.PgError
	if errors.As(err, &pgerr) {
		spew.Dump(pgerr)
	}
}