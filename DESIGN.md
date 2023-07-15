# Design

## Overview

The architecture of gotato is very simple. Fundamentally, it consists of two
goroutines, one that listens to Twitch IRC for commands, and another that
reacts to events sent through a channel by the aforementioned channel.

It is meant to be as simple as possible in its implementation so that it is
as simple as possible to scale by infrastructural means. A single instance
follows a single channel, and configuration happens via environment variables,
which tend to be easier to work with in serverless setups than files.

## Mechanics

### Configuration

Configuration fields are covered in the README, but it's worth noting that,
beyond essential authentication details, users can tweak all timing and reward/
punishment aspects of the game per their needs.

### Handling commands

We use `gempir/go-twitch-irc` to handle the tedious parts of IRC while we
listen for commands. The commands are as follows:

- `!gotato`: Starts the game
- `!join`: Adds the sender to the participants list
- `!bet/wager <number>`: Joins and registers a bet from the sender's point bank
- `!pass/toss`: Passes the potato to another participant
- `!points`: Displays the sender's point bank during wait/cooldown phases
- `!reset`: Resets the game to the join phase

### State, the main character

The `State` structure defined in `state.go` is the foundation of the game. It
consists of the following fields:

- `Timer`: A randomly-set (within configured bounds) game timer
- `Holder`: The user ID of the current potato holder
- `LastUpdate`: A timestamp updated on every cycle to assign holding scores
- `Participants`: A list of user IDs that have joined the game
- `Aliases`: A map of user IDs to usernames for chat messages
- `Scores`: A map of user IDs to scores
- `Bets`: A map of user IDs to bets
- `Reward`: The total reward pool value

Everything in state is ephemeral and reset after every game/cooldown period.
The only values that persist across games are scores/points, which are saved
in a local database.

### Points, scoring, and betting

A miminal Bolt DB interface is implemented in `points.go` for setting/getting
points in a local file. The setter is used in the end game to reward the winner
with the reward pool. The winner is determined by the user who held the potato
the longest and isn't holding it at the end.

A base reward (100 points by default) is given to the winner in addition to
whatever is bet by players during the join phase. Players can place a bet once
during the join phase and if it exceeds what is in their point bank, their
entire balance will be registered for betting.

At that point, bets are only *registered*, not actually *spent*. The balances
are updated immediately once the game phase starts and before the initial pass
is executed.

## Phases

The game goroutine is broken into four phases: wait, join, game, and cooldown.
Each handles events differently and contains the necessary functionality for
progressing state.

### Wait

The wait phase is the default, inactive phase. During this phase, the only
commands that can be issued are the starter `!gotato` command and the `!points`
command which sends a whisper to the user with their points bank value. This
phase runs indefinitely, and the `!gotato` command will move execution to the 
join phase.

### Join

In the join phase, chatters can register in the participant pool and place
bets with their earned points. Note that a bet is a join, but a join does not
require a bet. For example, a user can issue `!join` and offer no bet, or they
can issue `!bet 20` to join and register a bet of 20 points. These are wrapped
together for brevity and ease of use.

The join phase runs on a timer set by the `JOIN_DURATION` variable, which is
30 seconds by default. Once that timer (which is a native `time.Timer`) has run
out, execution will move to the game phase.

### Game

The game phase starts by applying registered bets to user point balances. Once
that is complete, an initial randomized pass is executed and the loop begins.
The only user able to progress state at this point is the one with the potato.
By issuing a `!pass` command, the potato is passed to another player and their
score is increased based on how many seconds they held the potato before
passing it.

The game phase also runs on a timer set randomly between the `GAME_DURATION_MIN`
and `GAME_DURATION_MAX` variables. Once that timer runs out, the user holding
the potato is timed out according to the `TIMEOUT_DURATION` variable, and the
user with the highest score (i.e. the one who held the potato the longest)
has the `state.Reward` value (i.e. the `REWARD_BASE` + all bets) added to their
point balance.

### Cooldown

Once those actions have completed, execution is moved to the cooldown phase,
the timer of which is set by the `COOLDOWN_DURATION` variable. The only command
that can be issued in this phase is the `!points` command. Once the timer is
depleted, execution moves back to the wait phase until another start command
is issued.
