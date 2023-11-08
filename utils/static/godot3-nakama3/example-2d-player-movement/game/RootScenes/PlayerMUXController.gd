extends "res://RootScenes/RootController.gd"

export var character_scene: PackedScene

var channel_name = "player"

var user_ids: Array = []


func _ready():
	user_ids.append(SessionManager.session.user_id)

	if channel_name != null && channel_name != '':
		yield(PlayerManager.connect_socket(), "completed")
		yield(_join_player(), "completed")
		var _rm = PlayerManager.socket.connect("received_match_presence", self, "_on_match_presence")



func _exit_tree():
	if PlayerManager.game_match != null:
		PlayerManager.socket.disconnect("received_match_presence", self, "_on_match_presence")
		PlayerManager.game_match = null


func _join_player():
	var player_channel_object = yield(PlayerManager.find_or_create_match(channel_name, player_entry_node.position), "completed")

	for presence in player_channel_object.presences:
		_handle_match_join_event(presence)


func _on_match_presence(match_event: NakamaRTAPI.MatchPresenceEvent):
	for presence in match_event.leaves:
		user_ids.erase(presence.user_id)
		var user_to_erase = find_node(presence.user_id, true, false)
		if user_to_erase != null:
			user_to_erase.queue_free()

	for presence in match_event.joins:
		if presence.user_id != SessionManager.session.user_id:
			_handle_match_join_event(presence)


func _handle_match_join_event(presence):
	user_ids.append(presence.user_id)
	call_deferred("_add_networked_player_to_scene", presence.user_id, player_entry_node.position)


func _add_networked_player_to_scene(user_id: String, position: Vector2):
	var player_node = character_scene.instance()

	player_node.user_id = user_id
	player_node.position = position

	environment_items.add_child(player_node)
