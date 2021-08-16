# album-streamer

An exersise in TDD and test of duplex communication via HTTP/2.

A centralized mp3 library which can stream an "Album" to a "Player" in a chromecast-like way.

Start a Streamer passing it a library of mp3 files with an Artist -> Album -> Song.mp3 folder hierarchy.
Then, choose which album to stream and which player to stream it to.

See Integration test for example usage.

I plan to use this to centrally host my mp3 library and stream albums to raspberry pi based amplifiers in different parts
of my home.

There is no seek, there is no skip, there is no pause. You pick an album and listen to it all the way through (or stop in the middle, I suppose).

I may end up implementing pause, but part of the intention of this project was to get more of a physical LP like listening experience with my digital collection.

I'd also like to implement album "stacks". This would just be a playlist of albums to play in sequence.

The duplex communication is currently pointless 😋.

