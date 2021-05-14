extends Node

var socket: NakamaSocket
var game_match: NakamaRTAPI.Match

const dungeon_label_mappings: Dictionary = {
	"town0": "Downtown", "town0CommunityFarm": "Community Farm"
}


func connect_socket():
	socket = Nakama.create_socket_from(SessionManager.client)
	yield(socket.connect_async(SessionManager.session), "completed")
	socket.connect("received_match_state", self, "_on_match_state")
	socket.connect("closed", self, "_on_socket_disconnect")


func find_or_create_match(label: String, starting_position: Vector2):
	var response = yield(
		SessionManager.rpc_async("find_or_create_dungeon", "dungeon_type"), "completed"
	)

	var match_id = response.payload.replace('"', "")

	var _match = yield(_join_match(match_id, label, starting_position), "completed")

	return _match


func leave_match():
	if game_match != null:
		yield(socket.leave_match_async(game_match.match_id), "completed")
		game_match = null
	else:
		yield()


func _join_match(match_id: String, label: String, starting_position: Vector2) -> NakamaRTAPI.Match:
	game_match = yield(
		socket.join_match_async(
			match_id,
			{
				"x": String(starting_position.x),
				"y": String(starting_position.y),
			}
		),
		"completed"
	)

	if game_match.is_exception():
		print("An error occured: %s" % game_match)
		return

	var dungeon_label = (
		dungeon_label_mappings[label]
		if dungeon_label_mappings.has(label)
		else label
	)

	return game_match


func _on_match_state(state: NakamaRTAPI.MatchData):
	MatchEvent.handle_match_state_update(state)


func _on_socket_disconnect():
	# TPLG.base_change_scene("res://RootScenes/Authentication/Authentication.tscn")
	pass
