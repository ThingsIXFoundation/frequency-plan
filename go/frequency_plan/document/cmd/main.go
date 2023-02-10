package main

import (
	"github.com/ThingsIXFoundation/frequency-plan/go/frequency_plan"
	"github.com/brocaar/lorawan/band"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

func main() {
	eu868band, err := frequency_plan.GetBand(string(frequency_plan.EU868))
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Infof("Printing Frequency Plan %s", eu868band.Name())

	logrus.Infof("Uplink channels")

	for _, channelIndex := range eu868band.GetEnabledUplinkChannelIndices() {
		channel, err := eu868band.GetUplinkChannel(channelIndex)
		if err != nil {
			logrus.Fatal(err)
		}

		minDr, err := eu868band.GetDataRate(channel.MinDR)
		if err != nil {
			logrus.Fatal(err)
		}

		maxDr, err := eu868band.GetDataRate(channel.MaxDR)
		if err != nil {
			logrus.Fatal(err)
		}
		if maxDr.Modulation == band.LRFHSSModulation {
			for dr := channel.MaxDR; dr > 0; dr-- {
				maxDr, err = eu868band.GetDataRate(dr)
				if err != nil {
					logrus.Fatal(err)
				}

				if maxDr.Modulation == band.LoRaModulation {
					break
				}
			}
		}

		if minDr != maxDr {
			logrus.Infof("Multi-SF LoRa channel %d, Frequency=%.1f MHz, SF%dBW%d-SF%dBW%d", channelIndex, float64(channel.Frequency)/(1000.0*1000.0), maxDr.SpreadFactor, maxDr.Bandwidth, minDr.SpreadFactor, minDr.Bandwidth)
		}

		if minDr == maxDr && minDr.Modulation == band.LoRaModulation {
			logrus.Infof("Std LoRa channel %d, Frequency=%.1f MHz, SF%dBW%d", channelIndex, float64(channel.Frequency)/(1000.0*1000.0), minDr.SpreadFactor, minDr.Bandwidth)
		}

		if minDr == maxDr && minDr.Modulation == band.FSKModulation {
			logrus.Infof("FSK channel %d, Frequency=%.1f MHz", channelIndex, float64(channel.Frequency)/(1000.0*1000.0))
		}
	}

	logrus.Infof("RX1 Downlink channels")
	var reportedIndexes []int
	for _, channelIndex := range eu868band.GetEnabledUplinkChannelIndices() {
		rx1ChannelIndex, err := eu868band.GetRX1ChannelIndexForUplinkChannelIndex(channelIndex)
		if err != nil {
			logrus.Fatal(err)
		}
		if slices.Contains(reportedIndexes, rx1ChannelIndex) {
			continue
		}

		reportedIndexes = append(reportedIndexes, rx1ChannelIndex)

		channel, err := eu868band.GetDownlinkChannel(rx1ChannelIndex)
		if err != nil {
			logrus.Fatal(err)
		}

		minDr, err := eu868band.GetDataRate(channel.MinDR)
		if err != nil {
			logrus.Fatal(err)
		}

		maxDr, err := eu868band.GetDataRate(channel.MaxDR)
		if err != nil {
			logrus.Fatal(err)
		}
		if maxDr.Modulation == band.LRFHSSModulation {
			for dr := channel.MaxDR; dr > 0; dr-- {
				maxDr, err = eu868band.GetDataRate(dr)
				if err != nil {
					logrus.Fatal(err)
				}

				if maxDr.Modulation == band.LoRaModulation {
					break
				}
			}
		}

		if minDr != maxDr {
			logrus.Infof("Multi-SF LoRa channel %d, Frequency=%.1f MHz, SF%dBW%d-SF%dBW%d", channelIndex, float64(channel.Frequency)/(1000.0*1000.0), maxDr.SpreadFactor, maxDr.Bandwidth, minDr.SpreadFactor, minDr.Bandwidth)
		}

		if minDr == maxDr && minDr.Modulation == band.LoRaModulation {
			logrus.Infof("Std LoRa channel %d, Frequency=%.1f MHz, SF%dBW%d", channelIndex, float64(channel.Frequency)/(1000.0*1000.0), minDr.SpreadFactor, minDr.Bandwidth)
		}

		if minDr == maxDr && minDr.Modulation == band.FSKModulation {
			logrus.Infof("FSK channel %d, Frequency=%.1f MHz", channelIndex, float64(channel.Frequency)/(1000.0*1000.0))
		}
	}

	logrus.Infof("RX2 Downlink channel")

	rx2dr, err := eu868band.GetDataRate(eu868band.GetDefaults().RX2DataRate)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Infof("Frequency=%d, SF%dBW%d", eu868band.GetDefaults().RX2Frequency, rx2dr.SpreadFactor, rx2dr.Bandwidth)
}
