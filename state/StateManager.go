package state

import "fmt"

func NewStateManager() *StateManager {
	return &StateManager{
		states: make(map[int64]*UserState),
	}
}
func (sm *StateManager) GetState(chatID int64) *UserState {
	if state, exists := sm.states[chatID]; exists {
		return state
	}
	state := &UserState{
		Current: Idle,
		Context: make(map[string]string),
	}
	sm.states[chatID] = state
	return state
}
func (sm *StateManager) SetState(chatID int64, newState State) {
	state := sm.GetState(chatID)
	state.Current = newState
}
func (sm *StateManager) GetContext(chatID int64, key string) (string, error) {
	state := sm.GetState(chatID)
	if value, exists := state.Context[key]; exists {
		return value, nil
	}
	return "", fmt.Errorf("no context for key: %s", key)
}
func (sm *StateManager) SetContext(chatID int64, key string, value string) {
	state := sm.GetState(chatID)
	state.Context[key] = value
}
