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
	sessionRepo repository.SessionRepository
	state       sessionState
}

func NewSessionUsecase(sessionRepo repository.SessionRepository) *SessionUsecase {
	return &SessionUsecase{
		sessionRepo: sessionRepo,
		state: sessionState{
			conditionCounts: make(map[string]int),
		},
	}
}

// ProcessCondition dipanggil setiap kali MQTT handler menerima telemetry.
// userID  : dari config/env
// condition: kondisi yang sudah dievaluasi di ESP32
// ts      : timestamp telemetry (gunakan time.Now() jika tidak ada)
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
	// Reset away timer setiap kali ada kehadiran
	u.state.awayStartedAt = nil

	// Jika belum ada sesi aktif, cek DB dulu (restart safety)
	if u.state.activeSessionID == "" {
		existing, err := u.sessionRepo.GetActiveByUserID(ctx, userID)
		if err != nil {
			return err
		}
		if existing != nil {
			// Resume sesi yang ada di DB (misal backend restart)
			u.state.activeSessionID = existing.ID
			log.Printf("[session] resumed existing session %s", existing.ID)
		} else {
			// Mulai sesi baru
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

	// Catat kondisi untuk kalkulasi dominant_condition
	u.state.conditionCounts[condition]++
	return nil
}

// handleAway — kondisi away
func (u *SessionUsecase) handleAway(ctx context.Context, userID string, ts time.Time) error {
	// Tidak ada sesi aktif — tidak ada yang perlu dilakukan
	if u.state.activeSessionID == "" {
		return nil
	}

	// Mulai / lanjutkan away timer
	if u.state.awayStartedAt == nil {
		u.state.awayStartedAt = &ts
		log.Printf("[session] away started at %s", ts.Format(time.RFC3339))
		return nil
	}

	awayDuration := ts.Sub(*u.state.awayStartedAt)
	timeoutDuration := time.Duration(defaultAwayTimeoutMinutes) * time.Minute

	if awayDuration < timeoutDuration {
		// Belum timeout — tunggu
		return nil
	}

	// Away timeout tercapai — akhiri sesi
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

	// Reset state
	u.state.activeSessionID = ""
	u.state.conditionCounts = make(map[string]int)
	u.state.awayStartedAt = nil

	return nil
}

// calcDominantCondition returns kondisi dengan frekuensi terbanyak.
// Jika kosong (tidak ada data), return "away".
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
		// Fallback — seharusnya tidak terjadi
		log.Printf("[session] WARNING: session %s not found in DB during end", sessionID)
		return time.Now(), nil
	}
	return session.StartedAt, nil
}

// ActiveSessionID returns ID sesi aktif saat ini, atau "" jika tidak ada.
// Digunakan oleh MQTT handler untuk inject ke WebSocket payload.
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
