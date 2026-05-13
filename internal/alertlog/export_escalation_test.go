package alertlog

import "time"

// NewEscalationPolicyExported exposes the constructor for white-box testing.
func NewEscalationPolicyExported(threshold int, window time.Duration) *EscalationPolicy {
	return NewEscalationPolicy(threshold, window)
}
