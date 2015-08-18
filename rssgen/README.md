Some sample RSS, full and stripped down:

- thecrashcast.rss
- simple.rss
- minimal.rss

Minimally simple RSS generator.  Useful for downloading a set of
mp3 files in order to a podcast player.  No frills.

- rssgen

Usage:

./rssgen title url files... >r2.rss

Example:

cd ~/html/book
./rssgen MyAudioBook http://example.com/book *.mp3 >contents.rss
