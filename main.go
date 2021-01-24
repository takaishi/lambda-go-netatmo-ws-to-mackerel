/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/exzz/netatmo-api-go"
	mkr "github.com/mackerelio/mackerel-client-go"
	"log"
	"os"
	"time"
)

var (
	netatmoEmail        = os.Getenv("NETATMO_EMAIL")
	netatmoPassword     = os.Getenv("NETATMO_PASSWORD")
	netatmoAppID        = os.Getenv("NETATMO_APP_ID")
	netatmoAppSecret    = os.Getenv("NETATMO_APP_SECRET")
	mackerelAPIKey      = os.Getenv("MACKEREL_APIKEY")
	mackerelServiceName = os.Getenv("MACKEREL_SERVICE_NAME")
)

var metricPrefixes = map[string]string{
	"BatteryPercent": "battery",
	"RFStatus":       "rf_status",
	"WifiStatus":     "wifi",
	"Temperature":    "temperature",
	"Humidity":       "humidity",
	"CO2":            "co2",
	"Noise":          "noise",
	"Pressure":       "pressure",
}


func HandleRequest(ctx context.Context) error {
	validateEnvironmentVariables()
	metrics, err := fetchWeatherStationMetrics()
	if err != nil {
		return err
	}
	mackerelClient := mkr.NewClient(mackerelAPIKey)
	return mackerelClient.PostServiceMetricValues(mackerelServiceName, metrics)
}

func main() {
	lambda.Start(HandleRequest)
}

func validateEnvironmentVariables() {
	if netatmoEmail == "" {
		log.Fatalf("NETATMO_EMAIL is empty")
	}
	if netatmoPassword  == "" {
		log.Fatalf("NETATMO_PASSWORD is empty")
	}
	if netatmoAppID == "" {
		log.Fatalf("NETATMO_APP_ID is empty")
	}
	if netatmoAppSecret == "" {
		log.Fatalf("NETATMO_APP_SECRET is empty")
	}
	if mackerelAPIKey == "" {
		log.Fatalf("MACKEREL_APIKEY is empty")
	}
	if mackerelServiceName == "" {
		log.Fatalf("MACKEREL_SERVICE_NAME")
	}
}



func fetchWeatherStationMetrics() ([]*mkr.MetricValue, error){
	metrics := make([]*mkr.MetricValue, 0)
	conf := netatmo.Config{
		ClientID:     netatmoAppID,
		ClientSecret: netatmoAppSecret,
		Username:     netatmoEmail,
		Password:     netatmoPassword,
	}
	netatmoClient, err := netatmo.NewClient(conf)
	if err != nil {
		return metrics, fmt.Errorf("failed: netatmo.NewClient (%s)", err)
	}
	dc, err := netatmoClient.Read()
	if err != nil {
		return metrics, err
	}

	for _, station := range dc.Stations() {
		for _, module := range station.Modules() {
			infoEpoch, info := module.Info()
			infoTs := time.Unix(infoEpoch, 0)
			for name, value := range info {
				log.Printf("Timestamp:%s Station:%s Module:%s Name:%s Value:%#v\n", infoTs, module.StationName, module.ModuleName, name, value)
				if metricPrefix, ok := metricPrefixes[name]; ok {
					if v, ok := float64of(value); ok {
						metrics = append(metrics, &mkr.MetricValue{
							Name:  "netatmo" + "." + metricPrefix + "." + module.ModuleName,
							Value: v,
							Time:  infoTs.Unix(),
						})
					}
				}
			}

			dataEpoch, data := module.Data()
			dataTs := time.Unix(dataEpoch, 0)
			for name, value := range data {
				log.Printf("Timestamp:%s Station:%s Module:%s Name:%s Value:%#v\n", dataTs, module.StationName, module.ModuleName, name, value)
				if metricPrefix, ok := metricPrefixes[name]; ok {
					if v, ok := float64of(value); ok {
						metrics = append(metrics, &mkr.MetricValue{
							Name:  "netatmo" + "." + metricPrefix + "." + module.ModuleName,
							Value: v,
							Time:  dataTs.Unix(),
						})
					}
				}
			}
		}
	}

	return metrics, nil
}


func float64of(value interface{}) (float64, bool) {
	switch t := value.(type) {
	case float64:
		return t, true
	case *float64:
		return *t, true
	case float32:
		return float64(t), true
	case *float32:
		return float64(*t), true
	case int32:
		return float64(t), true
	case *int32:
		return float64(*t), true
	case int64:
		return float64(t), true
	case *int64:
		return float64(*t), true
	default:
		return 0, false
	}
}
