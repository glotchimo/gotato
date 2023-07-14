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

## Structure

