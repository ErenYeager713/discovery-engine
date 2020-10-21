package libs

import (
	"strings"

	"github.com/accuknox/knoxAutoPolicy/types"
	pb "github.com/accuknox/knoxServiceFlowMgmt/src/proto"
)

func isSynFlagOnly(tcp *pb.TCP) bool {
	if tcp.Flags.SYN && !tcp.Flags.ACK {
		return true
	}
	return false
}

func getL4Ports(l4 *pb.Layer4) (int, int) {
	if l4.TCP != nil {
		return int(l4.TCP.SourcePort), int(l4.TCP.DestinationPort)
	} else if l4.UDP != nil {
		return int(l4.UDP.SourcePort), int(l4.UDP.DestinationPort)
	} else if l4.ICMPv4 != nil {
		return int(l4.ICMPv4.Type), int(l4.ICMPv4.Code)
	} else {
		return -1, -1
	}
}

func getProtocol(l4 *pb.Layer4) int {
	if l4.TCP != nil {
		return 6
	} else if l4.UDP != nil {
		return 17
	} else if l4.ICMPv4 != nil {
		return 1
	} else {
		return 0 // unknown?
	}
}

func ConvertTrafficToLog(flow *pb.TrafficFlow) types.NetworkLog {
	log := types.NetworkLog{}

	if flow.Source.Namespace == "" {
		log.SrcMicroserviceName = "external"
	} else {
		log.SrcMicroserviceName = flow.Source.Namespace
	}

	if flow.Source.Pod == "" {
		log.SrcContainerGroupName = flow.Ip.Source
	} else {
		log.SrcContainerGroupName = flow.Source.Pod
	}

	if flow.Destination.Namespace == "" {
		log.DstMicroserviceName = "external"
	} else {
		log.DstMicroserviceName = flow.Destination.Namespace
	}

	if flow.Destination.Pod == "" {
		log.DstContainerGroupName = flow.Ip.Destination
	} else {
		log.DstContainerGroupName = flow.Destination.Pod
	}

	log.SrcMac = flow.Ethernet.Source
	log.DstMac = flow.Ethernet.Destination

	log.Protocol = getProtocol(flow.L4)
	if log.Protocol == 6 { //
		log.SynFlag = isSynFlagOnly(flow.L4.TCP)
	}

	log.SrcIP = flow.Ip.Source
	log.DstIP = flow.Ip.Destination

	log.SrcPort, log.DstPort = getL4Ports(flow.L4)

	if flow.Verdict == "FORWARDED" {
		log.Action = "allow"
	} else if flow.Verdict == "DROPPED" {
		log.Action = "deny"
	} else { // default
		log.Action = "unknown"
	}

	log.Direction = flow.TrafficDirection

	return log
}

func ToCiliumNetworkPolicy(inPolicy types.KnoxNetworkPolicy) types.CiliumNetworkPolicy {
	ciliumPolicy := types.CiliumNetworkPolicy{}

	ciliumPolicy.APIVersion = "cilium.io/v2"
	ciliumPolicy.Kind = "CiliumNetworkPolicy"
	ciliumPolicy.Metadata = map[string]string{}
	for k, v := range inPolicy.Metadata {
		ciliumPolicy.Metadata[k] = v
	}

	// update selector
	ciliumPolicy.Spec.Selector.MatchLabels = map[string]string{}
	for k, v := range inPolicy.Spec.Selector.MatchLabels {
		ciliumPolicy.Spec.Selector.MatchLabels[k] = v
	}

	// update egress
	egress := types.CiliumEgress{}

	if inPolicy.Spec.Egress.MatchLabels != nil {
		matchLabels := map[string]string{}
		for k, v := range inPolicy.Spec.Egress.MatchLabels {
			matchLabels[k] = v
		}

		toEndpoints := []types.CiliumToEndpoints{types.CiliumToEndpoints{matchLabels}}
		egress.ToEndpoints = toEndpoints
	}

	// update toPorts
	for _, toPort := range inPolicy.Spec.Egress.ToPorts {
		if egress.ToPorts == nil {
			egress.ToPorts = []types.CiliumToPort{}
			ciliumPort := types.CiliumToPort{}
			ciliumPort.Ports = []types.CiliumPort{}
			egress.ToPorts = append(egress.ToPorts, ciliumPort)
		}

		port := types.CiliumPort{Port: toPort.Ports, Protocol: strings.ToUpper(toPort.Protocol)}
		egress.ToPorts[0].Ports = append(egress.ToPorts[0].Ports, port)
	}

	// update toCIDRs
	for _, toCIDR := range inPolicy.Spec.Egress.ToCIDRs {
		if egress.ToCIDRs == nil {
			egress.ToCIDRs = []types.ToCIDR{}
		}

		egress.ToCIDRs = append(egress.ToCIDRs, toCIDR)
	}

	ciliumPolicy.Spec.Egress = []types.CiliumEgress{}
	ciliumPolicy.Spec.Egress = append(ciliumPolicy.Spec.Egress, egress)

	return ciliumPolicy
}

func ConvertTrafficToLogs(flows []*pb.TrafficFlow) []types.NetworkLog {
	networkLogs := []types.NetworkLog{}
	for _, flow := range flows {
		log := ConvertTrafficToLog(flow)
		networkLogs = append(networkLogs, log)
	}

	return networkLogs
}