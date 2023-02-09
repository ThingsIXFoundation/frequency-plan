package main

import (
	"github.com/ThingsIXFoundation/frequency-plan/go/frequency_plan"
	"github.com/brocaar/lorawan/band"
	"github.com/sirupsen/logrus"
)

func main() {
	eu868band, err := frequency_plan.GetBand(string(frequency_plan.EU868))
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Infof("Printing Frequency Plan %s", eu868band.Name())

	logrus.Infof("Uplink channels")

	for _, channelIndex := range eu868band.GetUplinkChannelIndices() {
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

		if minDr != maxDr {
			logrus.Infof("Multi-SF LoRa channel %d, Frequency=%d, SF%dBW%d-SF%dBW%d", channelIndex+1, channel.Frequency, maxDr.SpreadFactor, maxDr.Bandwidth, minDr.SpreadFactor, minDr.Bandwidth)
		}

		if minDr == maxDr && minDr.Modulation == band.LoRaModulation {
			logrus.Infof("Std LoRa channel %d, Frequency=%d, SF%dBW%d", channelIndex+1, channel.Frequency, minDr.SpreadFactor, minDr.Bandwidth)
		}

		if minDr == maxDr && minDr.Modulation == band.FSKModulation {
			logrus.Infof("FSK channel %d, Frequency=%d", channelIndex+1, channel.Frequency)
		}
	}

	logrus.Infof("RX1 Downlink channels")

	for _, channelIndex := range eu868band.GetUplinkChannelIndices() {
		rx1ChannelIndex, err := eu868band.GetRX1ChannelIndexForUplinkChannelIndex(channelIndex)
		if err != nil {
			logrus.Fatal(err)
		}

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

		if minDr != maxDr {
			logrus.Infof("Multi-SF LoRa channel %d, Frequency=%d, SF%dBW%d-SF%dBW%d", channelIndex+1, channel.Frequency, maxDr.SpreadFactor, maxDr.Bandwidth, minDr.SpreadFactor, minDr.Bandwidth)
		}

		if minDr == maxDr && minDr.Modulation == band.LoRaModulation {
			logrus.Infof("Std LoRa channel %d, Frequency=%d, SF%dBW%d", channelIndex+1, channel.Frequency, minDr.SpreadFactor, minDr.Bandwidth)
		}

		if minDr == maxDr && minDr.Modulation == band.FSKModulation {
			logrus.Infof("FSK channel %d, Frequency=%d", channelIndex+1, channel.Frequency)
		}
	}

	logrus.Infof("RX2 Downlink channel")

	rx2dr, err := eu868band.GetDataRate(eu868band.GetDefaults().RX2DataRate)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Infof("Frequency=%d, SF%dBW%d", eu868band.GetDefaults().RX2Frequency, rx2dr.SpreadFactor, rx2dr.Bandwidth)
}
