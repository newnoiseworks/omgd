extends Node

enum {
	MOVEMENT = 0,
	MATCH_JOIN = 1,
}

signal movement(state, presence)
signal match_join(state, presence)

var signal_map = {
	MOVEMENT: "movement",
	MATCH_JOIN: "match_join",
}


func handle_match_state_update(state: NakamaRTAPI.MatchData):
	if state.presence != null && state.presence.user_id == SessionManager.session.user_id:
		return

	match state.op_code:
		MOVEMENT:
			emit_signal("movement", state.data, state.presence)
		MATCH_JOIN:
			emit_signal("match_join", state.data, state.presence)


func emit(op_code: int, payload: String):
	var event: String = signal_map[op_code]
	emit_signal(event, payload, {"user_id": SessionManager.session.user_id})

	if MatchManager.game_match != null:
		MatchManager.socket.send_match_state_async(
			MatchManager.game_match.match_id, op_code, payload
		)


func movement(payload):
	emit(MOVEMENT, JSON.print(payload))
