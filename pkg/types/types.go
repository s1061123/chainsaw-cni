package types

import (
  "encoding/json"
  "fmt"
  cnitypes "github.com/containernetworking/cni/pkg/types"
  current "github.com/containernetworking/cni/pkg/types/040"
  cniVersion "github.com/containernetworking/cni/pkg/version"
  "net"
)

// NetConf is our definition for the CNI configuration
type NetConf struct {
  cnitypes.NetConf
  PrevResult       *current.Result `json:"-"`
  Foo              string          `json:"foo"`
  FilterExpression string          `json:"filter_expression"`
  SocketEnabled    bool            `json:"socket_enabled"`
  SocketPath       string          `json:"socket_path"`
  Kubeconfig       string          `json:"kubeconfig"`
}

type K8sArgs struct {
  cnitypes.CommonArgs
  IP                         net.IP
  K8S_POD_NAME               cnitypes.UnmarshallableString
  K8S_POD_NAMESPACE          cnitypes.UnmarshallableString
  K8S_POD_INFRA_CONTAINER_ID cnitypes.UnmarshallableString
  K8S_POD_UID                cnitypes.UnmarshallableString
}

// LoadNetConf parses our cni configuration
func LoadNetConf(bytes []byte) (*NetConf, error) {
  conf := NetConf{
    SocketEnabled: true,
    SocketPath:    "/var/run/chainsaw-cni/chainsaw.sock",
    Kubeconfig:    "/etc/cni/net.d/chainsaw.d/chainsaw.kubeconfig",
  }
  if err := json.Unmarshal(bytes, &conf); err != nil {
    return nil, fmt.Errorf("failed to load netconf: %s", err)
  }

  // Parse previous result
  if conf.RawPrevResult != nil {
    resultBytes, err := json.Marshal(conf.RawPrevResult)
    if err != nil {
      return nil, fmt.Errorf("could not serialize prevResult: %v", err)
    }

    res, err := cniVersion.NewResult(conf.CNIVersion, resultBytes)

    if err != nil {
      return nil, fmt.Errorf("could not parse prevResult: %v", err)
    }

    conf.RawPrevResult = nil
    conf.PrevResult, err = current.NewResultFromResult(res)
    if err != nil {
      return nil, fmt.Errorf("could not convert result to current version: %v", err)
    }
  }

  return &conf, nil
}
