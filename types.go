package main

type HostsType struct {
	Global struct {
		Login    string `yaml:"login"`
		Password string `yaml:"password"`
	} `yaml:"global"`
	Hosts        []Host `yaml:"hosts"`
	SessionKeys  map[string]string
	ExpiresDates map[string]string
	Cookies      map[string]string
}

type Host struct {
	Address  string `yaml:"address"`
	Name     string `yaml:"name"`
	Login    string `yaml:"login"`
	Password string `yaml:"password"`
}

type iloCookie struct {
	SessionKey       string `json:"session_key"`
	UserName         string `json:"user_name"`
	UserAccount      string `json:"user_account"`
	UserDn           string `json:"user_dn"`
	UserType         string `json:"user_type"`
	UserIP           string `json:"user_ip"`
	UserExpires      string `json:"user_expires"`
	LoginPriv        int    `json:"login_priv"`
	RemoteConsPriv   int    `json:"remote_cons_priv"`
	VirtualMediaPriv int    `json:"virtual_media_priv"`
	ResetPriv        int    `json:"reset_priv"`
	ConfigPriv       int    `json:"config_priv"`
	UserPriv         int    `json:"user_priv"`
	HostAddress      string `json:"address"`
}

type HealthType struct {
	SelfTest                string `json:"self_test"`
	SystemHealth            string `json:"system_health"`
	HostpwrState            string `json:"hostpwr_state"`
	FansStatus              string `json:"fans_status"`
	FansRedundancy          string `json:"fans_redundancy"`
	TemperatureStatus       string `json:"temperature_status"`
	PowerSuppliesStatus     string `json:"power_supplies_status"`
	PowerSuppliesRedundancy string `json:"power_supplies_redundancy"`
	PowerSuppliesMismatch   int    `json:"power_supplies_mismatch"`
	StorageStatus           string `json:"storage_status"`
	NicStatus               string `json:"nic_status"`
	CPUStatus               string `json:"cpu_status"`
	MemStatus               string `json:"mem_status"`
	ExtHlthStatus           string `json:"ext_hlth_status"`
	BatteryStatus           string `json:"battery_status"`
	AmsReady                string `json:"ams_ready"`
	InPost                  int    `json:"in_post"`
}

type TemperaturesType struct {
	HostpwrState string `json:"hostpwr_state"`
	InPost       int    `json:"in_post"`
	Temperature  []struct {
		Label          string `json:"label"`
		Xposition      int    `json:"xposition"`
		Yposition      int    `json:"yposition"`
		Location       string `json:"location"`
		Status         string `json:"status"`
		Currentreading int    `json:"currentreading"`
		Caution        int    `json:"caution"`
		Critical       int    `json:"critical"`
		TempUnit       string `json:"temp_unit"`
	} `json:"temperature"`
}

