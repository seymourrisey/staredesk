package usecase

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/seymourrisey/staredesk/internal/entity"
	"github.com/seymourrisey/staredesk/internal/repository"
	"github.com/seymourrisey/staredesk/pkg/idgen"
)

const defaultAwayTimeoutMinutes = 5

type sessionState struct {
	mu               sync.Mutex
	activeSessionID  string
	sessionStartedAt *time.Time // tambah ini
	conditionCounts  map[string]int
	awayStartedAt    *time.Time
}

type SessionUsecase struct {
	sessionRepo      repository.SessionRepository
	deviceConfigRepo repository.DeviceConfigRepository
	state            sessionState
}

func NewSessionUsecase(
	sessionRepo repository.SessionRepository,
	deviceConfigRepo repository.DeviceConfigRepository,
) *SessionUsecase {
	return &SessionUsecase{
		sessionRepo:      sessionRepo,
		deviceConfigRepo: deviceConfigRepo,
		state: sessionState{
			conditionCounts: make(map[string]int),
		},
	}
}

// SessionEvent adalah event yang dikembalikan ProcessCondition ke caller.
type SessionEvent string

const (
	SessionEventNone  SessionEvent = ""
	SessionEventStart SessionEvent = "session_start"
	SessionEventEnd   SessionEvent = "session_end"
)

// ProcessCondition sekarang return SessionEvent selain error.
func (u *SessionUsecase) ProcessCondition(ctx context.Context, userID, condition string, ts time.Time) (SessionEvent, error) {
	u.state.mu.Lock()
	defer u.state.mu.Unlock()

	if condition == "away" {
		return u.handleAway(ctx, userID, ts)
	}
	return u.handlePresent(ctx, userID, condition, ts)
}

func (u *SessionUsecase) handlePresent(ctx context.Context, userID, condition string, ts time.Time) (SessionEvent, error) {
	u.state.awayStartedAt = nil
	event := SessionEventNone

	if u.state.activeSessionID == "" {
		existing, err := u.sessionRepo.GetActiveByUserID(ctx, userID)
		if err != nil {
			return SessionEventNone, err
		}
		if existing != nil {
			u.state.activeSessionID = existing.ID
			u.state.sessionStartedAt = &existing.StartedAt
			log.Printf("[session] resumed existing session %s", existing.ID)
		} else {
			newSession := &entity.Session{
				ID:        idgen.NewSessionID(),
				UserID:    userID,
				StartedAt: ts,
			}
			if err := u.sessionRepo.Create(ctx, newSession); err != nil {
				return SessionEventNone, err
			}
			u.state.activeSessionID = newSession.ID
			u.state.sessionStartedAt = &newSession.StartedAt
			u.state.conditionCounts = make(map[string]int)
			log.Printf("[session] started new session %s", newSession.ID)
			event = SessionEventStart
		}
	}

	u.state.conditionCounts[condition]++
	return event, nil
}

func (u *SessionUsecase) handleAway(ctx context.Context, userID string, ts time.Time) (SessionEvent, error) {
	if u.state.activeSessionID == "" {
		return SessionEventNone, nil
	}

	if u.state.awayStartedAt == nil {
		u.state.awayStartedAt = &ts
		log.Printf("[session] away started at %s", ts.Format(time.RFC3339))
		return SessionEventNone, nil
	}

	awayDuration := ts.Sub(*u.state.awayStartedAt)
	timeoutMinutes := u.getAwayTimeout(ctx, userID)
	timeoutDuration := time.Duration(timeoutMinutes) * time.Minute

	if awayDuration < timeoutDuration {
		return SessionEventNone, nil
	}

	log.Printf("[session] away timeout reached (%.1f min), ending session %s",
		awayDuration.Minutes(), u.state.activeSessionID)

	dominant := u.calcDominantCondition()
	endedAt := ts
	startedAt, err := u.getSessionStartedAt(ctx, u.state.activeSessionID, userID)
	if err != nil {
		return SessionEventNone, err
	}

	durationSec := int(ts.Sub(startedAt).Seconds())
	updated := &entity.Session{
		ID:                u.state.activeSessionID,
		EndedAt:           &endedAt,
		DurationSec:       &durationSec,
		DominantCondition: &dominant,
	}

	if err := u.sessionRepo.Update(ctx, updated); err != nil {
		return SessionEventNone, err
	}

	log.Printf("[session] session %s ended — duration %ds, dominant: %s",
		u.state.activeSessionID, durationSec, dominant)

	endedSessionID := u.state.activeSessionID
	u.state.activeSessionID = ""
	u.state.sessionStartedAt = nil
	u.state.conditionCounts = make(map[string]int)
	u.state.awayStartedAt = nil

	_ = endedSessionID // akan dipakai di handler untuk payload
	return SessionEventEnd, nil
}

