# (Y)et (A)nother (R)obber: Sail ye seas of git for booty is to be found
Sail ho! Yar is a tool for plunderin' organizations, users and/or repositories...

In all seriousness though, yar is an OSINT tool for reconnaissance of repositories/users/organizations on Github. Yar clones repositories of users/organizations given to it
and goes through the whole commit history in order of commit time, in search for secrets/tokens/passwords, essentially anything that shouldn't be there. Whenever yar finds a secret,
it will print it out for you to further assess.

Yar searches either by regex, entropy or both, the choice is yours. This idea is inspired by the infamous [truffleHog](https://github.com/dxa4481/truffleHog) tool and the code used
for entropy search is in fact borrowed from their repository, which in turn was borrowed from [this blog post](http://blog.dkbza.org/2007/05/scanning-data-for-entropy-anomalies.html).

You can think of yar as a bigger and better truffleHog, it does everything that truffleHog does and more! Yar also does it faster and even finds more secrets (yes I know bold statement with no data to show it, I'll hopefully have some data to show it in the near future).

Right now yar is in early development (v0.8.0), however it is still usable and I've found it to be performing better than truffleHog, both in performance and the number of secrets it finds
If you want to know more regarding the development of yar, consult the Future Plans section.

## Installation
To install this you can simply run the following commands.

```
git clone [link here]
cd [folder here]
go install
```

## Usage
Want to search for secrets in an organization?
```
yar -o orgname
```

Want to search for secrets for a user?
```
yar -u username
```

Want to search for secrets in a single repository?
```
yar -r reponame
```

Want to search for secrets in an organization, for a user and a repository?
```
yar -o orgname -u username -r reponame
```

Have your own predefined rules?
Rules are stored in a JSON file with the following format. If you are familiar with truffleHog you already know how to use this.
```
{
  "Reason": "regex-rule"
}
```

```
yar -u username --rules PATH_TO_JSON_FILE
```

Don't like regex?
```
yar -u username --entropy
```

Want the best of both worlds?
```
yar -u username --both
```


## Help
```
usage: yar [-h|--help] [-o|--org "<value>"] [-u|--user "<value>"] [-r|--repo
           "<value>"] [--rules <file>] [-c|--context <integer>] [-e|--entropy]
           [-b|--both] [-n|--no-context] [-d|--debug]

           Sail ye seas of git for booty is to be found

Arguments:

  -h  --help        Print help information
  -o  --org         Organization to plunder. Default: 
  -u  --user        User to plunder. Default: 
  -r  --repo        Repository to plunder. Default: 
      --rules       JSON file containing regex rulesets. Default: rules.json
  -c  --context     Show N number of lines for context. Default: 2
  -e  --entropy     Search for secrets using entropy analysis. Default: false
  -b  --both        Search by using both regex and entropy analysis. Default:
                    false
  -n  --no-context  Only show the secret itself, similar to trufflehog's regex
                    output. Default: false
  -d  --debug  
```

## Future Plans
Yar is in active development and there are big plans for the near future.

### v0.9.0
Features:
+ Filename of the file that contains the secret to the finding output.
+ Environment variables for color customization of output.
+ Ensure the project is go gettable.
+ Add cleanup flag 

Extras:
+ Add contribution guidelines so strangers can take part in the project.

### v1.0.0
Features:
+ Make yar concurrent and hopefully parallel.
+ Redesign the findings output.
+ Option to save the results in some specified format (JSON maybe?).

Extras:
+ Benchmarking.
+ Testing.

