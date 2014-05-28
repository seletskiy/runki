runki
=====

Manage Anki flashcards without a friction.


Typical usage
-------------

`go get github.com/seletskiy/runki` it first.

Make sure, that after `go get` command `runki` is available in your shell.

Create executable file with following contents (named `add-anki-word`):
```bash
#!/bin/bash

notify-send "$(echo $(xclip -o) \
    | runki --user=<USERNAME> --pass=<PASSWORD> --deck=<DECKNAME>)"
```

Then create shortcut in your window manager. i3, for example:
```
bindsym $mod+Escape exec add-anki-word
```

Reload configuration and it should work.

So, you just select the word, press shortcut ($mod+Escape in example), and get
word translation and transcription in popup. What will you see in popup will be
added to specified <DECKNAME> in your Anki account.

Every word will be checked on uniqueness prior to be added.
