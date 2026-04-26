package usecase

import (
	"context"
	"time"

	"github.com/seymourrisey/staredesk/internal/repository"
)

type AnalyticsUsecase struct {
	analyticsRepo repository.AnalyticsRepository
}

func NewAnalyticsUsecase(analyticsRepo repository.AnalyticsRepository) *AnalyticsUsecase {
	return &AnalyticsUsecase{analyticsRepo: analyticsRepo}
}

// --- Result structs ---

type PeakHoursResult struct {
	Range   string
	Entries []*repository.PeakHourEntry
}

type ConditionBreakdownResult struct {
	Range   string
	Entries []*repository.ConditionBreakdownEntry
}

type TimelineResult struct {
	Date    string
	Entries []*repository.TimelineEntry
}

// --- Helpers ---

// parseRange mengkonversi "today/week/month" ke from/to time.Time.
func parseRange(rangeParam string) (from, to time.Time, label string) {
	now := time.Now()
	switch rangeParam {
	case "week":
		return now.AddDate(0, 0, -7), now, "week"
	case "month":
		return now.AddDate(0, -1, 0), now, "month"
	default:
		dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		return dayStart, dayStart.Add(24 * time.Hour), "today"
	}
}

// --- Usecase methods ---

func (u *AnalyticsUsecase) GetPeakHours(ctx context.Context, userID, rangeParam string) (*PeakHoursResult, error) {
	from, to, label := parseRange(rangeParam)

	entries, err := u.analyticsRepo.GetPeakHours(ctx, userID, from, to)
	if err != nil {
		return nil, err
	}
	if entries == nil {
		entries = []*repository.PeakHourEntry{}
	}

	return &PeakHoursResult{
		Range:   label,
		Entries: entries,
	}, nil
}

func (u *AnalyticsUsecase) GetConditionBreakdown(ctx context.Context, userID, rangeParam string) (*ConditionBreakdownResult, error) {
	from, to, label := parseRange(rangeParam)

	entries, err := u.analyticsRepo.GetConditionBreakdown(ctx, userID, from, to)
	if err != nil {
		return nil, err
	}
	if entries == nil {
		entries = []*repository.ConditionBreakdownEntry{}
	}

	return &ConditionBreakdownResult{
		Range:   label,
		Entries: entries,
	}, nil
}

func (u *AnalyticsUsecase) GetTimeline(ctx context.Context, userID, dateParam string) (*TimelineResult, error) {
	var date time.Time
	var err error

	if dateParam == "" {
		date = time.Now()
	} else {
		date, err = time.Parse("2006-01-02", dateParam)
		if err != nil {
			date = time.Now()
		}
	}

	entries, err := u.analyticsRepo.GetTimeline(ctx, userID, date)
	if err != nil {
		return nil, err
	}
	if entries == nil {
		entries = []*repository.TimelineEntry{}
	}

	return &TimelineResult{
		Date:    date.Format("2006-01-02"),
		Entries: entries,
	}, nil
}