func (u *SessionUsecase) getAwayTimeout(ctx context.Context, userID string) int {
	config, err := u.deviceConfigRepo.GetByUserID(ctx, userID)
	if err != nil || config == nil {
		return defaultAwayTimeoutMinutes
	}
	return config.AwayTimeoutMinutes
}

func (u *SessionUsecase) calcDominantCondition() string {
	dominant := "away"
	maxCount := 0
	for condition, count := range u.state.conditionCounts {
		if count > maxCount {
			maxCount = count
			dominant = condition
		}
	}
	return dominant
}

func (u *SessionUsecase) getSessionStartedAt(ctx context.Context, sessionID, userID string) (time.Time, error) {
	session, err := u.sessionRepo.GetByID(ctx, sessionID, userID)
	if err != nil {
		return time.Time{}, err
	}
	if session == nil {
		log.Printf("[session] WARNING: session %s not found in DB during end", sessionID)
		return time.Now(), nil
	}
	return session.StartedAt, nil
}

func (u *SessionUsecase) ActiveSessionID() string {
	u.state.mu.Lock()
	defer u.state.mu.Unlock()
	return u.state.activeSessionID
}

func (u *SessionUsecase) IsSessionActive() bool {
	u.state.mu.Lock()
	defer u.state.mu.Unlock()
	return u.state.activeSessionID != ""
}

func (u *SessionUsecase) ActiveSessionStartedAt() *time.Time {
	u.state.mu.Lock()
	defer u.state.mu.Unlock()
	return u.state.sessionStartedAt
}

// --- REST API methods ---

// SessionListResult adalah hasil GetAll dengan pagination info.
type SessionListResult struct {
	Sessions []*entity.Session
	Total    int
	Limit    int
	Offset   int
}

// SessionSummaryResult adalah hasil GetSummary untuk range tertentu.
type SessionSummaryResult struct {
	Range        string
	TotalSec     int
	SessionCount int
	Sessions     []*entity.Session
}

// GetAll mengembalikan list sesi selesai dengan pagination.
func (u *SessionUsecase) GetAll(ctx context.Context, userID string, limit, offset int) (*SessionListResult, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	sessions, total, err := u.sessionRepo.GetAll(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	if sessions == nil {
		sessions = []*entity.Session{}
	}

	return &SessionListResult{
		Sessions: sessions,
		Total:    total,
		Limit:    limit,
		Offset:   offset,
	}, nil
}

// GetByID mengembalikan detail satu sesi.
func (u *SessionUsecase) GetByID(ctx context.Context, id, userID string) (*entity.Session, error) {
	return u.sessionRepo.GetByID(ctx, id, userID)
}

// GetSummary mengembalikan ringkasan sesi untuk range today/week/month.
func (u *SessionUsecase) GetSummary(ctx context.Context, userID, rangeParam string) (*SessionSummaryResult, error) {
	now := time.Now()
	var from, to time.Time

	switch rangeParam {
	case "week":
		from = now.AddDate(0, 0, -7)
		to = now
	case "month":
		from = now.AddDate(0, -1, 0)
		to = now
	default: // "today"
		from = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		to = from.Add(24 * time.Hour)
		rangeParam = "today"
	}

	sessions, err := u.sessionRepo.GetByDateRange(ctx, userID, from, to)
	if err != nil {
		return nil, err
	}
	if sessions == nil {
		sessions = []*entity.Session{}
	}

	totalSec := 0
	for _, s := range sessions {
		if s.DurationSec != nil {
			totalSec += *s.DurationSec
		}
	}

	return &SessionSummaryResult{
		Range:        rangeParam,
		TotalSec:     totalSec,
		SessionCount: len(sessions),
		Sessions:     sessions,
	}, nil
}
