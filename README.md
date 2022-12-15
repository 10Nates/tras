# Text Response Automation System
```                                     
  _|_|_|_|_|  _|_|_|_|       _|       _|_|_|_|
 /////_|///  /_|////_|      _|_|     _|////// 
     /_|     /_|   /_|     _|//_|   /_|       
     /_|     /_|_|_|_|    _|  //_|  /_|_|_|_| 
     /_|     /_|///_|    _|_|_|_|_| ////////_|
     /_|     /_|  //_|  /_|//////_|        /_|
     /_|     /_|   //_| /_|     /_|  _|_|_|_| 
     //      //     //  //      //  ////////                                          
```
----------------------------------

> Version 3.0.0

> Made by Nathan Hedge @ https://almostd.one/

----------------------------------

__A Discord bot for text-based commands and responses.__

This may sound like any other bot at first, but this is **much** more than basic text.

----------------------------------

<br>
<br>

__LIST OF COMMANDS__
---
### `@TRAS help`
> Summons a help list.
 
### `@TRAS about`
> Gives information about the bot.
 
### `@TRAS oof`
> Mega OOF
 
### `@TRAS f`
> Mega F
 
### `@TRAS pi`
> First 1 million digits of Pi

### `@TRAS big`
> Make a larger verison of word/text made of the letter. 
> Starts getting wonky with emojis. Becomes file over 520 characters. 
> You can enable thin letters with -t or --thin.
> 
> Format: `@TRAS big (-t/--thin) [letter] [text]`
 
### `@TRAS jumble`
> Jumbles the words in a sentence so it's confusing to read.
> 
> Format: `@TRAS jumble [text]`

### `@TRAS emojify`
> Turn all characters into emojis.
> 
> Format: `@TRAS emojify [text]`
 
### `@TRAS flagify`
> Turn all letters into regional indicators, creating flags.
> 
> Format: `@TRAS flagify [text]`

### `@TRAS superscript`
> Turn all numbers and letters plus a few math symbols into superscript. 
> Some letters are always lowercase or replaced with something similar due to Unicode limitations.
> 
> Format: `@TRAS superscript [text]`
 
### `@TRAS unicodify`
> Turn all numbers and letters into a non-Latin equivilant.
> 
> Format: `@TRAS unicodify [text]`
 
### `@TRAS bold`
> Bolds all Latin letters and numbers using Unicode.
> 
> Format: `@TRAS bold [text]`
 
### `@TRAS replace`
> Replaces every appearance of a set item with a set replacement.
> 
> Format: `@TRAS replace [item] [replacement] [text]`
 
### `@TRAS overcomplicate`
> Replaces all words with synonyms of the word.
> 
> Format: `@TRAS overcomplicate [text]`
 
### `@TRAS word info`
> Get the definition or Part-of-Speech of a word.
> 
> Format: `@TRAS word info [definition/pos] [word]`
 
### `@TRAS ascii art`
> Generate ascii art. Over 15 characters responds with a file.
> 
> Format: `@TRAS ascii art [text/{font:[Font (use "\ " as space)]}/{getFonts}] [(font)text]`
 
### `@TRAS commands`
> View and manage custom server commands, managing requires "Manage Messages" perms.
> Custom commands feature requires TRAS Deluxe TBS.
> 
> Format: `@TRAS commands [manage/view] [(manage)...]`
> Format (manage): `@TRAS commands manage [set/delete/schedule] [(set/delete)trigger//(schedule)time of day (hh:mm:ss)] [(set/schedule)reply]`

### `@TRAS rank` 
> Shows your rank, lets your reset your rank, and allows you to roll dice for a new rank if it's enabled. 
> Admins get other commands as well. Dice rolling disabled by default.
> 
> Format: `@TRAS rank [info/checkDice/dice/reset] [user(4resetORset,admin)/amount(4set,admin)/-real(4info)] [amount(4set,admin)]`
 
### `@TRAS set nickname`
> Set the bot's Nickname on the server. Reset with "{RESET}". 
> Requires "Manage Messages" or "Change Nicknames".
> 
> Format: `@TRAS set nickname [nickname/{RESET}]`
 
### `@TRAS speak`
> Generate a sentence, plus toggle and get the status of random generated messages. 
> Toggling requires "Manage Messages" permissions. Random messages are off by default.
> 
> Format: `@TRAS speak [generate/random speak] [(random speak)on/off/status//(generate)starter]`
 
### `@TRAS combinations`
> Sends file with all possible combinations of the units you have selected and given.
> 
> Format: `@TRAS combinations [words/characters] [items]`
 
### Default response
> Responds "What's up?"
 
### Generated messages
> Random generated messages that can be toggled to randomly to be said in response to messages. 
> Random messages will not reply to commands.
> Toggle random messages with `@TRAS speak random speak [on/off]`

<br>
<br>

----------------------------------

__GENERAL DETAILS__
---
> All data is stored in MongoDB.

> Large items are stored in files.

> Controversial features are toggleable.

> Custom commands are activated like normal commands. / Ex: `@TRAS [trigger]`

> ALL commands (except for some modifiers) work with aNY CapItaLIzATIoN.

> All commands work in DMs. (There are some exeptions due to Discord's limits, or having no purpose)

> Ranking levels are base 2 logirithms of your progress.

> Spaces can be included in individual arguments by using backspaces / Ex: `This\ is\ one\ argument, but\ this\ is\ the\ second`

----------------------------------

<br>
<br>

## Copyright (c) 2022 Nathan B. (10Nates)