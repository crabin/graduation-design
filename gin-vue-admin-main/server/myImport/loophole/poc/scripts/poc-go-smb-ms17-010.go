package scripts

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/myImport/loophole/pkg/util"
	"net"
	"time"
)

var (
	negotiateProtocolRequest, _  = hex.DecodeString("00000085ff534d4272000000001853c00000000000000000000000000000fffe00004000006200025043204e4554574f524b2050524f4752414d20312e3000024c414e4d414e312e30000257696e646f777320666f7220576f726b67726f75707320332e316100024c4d312e325830303200024c414e4d414e322e3100024e54204c4d20302e313200")
	sessionSetupRequest, _       = hex.DecodeString("00000088ff534d4273000000001807c00000000000000000000000000000fffe000040000dff00880004110a000000000000000100000000000000d40000004b000000000000570069006e0064006f007700730020003200300030003000200032003100390035000000570069006e0064006f007700730020003200300030003000200035002e0030000000")
	treeConnectRequest, _        = hex.DecodeString("00000060ff534d4275000000001807c00000000000000000000000000000fffe0008400004ff006000080001003500005c005c003100390032002e003100360038002e003100370035002e003100320038005c00490050004300240000003f3f3f3f3f00")
	transNamedPipeRequest, _     = hex.DecodeString("0000004aff534d42250000000018012800000000000000000000000000088ea3010852981000000000ffffffff0000000000000000000000004a0000004a0002002300000007005c504950455c00")
	trans2SessionSetupRequest, _ = hex.DecodeString("0000004eff534d4232000000001807c00000000000000000000000000008fffe000841000f0c0000000100000000000000a6d9a40000000c00420000004e0001000e000d0000000000000000000000000000")
)

func checkMs17010(args *ScriptScanArgs) (*util.ScanResult, error) {
	// connecting to a host in LAN if reachable should be very quick
	ip := args.Host
	conn, err := net.DialTimeout("tcp", ip+":445", time.Second*3)
	if err != nil {
		// fmt.Printf("failed to connect to %s\n", ip)
		return nil, err
	}
	defer conn.Close()

	conn.Write(negotiateProtocolRequest)
	reply := make([]byte, 1024)
	// let alone half packet
	if n, err := conn.Read(reply); err != nil || n < 36 {
		return nil, err
	}

	if binary.LittleEndian.Uint32(reply[9:13]) != 0 {
		// status != 0
		return nil, err
	}

	conn.Write(sessionSetupRequest)

	n, err := conn.Read(reply)
	if err != nil || n < 36 {
		return nil, err
	}

	if binary.LittleEndian.Uint32(reply[9:13]) != 0 {
		// status != 0
		fmt.Printf("can't determine whether %s is vulnerable or not\n", ip)
		return nil, err
	}

	userID := reply[32:34]
	treeConnectRequest[32] = userID[0]
	treeConnectRequest[33] = userID[1]
	// TODO change the ip in tree path though it doesn't matter
	conn.Write(treeConnectRequest)

	if n, err := conn.Read(reply); err != nil || n < 36 {
		return nil, err
	}

	treeID := reply[28:30]
	transNamedPipeRequest[28] = treeID[0]
	transNamedPipeRequest[29] = treeID[1]
	transNamedPipeRequest[32] = userID[0]
	transNamedPipeRequest[33] = userID[1]

	conn.Write(transNamedPipeRequest)
	if n, err := conn.Read(reply); err != nil || n < 36 {
		return nil, err
	}

	if reply[9] == 0x05 && reply[10] == 0x02 && reply[11] == 0x00 && reply[12] == 0xc0 {
		//fmt.Printf("%s(%s) is likely VULNERABLE to MS17-010!\n", ip, os)
		return util.VulnerableTcpOrUdpResult(ip, "", nil, nil), nil
	} else {
		//fmt.Printf("%s(%s) stays in safety\n", ip, os)
		return &util.InVulnerableResult, nil
	}

}

func init() {
	ScriptRegister("poc-go-smb-ms17-010", checkMs17010)
}
