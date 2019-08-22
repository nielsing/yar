package robber

import (
	"errors"
	"fmt"
	"github.com/akamensky/argparse"
	"os"
	"path/filepath"
	"strconv"
)

const (
	maxInt = int(^uint(0) >> 1)
)

// Flags struct keeps a hold of all of the CLI arguments that were given.
type Flags struct {
	Org         *string
	User        *string
	Repo        *string
	Rules       *os.File
	Context     *int
	Entropy     *bool
	Both        *bool
	NoContext   *bool
	Forks       *bool
	Verbose     *bool
	CleanUp     *bool
	CommitDepth *int
	Noise       *int
}

type bound struct {
	upper int
	lower int
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

// ParseFlags parses CLI arguments and returns them.
func ParseFlags() *Flags {
	parser := argparse.NewParser("yar", "Sail ye seas of git for booty is to be found")
	flags := &Flags{
		Org: parser.String("o", "org", &argparse.Options{
			Required: false,
			Help:     "Organization to plunder",
			Default:  "",
		}),

		User: parser.String("u", "user", &argparse.Options{
			Required: false,
			Help:     "User to plunder",
			Default:  "",
		}),

		Repo: parser.String("r", "repo", &argparse.Options{
			Required: false,
			Help:     "Repository to plunder",
			Default:  "",
		}),

		Rules: parser.File("", "rules", os.O_RDONLY, 0600, &argparse.Options{
			Required: false,
			Help:     "JSON file containing regex rulesets",
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

		// Overrides context flag
		NoContext: parser.Flag("", "no-context", &argparse.Options{
			Required: false,
			Help:     "Only show the secret itself, similar to trufflehog's regex output. Overrides context flag",
			Default:  false,
		}),

		Forks: parser.Flag("f", "forks", &argparse.Options{
			Required: false,
			Help:     "Specifies whether forked repos are included or not",
			Default:  false,
		}),

		Verbose: parser.Flag("v", "verbose", &argparse.Options{
			Required: false,
			Default:  false,
		}),

		// If cleanup is set, yar will ignore all other flags and only perform cleanup
		CleanUp: parser.Flag("", "cleanup", &argparse.Options{
			Required: false,
			Help:     "Remove all temporary directories used for cloning",
			Default:  false,
		}),

		CommitDepth: parser.Int("", "depth", &argparse.Options{
			Required: false,
			Help:     "Specify the depth limit of commits fetched when cloning",
			Default:  100000,
			Validate: func(args []string) error {
				return validateInt("Depth", args[0], &bound{0, maxInt})
			},
		}),

		Noise: parser.Int("n", "noise", &argparse.Options{
			Required: false,
			Help:     "Specify the maximum noise level of findings to output",
			Default:  3,
			Validate: func(args []string) error {
				return validateInt("Noiselevel", args[0], &bound{0, 5})
			},
		}),
	}

	if err := parser.Parse(os.Args); err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}
	validateFlags(flags, parser)
	return flags
}

func validateFlags(flags *Flags, parser *argparse.Parser) {
	if *flags.User == "" && *flags.Repo == "" && *flags.Org == "" && !*flags.CleanUp {
		fmt.Print(parser.Usage("Must give atleast one of org/user/repo"))
		os.Exit(1)
	}
}
