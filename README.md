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
