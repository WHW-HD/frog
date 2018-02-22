package ads1115


import "time"
import "golang.org/x/exp/io/i2c"

const (
	ADS1015_REG_CONFIG_CQUE_NONE 	= 0x0003 // Disable the comparator and put ALERT/RDY in high state (default)
	ADS1015_REG_CONFIG_CLAT_NONLAT 	= 0x0000 // Non-latching comparator (default)
	ADS1015_REG_CONFIG_CPOL_ACTVLOW = 0x0000 // ALERT/RDY pin is low when active (default)
	ADS1015_REG_CONFIG_CMODE_TRAD 	= 0x0000 // Traditional comparator with hysteresis (default)
	ADS1015_REG_CONFIG_MODE_SINGLE 	= 0x0100 // Power-down single-shot mode (default)
	ADS1115_REG_CONFIG_DR_8SPS 	= 0x0000 // 8 samples per second
	ADS1015_REG_CONFIG_PGA_6_144V 	= 0x0000 // +/-6.144V range
	ADS1015_REG_CONFIG_MUX_SINGLE_0 = 0x4000 // Single-ended AIN0
	ADS1015_REG_CONFIG_OS_SINGLE 	= 0x8000 // Write: Set to start a single-conversion
	ADS1115_Config =
		ADS1015_REG_CONFIG_CQUE_NONE |
		ADS1015_REG_CONFIG_CLAT_NONLAT |
		ADS1015_REG_CONFIG_CPOL_ACTVLOW |
		ADS1015_REG_CONFIG_CMODE_TRAD |
		ADS1015_REG_CONFIG_MODE_SINGLE |
		ADS1115_REG_CONFIG_DR_8SPS |
		ADS1015_REG_CONFIG_PGA_6_144V |
		ADS1015_REG_CONFIG_MUX_SINGLE_0 |
		ADS1015_REG_CONFIG_OS_SINGLE

	ADS1015_REG_POINTER_CONFIG = 0x01
	ADS1015_REG_POINTER_CONVERT = 0x00
)

type Ads1115 struct {
	device *i2c.Device
}

func New() (Ads1115, error) {
	result := Ads1115{}
	device, err := i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, 0x48)
	if err != nil {
		return result, err
	}
	result.device = device
	return result, nil
}

func (ads *Ads1115) Read() (float32, error) {
	bytes := []byte{(ADS1115_Config >> 8) & 0xFF, ADS1115_Config & 0xFF}
	err := ads.device.WriteReg(ADS1015_REG_POINTER_CONFIG, bytes)
	if err != nil {
		return float32(-1), err
	}
	delay := 1000 / 128 + 1
	time.Sleep(time.Duration(delay) * time.Millisecond)
	var result = make([]byte, 2)
	err = ads.device.ReadReg(ADS1015_REG_POINTER_CONVERT, result)
	if err != nil {
		return float32(-1), err
	}
	val := (uint32(result[0]) << 8) | uint32(result[1])
	var data float32
	if val > 0x7FFF	{
		data = float32((val - 0xFFFF) * 6144) / 32768.0
	} else {
		data = float32(val * 6144) / 32768.0
	}

	return data, nil
}

func (ads *Ads1115) Close() error {
	return ads.device.Close()
}
