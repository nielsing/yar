package robber

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"regexp"
)

const regexErrorMessage = `Invalid regex rule in config file!
    Key: %s
    Rule: %s
    Error: %s
Read here for more information: https://golang.org/pkg/regexp/syntax/
`
const numberOfDefaultRules = 100

var allRules = []*Rule{
	{
		"AWS Access Key ID Value",
		regexp.MustCompile("(A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16}"),
		1,
	},
	{
		"AWS Access Key ID Value Base64",
		regexp.MustCompile("(QTNU|QUtJQ|QUdQQ|QUlEQ|QVJPQ|QUlQQ|QU5QQ|QU5WQ|QVNJQ)[%a-zA-Z0-9+/]{20,24}={0,2}"),
		3,
	},
	{
		"AWS Account ID",
		regexp.MustCompile("((\\\"|'|`)?((?i)aws)?_?((?i)account)_?((?i)id)?(\\\"|'|`)?\\\\s{0,50}(:|=>|=)\\\\s{0,50}(\\\"|'|`)?[0-9]{4}-?[0-9]{4}-?[0-9]{4}(\\\"|'|`)?)"),
		3,
	},
	{
		"AWS Secret Access Key",
		regexp.MustCompile("((\\\"|'|`)?((?i)aws)?_?((?i)secret)_?((?i)access)?_?((?i)key)?_?((?i)id)?(\\\"|'|`)?\\\\s{0,50}(:|=>|=)\\\\s{0,50}(\\\"|'|`)?[A-Za-z0-9/+=]{40}(\\\"|'|`)?)"),
		3,
	},
	{
		"AWS Session Token",
		regexp.MustCompile("((\\\"|'|`)?((?i)aws)?_?((?i)session)?_?((?i)token)?(\\\"|'|`)?\\\\s{0,50}(:|=>|=)\\\\s{0,50}(\\\"|'|`)?[A-Za-z0-9/+=]{16,}(\\\"|'|`)?)"),
		3,
	},
	{
		"Amazon MWS Auth Token",
		regexp.MustCompile("amzn\\.mws\\.[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}"),
		3,
	},
	{
		"Artifactory",
		regexp.MustCompile("(?i)artifactory.{0,50}(\\\"|'|`)?[a-zA-Z0-9=]{112}(\\\"|'|`)?"),
		3,
	},
	{
		"CC",
		regexp.MustCompile("(?i)codeclima.{0,50}(\\\"|'|`)?[0-9a-f]{64}(\\\"|'|`)?"),
		1,
	},
	{
		"Facebook Access Token",
		regexp.MustCompile("EAACEdEose0cBA[0-9A-Za-z]+"),
		3,
	},
	{
		"Facebook Access Token Base64",
		regexp.MustCompile("RUFBQ0VkRW9zZTBjQk[%a-zA-Z0-9+/]+={0,2}"),
		3,
	},
	{
		"Facebook Oauth",
		regexp.MustCompile("(?i)facebook[^/]{0,50}(\\\"|'|`)?[0-9a-f]{32}(\\\"|'|`)?"),
		3,
	},
	{
		"Google (GCP) Service-account",
		regexp.MustCompile("((\\\"|'|`)?type(\\\"|'|`)?\\\\s{0,50}(:|=>|=)\\\\s{0,50}(\\\"|'|`)?service_account(\\\"|'|`)?,?)"),
		3,
	},
	{
		"Google API Key",
		regexp.MustCompile("AIza[0-9A-Za-z\\-_]{35}"),
		1,
	},
	{
		"Google API Key Base64",
		regexp.MustCompile("QUl6Y[%a-zA-Z0-9+/]{47}"),
		1,
	},
	{
		"Google OAuth",
		regexp.MustCompile("[0-9]+-[0-9A-Za-z_]{32}\\.apps\\.googleusercontent\\.com"),
		3,
	},
	{
		"Google OAuth Access Token",
		regexp.MustCompile("ya29\\.[0-9A-Za-z\\-_]+"),
		3,
	},
	{
		"Google Oauth",
		regexp.MustCompile("((\\\"|'|`)?client_secret(\\\"|'|`)?\\\\s{0,50}(:|=>|=)\\\\s{0,50}(\\\"|'|`)?[a-zA-Z0-9-_]{24}(\\\"|'|`)?)"),
		3,
	},
	{
		"Heroku API Key",
		regexp.MustCompile("(?i)heroku.{0,50}[0-9A-F]{8}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{12}"),
		2,
	},
	{
		"Hockeyapp",
		regexp.MustCompile("(?i)hockey.{0,50}(\\\"|'|`)?[0-9a-f]{32}(\\\"|'|`)?"),
		3,
	},
	{
		"MailChimp API Key",
		regexp.MustCompile("[0-9a-f]{32}-us[0-9]{1,2}"),
		3,
	},
	{
		"Mailgun API Key",
		regexp.MustCompile("key-[0-9a-zA-Z]{32}"),
		3,
	},
	{
		"Outlook team",
		regexp.MustCompile("https\\://outlook\\.office.com/webhook/[0-9a-f-]{36}\\@"),
		3,
	},
	{
		"PayPal Braintree Access Token",
		regexp.MustCompile("access_token\\$production\\$[0-9a-z]{16}\\$[0-9a-f]{32}"),
		1,
	},
	{
		"PGP private key block",
		regexp.MustCompile("-----BEGIN PGP PRIVATE KEY BLOCK-----"),
		1,
	},
	{
		"PGP private key block Base64",
		regexp.MustCompile("LS0tLS1CRUdJTiBQR1AgUFJJVkFURSBLRVkgQkxPQ0stLS0tL[%a-zA-Z0-9+/]+={0,2}"),
		1,
	},
	{
		"RSA private key",
		regexp.MustCompile("-----BEGIN RSA PRIVATE KEY-----"),
		1,
	},
	{
		"RSA private key Base64",
		regexp.MustCompile("LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tL[%a-zA-Z0-9+/]+={0,2}"),
		1,
	},
	{
		"SSH (DSA) private key",
		regexp.MustCompile("-----BEGIN DSA PRIVATE KEY-----"),
		1,
	},
	{
		"SSH (DSA) private key Base64",
		regexp.MustCompile("LS0tLS1CRUdJTiBEU0EgUFJJVkFURSBLRVktLS0tL[%a-zA-Z0-9+/]+={0,2}"),
		1,
	},
	{
		"SSH (EC) private key",
		regexp.MustCompile("-----BEGIN EC PRIVATE KEY-----"),
		1,
	},
	{
		"SSH (EC) private key Base64",
		regexp.MustCompile("LS0tLS1CRUdJTiBFQyBQUklWQVRFIEtFWS0tLS0t[%a-zA-Z0-9+/]+={0,2}"),
		1,
	},
	{
		"SSH (OPENSSH) private key",
		regexp.MustCompile("-----BEGIN OPENSSH PRIVATE KEY-----"),
		1,
	},
	{
		"SSH (OPENSSH) private key Base64",
		regexp.MustCompile("LS0tLS1CRUdJTiBPUEVOU1NIIFBSSVZBVEUgS0VZLS0tLS[%a-zA-Z0-9+/]+={0,2}"),
		1,
	},
	{
		"Sauce",
		regexp.MustCompile("(?i)sauce.{0,50}(\\\"|'|`)?[0-9a-f-]{36}(\\\"|'|`)?"),
		3,
	},
	{
		"Slack Token",
		regexp.MustCompile("(xox[pboa]-[0-9]{12}-[0-9]{12}-[0-9]{12}-[a-z0-9]{32})"),
		3,
	},
	{
		"Slack Webhook",
		regexp.MustCompile("https://hooks.slack.com/services/T[a-zA-Z0-9_]{8}/B[a-zA-Z0-9_]{8}/[a-zA-Z0-9_]{24}"),
		3,
	},
	{
		"Sonar",
		regexp.MustCompile("(?i)sonar.{0,50}(\\\"|'|`)?[0-9a-f]{40}(\\\"|'|`)?"),
		3,
	},
	{
		"Square Access Token",
		regexp.MustCompile("sq0atp-[0-9A-Za-z\\-_]{22}"),
		3,
	},
	{
		"Square OAuth Secret",
		regexp.MustCompile("sq0csp-[0-9A-Za-z\\-_]{43}"),
		3,
	},
	{
		"Stripe API Key",
		regexp.MustCompile("sk_live_[0-9a-zA-Z]{24}"),
		3,
	},
	{
		"Stripe Restricted API Key",
		regexp.MustCompile("rk_live_[0-9a-zA-Z]{24}"),
		1,
	},
	{
		"Surge",
		regexp.MustCompile("(?i)surge.{0,50}(\\\"|'|`)?[0-9a-f]{32}(\\\"|'|`)?"),
		3,
	},
	{
		"Twilio API Key",
		regexp.MustCompile("SK[0-9a-fA-F]{32}"),
		2,
	},
	{
		"Twitter Oauth",
		regexp.MustCompile("(?i)twitter[^/]{0,50}[0-9a-zA-Z]{35,44}"),
		3,
	},
	{
		"Password in URL",
		regexp.MustCompile("[a-zA-Z]{3,10}://[^/\\s:@]{3,20}:[^/\\s:@]{3,20}@.{1,100}[\"'\\s]"),
		2,
	},
	{
		"S3 Buckets",
		regexp.MustCompile("[a-z0-9.-]+\\.s3\\.amazonaws\\.com|[a-z0-9.-]+\\.s3-[a-z0-9-]\\.amazonaws\\.com|[a-z0-9.-]+\\.s3-website[.-](eu|ap|us|ca|sa|cn)|//s3\\.amazonaws\\.com/[a-z0-9._-]+|//s3-[a-z0-9-]+\\.amazonaws\\.com/[a-z0-9._-]+"),
		4,
	},
	{
		"Generic Private Key",
		regexp.MustCompile("-----BEGIN [ A-Za-z0-9]*PRIVATE KEY[ A-Za-z0-9]*-----"),
		4,
	},
	{
		"Generic certificate header",
		regexp.MustCompile("-----BEGIN .{3,100}-----"),
		5,
	},
	{
		"Generic certificate header Base64",
		regexp.MustCompile("LS0tLS1CRUdJT[%a-zA-Z0-9+/]+={0,2}"),
		5,
	},
	{
		"Generic Password",
		regexp.MustCompile("(?i)pass(word)?[\\w-]*\\\\s*[=:>|]+\\s*['\"`][^'\"`]{3,100}['\"`]"),
		4,
	},
	{
		"Generic Secret",
		regexp.MustCompile("(?i)secret[\\w-]*\\s*[=:>|]+\\s*['\"`][^'\"`]{3,100}['\"`]"),
		4,
	},
	{
		"Generic Token",
		regexp.MustCompile("(?i)token[\\w-]*\\s*[=:>|]+\\s*['\"`][^'\"`]{3,100}['\"`]"),
		4,
	},
	{
		"Common secret names",
		regexp.MustCompile("(?i)aws_access|aws_secret|api[_-]?key|listbucketresult|s3_access_key|authorization:|ssh-rsa AA|pass(word)?|secret|token"),
		4,
	},
	{
		"PHP Things Base64",
		regexp.MustCompile("(eyJ|YTo|Tzo|PD[89]|aHR0cHM6L|aHR0cDo|rO0)[%a-zA-Z0-9+/]+={0,2}"),
		7,
	},
	{
		"IP Address",
		regexp.MustCompile("\\b(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\b"),
		8,
	},
	{
		"URL",
		regexp.MustCompile("[a-zA-Z]{2,10}://[-a-zA-Z0-9@:%._\\+~#=]{1,256}\\.[a-zA-Z0-9]{1,6}\\b([-a-zA-Z0-9()@:%_\\+.~#?&//=]*)"),
		9,
	},
	{
		"HTTP URL Base64",
		regexp.MustCompile("(aHR0cD|aHR0cHM6)[%a-zA-Z0-9+/]+={0,2}"),
		9,
	},
	{
		"Email",
		regexp.MustCompile("[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+"),
		9,
	},
	{
		"Hostname",
		regexp.MustCompile("\\b(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9])\\.)+(com|org|net)\\b"),
		9,
	},
	{
		"Suspicious Comments",
		regexp.MustCompile("(?i)\\b(hack|hax|fix|oo+ps|fuck|ugly|todo|shit)\\b"),
		10,
	},
}

