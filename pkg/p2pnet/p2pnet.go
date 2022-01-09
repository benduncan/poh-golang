package p2pnet

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"net"

	"github.com/benduncan/poh-golang/pkg/poh_hash"
	"github.com/benduncan/poh-golang/pkg/wallet"
)

// "Safe" UDP packet size up to 508 bytes
type Packet struct {
	Version            [1]byte
	Reserved           [3]byte
	SenderPublicKey    [32]byte // TODO: Consider Base36 (with ICAP) or Base56 encoding w/ unique identifier
	RecipientPublicKey [32]byte
	Payload            [376]byte
	SenderSignature    [64]byte
}

type P2P struct {
	Host            string
	Port            string
	Version         uint8
	maxDatagramSize int
	POH             *poh_hash.POH
}

func New() *P2P {

	return &P2P{Port: "15500", Version: 1, Host: "224.0.0.1", maxDatagramSize: 8192}
}

func (this *P2P) BroadcastListen(h func(*net.UDPAddr, int, []byte)) {

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", this.Host, this.Port))

	fmt.Println("BroadcastListen")

	if err != nil {
		log.Fatal(err)
	}
	l, err := net.ListenMulticastUDP("udp", nil, addr)
	l.SetReadBuffer(8192) //this.maxDatagramSize)
	for {
		b := make([]byte, this.maxDatagramSize)
		n, src, err := l.ReadFromUDP(b)
		if err != nil {
			log.Fatal("ReadFromUDP failed:", err)
		}
		h(src, n, b)
	}

}

func (this *P2P) BroadcastSend(data string) {

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%s", this.Port))

	if err != nil {
		log.Fatal(err)
	}

	c, err := net.DialUDP("udp", nil, addr)

	mywallet := wallet.New()
	mywallet.GenerateWallet()

	senderwallet := wallet.New()
	senderwallet.GenerateWallet()

	packet := Packet{}

	packet.Version = [1]byte{1}
	packet.Reserved = [3]byte{0, 0, 0}

	copy(packet.SenderPublicKey[:], mywallet.PublicKey)
	copy(packet.RecipientPublicKey[:], mywallet.PublicKey) //senderwallet.PublicKey)

	msg := []byte(fmt.Sprintf("WOW A super important message here for you to consider %s", data))

	copy(packet.Payload[:], msg)

	signature, _ := mywallet.Sign(packet.Payload[:])

	copy(packet.SenderSignature[:], signature)

	/*
		fmt.Println("Len PublicKey => ", len(mywallet.PublicKey))
		fmt.Println("Len PrivateKey => ", len(mywallet.PrivateKey))
		fmt.Println("Len Signature => ", len(signature))
		fmt.Println("Len Msg => ", len(msg))
		fmt.Println("Len Payload => ", len(packet.Payload[:]))

		fmt.Printf("Orig PubKey => %v\n", mywallet.PublicKey)
		fmt.Printf("Orig Signature => %v\n", packet.SenderSignature[:])
		fmt.Printf("Orig Payload => %v\n", packet.Payload[:])
	*/

	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, packet); err != nil {
		fmt.Println(err)
		return
	}

	verify, _ := mywallet.Verify(packet.Payload[:], signature)

	if verify == true {
		c.Write(buf.Bytes())
	}

	defer c.Close()

}

func (this *P2P) MsgHandler(src *net.UDPAddr, n int, b []byte) {
	log.Println(n, "bytes read from", src)
	log.Println(hex.Dump(b[:n]))

	packet := Packet{}

	r := bytes.NewReader(b)

	if err := binary.Read(r, binary.BigEndian, &packet); err != nil {
		fmt.Println("failed to Read:", err)
		return
	}

	fmt.Printf("Payload => %v\n", packet.Payload[:])
	//fmt.Printf("Signature => %v\n", packet.SenderSignature[:])

	// Confirm if validated
	mywallet := wallet.New()
	verify := mywallet.VerifyRaw(packet.SenderPublicKey[:], packet.Payload[:], packet.SenderSignature[:])

	//fmt.Printf("Pub key => %v\n", mywallet.PublicKey)
	fmt.Println("\nVerify Check => ", verify)

	//fmt.Printf("Verification %s\n", verify)

	if verify == true {
		queuedata := poh_hash.QueueData{Data: packet.Payload[:], Sender: packet.SenderPublicKey[:]}
		this.POH.QueueSync.State = append(this.POH.QueueSync.State, queuedata)
	}

}
