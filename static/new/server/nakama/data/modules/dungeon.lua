local nk = require("nakama")
local farm = require("farm_grid")
local DM = require("dungeon_manager")

local M = {}

local empty_room_ticks = 0
local max_empty_room_ticks = 150

function M.match_init(context, setupstate)
  local gamestate = {
    presences = {},
    plot_map = {},
    avatar_map = {},
    movement_targets = {},
    positions = {}
  }

  local tickrate = 15 -- per sec

  nk.logger_info("Dungeon started")

  return gamestate, tickrate, setupstate.label
end

function M.match_join_attempt(context, dispatcher, tick, state, presence, metadata)
  local acceptuser = true

  local presence_count = get_presence_count(state)

  if (presence_count >= DM.max_dungeon_size) then
    acceptuser = false
  else
    set_plot_id(state, presence)
    state.presences[presence.session_id] = presence
    state.avatar_map[presence.user_id] = metadata.avatarName
    state.positions[presence.user_id] = {
      x = metadata.x,
      y = metadata.y
    }
  end

  return state, acceptuser
end

function M.match_join(context, dispatcher, tick, state, presences)
  --for _, presence in pairs(presences) do
  --end

  broadcast_state(dispatcher, state)

  return state
end

function M.match_leave(context, dispatcher, tick, state, presences)
  for _, presence in pairs(presences) do
    state.presences[presence.session_id] = nil
    state.plot_map[presence.user_id] = nil
    state.avatar_map[presence.user_id] = nil
    state.movement_targets[presence.user_id] = nil
    state.positions[presence.user_id] = nil
  end

  broadcast_state(dispatcher, state)

  return state
end

function broadcast_state(dispatcher, state)
  dispatcher.broadcast_message(
    6,
    nk.json_encode(
      {
        plotMap = state.plot_map,
        avatarMap = state.avatar_map,
        movementTargets = state.movement_targets,
        positions = state.positions
      }
    ),
    nil
  )
end

function set_plot_id(state, presence)
  local plot_id = 0
  local taken_plots = {}

  for user, plot in pairs(state.plot_map) do
    taken_plots[plot] = plot
  end

  for i = 0, DM.max_dungeon_size, 1 do
    if (taken_plots[i] == nil) then
      plot_id = i
      break
    end
  end

  state.plot_map[presence.user_id] = plot_id
end

function get_presence_count(state)
  local presence_count = 0
  for _, presence in pairs(state.presences) do
    -- nk.logger_info(("Presence %s named %s"):format(presence.user_id, presence.username))
    presence_count = presence_count + 1
  end

  return presence_count
end

function should_room_be_open(state)
  presence_count = get_presence_count(state)

  if presence_count == 0 then
    -- nk.logger_info(("match room empty, tick count: %s"):format(empty_room_ticks))
    empty_room_ticks = empty_room_ticks + 1

    if empty_room_ticks >= max_empty_room_ticks then
      nk.logger_info("match room empty for too long, terminating...")
      return false
    end
  elseif empty_room_ticks > 0 then
    -- nk.logger_info("match room not empty, zeroing out empty_room_ticks")
    empty_room_ticks = 0
  end

  return true
end

function M.match_loop(context, dispatcher, tick, state, messages)
  if should_room_be_open(state) == false then
    return nil
  end

  for _, message in ipairs(messages) do
    dispatcher.broadcast_message(message.op_code, message.data, nil, message.sender)
  end

  for _, message in ipairs(messages) do
    state = perform_message_validation(message, dispatcher, state)
  end

  return state
end

function perform_message_validation(message, dispatcher, state)
  if message.op_code == 1 then -- 1 == farming event
    if farm.handle_match_farm_action(message) == false then
      dispatcher.broadcast_message(
        4, -- data reset game event enum ID in client
        nk.json_encode(
          {
            offendingUUID = message.sender.user_id
          }
        ),
        nil
      )
    end
  elseif message.op_code == 0 then -- 0 == movement event
    local data = nk.json_decode(message.data)

    if (data.ping == "1") then -- 1 == "ping"
      state.positions[message.sender.user_id] = {
        x = data.x,
        y = data.y
      }
    else
      state.movement_targets[message.sender.user_id] = {
        x = data.x,
        y = data.y
      }
    end
  end

  return state
end

function M.match_terminate(context, dispatcher, tick, state, grace_seconds)
  nk.logger_info("terminating match")
  local message = "Server shutting down match in " .. grace_seconds .. " seconds"
  -- TODO: implement message and opcode for server messaging
  -- dispatcher.broadcast_message(2, message)
  return nil
end

return M
