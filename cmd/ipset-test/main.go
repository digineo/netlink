package main

import (
	"flag"
	"log"
	"net"

	"github.com/vishvananda/netlink"
)

func main() {
	timeout := flag.Int("timeout", -1, "timeout, negative means omit the argument")
	comment := flag.String("comment", "", "comment")
	withComments := flag.Bool("with-comments", false, "create set with comment support")
	withCounters := flag.Bool("with-counters", false, "create set with counters support")
	withSkbinfo := flag.Bool("with-skbinfo", false, "create set with skbinfo support")
	replace := flag.Bool("replace", false, "replace existing set/entry")
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		panic("invalid arguments")
	}

	var timeoutVal *uint32
	if *timeout >= 0 {
		v := uint32(*timeout)
		timeoutVal = &v
	}

	log.SetFlags(log.Lshortfile)

	cmd := args[0]
	args = args[1:]

	switch cmd {
	case "protocol":
		protocol, err := netlink.IpsetProtocol()
		if err != nil {
			panic(err)
		}
		log.Println("Protocol:", protocol)

	case "create":
		if len(args) != 2 {
			panic("invalid arguments")
		}

		err := netlink.IpsetCreate(args[0], args[1], netlink.IpsetCreateOptions{
			Replace:  *replace,
			Timeout:  timeoutVal,
			Comments: *withComments,
			Counters: *withCounters,
			Skbinfo:  *withSkbinfo,
		})
		if err != nil {
			panic(err)
		}

	case "destroy":
		if len(args) != 1 {
			panic("invalid arguments")
		}
		err := netlink.IpsetDestroy(args[0])
		if err != nil {
			panic(err)
		}

	case "list":
		if len(args) != 1 {
			panic("invalid arguments")
		}

		result, err := netlink.IpsetList(args[0])
		if err != nil {
			panic(err)
		}
		log.Printf("%+v", result)

	case "listall":
		result, err := netlink.IpsetListAll()
		if err != nil {
			panic(err)
		}
		for _, ipset := range result {
			log.Printf("%+v", ipset)
		}

	case "add", "del":
		if len(args) != 2 {
			panic("invalid arguments")
		}

		setName := args[0]
		element := args[1]

		mac, _ := net.ParseMAC(element)
		entry := netlink.IPSetEntry{
			Timeout: timeoutVal,
			MAC:     mac,
			Comment: *comment,
			Replace: *replace,
		}

		var err error
		if cmd == "add" {
			err = netlink.IpsetAdd(setName, &entry)
		} else {
			err = netlink.IpsetDel(setName, &entry)
		}

		if err != nil {
			panic(err)
		}
	default:
		panic("invalid command")
	}
}
