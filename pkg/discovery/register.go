package discovery

import (
	"context"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// Register 结构体用于保存注册服务所需的配置和状态。
type Register struct {
	EtcdAddrs   []string // Etcd服务器的地址列表。
	DialTimeout int      // 连接Etcd服务器的超时时间（秒）。

	closeCh     chan struct{}                           // 用于通知关闭服务的通道。
	leasesID    clientv3.LeaseID                        // etcd租约ID。
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse // 租约续约的响应通道。

	srvInfo Server           // 服务信息。
	srvTTL  int64            // 服务的TTL值，用于设置etcd中的租约时间。
	cli     *clientv3.Client // etcd客户端实例。
}

// NewRegister 创建一个新的Register实例。
func NewRegister(etcdAddrs []string) *Register {
	return &Register{
		EtcdAddrs:   etcdAddrs,
		DialTimeout: 3, // 默认超时时间设置为3秒。
	}
}

// Register 方法用于注册服务到etcd。
func (r *Register) Register(srvInfo Server, ttl int64) (chan<- struct{}, error) {
	var err error

	// 检查服务地址是否有效。
	if strings.Split(srvInfo.Addr, ":")[0] == "" {
		return nil, errors.New("invalid ip address")
	}

	// 创建etcd客户端。
	if r.cli, err = clientv3.New(clientv3.Config{
		Endpoints:   r.EtcdAddrs,
		DialTimeout: time.Duration(r.DialTimeout) * time.Second,
	}); err != nil {
		return nil, err
	}

	r.srvInfo = srvInfo
	r.srvTTL = ttl

	// 注册服务。
	if err = r.register(); err != nil {
		return nil, err
	}

	r.closeCh = make(chan struct{})

	// 启动保持租约的协程。
	go r.keepAlive()

	return r.closeCh, nil
}

// register 方法用于在etcd中注册服务。
func (r *Register) register() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.DialTimeout)*time.Second)
	defer cancel()

	// 获取etcd租约。
	leaseResp, err := r.cli.Grant(ctx, r.srvTTL)
	if err != nil {
		return err
	}

	r.leasesID = leaseResp.ID

	// 开始租约续约。
	if r.keepAliveCh, err = r.cli.KeepAlive(context.Background(), r.leasesID); err != nil {
		return err
	}

	// 将服务信息序列化后存储到etcd。
	data, err := json.Marshal(r.srvInfo)
	if err != nil {
		return err
	}

	_, err = r.cli.Put(context.Background(), BuildRegisterPath(r.srvInfo), string(data), clientv3.WithLease(r.leasesID))

	return err
}

// Stop 方法用于停止服务注册。
func (r *Register) Stop() {
	r.closeCh <- struct{}{}
}

// unregister 方法用于从etcd中删除服务。
func (r *Register) unregister() error {
	_, err := r.cli.Delete(context.Background(), BuildRegisterPath(r.srvInfo))
	return err
}

// keepAlive 方法用于维持etcd租约。
func (r *Register) keepAlive() {
	ticker := time.NewTicker(time.Duration(r.srvTTL) * time.Second)

	for {
		select {
		case <-r.closeCh:
			if err := r.unregister(); err != nil {
				zap.S().Error("unregister failed, error: ", err)
			}

			if _, err := r.cli.Revoke(context.Background(), r.leasesID); err != nil {
				zap.S().Error("revoke failed, error: ", err)
			}
		case res := <-r.keepAliveCh:
			if res == nil {
				if err := r.register(); err != nil {
					zap.S().Error("register failed, error: ", err)
				}
			}
		case <-ticker.C:
			if r.keepAliveCh == nil {
				if err := r.register(); err != nil {
					zap.S().Error("register failed, error: ", err)
				}
			}
		}
	}
}

// UpdateHandler 方法用于处理HTTP请求，更新服务权重。
func (r *Register) UpdateHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		weightstr := req.URL.Query().Get("weight")
		weight, err := strconv.Atoi(weightstr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		var update = func() error {
			r.srvInfo.Weight = int64(weight)
			data, err := json.Marshal(r.srvInfo)
			if err != nil {
				return err
			}

			_, err = r.cli.Put(context.Background(), BuildRegisterPath(r.srvInfo), string(data), clientv3.WithLease(r.leasesID))
			return err
		}

		if err := update(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		_, _ = w.Write([]byte("update service weight success"))
	})
}

// GetServerInfo 方法用于从etcd获取服务信息。
func (r *Register) GetServerInfo() (Server, error) {
	resp, err := r.cli.Get(context.Background(), BuildRegisterPath(r.srvInfo))
	if err != nil {
		return r.srvInfo, err
	}

	server := Server{}
	if resp.Count >= 1 {
		if err := json.Unmarshal(resp.Kvs[0].Value, &server); err != nil {
			return server, err
		}
	}

	return server, err
}
