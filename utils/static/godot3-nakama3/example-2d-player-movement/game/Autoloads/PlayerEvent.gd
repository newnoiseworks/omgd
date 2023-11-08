extends Node

enum {
  MOVEMENT = 0,
}


signal movement(state, presence)


var signal_map = {
	MOVEMENT: "movement",
}


func handle_match_state_update(state: NakamaRTAPI.MatchData):
	if state.presence != null && state.presence.user_id == SessionManager.session.user_id:
		return

	match state.op_code:
		MOVEMENT:
			emit_signal("movement", state.data, state.presence)


func emit(op_code: int, payload: String):
	var event: String = signal_map[op_code]
	emit_signal(event, payload, {"user_id": SessionManager.session.user_id})

	if PlayerManager.game_match != null:
		PlayerManager.socket.send_match_state_async(
			PlayerManager.game_match.match_id, op_code, payload
		)




func movement(payload):
	emit(MOVEMENT, JSON.print(payload))
