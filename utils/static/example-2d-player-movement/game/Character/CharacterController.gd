extends KinematicBody2D

export var user_id: String

var next_movement_point: Vector2


func _enter_tree():
	next_movement_point = position
	var _mc = PlayerEvent.connect("movement", self, "_handle_move_event")


func _exit_tree():
	PlayerEvent.disconnect("movement", self, "_handle_move_event")


func _handle_move_event(msg, presence):
	if msg == null:
		return

	if presence.user_id != user_id:
		return

	var args = JSON.parse(msg).result
	_move_event(args)


func _move_event(args):
	next_movement_point = Vector2(args.x, args.y)
	var _ms = move_and_slide(next_movement_point)

	update()