var defaultConfig = &Config{
	Rules: make([]*Rule, numberOfDefaultRules, numberOfDefaultRules),
	Blacklist: []*regexp.Regexp{
		regexp.MustCompile("\\.min\\.js$"),
		regexp.MustCompile("\\.(less|s?css|html|lock|pbxproj)$"),
		regexp.MustCompile("node_modules/"),
		regexp.MustCompile("package-lock\\.json$"),
		regexp.MustCompile("bower\\.json$"),
		regexp.MustCompile("\\.pdf$"),
		regexp.MustCompile("npm-debug\\.log"),
	},
}

// Config struct contains the parsed contents of the JSON config file
type Config struct {
	Rules     []*Rule
	Blacklist []*regexp.Regexp
}

func newConfig(r *Robber) *Config {
	if r.Args.Config == "" {
		return defaultConfig
	}
	return parseConfig(r)
}

// jsonConfig struct for reading a given JSON config file
type jsonConfig struct {
	Rules []struct {
		Reason string `json:"Reason"`
		Rule   string `json:"Rule"`
		Noise  int    `json:"Noise"`
	} `json:"Rules"`
	FileBlacklist []string `json:"FileBlacklist"`
	Colors        []struct {
		Name  string `json:"Name"`
		Value string `json:"Value"`
	} `json:"Colors"`
}

