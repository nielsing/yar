package robber

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/akamensky/argparse"
)

const (
	maxInt = int(^uint(0) >> 1)
)

// Flags struct keeps a hold of all of the CLI arguments that were given.
type Flags struct {
	Org            *string
	User           *string
	Repo           *string
	Save           *string
	CleanUp        *string
	Noise          *int
	Config         *os.File
	Entropy        *bool
	Both           *bool
	NoContext      *bool
	Forks          *bool
	NoBare         *bool
	IncludeMembers *bool
	NoCache        *bool
	Context        *int
	CommitDepth    *int

	SavePresent    bool
	CleanUpPresent bool
}

type bound struct {
	lower int
	upper int
}

func validateInt(argname string, arg string, bound *bound) error {
	num, err := strconv.Atoi(arg)
	if err != nil || num < bound.lower {
		return fmt.Errorf("%s must be a non-negative integer", argname)
	}
	if num > bound.upper {
		return fmt.Errorf("%s must be a number between %d and %d", argname, bound.lower, bound.upper)
	}
	return nil
}

func validErr(err error) bool {
	return err.Error() != "not enough arguments for -s|--save" && err.Error() != "not enough arguments for --cleanup"
}

func flagPresent(shortHand string, name string) bool {
	for _, val := range os.Args {
		if val == shortHand || val == name {
			return true
		}
	}
	return false
}

// ParseFlags parses CLI arguments and returns them.
func ParseFlags() *Flags {
	parser := argparse.NewParser("yar", "Sail ye seas of git for booty is to be found")
	flags := &Flags{
		Org: parser.String("o", "org", &argparse.Options{
			Required: false,
			Help:     "Organization to plunder",
		}),

		User: parser.String("u", "user", &argparse.Options{
			Required: false,
			Help:     "User to plunder",
		}),

		Repo: parser.String("r", "repo", &argparse.Options{
			Required: false,
			Help:     "Repository to plunder",
		}),

		Context: parser.Int("c", "context", &argparse.Options{
			Required: false,
			Help:     "Show N number of lines for context",
			Default:  2,
			Validate: func(args []string) error {
				return validateInt("Context", args[0], &bound{0, 10})
			},
		}),

		Entropy: parser.Flag("e", "entropy", &argparse.Options{
			Required: false,
			Help:     "Search for secrets using entropy analysis",
			Default:  false,
		}),

		// Overrides entropy flag
		Both: parser.Flag("b", "both", &argparse.Options{
			Required: false,
			Help:     "Search by using both regex and entropy analysis. Overrides entropy flag",
			Default:  false,
		}),

		Forks: parser.Flag("f", "forks", &argparse.Options{
			Required: false,
			Help:     "Specifies whether forked repos are included or not",
			Default:  false,
		}),

		Noise: parser.Int("n", "noise", &argparse.Options{
			Required: false,
			Help:     "Specify the maximum noise level of findings to output",
			Default:  3,
			Validate: func(args []string) error {
				return validateInt("Noiselevel", args[0], &bound{1, 10})
			},
		}),

		CommitDepth: parser.Int("", "depth", &argparse.Options{
			Required: false,
			Help:     "Specify the depth limit of commits fetched when cloning",
			Default:  100000,
			Validate: func(args []string) error {
				return validateInt("Depth", args[0], &bound{0, maxInt})
			},
		}),

		NoBare: parser.Flag("", "nobare", &argparse.Options{
			Required: false,
			Help:     "Clone the whole repository",
			Default:  false,
		}),

		Config: parser.File("", "config", os.O_RDONLY, 0600, &argparse.Options{
			Required: false,
			Help:     "JSON file containing yar config",
			Default:  filepath.Join(GetGoPath(), "src", "github.com", "Furduhlutur", "yar", "config", "yarconfig.json"),
			Validate: func(args []string) error {
				filename := args[0]
				_, err := os.Stat(filename)
				if os.IsNotExist(err) {
					return errors.New("Rules file does not exist")
				} else if os.IsPermission(err) {
					return errors.New("You do not have permission to read the rules file")
				} else if err != nil {
					return errors.New("Unable to read rules file")
				}
				return nil
			},
		}),

		// Overrides context flag
		NoContext: parser.Flag("", "no-context", &argparse.Options{
			Required: false,
			Help:     "Only show the secret itself, similar to trufflehog's regex output. Overrides context flag",
			Default:  false,
		}),

		NoCache: parser.Flag("", "--nocache", &argparse.Options{
			Required: false,
			Help:     "Don't load from cache",
			Default:  false,
		}),

		IncludeMembers: parser.Flag("", "--includemembers", &argparse.Options{
			Required: false,
			Help:     "Include an organization's members for plunderin'",
			Default:  false,
		}),

		SavePresent:    flagPresent("-s", "--save"),
		CleanUpPresent: flagPresent("", "--cleanup"),

		// If cleanup is set, yar will ignore all other flags and only perform cleanup
		CleanUp: parser.String("", "cleanup", &argparse.Options{
			Required: false,
			Help:     "Remove specified cloned directory within yar cache folder. Leave blank to remove the cache folder completely",
			Default:  "",
		}),

		Save: parser.String("s", "save", &argparse.Options{
			Required: false,
			Help:     "Yar will save all findings to a specified file",
			Default:  "findings.json",
		}),
	}

	if err := parser.Parse(os.Args); err != nil && validErr(err) {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}
	validateFlags(flags, parser)
	return flags
}

func validateFlags(flags *Flags, parser *argparse.Parser) {
	if *flags.User == "" && *flags.Repo == "" && *flags.Org == "" && !flags.CleanUpPresent {
		fmt.Print(parser.Usage("Must give atleast one of org/user/repo"))
		os.Exit(1)
	}
	if *flags.Save == "" {
		*flags.Save = "findings.json"
	}
}
