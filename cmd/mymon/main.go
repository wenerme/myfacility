package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"time"
)

func main() {
	pcapLive()
}
func decode(data []byte) {
	// Decode a packet
	packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Default)
	// Get the TCP layer from this packet
	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		fmt.Println("This is a TCP packet!")
		// Get actual TCP data from this layer
		tcp, _ := tcpLayer.(*layers.TCP)
		fmt.Printf("From src port %d to dst port %d\n", tcp.SrcPort, tcp.DstPort)
	}
	// Iterate over all layers, printing out each layer type
	for _, layer := range packet.Layers() {
		fmt.Println("PACKET LAYER:", layer.LayerType())
	}
}
func pcapInactive() {
	inactive, err := pcap.NewInactiveHandle("")
	if err != nil {
		log.Fatal(err)
	}
	defer inactive.CleanUp()

	// Call various functions on inactive to set it up the way you'd like:
	if err = inactive.SetTimeout(time.Minute); err != nil {
		log.Fatal(err)
	} else if err = inactive.SetTimestampSource(pcap.TimestampSource(0)); err != nil {
		log.Fatal(err)
	}

	// Finally, create the actual handle by calling Activate:
	handle, err := inactive.Activate() // after this, inactive is no longer valid
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// Now use your handle as you see fit.
}
func pcapLive() {
	if handle, err := pcap.OpenLive("lo0", 1600, true, pcap.BlockForever); err != nil {
		panic(err)
	} else if err := handle.SetBPFFilter("tcp and port 3306"); err != nil { // optional
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			_ = packet
			//			fmt.Println(hex.Dump(packet.Data()))
			fmt.Println(packet.Dump())
			//			handlePacket(packet)  // Do something with a packet here.
		}
	}
}
func pcapFile() {
	if handle, err := pcap.OpenOffline("/path/to/my/file"); err != nil {
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			_ = packet
			//			handlePacket(packet)  // Do something with a packet here.
		}
	}
}
