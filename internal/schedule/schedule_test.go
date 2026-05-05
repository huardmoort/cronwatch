package schedule_test

import (
	"testing"
	"time"

	"github.com/cronwatch/cronwatch/internal/schedule"
)

func TestNextRun_ValidExpression(t *testing.T) {
	// Every minute: next run should be within 60 seconds
	from := time.Now()
	next, err := schedule.NextRun("* * * * *", from)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !next.After(from) {
		t.Errorf("expected next run to be after %v, got %v", from, next)
	}
	if next.Sub(from) > 60*time.Second {
		t.Errorf("next run too far in the future: %v", next.Sub(from))
	}
}

func TestNextRun_InvalidExpression(t *testing.T) {
	_, err := schedule.NextRun("not-a-cron", time.Now())
	if err == nil {
		t.Error("expected error for invalid cron expression, got nil")
	}
}

func TestIsMissed_NotMissed(t *testing.T) {
	// Job ran just now; it should not be missed.
	missed, err := schedule.IsMissed("* * * * *", time.Now(), 30*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if missed {
		t.Error("expected job not to be missed when lastRun is now")
	}
}

func TestIsMissed_Missed(t *testing.T) {
	// Job last ran 10 minutes ago on a every-minute schedule with no grace.
	lastRun := time.Now().Add(-10 * time.Minute)
	missed, err := schedule.IsMissed("* * * * *", lastRun, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !missed {
		t.Error("expected job to be missed")
	}
}

func TestDurationUntilNext_Positive(t *testing.T) {
	d, err := schedule.DurationUntilNext("* * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d <= 0 {
		t.Errorf("expected positive duration, got %v", d)
	}
	if d > 60*time.Second {
		t.Errorf("duration unexpectedly large: %v", d)
	}
}
