extends "res://Character/CharacterController.gd"

# icon from https://www.reddit.com/r/godot/comments/icagss/i_made_a_claylike_3d_desktop_icon_for_godot_ico/

var map_constraints: Dictionary = {}

const MOVE_POSITION_FROM_BUTTON: int = 32
const MOVE_DELAY_MINIMUM: float = .25
const SPEED: int = 200

onready var camera: Camera2D = $Camera2D


func _physics_process(_delta: float):
	var velocity = Vector2()

	if Input.is_action_pressed("move_right"):
		velocity.x += 1
	if Input.is_action_pressed("move_left"):
		velocity.x -= 1
	if Input.is_action_pressed("move_down"):
		velocity.y += 1
	if Input.is_action_pressed("move_up"):
		velocity.y -= 1

	var target = velocity.normalized() * SPEED

	if velocity != Vector2.ZERO:
		PlayerEvent.movement({"x": target.x, "y": target.y})


func restrict_camera_to_tile_map(map: TileMap):
	var limits: Rect2 = map.get_used_rect()
	var cell_size: Vector2 = map.cell_size

	camera.limit_left = int(round(limits.position.x * cell_size.x))
	camera.limit_right = int(round(limits.end.x * cell_size.x))
	camera.limit_top = int(round(limits.position.y * cell_size.y))
	camera.limit_bottom = int(round(limits.end.y * cell_size.y))

	map_constraints["left"] = camera.limit_left
	map_constraints["right"] = camera.limit_right
	map_constraints["top"] = camera.limit_top
	map_constraints["bottom"] = camera.limit_bottom
