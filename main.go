package main

import (
	"os"
	"regexp"
	"text/template"
)

func proto2Func(protoFilePath, outFilePath string) []RpcMethod {
	content, err := os.ReadFile(protoFilePath)
	if err != nil {
		panic(err)
	}

	str := string(content)
	re := regexp.MustCompile(must)
	matches := re.FindAllStringSubmatch(str, -1)

	methods := make([]RpcMethod, 0)
	for _, match := range matches {
		log.Debug().Interface("sss", match).Msg("ss")
		if len(match) >= 3 {
			if match[2] == "google.protobuf.StringValue" {
				match[2] = "wrapperspb.StringValue"
			}
			if match[3] == "google.protobuf.StringValue" {
				match[3] = "wrapperspb.StringValue"
			}
			if match[2] == "google.protobuf.Empty" {
				match[2] = "emptypb.Empty"
			}
			if match[3] == "google.protobuf.Empty" {
				match[3] = "emptypb.Empty"
			}
			methods = append(methods, RpcMethod{
				Name:     match[1],
				Request:  match[2],
				Response: match[3],
			})
		}
	}

	tmpl, err := template.New("rpc").Parse(rpcTemplate)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(outFilePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for _, v := range methods {
		err = tmpl.Execute(file, v)
		if err != nil {
			panic(err)
		}
	}
	return methods
}

type RpcMethod struct {
	Name     string `json:"name"`
	Request  string `json:"request"`
	Response string `json:"response"`
}

const must = `rpc (\w+)\(([\w\.]+)\) returns \(([\w\.]+)\) \{\};`

const rpcTemplate = `
func (h *GrpcHandler) {{.Name}}(ctx context.Context, in *pb.{{.Request}}) (*pb.{{.Response}}, error) {
  
	c_list, err := h.Svc.{{.Name}}(ctx)
    if err != nil {
        return nil, err
    }

    return nil, nil
}
`

