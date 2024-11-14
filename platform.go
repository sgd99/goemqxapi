package goemqxapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type ClientData struct {
	Node           string `json:"node"`            // 客户端所连接的节点名称
	Clientid       string `json:"clientid"`        // 客户端标识符
	Username       string `json:"username"`        // 客户端连接时使用的用户名
	ProtoName      string `json:"proto_name"`      // 客户端协议名称
	ProtoVer       int    `json:"proto_ver"`       // 客户端使用的协议版本
	IpAdress       string `json:"ip_address"`      // 客户端的 IP 地址
	Port           int    `json:"port"`            // 客户端的端口
	IsBridge       bool   `json:"is_bridge"`       // 指示客户端是否通过桥接方式连接
	ConnectedAt    string `json:"connected_at"`    // 客户端连接时间，格式为 "YYYY-MM-DD HH:mm:ss"
	DisconnectedAt string `json:"disconnected_at"` // 客户端离线时间，格式为 "YYYY-MM-DD HH:mm:ss"， 此字段仅在 connected 为 false 时有效并被返回
	Connected      bool   `json:"connected"`       // 客户端是否处于连接状态
	Zone           string `json:"zone"`            // 指示客户端使用的配置组
	KeepAlive      int    `json:"keepalive"`       // 保持连接时间，单位：秒

	CleanStart     bool   `json:"clean_start"`     // 指示客户端是否使用了全新的会话
	ExpiryInterval int    `json:"expiry_interval"` // 会话过期间隔，单位：秒
	CreateAt       string `json:"create_at"`       // 	会话创建时间，格式为 "YYYY-MM-DD HH:mm:ss"

	SubscriptionsCnt int `json:"subscriptions_cnt"` // 此客户端已建立的订阅数量
	MaxSubscriptions int `json:"max_subscriptions"` // 此客户端允许建立的最大订阅数量

	Inflight      int `json:"inflight"`       // 飞行队列当前长度
	MaxInflight   int `json:"max_inflight"`   // 飞行队列最大长度
	MqueueLen     int `json:"mqueue_len"`     // 消息队列当前长度
	MaxMqueue     int `json:"max_mqueue"`     // 消息队列最大长度
	MqueueDropped int `json:"mqueue_dropped"` // 消息队列因超出长度而丢弃的消息数量

	AwaitingRel    int `json:"awaiting_rel"`     // 未确认的 PUBREC 报文数量
	MaxAwaitingRel int `json:"max_awaiting_rel"` // 允许存在未确认的 PUBREC 报文的最大数量

	RecvOct int `json:"recv_oct"` // EMQX Broker（下同）接收的字节数量
	RecvCnt int `json:"recv_cnt"` // 接收的 TCP 报文数量
	RecvPkt int `json:"recv_pkt"` // 接收的 MQTT 报文数量
	RecvMsg int `json:"recv_msg"` // 接收的 PUBLISH 报文数量

	SendOct int `json:"send_oct"` // 发送的字节数量
	SendCnt int `json:"send_cnt"` // 发送的 TCP 报文数量
	SendPkt int `json:"send_pkt"` // 发送的 MQTT 报文数量
	SendMsg int `json:"send_msg"` // 发送的 PUBLISH 报文数量

	MailboxLen int `json:"mailbox_len"` // 进程邮箱大小
	HeapSize   int `json:"heap_size"`   // 进程堆栈大小，单位：字节
	Reductions int `json:"reductions"`  // Erlang reduction
}

type ClientsMetaData struct {
	Page  int `json:"page"`  // 页码
	Limit int `json:"limit"` // 每页显示的数据条数
	Count int `json:"count"` // 数据总条数
}

type CLientsData struct {
	Data []ClientData    `json:"data"`
	Meta ClientsMetaData `json:"meta"`
}

func (c *CLientsData) GetClient(clientId string) *ClientData {
	for _, client := range c.Data {
		if client.Clientid == clientId {
			return &client
		}
	}
	return nil
}

func (c *CLientsData) GetClientIds() []string {
	ids := make([]string, 0, len(c.Data))
	for _, client := range c.Data {
		ids = append(ids, client.Clientid)
	}
	return ids
}

