package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v2"
)

func checkFileMod(filename string) bool {
	file, err := os.Stat(filename)

	if err != nil {
		fmt.Println(err)
		return false
	}

	if LastHostsModTime.IsZero() {
		LastHostsModTime = file.ModTime()
		fmt.Println("read hosts file")
		return true
	}

	if -LastHostsModTime.Sub(file.ModTime()).Seconds() > 0 {
		LastHostsModTime = file.ModTime()
		fmt.Println("hosts file was modified. re-read")
		return true
	}

	return false
}

func readHosts() {
	for {
		if checkFileMod(*HostsFile) {
			content, err := ioutil.ReadFile(*HostsFile)
			if err != nil {
				panic(err)
			}
			err = yaml.Unmarshal([]byte(content), &HostsArray)
			if err != nil {
				panic(err)
			}
			for i, host := range HostsArray.Hosts {
				if host.Login == "" {
					HostsArray.Hosts[i].Login = HostsArray.Global.Login
				}
				if host.Password == "" {
					HostsArray.Hosts[i].Password = HostsArray.Global.Password
				}

				if HostsArray.SessionKeys == nil {
					HostsArray.SessionKeys = map[string]string{}
				}
				if HostsArray.ExpiresDates == nil {
					HostsArray.ExpiresDates = map[string]string{}
				}
				if HostsArray.Cookies == nil {
					HostsArray.Cookies = map[string]string{}
				}

				HostsArray.Cookies[host.Address] = ""
				HostsArray.ExpiresDates[host.Address] = ""
				HostsArray.SessionKeys[host.Address] = ""
			}
		}

		time.Sleep(30 * time.Second)
	}
}

func currentTimestamp() string {
	retTime := time.Now().Add(time.Duration(-5) * time.Second)
	return string(retTime.Unix())
}

func checkUserExpire(host Host) int64 {
	t, _ := time.Parse(time.ANSIC, HostsArray.ExpiresDates[host.Address])
	return int64(
		math.Round(-time.Now().Sub(t).Seconds()))
}

func labelsList(parameter string, ilo string) []prometheus.Labels {
	var outVar []prometheus.Labels
	metrics, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		return nil
	}
	for _, metric := range metrics {
		if len(metric.Metric) > 0 && *metric.Name == parameter {
			for _, subMetric := range metric.Metric {
				if len(subMetric.Label) > 0 {
					tmpVar := make(prometheus.Labels, 20)
					for _, label := range subMetric.Label {
						tmpVar[*label.Name] = *label.Value
					}
					if tmpVar["ilo"] == ilo {
						outVar = append(outVar, tmpVar)
					}
				}
			}
		}
	}
	return outVar
}

func deletePrometheusMetrics(ilo string) {
	for _, parameter := range []string{"temperature_current", "temperature_critical"} {
		for _, labels := range labelsList(parameter, ilo) {
			PromTemperaturesCurrent.Delete(labels)
			PromTemperaturesCritical.Delete(labels)
		}
	}
	for _, parameter := range []string{"power_state", "in_post"} {
		for _, labels := range labelsList(parameter, ilo) {
			PromPowerState.Delete(labels)
			PromInPost.Delete(labels)
		}
	}
}

