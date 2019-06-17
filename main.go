package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	ginprometheus "github.com/mcuadros/go-gin-prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	IloHealth               HealthType
	HostsArray              HostsType
	ListenAddr              = getEnv("LISTEN_ADDR", "0.0.0.0")
	ListenPort              = getEnv("LISTEN_PORT", "10011")
	LastHostsModTime        time.Time
	PromTemperaturesCurrent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "temperature_current",
			Help: "temperatures from ilo",
		},
		[]string{"ilo", "label", "location"},
	)
	PromTemperaturesCritical = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "temperature_critical",
			Help: "temperatures from ilo",
		},
		[]string{"ilo", "label", "location"},
	)
	PromPowerState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "power_state",
			Help: "power state",
		},
		[]string{"ilo"},
	)

	PromInPost = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "in_post",
			Help: "In POST",
		},
		[]string{"ilo"},
	)

	PromGatheringState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gathering_state",
			Help: "power state",
		},
		[]string{"ilo"},
	)
)

func getCookie(host Host) {

	var cookie iloCookie

	url := "https://" + host.Address + "/json/login_session"

	payload := strings.NewReader("{\"method\": \"login\", \"user_login\": \"" + host.Login + "\", \"password\": \"" + host.Password + "\"}")

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		fmt.Println(err)
		PromGatheringState.WithLabelValues(host.Address).Set(0)
		deletePrometheusMetrics(host.Address)
		return
	}

	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("Cache-Control", "no-cache")

	timeout := time.Duration(5 * time.Second)
	client := &http.Client{Timeout: timeout}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		PromGatheringState.WithLabelValues(host.Address).Set(0)
		deletePrometheusMetrics(host.Address)
		return
	}

	err = json.NewDecoder(res.Body).Decode(&cookie)
	if cookie.SessionKey != "" {
		HostsArray.SessionKeys[host.Address] = cookie.SessionKey
		HostsArray.ExpiresDates[host.Address] = cookie.UserExpires
		HostsArray.Cookies[host.Address] = "sessionKey=" + HostsArray.SessionKeys[host.Address] + "; path=/; domain=" + host.Address + "; Secure; Expires=" + HostsArray.ExpiresDates[host.Address] + ";"
	}

	if err != nil {
		fmt.Println(err)
		PromGatheringState.WithLabelValues(host.Address).Set(0)
		deletePrometheusMetrics(host.Address)
	} else {
		PromGatheringState.WithLabelValues(host.Address).Set(1)
	}
	err = res.Body.Close()
	if err != nil {
		fmt.Println(err)
		PromGatheringState.WithLabelValues(host.Address).Set(0)
		deletePrometheusMetrics(host.Address)
	}
}

func getMetrics(host Host) {
	if HostsArray.Cookies[host.Address] == "" {
		PromGatheringState.WithLabelValues(host.Address).Set(0)
		return
	}
	url := "https://" + host.Address + "/json/health_summary?_=" + currentTimestamp()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		PromGatheringState.WithLabelValues(host.Address).Set(0)
		return
	}

	timeout := time.Duration(5 * time.Second)
	client := &http.Client{Timeout: timeout}

	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Cookie", HostsArray.Cookies[host.Address])
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		PromGatheringState.WithLabelValues(host.Address).Set(0)
		return
	}

	err = json.NewDecoder(res.Body).Decode(&IloHealth)
	if err != nil {
		fmt.Println(err)
		PromGatheringState.WithLabelValues(host.Address).Set(0)
	} else {
		PromGatheringState.WithLabelValues(host.Address).Set(1)
	}
	err = res.Body.Close()
	if err != nil {
		fmt.Println(err)
		PromGatheringState.WithLabelValues(host.Address).Set(0)
	}
}

func getTemperatures(host Host) {
	var IloTemperature TemperaturesType
	if HostsArray.Cookies[host.Address] == "" {
		PromGatheringState.WithLabelValues(host.Address).Set(0)
		deletePrometheusMetrics(host.Address)
		return
	}
	url := "https://" + host.Address + "/json/health_temperature?_=" + currentTimestamp()

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println(err)
		PromGatheringState.WithLabelValues(host.Address).Set(0)
		deletePrometheusMetrics(host.Address)
		return
	}

	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Cookie", HostsArray.Cookies[host.Address])

	timeout := time.Duration(5 * time.Second)
	client := &http.Client{Timeout: timeout}

	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		PromGatheringState.WithLabelValues(host.Address).Set(0)
		deletePrometheusMetrics(host.Address)
		return
	}

	err = json.NewDecoder(res.Body).Decode(&IloTemperature)
	if err != nil {
		fmt.Println(err)
		PromGatheringState.WithLabelValues(host.Address).Set(0)
		deletePrometheusMetrics(host.Address)
	} else {
		PromGatheringState.WithLabelValues(host.Address).Set(1)
	}

	if IloTemperature.HostpwrState == "ON" {
		PromPowerState.WithLabelValues(host.Address).Set(1)
	}

	PromInPost.WithLabelValues(host.Address).Set(float64(IloTemperature.InPost))

	for _, temp := range IloTemperature.Temperature {
		PromTemperaturesCurrent.WithLabelValues(host.Address, temp.Label, temp.Location).Set(float64(temp.Currentreading))
		PromTemperaturesCritical.WithLabelValues(host.Address, temp.Label, temp.Location).Set(float64(temp.Critical))
	}

	err = res.Body.Close()

	if err != nil {
		fmt.Println(err)
		PromGatheringState.WithLabelValues(host.Address).Set(0)
		deletePrometheusMetrics(host.Address)
	}
}

func init() {
	if PromErr := prometheus.Register(PromTemperaturesCurrent); PromErr != nil {
		panic(PromErr)
	}

	if PromErr := prometheus.Register(PromTemperaturesCritical); PromErr != nil {
		panic(PromErr)
	}

	if PromErr := prometheus.Register(PromPowerState); PromErr != nil {
		panic(PromErr)
	}

	if PromErr := prometheus.Register(PromInPost); PromErr != nil {
		panic(PromErr)
	}

	if PromErr := prometheus.Register(PromGatheringState); PromErr != nil {
		panic(PromErr)
	}

}

func main() {
	go readHosts()
	go func() {
		for {
			for _, host := range HostsArray.Hosts {
				if checkUserExpire(Host{Address: host.Address, Name: host.Name, Login: host.Login, Password: host.Password}) < 180 {
					fmt.Println("User cookie updating for host: " + host.Address)
					go getCookie(Host{Address: host.Address, Name: host.Name, Login: host.Login, Password: host.Password})
				}
			}
			time.Sleep(15 * time.Second)
		}
	}()

	go func() {
		for {
			for _, host := range HostsArray.Hosts {
				go getMetrics(Host{Address: host.Address, Name: host.Name, Login: host.Login, Password: host.Password})
				go getTemperatures(Host{Address: host.Address, Name: host.Name, Login: host.Login, Password: host.Password})
			}

			time.Sleep(5 * time.Second)
		}
	}()

	r := gin.Default()
	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)
	err := r.Run(ListenAddr + ":" + ListenPort)

	if err != nil {
		panic(err)
	}
}
