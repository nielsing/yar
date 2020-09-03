package robber

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	arg "github.com/alexflint/go-arg"
)

const (
	maxInt   = int(^uint(0) >> 1)
	minNoise = 0
	maxNoise = 9
)

// TODO: Comment
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

// TODO: Comment
func validateFile(argname, filename string) error {
	if filename == "" {
		return nil
	}
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return errors.New("Config file does not exist")
	} else if os.IsPermission(err) {
		return errors.New("You do not have permission to read the config file")
	} else if err != nil {
		return errors.New("Unable to read the config file")
	}
	return nil
}

// ClearCmd handles the 'clear' subcommand which allows a user to delete a specific directory within
// the designated yar cache folder.
type ClearCmd struct {
	Cache string `arg:"positional" help:"Remove specified directory within yar cache folder. Leave blank to remove the cache folder completely" default:""`
}

type GitCmd struct {
	Repo   []string `arg:"-r" help:"repository to plunder"`
	Depth  int      `arg:"-d" help:"Specify the depth limit of commits fetched when cloning" default:"10000"`
	NoBare bool     `arg:"--no-bare" help:"Clone the whole repository"`
}

type GithubCmd struct {
	Org            []string `arg:"-o" help:"organization to plunder"`
	User           []string `arg:"-u" help:"user to plunder"`
	Repo           []string `arg:"-r" help:"repository to plunder"`
	Forks          bool     `arg:"-f" help:"specifies whether forked repos are included or not"`
	Depth          int      `arg:"-d" help:"Specify the depth limit of commits fetched when cloning" default:"10000"`
	NoBare         bool     `arg:"--no-bare" help:"Clone the whole repository"`
	IncludeMembers bool     `arg:"--include-members" help:"Include an organization's members for plunderin'"`
}

// TODO: Implemented!
type GitlabCmd struct{}

// TODO: Implemented!
type BitbucketCmd struct{}

// TODO: Explain Noise command
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
	Both      bool   `arg:"-b" help:"search by using both regex and entropy analysis"`
	Save      string `arg:"-s" help:"yar will save all findings to a specified file" default:"findings.json" placeholder:"FILE"`
	Noise     Noise  `arg:"-n" help:"specify the range of the noise for rules. Can be specified as up to (and including) a certain value (-4), from a certain value (5-), between two values (3-5), just a single value (4) or the whole range (-)" default:"-5" placeholder:"X-Y"`
	Context   int    `arg:"-c" help:"show N number of lines for context" default:"2" placeholder:"N"`
	Entropy   bool   `arg:"-e" help:"search for secrets using entropy analysis"`
	NoCache   bool   `arg:"--no-cache" help:"don't load from cache"`
	NoContext bool   `arg:"--no-context" help:"only show the secret itself. Overrides context flag"`

	// Subcommands
	Clear     *ClearCmd     `arg:"subcommand:clear" help:"unimplemented!"`
	Git       *GitCmd       `arg:"subcommand:git" help:"unimplemented!"`
	Github    *GithubCmd    `arg:"subcommand:github" help:"unimplemented!"`
	Gitlab    *GitlabCmd    `arg:"subcommand:gitlab" help:"unimplemented!"`
	Bitbucket *BitbucketCmd `arg:"subcommand:bitbucket" help:"unimplemented!"`

	// Environment commands
	Config string `arg:"env" help:"JSON file containing yar config" placeholder:"FILE"`
}

func (Args) Version() string {
	return "yar 2.0.0"
}

func (Args) Description() string {
	return "Sail ye seas of internets for booty is to be found"
}

func parseArgs() *Args {
	parsedArgs := &Args{}
	parser := arg.MustParse(parsedArgs)

	// Start Validation of commands
	if parser.Subcommand() == nil {
		parser.Fail("Missing subcommand")
	}

	_, err := validateInt("context", strconv.Itoa(parsedArgs.Context), 0, maxInt)
	if err != nil {
		parser.Fail(err.Error())
	}

	err = validateFile("config", parsedArgs.Config)
	if err != nil {
		parser.Fail(err.Error())
	}

	return parsedArgs
}
