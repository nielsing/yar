package args

import (
    "fmt"
    "errors"
    "strconv"

    arg "github.com/alexflint/go-arg"
)

const (
    minNoise = 0
    maxNoise = 9
)

func validateInt(argname, arg string, lower, upper int) (int, error) {
    num, err := strconv.Atoi(arg)
    if err != nil {
        return 0, fmt.Errorf("%s is not a number", arg)
    } else if num < lower {
        return 0, fmt.Errorf("%s can not be smaller than %d", argname, lower)
    } else if num > upper {
        return 0, fmt.Errorf("%s can not be larger than %d", argname, upper)
    }
    return num, nil
}

type CleanCmd struct {
    Cache string `arg:"positional"`
}

type GithubCmd struct {
    Org string `arg:"-o"`
    User string `arg:"-u"`
    Repo string `arg:"-r"`
    Forks bool `arg:"-f"`
    Depth int `arg:"-d"`
    NoBare bool `arg:"--no-bare"`
    IncludeMembers bool `arg:"--include-members"`
}

// TODO: Implemented!
type GitlabCmd struct {}

// TODO: Implemented!
type BitbucketCmd struct {}

// TODO: Explain s command
type Noise struct {
    Lower, Upper int
}

func (n *Noise) UnmarshalText(b []byte) error {
    s := string(b)
    switch length := len(s); length {
	case 3:
		lower, err1 := validateInt("noise level", string(s[0]), minNoise, maxNoise)
		upper, err2 := validateInt("noise level", string(s[2]), minNoise, maxNoise)
		if err1 != nil {
			return err1
		} else if err2 != nil {
			return err2
		}
		if upper < lower {
			return errors.New("noise level must be X-Y such that X <= Y")
		}
        n.Lower = lower
        n.Upper = upper
		return nil
	case 2:
		if string(s[0]) == "-" {
			num, err := validateInt("noise level", string(s[1]), minNoise, maxNoise)
			if err != nil {
				return err
			}
            n.Lower = 0
            n.Upper = num
			return nil
		}
		num, err := validateInt("noise level", string(s[0]), minNoise, maxNoise)
		if err != nil {
			return err
		}
        n.Lower = num
        n.Upper = 9
		return nil
	case 1:
		if s == "-" {
            n.Lower = minNoise
            n.Upper = maxNoise
			return nil
		}
		num, err := validateInt("noise level", string(s[0]), minNoise, maxNoise)
		if err != nil {
			return err
		}
        n.Lower = num
        n.Upper = num
		return nil
	default:
		return errors.New("noise argument must be in any of these forms:\nX, -X, X-, X-Y")
	}
}

type Args struct {
    // General flags
    Both bool `arg:"-b"`
    Save bool `arg:"-s"`
    Noise Noise `arg:"-n" default:"-5"`
    Config string `arg:"-C"`
    Context string `arg:"-c"`
    Entropy string `arg:"-e"`
    NoCache bool `arg:"--no-cache"`
    NoContext bool `arg:"--no-context"`

    // Commands
    Clean *CleanCmd `arg:"subcommand:clean" help:"Unimplemented!"`
    Github *GithubCmd `arg:"subcommand:github" help:"Unimplemented!"`
    Gitlab *GitlabCmd `arg:"subcommand:gitlab" help:"Unimplemented!"`
    Bitbucket *BitbucketCmd `arg:"subcommand:bitbucket" help:"Unimplemented!"`
}

func (Args) Version() string {
    return "yar 2.0.0"
}

func (Args) Description() string {
    return "Sail ye seas of internets for booty is to be found"
}

func ParseArgs() Args {
    parsedArgs := Args{}
    parser := arg.MustParse(&parsedArgs)
    // Start Validation of commands
    if parser.Subcommand() == nil {
        parser.Fail("Missing subcommand")
    }
    return parsedArgs
}