type ClientsRequest struct {
	Page  int `json:"_page,omitempty"`  // 页码
	Limit int `json:"_limit,omitempty"` // 每页显示的数据条数

	Clientid string `json:"clientid,omitempty"` // 客户端标识符
	Username string `json:"username,omitempty"` // 客户端用户名
	Zone     string `json:"zone,omitempty"`     // 客户端配置组名称

	IpAdress   string `json:"ip_address,omitempty"`      // 客户端 IP 地址
	ConnState  string `json:"connected_state,omitempty"` // 客户端当前连接状态， 可取值有：connected,idle,disconnected
	CleanStart bool   `json:"clean_start,omitempty"`     // 客户端是否使用了全新的会话
	ProtoName  string `json:"proto_name,omitempty"`      // 客户端协议名称， 可取值有：MQTT,CoAP,MQTT-SN
	ProtoVer   string `json:"proto_ver,omitempty"`       // 客户端协议版本

	GteCreateAt int `json:"_gte_create_at,omitempty"` // 客户端会话创建时间，小于等于查找
	LteCreateAt int `json:"_lte_create_at,omitempty"` // 客户端会话创建时间，大于等于查找

	GteConnectedAt int `json:"_gte_connected_at,omitempty"` // 客户端连接创建时间，小于等于查找
	LteConnectedAt int `json:"_lte_connected_at,omitempty"` // 客户端连接创建时间，大于等于查找
}

func (c *ClientsRequest) QueryString() string {
	query := url.Values{}

	if c.Page != 0 {
		query.Set("_page", strconv.Itoa(c.Page))
	}
	if c.Limit != 0 {
		query.Set("_limit", strconv.Itoa(c.Limit))
	}
	if c.Clientid != "" {
		query.Set("clientid", c.Clientid)
	}
	if c.Username != "" {
		query.Set("username", c.Username)
	}
	if c.Zone != "" {
		query.Set("zone", c.Zone)
	}
	if c.IpAdress != "" {
		query.Set("ip_address", c.IpAdress)
	}
	if c.ConnState != "" {
		query.Set("connected_state", c.ConnState)
	}
	if c.CleanStart {
		query.Set("clean_start", strconv.FormatBool(c.CleanStart))
	}
	if c.ProtoName != "" {
		query.Set("proto_name", c.ProtoName)
	}
	if c.ProtoVer != "" {
		query.Set("proto_ver", c.ProtoVer)
	}
	if c.GteCreateAt != 0 {
		query.Set("_gte_create_at", strconv.Itoa(c.GteCreateAt))
	}
	if c.LteCreateAt != 0 {
		query.Set("_lte_create_at", strconv.Itoa(c.LteCreateAt))
	}
	if c.GteConnectedAt != 0 {
		query.Set("_gte_connected_at", strconv.Itoa(c.GteConnectedAt))
	}
	if c.LteConnectedAt != 0 {
		query.Set("_lte_connected_at", strconv.Itoa(c.LteConnectedAt))
	}
	return query.Encode()
}

func (g *Goemq) PlatformGetClients(q ClientsRequest) (*CLientsData, error) {
	url, _ := url.ParseRequestURI(g.BaseURL + "/clients")
	url.RawQuery = q.QueryString()

	req, _ := http.NewRequest("GET", url.String(), nil)

	req.Header.Set("Authorization", g.getBasicAuthHeader())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data CLientsData

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (g *Goemq) PlatformGetClient(clientId string) (*ClientData, error) {
	url, _ := url.ParseRequestURI(g.BaseURL + "/clients/" + clientId)

	req, _ := http.NewRequest("GET", url.String(), nil)
	req.Header.Set("Authorization", g.getBasicAuthHeader())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data ClientData
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (g *Goemq) PlatformDeleteClient(clientId string) error {
	url, _ := url.ParseRequestURI(g.BaseURL + "/clients/" + clientId)

	req, _ := http.NewRequest("DELETE", url.String(), nil)
	req.Header.Set("Authorization", g.getBasicAuthHeader())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete client: %s", string(body))
	}
	return nil
}

func (g *Goemq) PlatformSubscribe(clientId, topic string, qos byte) error {
	url, _ := url.ParseRequestURI(g.BaseURL + "/clients/" + clientId + "/subscribe")

	payload := fmt.Sprintf(`{"topic":"%s","qos":%d}`, topic, qos)
	req, _ := http.NewRequest("POST", url.String(), bytes.NewBufferString(payload))
	req.Header.Set("Authorization", g.getBasicAuthHeader())
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to subscribe client: %s", string(body))
	}
	return nil
}

func (g *Goemq) PlatformUnsubscribe(clientId, topic string) error {
	url, _ := url.ParseRequestURI(g.BaseURL + "/clients/" + clientId + "/unsubscribe")

	payload := fmt.Sprintf(`{"topic":"%s"}`, topic)
	req, _ := http.NewRequest("POST", url.String(), bytes.NewBufferString(payload))
	req.Header.Set("Authorization", g.getBasicAuthHeader())
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to unsubscribe client: %s", string(body))
	}
	return nil
}
