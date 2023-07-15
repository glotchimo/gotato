# gotato

A Twitch chatbot that allows chatters to play hot potato for points that they
can then gamble (for fun and non-channel points because Twitch doesn't allow bots
to interact with channel points)!

## Usage

Building is of course the first priority, so assuming you have Go installed on
your computer, you can run `go build` in this repo and you'll get a binary.

The second thing you'll need is some environment variables for authenication:

- `GOTATO_CHANNEL`: The channel to connect to
- `GOTATO_USERNAME`: Twitch username to connect with
- `GOTATO_CLIENT_ID`: Twitch app client ID
- `GOTATO_CLIENT_SECRET`: Twitch app client secret

*You can get client ID/secret pairs by creating an application at
https://dev.twitch.tv/console*

The third part, which is optional and just changes the existing defaults, are
some more environment variables that tweak game options:

- `GOTATO_JOIN_DURATION`: Join phase duration in seconds (default: 30)
- `GOTATO_DURATION_MIN`: Minimum game length in seconds (default: 30)
- `GOTATO_DURATION_MAX`: Maximum game length in seconds (default: 60)
- `GOTATO_TIMEOUT_DURATION`: Loss timeout in seconds (default: 30)
- `GOTATO_REWARD_BASE`: Base points rewarded to winner (default: 100)
- `GOTATO_COOLDOWN_DURATION`: Cooldown between games in seconds (default: 120)

Once those are set (or not if you want the defaults), just run the binary:

	./gotato

This will open up a browser tab that will prompt you to authorize the app for
chat read/write/ban permissions (ban for timeout punishments). Once you
complete that flow, the goroutines will boot up and you'll see a message in
your chat from gotato.

**Note that the client user needs to type !enable in the chat to turn on
timeouts and points whispers.** This is necessary for the API client to 
receive the necessary user IDs to execute requests.

## Design

See `DESIGN.md` for design/implementation details.
