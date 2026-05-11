package alertlog_test

import (
	"testing"
	"time"

	"github.com/example/cronwatch/internal/alertlog"
)

func TestMute_IsMuted_BeforeExpiry(t *testing.T) {
	p := alertlog.NewMutePolicy()
	p.Mute("backup", 5*time.Minute)

	if !p.IsMuted("backup") {
		t.Fatal("expected job to be muted")
	}
}

func TestMute_NotMuted_WhenNeverSet(t *testing.T) {
	p := alertlog.NewMutePolicy()

	if p.IsMuted("backup") {
		t.Fatal("expected job to not be muted")
	}
}

func TestMute_Expired_ReturnsFalse(t *testing.T) {
	p := alertlog.NewMutePolicy()
	p.Mute("backup", -1*time.Second) // already expired

	if p.IsMuted("backup") {
		t.Fatal("expected expired mute to return false")
	}
}

func TestUnmute_ClearsActiveMute(t *testing.T) {
	p := alertlog.NewMutePolicy()
	p.Mute("backup", 5*time.Minute)
	p.Unmute("backup")

	if p.IsMuted("backup") {
		t.Fatal("expected job to be unmuted after explicit Unmute")
	}
}

func TestMute_ExtendsDuration(t *testing.T) {
	p := alertlog.NewMutePolicy()
	p.Mute("backup", 1*time.Millisecond)
	p.Mute("backup", 5*time.Minute) // extend

	if !p.IsMuted("backup") {
		t.Fatal("expected extended mute to still be active")
	}
}

func TestActiveMutes_ReturnsOnlyActive(t *testing.T) {
	p := alertlog.NewMutePolicy()
	p.Mute("job-a", 5*time.Minute)
	p.Mute("job-b", -1*time.Second) // expired

	active := p.ActiveMutes()

	if _, ok := active["job-a"]; !ok {
		t.Error("expected job-a in active mutes")
	}
	if _, ok := active["job-b"]; ok {
		t.Error("expected job-b to be absent (expired)")
	}
}

func TestActiveMutes_EmptyWhenNone(t *testing.T) {
	p := alertlog.NewMutePolicy()

	active := p.ActiveMutes()
	if len(active) != 0 {
		t.Fatalf("expected empty map, got %d entries", len(active))
	}
}

func TestMute_DifferentJobsIndependent(t *testing.T) {
	p := alertlog.NewMutePolicy()
	p.Mute("job-a", 5*time.Minute)

	if p.IsMuted("job-b") {
		t.Fatal("muting job-a should not affect job-b")
	}
}
