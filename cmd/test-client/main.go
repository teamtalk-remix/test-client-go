package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gogo/protobuf/proto"
	"github.com/teamtalk-remix/test-client-go/cmd/test-client/help"
	"github.com/teamtalk-remix/test-client-go/cmd/test-client/utils"
	"github.com/teamtalk-remix/test-client-go/pkg/pdubase"
	b "github.com/teamtalk-remix/test-client-go/proto/IM_BaseDefine"
	"github.com/teamtalk-remix/test-client-go/proto/IM_Login"
	"github.com/teamtalk-remix/test-client-go/proto/IM_Message"
	"github.com/teamtalk-remix/test-client-go/proto/IM_Other"
)

var (
	logfileStr = fmt.Sprintf("/tmp/testc-%s.log", time.Now().Format("20060102150405"))
	logpath    = flag.String("logpath", logfileStr, "Log Path")
)

var globalConn net.Conn
var msgid uint32

const CLIENT_HEARTBEAT_INTERVAL = 30000

func Connect(serverIP, serverPort string) (net.Conn, error) {
	c, err := net.Dial("tcp", serverIP+":"+serverPort)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func SendMsg(conn net.Conn, from, to uint32, s string) {
	msgid = msgid + 1
	mt := b.MsgType_MSG_TYPE_SINGLE_TEXT
	msg := &IM_Message.IMMsgData{
		FromUserId:  proto.Uint32(from),
		ToSessionId: proto.Uint32(to),
		MsgId:       proto.Uint32(msgid),
		MsgData:     []byte(s),
		CreateTime:  proto.Uint32(uint32(time.Now().Unix())),
		MsgType:     &mt,
	}
	var pdu pdubase.CImPdu
	pdu.SetServiceId(b.ServiceID_SID_MSG)
	pdu.SetCommandId(b.CommandID_CID_MSG_DATA)
	pdu.SetPB(msg)

	_, err := conn.Write(pdu.Buffer().Bytes())
	if err != nil {
		checkErr(err)
	}
}

func SendLogin(conn net.Conn, userName, userPass string) {
	ols := b.UserStatType_USER_STATUS_ONLINE
	ct := b.ClientType_CLIENT_TYPE_WINDOWS

	msg := &IM_Login.IMLoginReq{
		UserName:      proto.String(userName),
		Password:      proto.String(getMD5Hash(userPass)),
		OnlineStatus:  &ols,
		ClientType:    &ct,
		ClientVersion: proto.String("1.0"),
	}
	var pdu pdubase.CImPdu
	pdu.SetServiceId(b.ServiceID_SID_LOGIN)
	pdu.SetCommandId(b.CommandID_CID_LOGIN_REQ_USERLOGIN)
	pdu.SetPB(msg)

	_, err := conn.Write(pdu.Buffer().Bytes())
	if err != nil {
		checkErr(err)
	}
}

func SendHeartBeat(conn net.Conn, d time.Duration) {
	for x := range time.Tick(d) {
		if conn != nil {
			utils.Log.Println("Time: ", x)
			msg := &IM_Other.IMHeartBeat{}
			var pdu pdubase.CImPdu
			pdu.SetServiceId(b.ServiceID_SID_OTHER)
			pdu.SetCommandId(b.CommandID_CID_OTHER_HEARTBEAT)
			pdu.SetPB(msg)
			_, err := conn.Write(pdu.Buffer().Bytes())
			if err != nil {
				checkErr(err)
			}
			utils.Log.Println("Sending heartbeat succeed")
		} else {
			utils.Log.Println("globalConn is null exit....")
			break
		}
	}
}

func HandleLoginResponse(conn net.Conn, pdu pdubase.CImPdu) {
	msg := IM_Login.IMLoginRes{}
	pdu.GetPBMsg(pdu.Buffer().Bytes(), &msg)
	spew.Dump(msg.String())
	spew.Dump(msg.GetResultCode())
	spew.Dump(msg.GetResultString())
	spew.Dump(msg.GetUserInfo())

	if msg.GetResultCode() == b.ResultType_REFUSE_REASON_NONE {
		go SendHeartBeat(conn, CLIENT_HEARTBEAT_INTERVAL*time.Millisecond)
	}
}

func HandleMsgDataResponse(conn net.Conn, pdu pdubase.CImPdu) {
	msg := IM_Message.IMMsgData{}
	pdu.GetPBMsg(pdu.Buffer().Bytes(), &msg)
	spew.Dump("receive HandleMsgDataResponse from:", msg.FromUserId, msg.MsgType)
	spew.Dump(msg.String())
}

func HandleListAllUserResponse(conn net.Conn, pdu pdubase.CImPdu) {
	//IM::Buddy::IMAllUserRsp msgResp;
	//uint32_t nSeqNo = pPdu->GetSeqNum();
	//if(msgResp.ParseFromArray(pPdu->GetBodyData(), pPdu->GetBodyLength()))
	//{
	//	uint32_t userCnt = msgResp.user_list_size();
	//	printf("get %d users\n", userCnt);
	//	list<IM::BaseDefine::UserInfo> lsUsers;
	//	for(uint32_t i=0; i<userCnt; ++i)
	//	{
	//	IM::BaseDefine::UserInfo cUserInfo = msgResp.user_list(i);
	//	lsUsers.push_back(cUserInfo);
	//	}
	//	m_pCallback->onGetChangedUser(nSeqNo, lsUsers);
	//}
	//else
	//{
	//	m_pCallback->onError(nSeqNo, pPdu->GetCommandId(), "parse pb error");
	//}
}

func HandleListUserInfoResponse(conn net.Conn, pdu pdubase.CImPdu) {
	//IM::Buddy::IMUsersInfoRsp msgResp;
	//uint32_t nSeqNo = pPdu->GetSeqNum();
	//if(msgResp.ParseFromArray(pPdu->GetBodyData(), pPdu->GetBodyLength()))
	//{
	//	uint32_t userCnt = msgResp.user_info_list_size();
	//	list<IM::BaseDefine::UserInfo> lsUser;
	//	for (uint32_t i=0; i<userCnt; ++i) {
	//IM::BaseDefine::UserInfo userInfo = msgResp.user_info_list(i);
	//	lsUser.push_back(userInfo);
	//}
	//	m_pCallback->onGetUserInfo(nSeqNo, lsUser);
	//}
	//else
	//{
	//	m_pCallback->onError(nSeqNo, pPdu->GetCommandId(), "parse pb error");
	//}
}

func HandleMsgDataAck(conn net.Conn, pdu pdubase.CImPdu) {
	//IM::Message::IMMsgDataAck msgResp;
	//uint32_t nSeqNo = pPdu->GetSeqNum();
	//if(msgResp.ParseFromArray(pPdu->GetBodyData(), pPdu->GetBodyLength()))
	//{
	//uint32_t nSendId = msgResp.user_id();
	//uint32_t nRecvId = msgResp.session_id();
	//uint32_t nMsgId = msgResp.msg_id();
	//IM::BaseDefine::SessionType nType = msgResp.session_type();
	//m_pCallback->onSendMsg(nSeqNo, nSendId, nRecvId, nType, nMsgId);
	//}
	//else
	//{
	//m_pCallback->onError(nSeqNo, pPdu->GetCommandId(), "parse pb error");
	//}
}

func HandleMsgListResponse(conn net.Conn, pdu pdubase.CImPdu) {
	//IM::Message::IMGetMsgListRsp msgResp;
	//uint32_t nSeqNo = pPdu->GetSeqNum();
	//if(msgResp.ParseFromArray(pPdu->GetBodyData(), pPdu->GetBodyLength()))
	//{
	//	uint32_t nUserId= msgResp.user_id();
	//IM::BaseDefine::SessionType nSessionType = msgResp.session_type();
	//	uint32_t nPeerId = msgResp.session_id();
	//	uint32_t nMsgId = msgResp.msg_id_begin();
	//	uint32_t nMsgCnt = msgResp.msg_list_size();
	//	list<IM::BaseDefine::MsgInfo> lsMsg;
	//	for(uint32_t i=0; i<nMsgCnt; ++i)
	//	{
	//	IM::BaseDefine::MsgInfo msgInfo = msgResp.msg_list(i);
	//	lsMsg.push_back(msgInfo);
	//	}
	//	m_pCallback->onGetMsgList(nSeqNo, nUserId, nPeerId, nSessionType, nMsgId, nMsgCnt, lsMsg);
	//}
	//else
	//{
	//	m_pCallback->onError(nSeqNo, pPdu->GetCommandId(), "parse pb falied");
	//}
}

func HandleMsgUnreadCntResponse(conn net.Conn, pdu pdubase.CImPdu) {
	//IM::Message::IMUnreadMsgCntRsp msgResp;
	//uint32_t nSeqNo = pPdu->GetSeqNum();
	//if(msgResp.ParseFromArray(pPdu->GetBodyData(), pPdu->GetBodyLength()))
	//{
	//	list<IM::BaseDefine::UnreadInfo> lsUnreadInfo;
	//	uint32_t nUserId = msgResp.user_id();
	//	uint32_t nTotalCnt = msgResp.total_cnt();
	//	uint32_t nCnt = msgResp.unreadinfo_list_size();
	//	for (uint32_t i=0; i<nCnt; ++i) {
	//IM::BaseDefine::UnreadInfo unreadInfo = msgResp.unreadinfo_list(i);
	//	lsUnreadInfo.push_back(unreadInfo);
	//}
	//	m_pCallback->onGetUnreadMsgCnt(nSeqNo, nUserId, nTotalCnt, lsUnreadInfo);
	//}
	//else
	//{
	//	m_pCallback->onError(nSeqNo, pPdu->GetCommandId(), "parse pb fail");
	//}
}
func HandleRecentContactSessionResponse(conn net.Conn, pdu pdubase.CImPdu) {
	//IM::Buddy::IMRecentContactSessionRsp msgResp;
	//uint32_t nSeqNo = pPdu->GetSeqNum();
	//if(msgResp.ParseFromArray(pPdu->GetBodyData(), pPdu->GetBodyLength()))
	//{
	//	list<IM::BaseDefine::ContactSessionInfo> lsSession;
	//	uint32_t nUserId = msgResp.user_id();
	//	uint32_t nCnt = msgResp.contact_session_list_size();
	//	for (uint32_t i=0; i<nCnt; ++i) {
	//IM::BaseDefine::ContactSessionInfo session = msgResp.contact_session_list(i);
	//	lsSession.push_back(session);
	//}
	//	m_pCallback->onGetRecentSession(nSeqNo, nUserId, lsSession);
	//}
	//else
	//{
	//	m_pCallback->onError(nSeqNo, pPdu->GetCommandId(), "parse pb error");
	//}
}

//func Timer
func HandlePdu(conn net.Conn, pdu pdubase.CImPdu) {
	switch pdu.CommandId() {
	case b.CommandID_CID_OTHER_HEARTBEAT:
		utils.Log.Println("HandlePdu:  heartbeat received....")
	case b.CommandID_CID_LOGIN_RES_USERLOGIN:
		HandleLoginResponse(conn, pdu)
	case b.CommandID_CID_BUDDY_LIST_ALL_USER_RESPONSE:
		HandleListAllUserResponse(conn, pdu)
	case b.CommandID_CID_BUDDY_LIST_USER_INFO_RESPONSE:
		HandleListUserInfoResponse(conn, pdu)
	case b.CommandID_CID_MSG_UNREAD_CNT_RESPONSE:
		HandleMsgUnreadCntResponse(conn, pdu)
	case b.CommandID_CID_BUDDY_LIST_RECENT_CONTACT_SESSION_RESPONSE:
		HandleRecentContactSessionResponse(conn, pdu)
	case b.CommandID_CID_MSG_LIST_RESPONSE:
		HandleMsgListResponse(conn, pdu)
	case b.CommandID_CID_MSG_DATA_ACK:
		HandleMsgDataAck(conn, pdu)
	case b.CommandID_CID_MSG_DATA:
		HandleMsgDataResponse(conn, pdu)
	case b.CommandID_CID_BUDDY_LIST_STATUS_NOTIFY:

	default:
		utils.Log.Println("wrong msg_type", pdu.CommandId(), pdu.ServiceId())
	}
}

func checkErr(err error) {
	if err == nil {
		utils.Log.Println("Ok")
		return

	} else if netError, ok := err.(net.Error); ok && netError.Timeout() {
		utils.Log.Println("Timeout")
		return
	}

	switch t := err.(type) {
	case *net.OpError:
		if t.Op == "dial" {
			utils.Log.Println("Unknown host")
		} else if t.Op == "read" {
			utils.Log.Println("Connection refused")
		}

	case syscall.Errno:
		if t == syscall.ECONNREFUSED {
			utils.Log.Println("Connection refused")
		}
	}
}

func testReceiveLoginResponse(conn net.Conn) {
	for {
		headerBuf := make([]byte, pdubase.PduHeaderSize)
		_, err := conn.Read(headerBuf)
		if err != nil {
			utils.Log.Println("Breaking......:", err.Error())
			checkErr(err)
			break
		}

		//deal with heartbeat
		var pdu pdubase.CImPdu
		pdu.GetHeader(headerBuf)
		pdu.Length()
		bodyBuf := make([]byte, pdu.Length()-pdubase.PduHeaderSize)
		_, err = conn.Read(bodyBuf[:pdu.Length()-pdubase.PduHeaderSize])
		if err != nil {
			utils.Log.Println("Breaking......:", err.Error())
			checkErr(err)
			break
		}
		//copy
		pdu.WriteBuffer(bodyBuf)
		HandlePdu(conn, pdu)
	}
}

func recvMsg(conn net.Conn) {
	utils.Log.Println("\n-> start lisiten.....")
	testReceiveLoginResponse(conn)
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func login(conn net.Conn, userName, userPass string) {
	SendLogin(conn, userName, userPass)
}

type httpResponse struct {
	BackupIP   string `json:"backupIP"`
	Code       int    `json:"code"`
	Discovery  string `json:"discovery"`
	MsfsBackup string `json:"msfsBackup"`
	MsfsPrior  string `json:"msfsPrior"`
	Msg        string `json:"msg"`
	Port       string `json:"port"`
	PriorIP    string `json:"priorIP"`
}

func getMsgServerIPPort() (string, string) {
	resp, err := http.Get("http://teamtalk-ng-login.cloud.matrix.int.iwarp.org/msg_server")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	res := httpResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		panic(err)
	}
	if res.Code != 0 {
		utils.Log.Panic("res.Code is not 0")
	}

	if res.PriorIP == "" {
		return res.BackupIP, res.Port
	}
	return "10.110.196.118", "8000"
	//return res.PriorIP, res.Port
}

func main() {

	utils.NewLog(*logpath)
	fmt.Println("Simple Client, type 'help' to show commands")
	fmt.Println("---------------------")

	for {
		fmt.Print("-> ")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)

		cmds := strings.Fields(text)
		if len(cmds) <= 0 {
			continue
		}

		switch cmds[0] {
		case "help":
			help.PrintHelp()
		case "login":
			if globalConn != nil {
				fmt.Println("client already login in...")
				continue
			}
			userName, userPass := cmds[1], cmds[2]
			serverIP, serverPort := getMsgServerIPPort()
			var err error
			globalConn, err = Connect(serverIP, serverPort)
			if err != nil {
				panic(err)
			}
			login(globalConn, userName, userPass)
			//TODO: this recvMsg could move out of switch block
			go recvMsg(globalConn)
		case "sendmsg":
			if globalConn == nil {
				fmt.Println("globalConn does not exist... please run connect command first...")
				os.Exit(1)
			}
			from, err := strconv.Atoi(cmds[1])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			to, err := strconv.Atoi(cmds[2])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			go SendMsg(globalConn, uint32(from), uint32(to), cmds[3])
		case "exit":
			if globalConn != nil {
				globalConn.Close()
			}
			fmt.Println("haha exit...")
			os.Exit(0)
		default:
			help.PrintHelp()
		}
	}
}
