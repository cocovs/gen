package main

import (
	"os"
	"regexp"
	"text/template"
)

type RpcMethod struct {
	Name       string `json:"name"`
	Request    string `json:"request"`
	Response   string `json:"response"`
	HTTPMethod string `json:"httpMethod"`
	URL        string `json:"url"`
}

const rpcTemplate = `
func (h *GrpcHandler) {{.Name}}(ctx context.Context, in *pb.{{.Request}}) (*pb.{{.Response}}, error) {
  
	c_list, err := h.CaseSvc.{{.Name}}(ctx, in.DomainID, in.MemberPath, in.NoCustomer)
    if err != nil {
        return nil, err
    }

    return nil, nil
}
`

func main() {
	content, err := os.ReadFile("example.proto")
	if err != nil {
		panic(err)
	}

	str := string(content)

	re := regexp.MustCompile(`rpc (\w+)\((\w+)\) returns \((\w+|\w+ \w+)\) {\s+option \(google.api.http\) = {\s+(\w+):"([^"]+)",`)
	matches := re.FindAllStringSubmatch(str, -1)

	methods := make([]RpcMethod, 0)
	for _, match := range matches {
		if len(match) > 5 {
			methods = append(methods, RpcMethod{
				Name:       match[1],
				Request:    match[2],
				Response:   match[3],
				HTTPMethod: match[4],
				URL:        match[5],
			})
		}
	}

	tmpl, err := template.New("rpc").Parse(rpcTemplate)
	if err != nil {
		panic(err)
	}

	filePath := "./handler.go"
	file, err := os.Create(filePath)
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
}
