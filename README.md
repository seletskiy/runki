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


Batch usage
-----------

Just pipe words from some source to `runki`, one word/phrase per line. Like
this:

```bash
cat words-list.txt | runki
```


Kindle usage
------------

It is possible to learn foreign words directly from Kindle. When you encounter
unknown word, just tap on it, then select "Highlight". All this highlights will
be saved to the file "/documents/My Clippings.txt" in the Kindle FS.

Afterwards this file can be filtered to extract single word highlights (see
example script here: https://github.com/seletskiy/dotfiles/blob/92a0e8fdb533106d700dc4c7415f4d41b232edff/bin/kindle-filter-words).

So, all highlighted words can be easily added to Anki like this:

```bash
kindle-filter-words "/mnt/documents/My Clippings.txt" | tac | runki --cut 3
```

Note, that `tac` command will reverse incoming list so most recent items will
came first. `--cut 3` flag tells runki to exit if it detect at least `3` words
that are already added to Anki. `tac` and `--cut` combination allows you to
sync words from Kindle in seamless way just running one single command again
and again.

Kindle + udev
-------------

Ok, let's go deeper. I want to sync words that I highlight on the kindle
automatically when I connect kindle to the computer.

So, we need to create udev rule like this:

```
ACTION=="add", SUBSYSTEM=="block", ENV{DEVTYPE}=="partition", ENV{ID_VENDOR_ID}=="1949", RUN+="/usr/bin/su <USERNAME> -lc 'DISPLAY=:0 kindle-to-anki $env{DEVNAME}'"
```

So, after kindle is connected, `kindle-to-anki` program will have to be runned.
It will add new words directly to Anki and show nice notification about how
many new words has been added to.

Example of this `kindle-to-anki` program can be found there: https://github.com/seletskiy/dotfiles/blob/1c9da6d347cc658c9d6d177a61ef94423a3c36d4/bin/kindle-to-anki
