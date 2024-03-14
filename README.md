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

> Version 3.0.1

> Made by Nathan Hedge @ https://almostd.one/

----------------------------------

__A Discord bot for text-based commands and responses.__

This may sound like any other bot at first, but this is **much** more than basic text.

----------------------------------

<br>
<br>

__LIST OF COMMANDS__
---
### `@TRAS help` or `/help`
> Summons a help list.
> Slash command doesn't require DMs
 
### `@TRAS about` or `/about`
> Gives information about the bot. 
> Add "NoCB" for devices that don't support links with command blocks.
> Slash command doesn't require DMs
 
### `@TRAS oof`
> Mega OOF
 
### `@TRAS f`
> Mega F
 
### `@TRAS pi`
> First 1 million digits of Pi

### `@TRAS big`
> Make a larger verison of word/text made of the letter. 
> Starts getting wonky with emojis. Becomes file over 400 characters. 
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
> Generate ascii art. Over 400 characters returns in a file.
> 
> Format: `@TRAS ascii art [font/getFonts] [text]`
 
### `@TRAS commands`
> View and manage custom server commands, managing requires administrator perms.
> Custom commands feature may require TRAS Deluxe in the future (TBD, currently not a thing).
> Schedule feature not currently implemented.
> 
> Format: `@TRAS commands [manage/view] [(manage)...]`
> Format (manage): `@TRAS commands manage [set/delete/schedule] [(set/delete)trigger//(schedule)time of day (hh:mm:ss)] [(set/schedule)reply]`

### `@TRAS rank` 
> Shows your rank, lets your reset your rank, and allows you to roll dice for a new rank if it's enabled. 
> Admins get other commands as well. Dice rolling disabled by default.
> 
> Format: `@TRAS rank [info/checkDice/dice] [(info)-real]`
> Format (admin): `@TRAS rank [set/reset/toggleDice] [(set/reset)user] [value]`

### `@TRAS set nickname`
> Set the bot's Nickname on the server. Reset with "{RESET}". 
> Requires "Manage Messages" or "Change Nicknames".
> 
> Format: `@TRAS set nickname [nickname/{RESET}]`
> Reset alternative: `@TRAS reset nickname`
 
### `/mydata`
> Download and delete data you own from the TRAS database. For data protection compliance. Slash only.
> Server data management requires administrator.
>
> Format: `/mydata target:[me/server] action:[download/delete]`
 
### `@TRAS speak`
> Generate a sentence, plus toggle and get the status of random generated messages. 
> Toggling requires "Manage Messages" permissions. Random messages are off by default.
> Frequent generation may require TRAS Deluxe in the future (TBD, currently not a thing).
> 
> Format: `@TRAS speak [generate/randomspeak] [(randomspeak)on/off/status]`
 
### `@TRAS combinations`
> Sends file with all possible combinations of the units you have selected and given.
> 
> Format: `@TRAS combinations [words/characters] [items]`
 
### `@TRAS ping`
> Check if the bot is alive. 
> Add 'info' or 'information' for latency data.

### Default response
> Responds "What's up?"
 
### Generated messages
> Fully generated messages *(not an AI so they're completely nonsensical)* can be toggled to randomly
> say them in response to  user messages. Random messages will not reply to commands.
> Toggle random messages with `@TRAS speak randomspeak [on/off]`

<br>
<br>

----------------------------------

__GENERAL DETAILS__
---
> All data is stored in Postgres.

> Large items are stored in files.

> Controversial features are toggleable.

> Custom commands are activated like normal commands. / Ex: `@TRAS [trigger]`

> ALL commands (except for some modifiers) work with aNY CapItaLIzATIoN.

> All commands should work in DMs.

> Ranking levels are base 2 logarithms of your progress.

> Spaces can be included in individual arguments by using escapes / Ex: `This\ is\ one\ argument, but\ this\ is\ the\ second`

> Escaped slashes can still be used at the end of individual arguments / Ex: `This\ is\ one\ argument\\ but\ this\\\ is\ the\ second` 

----------------------------------

<br>
<br>

![TRAS logo](src/traslogo.png)

<br>
<br>

## Copyright (c) 2024 Nathan B. (10Nates)