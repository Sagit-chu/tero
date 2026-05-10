package monitor

type GFWChecker interface {
	IsBlocked(ip string) (bool, error)
}
