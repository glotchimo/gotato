# gotato

A Twitch chat interaction that simulates a game of hot potato with chatters.

## Setup

Environment variables are used to configure things like the timer, rewards, timeouts, etc.
The following are available:

- `GOTATO_CHANNEL`: The channel to connect to
- `GOTATO_USERNAME`: Twitch username to connect with
- `GOTATO_PASSWORD`: Twitch password/token to connect with
- `GOTATO_TIMER_MIN`: Minimum game length in seconds
- `GOTATO_TIMER_MAX`: Maximim game length in seconds
- `GOTATO_TIMEOUT`: Loss timeout in seconds
- `GOTATO_REWARD`: Channel points rewarded to winner
- `GOTATO_COOLDOWN`: Cooldown between games in seconds

Once those are set (or not if you want the defaults), just run the binary:

	./gotato

## Design

### Interaction

We use Twitch's chat IRC to interact with chat. This is obviously necessary for interactions:

	glotchimo: !gotato
	gotato: The potato's been heated and is in @somebody's hands!
	somebody: !gotato
	gotato: The potato's been passed to @other!
	...
	gotato: Time's up! @other lost to potato LUL -5m
	gotato: @glotchimo held the potato the longest! +100 channel points PogChamp

By listening for the `!gotato` command we can progress state easily & atomically - whatever
IRC feeds us, we act on.

### Configuration

To give users the freedom to tweak the experience, we have a series of environment variables
that set the timer range, cooldown between games, timeout length, etc.

### State

Initially, state only requires two things. First, we need a random `timer` that gets set on a
first-cycle call (i.e. `!gotato` issued in chat) and gets decremented once every second or so.
Then we need a `holder` string that contains the ID of the chatter holding the potato, which
would be randomly set from the list of active viewers upon every subsequent call, potentially
with some additional logic to limit that list to recently active chatters.

An additional state element that would facilitate more competition and interesting interactions
would be a map of participants that gets populated and its values incremented whenever a given
chatter has the potato. For example, if somebody calls `!gotato` and I get the potato, I can hold it for as long as I want, and *then* pass it, and that duration would be counted. At the
end of the game, the winner would be the chatter who held onto the potato the longest.

### Conditions

In order to prevent multiple potatoes from being passed around, we need to make sure `holder` is
empty before starting a new game. That also requires that we clear `holder` after the timer has
run out and the timeout has been executed.

We also need to prevent chatters who are still active but have been timed out from
being included in the list of potato catchers. This could be achieved by watching for
"clear chat" messages from the IRC.

### Concurrency

We need three goroutines to progress state cleanly. First, we need a goroutine listening for
`!gotato` messages. When it receives one, it sends that event to a channel that gets consumed
by the state goroutine. Having a separate goroutine for state allows us to keep receiving calls
without breaking state up with the runtime conditions. Then we have a third that handles things
like timeouts and channel point rewards since those aren't that time-sensitive.

### Advanced Features

These are features that don't need to be in at first but would be cool to add in the future:

- **Channel point wager pool**: Chatters spend channel points to join, winner gets the pot.
- **Leaderboard**: A database is maintained with scores and points won over time.