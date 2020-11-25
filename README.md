# Yet Another Robber

<p align="center">
  <img src="https://raw.githubusercontent.com/nielsing/yar/master/images/yargopher3.png" alt="Yar the pirate gopher"/>
</p>

Yar searches for secrets either by regex, entropy or both, the choice is yours! Inspired by other secret grabbers.

## Installation
1. You can install Yar by running `go get github.com/nielsing/yar`
2. Or you can download the latest release of Yar for your operating system [here](https://github.com/nielsing/yar/releases).

## Usage
TODO: Add usage, the Yar CLI is going to be completely different in v2.

### Have your own predefined rules?
Rules are stored in a JSON file with the following format:
```
{
    "Rules": [
        {
            "Reason": "The reason for the match",
            "Rule": "The regex rule",
            "Noise": 3
        },
        {
            "Reason": "Super secret token",
            "Rule": "^Token: .*$",
            "Noise": 2
        }
    ]
    "FileBlacklist": [
        "Regex rule here"
        "^.*\\.lock"
    ]
}
```

You can then load your own rule set with the following command:
```
yar --rules PATH_TO_JSON_FILE
```

If you already have a truffleHog config and want to port it over to a yar config there is a script in the config folder that does it for you.
Simply run `python3 trufflestoconfig.py PATH_TO_TRUFFLEHOG_CONFIG` and the script will give you a file named `yarconfig.json`.

### Don't like regex?
```
yar --entropy
```

### Want the best of both worlds?
```
yar --both
```

### Don't like the default colors and want to add your own color settings?
It is possible to customize the colors of the output for Yar through environment variables.
The possible colors to choose from are the following:
```
black
blue
cyan
green
magenta
red
white
yellow
hiBlack
hiBlue
hiCyan
hiGreen
hiMagenta
hiRed
hiWhite
hiYellow
```
Each color can then be suffixed with `bold`, i.e. `blue bold` to make the letters bold.

This is done through the following env variables:
```
YAR_COLOR_VERBOSE -> Color of verbose lines.
YAR_COLOR_SECRET  -> Color of the highlighted secret.
YAR_COLOR_INFO    -> Color of info, that is, simple strings that tell you something.
YAR_COLOR_DATA    -> Color of data, i.e. commit message, reason, etc.
YAR_COLOR_SUCC    -> Color of succesful messages.
YAR_COLOR_WARN    -> Color of warnings.
YAR_COLOR_FAIL    -> Color of fatal warnings.
```
Like so `export YAR_COLOR_SECRET="hiRed bold"`.

## Extra Knowledge
There are some design decisions which might be good to know about. Yar saves all cloned github repos
in a folder named yar within the temp directory. Yar then tries to load github repos from this cache
by default, if you don't want to load from cache then you can add the `--no-cache` flag.

Yar also clones bare repos by default, if you want to get all files within a repo and not just the 
metadata then you can add the `--no-bare` flag.

If you want to remove repos from cache then you can use the `--cleanup` flag. This flag 
either removes the whole cache if no folder was specified or just removes the specified folder. The
folder structure within the cache folder is like so:
```
/yar
|--- /User1
|  |--- /Repo1
|  |--- /Repo2
|
|--- /User2
|  |--- /Repo1
|  |--- /Repo2

```
So you can run `--cleanup User1` to remove the cache of User1 or `--cleanup User1/Repo1` to clean up
Repo1 of User1. You can think of the flag as a wrapper around `rm -r /tmp/yar/{USER_INPUT}`.

Finally yar goes 10000 commits deep by default and goes through them in order of time
(oldest to newest). This depth is configurable so if you ever want to cover more or fewer commits
simply add the `--depth` flag with the depth you want.
