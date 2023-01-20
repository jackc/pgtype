package pgtype

type (
	IConfig interface {
		SetNaNInfinityNegativeInfinityAware(enable bool)
		GetNaNInfinityNegativeInfinityAware() bool
	}
	config struct {
		NaN_Infinity_NegativeInfinity_Aware bool
	}
)

var _config *config

func newConfig() *config {
	return &config{
		NaN_Infinity_NegativeInfinity_Aware: false,
	}
}
func GetConfig() IConfig {
	return _config
}
func (c *config) SetNaNInfinityNegativeInfinityAware(enable bool) {
	c.NaN_Infinity_NegativeInfinity_Aware = enable
}
func (c *config) GetNaNInfinityNegativeInfinityAware() bool {
	return c.NaN_Infinity_NegativeInfinity_Aware
}

func init() {
	// Register the default config
	_config = newConfig()
}
