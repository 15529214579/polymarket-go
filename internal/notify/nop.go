package notify

import "context"

// Nop discards every event. Used in tests, offline runs, and any environment
// where TELEGRAM_BOT_TOKEN / TELEGRAM_CHAT_ID are not set.
type Nop struct{}

func (Nop) RiskTrip(RiskTripEvent)         {}
func (Nop) RiskResume(RiskResumeEvent)     {}
func (Nop) LargeFill(LargeFillEvent)       {}
func (Nop) SignalPrompt(SignalPromptEvent) {}
func (Nop) Close(context.Context) error    { return nil }