// Rule struct holds a given regex rule with its reason for matching
type Rule struct {
	Reason string
	Regex  *regexp.Regexp
	Noise  int
}

// TODO: Comment
func loadDefaultRules(r *Robber, c *Config) {
	for _, rule := range allRules {
		if rule.Noise > r.Args.Noise.Upper || rule.Noise < r.Args.Noise.Lower {
			c.Rules = append(c.Rules, rule)
		}
	}
}

// TODO: Comment
func parseRules(r *Robber, config jsonConfig) []*Rule {
	var rules []*Rule
	for _, rule := range config.Rules {
		if rule.Noise > r.Args.Noise.Upper || rule.Noise < r.Args.Noise.Lower {
			continue
		}
		regex, err := regexp.Compile(rule.Rule)
		if err != nil {
			r.Logger.LogFail(regexErrorMessage, rule.Reason, rule.Rule, err)
		}
		// There is no true need for the Noise field in the Rule struct, it is just there for
		// parseing the default rules.
		rule := &Rule{
			Reason: rule.Reason,
			Regex:  regex,
		}
		rules = append(rules, rule)
	}
	return rules
}

// TODO: Comment
func parseBlacklist(r *Robber, config jsonConfig) []*regexp.Regexp {
	var blacklist []*regexp.Regexp
	for _, fileRule := range config.FileBlacklist {
		regex, err := regexp.Compile(fileRule)
		if err != nil {
			r.Logger.LogFail(regexErrorMessage, "File blacklist", fileRule, err)
		}
		blacklist = append(blacklist, regex)
	}
	return blacklist
}

// TODO: Comment
// User can specify a config file that contains only a blacklist, only color settings, only
// blacklist and color settings, and completely empty.
// User can not specify a config file with no rules, then the default rules will simply be loaded
func parseConfig(r *Robber) *Config {
	var config jsonConfig

	// If the Config argument is not set, add rules within the specified noise range to the
	// defaultConfig.
	if r.Args.Config == "" {
		loadDefaultRules(r, defaultConfig)
		return defaultConfig
	}

	// Read contents of JSON file
	file, _ := os.Open(r.Args.Config)
	defer file.Close()
	reader := bufio.NewReader(file)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		r.Logger.LogFail("Unable to read file %s: %s", file.Name(), err)
	}

	// Parse JSON file and compile regex rules
	json.Unmarshal([]byte(content), &config)
	parsedConfig := &Config{}
	if len(config.Rules) == 0 {
		loadDefaultRules(r, parsedConfig)
	} else {
		parsedConfig.Rules = parseRules(r, config)
	}
	parsedConfig.Blacklist = parseBlacklist(r, config)
	if len(config.Colors) > 0 {
		parseColors(r, config)
	}
	return parsedConfig
}
