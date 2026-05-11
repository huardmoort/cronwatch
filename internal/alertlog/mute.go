package alertlog

import (
	"sync"
	"time"
)

// MutePolicy suppresses alerts for specific jobs until a configured expiry time.
// This allows operators to silence noisy jobs during maintenance windows.
type MutePolicy struct {
	mu    sync.Mutex
	mutes map[string]time.Time // job name -> mute expiry
}

// NewMutePolicy returns an initialised MutePolicy with no active mutes.
func NewMutePolicy() *MutePolicy {
	return &MutePolicy{
		mutes: make(map[string]time.Time),
	}
}

// Mute silences alerts for the named job until now+duration.
// Calling Mute again before expiry extends the window.
func (p *MutePolicy) Mute(job string, duration time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.mutes[job] = time.Now().Add(duration)
}

// Unmute removes any active mute for the named job immediately.
func (p *MutePolicy) Unmute(job string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.mutes, job)
}

// IsMuted reports whether the named job is currently muted.
// Expired mutes are pruned lazily on each call.
func (p *MutePolicy) IsMuted(job string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	expiry, ok := p.mutes[job]
	if !ok {
		return false
	}
	if time.Now().After(expiry) {
		delete(p.mutes, job)
		return false
	}
	return true
}

// ActiveMutes returns a snapshot of job names and their expiry times.
func (p *MutePolicy) ActiveMutes() map[string]time.Time {
	p.mu.Lock()
	defer p.mu.Unlock()
	now := time.Now()
	out := make(map[string]time.Time)
	for job, expiry := range p.mutes {
		if now.Before(expiry) {
			out[job] = expiry
		}
	}
	return out
}
