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

// sessionState holds in-memory state for the current session.
// Single user — satu instance cukup.
type sessionState struct {
	mu              sync.Mutex
	activeSessionID string
	conditionCounts map[string]int
	awayStartedAt   *time.Time
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

// ProcessCondition dipanggil setiap kali MQTT handler menerima telemetry.
func (u *SessionUsecase) ProcessCondition(ctx context.Context, userID, condition string, ts time.Time) error {
	u.state.mu.Lock()
	defer u.state.mu.Unlock()

	if condition == "away" {
		return u.handleAway(ctx, userID, ts)
	}
	return u.handlePresent(ctx, userID, condition, ts)
}

// handlePresent — kondisi bukan away
func (u *SessionUsecase) handlePresent(ctx context.Context, userID, condition string, ts time.Time) error {
	u.state.awayStartedAt = nil

	if u.state.activeSessionID == "" {
		existing, err := u.sessionRepo.GetActiveByUserID(ctx, userID)
		if err != nil {
			return err
		}
		if existing != nil {
			u.state.activeSessionID = existing.ID
			log.Printf("[session] resumed existing session %s", existing.ID)
		} else {
			newSession := &entity.Session{
				ID:        idgen.NewSessionID(),
				UserID:    userID,
				StartedAt: ts,
			}
			if err := u.sessionRepo.Create(ctx, newSession); err != nil {
				return err
			}
			u.state.activeSessionID = newSession.ID
			u.state.conditionCounts = make(map[string]int)
			log.Printf("[session] started new session %s", newSession.ID)
		}
	}

	u.state.conditionCounts[condition]++
	return nil
}

// handleAway — kondisi away
func (u *SessionUsecase) handleAway(ctx context.Context, userID string, ts time.Time) error {
	if u.state.activeSessionID == "" {
		return nil
	}

	if u.state.awayStartedAt == nil {
		u.state.awayStartedAt = &ts
		log.Printf("[session] away started at %s", ts.Format(time.RFC3339))
		return nil
	}

	awayDuration := ts.Sub(*u.state.awayStartedAt)
	timeoutMinutes := u.getAwayTimeout(ctx, userID)
	timeoutDuration := time.Duration(timeoutMinutes) * time.Minute

	if awayDuration < timeoutDuration {
		return nil
	}

	log.Printf("[session] away timeout reached (%.1f min), ending session %s",
		awayDuration.Minutes(), u.state.activeSessionID)

	dominant := u.calcDominantCondition()
	endedAt := ts
	startedAt, err := u.getSessionStartedAt(ctx, u.state.activeSessionID)
	if err != nil {
		return err
	}

	durationSec := int(ts.Sub(startedAt).Seconds())

	updated := &entity.Session{
		ID:                u.state.activeSessionID,
		EndedAt:           &endedAt,
		DurationSec:       &durationSec,
		DominantCondition: &dominant,
	}

	if err := u.sessionRepo.Update(ctx, updated); err != nil {
		return err
	}

	log.Printf("[session] session %s ended — duration %ds, dominant: %s",
		u.state.activeSessionID, durationSec, dominant)

	u.state.activeSessionID = ""
	u.state.conditionCounts = make(map[string]int)
	u.state.awayStartedAt = nil

	return nil
}

// getAwayTimeout baca dari DB, fallback ke default jika row kosong atau error.
func (u *SessionUsecase) getAwayTimeout(ctx context.Context, userID string) int {
	config, err := u.deviceConfigRepo.GetByUserID(ctx, userID)
	if err != nil || config == nil {
		return defaultAwayTimeoutMinutes
	}
	return config.AwayTimeoutMinutes
}

// calcDominantCondition returns kondisi dengan frekuensi terbanyak.
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

// getSessionStartedAt fetch started_at dari DB untuk kalkulasi duration.
func (u *SessionUsecase) getSessionStartedAt(ctx context.Context, sessionID string) (time.Time, error) {
	session, err := u.sessionRepo.GetActiveByUserID(ctx, u.state.activeSessionID)
	if err != nil {
		return time.Time{}, err
	}
	if session == nil {
		log.Printf("[session] WARNING: session %s not found in DB during end", sessionID)
		return time.Now(), nil
	}
	return session.StartedAt, nil
}

// ActiveSessionID returns ID sesi aktif saat ini, atau "" jika tidak ada.
func (u *SessionUsecase) ActiveSessionID() string {
	u.state.mu.Lock()
	defer u.state.mu.Unlock()
	return u.state.activeSessionID
}

// IsSessionActive returns true jika ada sesi yang sedang berjalan.
func (u *SessionUsecase) IsSessionActive() bool {
	u.state.mu.Lock()
	defer u.state.mu.Unlock()
	return u.state.activeSessionID != ""
}
