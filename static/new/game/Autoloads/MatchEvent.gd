extends Node

enum {
	MOVEMENT = 0,
	FARMING = 1,
	FARM_PERMISSION_UPDATE = 2,
	FISHING_LURE = 3,
	DATA_RESET = 4,
	AVATAR_UPDATE = 5,
	MATCH_JOIN = 6,
	FISHING_VICTORY = 7
}

signal movement(state, presence)
signal farming(state, presence)
signal farm_permission_update(state, presence)
signal data_reset(state)
signal avatar_update(state, presence)
signal match_join(state, presence)
signal fishing_lure(state, presence)
signal fishing_victory(state, presence)

var signal_map = {
	MOVEMENT: "movement",
	FARMING: "farming",
	FARM_PERMISSION_UPDATE: "farm_permission_update",
	DATA_RESET: "data_reset",
	AVATAR_UPDATE: "avatar_update",
	MATCH_JOIN: "match_join",
	FISHING_LURE: "fishing_lure",
	FISHING_VICTORY: "fishing_victory"
}


func handle_match_state_update(state: NakamaRTAPI.MatchData):
	if state.presence != null && state.presence.user_id == SessionManager.session.user_id:
		return

	match state.op_code:
		MOVEMENT:
			emit_signal("movement", state.data, state.presence)
		FARMING:
			emit_signal("farming", state.data, state.presence)
		FARM_PERMISSION_UPDATE:
			emit_signal("farm_permission_update", state.data, state.presence)
		DATA_RESET:
			emit_signal("data_reset", state.data, state.presence)
		AVATAR_UPDATE:
			emit_signal("avatar_update", state.data, state.presence)
		MATCH_JOIN:
			emit_signal("match_join", state.data, state.presence)
		FISHING_LURE:
			emit_signal("fishing_lure", state.data, state.presence)
		FISHING_VICTORY:
			emit_signal("fishing_victory", state.data, state.presence)


func emit(op_code: int, payload: String):
	var event: String = signal_map[op_code]
	emit_signal(event, payload, {"user_id": SessionManager.session.user_id})

	if MatchManager.game_match != null:
		MatchManager.socket.send_match_state_async(
			MatchManager.game_match.match_id, op_code, payload
		)


func movement(payload):
	emit(MOVEMENT, JSON.print(payload))


func farming(payload):
	emit(FARMING, JSON.print(payload))


func farm_permission_update(payload):
	emit(FARM_PERMISSION_UPDATE, JSON.print(payload))


func avatar_update(payload):
	emit(AVATAR_UPDATE, JSON.print(payload))


func match_join(payload):
	emit(MATCH_JOIN, JSON.print(payload))


func fishing_lure():
	emit(FISHING_LURE, "fishing_lure")


func fishing_victory(payload):
	emit(FISHING_VICTORY, JSON.print(payload))
