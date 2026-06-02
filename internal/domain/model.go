package domain

type DomainConfig struct {
	ExpectedCodes []int `json:"expected_codes"`
	Retries       int   `json:"retries"`
}

type CheckResult struct {
	Domain string
	Up     bool
	Code   int
	Error  string
}

type DomainState struct {
	Status string
	Code   int
	Error  string
}
