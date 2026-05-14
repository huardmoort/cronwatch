package alertlog

import (
	"testing"
	"time"
)

func TestAck_NotAcked_WhenNeverSet(t *testing.T) {
	p := NewAckPolicy()
	if p.IsAcked("backup") {
		t.Fatal("expected not acked")
	}
}

func TestAck_IsAcked_BeforeExpiry(t *testing.T) {
	p := NewAckPolicy()
	p.Acknowledge("backup", 10*time.Minute)
	if !p.IsAcked("backup") {
		t.Fatal("expected acked")
	}
}

func TestAck_Expired_ReturnsFalse(t *testing.T) {
	p := NewAckPolicy()
	p.Acknowledge("backup", 1*time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	if p.IsAcked("backup") {
		t.Fatal("expected ack to have expired")
	}
}

func TestAck_ZeroDuration_IsPermanent(t *testing.T) {
	p := NewAckPolicy()
	p.Acknowledge("backup", 0)
	if !p.IsAcked("backup") {
		t.Fatal("expected permanent ack")
	}
}

func TestAck_Clear_RemovesAck(t *testing.T) {
	p := NewAckPolicy()
	p.Acknowledge("backup", 10*time.Minute)
	p.Clear("backup")
	if p.IsAcked("backup") {
		t.Fatal("expected ack to be cleared")
	}
}

func TestAck_DifferentJobs_Independent(t *testing.T) {
	p := NewAckPolicy()
	p.Acknowledge("job-a", 10*time.Minute)
	if p.IsAcked("job-b") {
		t.Fatal("job-b should not be acked")
	}
	if !p.IsAcked("job-a") {
		t.Fatal("job-a should be acked")
	}
}

func TestAck_AckedJobs_ReturnsList(t *testing.T) {
	p := NewAckPolicy()
	p.Acknowledge("job-a", 10*time.Minute)
	p.Acknowledge("job-b", 0)
	jobs := p.AckedJobs()
	if len(jobs) != 2 {
		t.Fatalf("expected 2 acked jobs, got %d", len(jobs))
	}
}
