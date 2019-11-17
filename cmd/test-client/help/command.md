   
    case CID_OTHER_HEARTBEAT:
     _HandleHeartBeat(pPdu);
     break;
    case CID_LOGIN_REQ_USERLOGIN:
     _HandleLoginRequest(pPdu );
     break;
 case CID_LOGIN_REQ_LOGINOUT:
     _HandleLoginOutRequest(pPdu);
     break;
 case CID_LOGIN_REQ_DEVICETOKEN:
     _HandleClientDeviceToken(pPdu);
     break;
 case CID_LOGIN_REQ_KICKPCCLIENT:
     _HandleKickPCClient(pPdu);
     break;
 case CID_LOGIN_REQ_PUSH_SHIELD:
     _HandlePushShieldRequest(pPdu);
     break;
 case CID_LOGIN_REQ_QUERY_PUSH_SHIELD:
     _HandleQueryPushShieldRequest(pPdu);
     break;
 case CID_MSG_DATA:
     _HandleClientMsgData(pPdu);
     break;
 case CID_MSG_DATA_ACK:
     _HandleClientMsgDataAck(pPdu);
     break;
 case CID_MSG_TIME_REQUEST:
     _HandleClientTimeRequest(pPdu);
     break;
 case CID_MSG_LIST_REQUEST:
     _HandleClientGetMsgListRequest(pPdu);
     break;
 case CID_MSG_GET_BY_MSG_ID_REQ:
     _HandleClientGetMsgByMsgIdRequest(pPdu);
     break;
 case CID_MSG_UNREAD_CNT_REQUEST:
     _HandleClientUnreadMsgCntRequest(pPdu );
     break;
 case CID_MSG_READ_ACK:
     _HandleClientMsgReadAck(pPdu);
     break;
 case CID_MSG_GET_LATEST_MSG_ID_REQ:
     _HandleClientGetLatestMsgIDReq(pPdu);
     break;
 case CID_SWITCH_P2P_CMD:
     _HandleClientP2PCmdMsg(pPdu );
     break;
 case CID_BUDDY_LIST_RECENT_CONTACT_SESSION_REQUEST:
     _HandleClientRecentContactSessionRequest(pPdu);
     break;
 case CID_BUDDY_LIST_USER_INFO_REQUEST:
     _HandleClientUserInfoRequest( pPdu );
     break;
 case CID_BUDDY_LIST_REMOVE_SESSION_REQ:
     _HandleClientRemoveSessionRequest( pPdu );
     break;
 case CID_BUDDY_LIST_ALL_USER_REQUEST:
     _HandleClientAllUserRequest(pPdu );
     break;
 case CID_BUDDY_LIST_CHANGE_AVATAR_REQUEST:
     _HandleChangeAvatarRequest(pPdu);
     break;
 case CID_BUDDY_LIST_CHANGE_SIGN_INFO_REQUEST:
     _HandleChangeSignInfoRequest(pPdu);
     break;
 case CID_BUDDY_LIST_USERS_STATUS_REQUEST:
     _HandleClientUsersStatusRequest(pPdu);
     break;
 case CID_BUDDY_LIST_DEPARTMENT_REQUEST:
     _HandleClientDepartmentRequest(pPdu);
     break;
 // for group process
 case CID_GROUP_NORMAL_LIST_REQUEST:
     s_group_chat->HandleClientGroupNormalRequest(pPdu, this);
     break;
 case CID_GROUP_INFO_REQUEST:
     s_group_chat->HandleClientGroupInfoRequest(pPdu, this);
     break;
 case CID_GROUP_CREATE_REQUEST:
     s_group_chat->HandleClientGroupCreateRequest(pPdu, this);
     break;
 case CID_GROUP_CHANGE_MEMBER_REQUEST:
     s_group_chat->HandleClientGroupChangeMemberRequest(pPdu, this);
     break;
 case CID_GROUP_SHIELD_GROUP_REQUEST:
     s_group_chat->HandleClientGroupShieldGroupRequest(pPdu, this);
     break;
 case CID_FILE_REQUEST:
     s_file_handler->HandleClientFileRequest(this, pPdu);
     break;
 case CID_FILE_HAS_OFFLINE_REQ:
     s_file_handler->HandleClientFileHasOfflineReq(this, pPdu);
     break;
 case CID_FILE_ADD_OFFLINE_REQ:
     s_file_handler->HandleClientFileAddOfflineReq(this, pPdu);
     break;
 case CID_FILE_DEL_OFFLINE_REQ:
     s_file_handler->HandleClientFileDelOfflineReq(this, pPdu);
     break;
 default:
     log("wrong msg, cmd id=%d, user id=%u. ", pPdu->GetCommandId(), GetUserId());
     break;