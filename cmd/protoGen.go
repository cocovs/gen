/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	gg "github.com/Xuanwo/gg"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// protoGenCmd represents the protoGen command
var protoGenCmd = &cobra.Command{
	Use:   "proto-gen",
	Short: "proto2HandlerSvc gen ",
	Long:  `proto2HanderSvc gen`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Str("proto file path", protoConfig.ProtoFilePath).Str("out floder", protoConfig.OutHandlerFolder).Str("out svc", protoConfig.OutSvcFolder).Msg("generating proto files")

		methods, err := protoConfig.In()
		if err != nil {
			log.Err(err).Msg("unable to generate proto files")
			return
		}

		if protoConfig.OutHandlerFolder != "" {
			if err = protoConfig.OutHandlerFrame(methods); err != nil {
				log.Err(err).Msg("unable to generate proto files")
				return
			}
		}

		if protoConfig.OutSvcFolder != "" {
			if err = protoConfig.OutSvcFrame(methods); err != nil {
				log.Err(err).Msg("unable to generate proto files")
				return
			}
		}

		log.Info().Msg("proto files generated successfully")
	},
}

func init() {
	rootCmd.AddCommand(protoGenCmd)

	protoGenCmd.Flags().StringVar(&protoConfig.ProtoFilePath, "proto-file", "", "path to the proto file. proto 文件路径")
	protoGenCmd.Flags().StringVar(&protoConfig.OutHandlerFolder, "out-handler-folder", "", "path to the  handler output path. handler 模板函数输出路径")
	protoGenCmd.Flags().StringVar(&protoConfig.OutSvcFolder, "out-svc-folder", "", "path to the  svc output path. svc 模板函数输出路径")

	protoGenCmd.Flags().Bool("append", false, "option to append or truncate the file 是否追加式，默认覆盖式")
}

var protoConfig ProtoConfig

type ProtoConfig struct {
	ProtoFilePath    string //path to the proto file . example: /path/to/file.proto
	OutHandlerFolder string //path to the hander output folder . example: /path/to/output
	OutSvcFolder     string //path to the svc output folder . example: /path/to/output
	AppendOption     bool   //option to append or truncate the file
	PackageName      string //package name
}

type rpcMethod struct {
	Name         string
	RequestName  string
	ResponseName string
}

var defaultProtobufType map[string]string = map[string]string{
	"google.protobuf.StringValue": "wrapperspb.StringValue",
	"google.protobuf.Empty":       "emptypb.Empty",
}

func (c *ProtoConfig) OutSvcFrame(methods map[string][]*rpcMethod) error {
	for rpcServiceName, methods := range methods {
		f := gg.NewGroup()
		f.AddPackage("svc")

		f.NewImport().
			AddPath("context")

		f.NewStruct(rpcServiceName + "GrpcService")

		for _, method := range methods {
			f.AddLine()
			f.NewFunction(method.Name).
				WithReceiver("svc", gg.String("*"+rpcServiceName+"GrpcService")).
				AddParameter("ctx", "context.Context").
				AddResult("", "error").
				AddBody(gg.String("return  nil"))
		}

		file, err := os.OpenFile(c.OutSvcFolder+strings.ToLower(rpcServiceName)+"_svc.go", func() int {
			if c.AppendOption {
				return os.O_APPEND | os.O_CREATE | os.O_WRONLY
			}
			return os.O_CREATE | os.O_WRONLY | os.O_TRUNC

		}(), 0644)
		if err != nil {
			log.Err(err).Msg("open file")
		}

		if _, err := file.WriteString(f.String()); err != nil {
			log.Err(err).Msg("write err")
			return err
		}
		file.Close()
	}

	return nil
}

func (c *ProtoConfig) OutHandlerFrame(methods map[string][]*rpcMethod) error {

	for rpcServiceName, methods := range methods {
		f := gg.NewGroup()
		f.AddPackage("handler")
		f.NewImport().
			AddPath("google.golang.org/protobuf/types/known/emptypb").
			AddPath("google.golang.org/protobuf/types/known/wrapperspb").
			AddPath("context").
			AddAlias(c.PackageName+"/proto", "pb").
			AddAlias(c.PackageName+"/svc", "svc")

		f.NewStruct(rpcServiceName+"GrpcHandler").
			AddField(rpcServiceName+"Svc", "svc."+rpcServiceName+"GrpcService")

		for _, method := range methods {
			f.AddLine()
			f.NewFunction(method.Name).
				WithReceiver("h", gg.String("*"+rpcServiceName+"GrpcHandler")).
				AddParameter("ctx", "context.Context").
				AddParameter("in", func() string {
					if v, ok := defaultProtobufType[method.RequestName]; ok {
						return "*" + v
					}
					return "*pb." + method.RequestName
				}()).
				AddResult("", func() string {
					if v, ok := defaultProtobufType[method.RequestName]; ok {
						return "*" + v
					}
					return "*pb." + method.RequestName
				}()).
				AddResult("", "error").
				AddBody(gg.Call(method.Name).WithOwner("h." + rpcServiceName + "Svc").AddParameter("ctx")).
				AddBody(gg.String("return nil, nil"))
		}
		file, err := os.OpenFile(c.OutHandlerFolder+strings.ToLower(rpcServiceName)+"_handler.go", func() int {
			if c.AppendOption {
				return os.O_APPEND | os.O_CREATE | os.O_WRONLY
			}
			return os.O_CREATE | os.O_WRONLY | os.O_TRUNC

		}(), 0644)
		if err != nil {
			log.Err(err).Msg("open file")
		}
		if _, err := file.WriteString(f.String()); err != nil {
			log.Err(err).Msg("write err")
			return err
		}
		file.Close()
	}

	return nil
}

func (c *ProtoConfig) In() (map[string][]*rpcMethod, error) {
	file, err := os.Open(c.ProtoFilePath)
	if err != nil {
		log.Err(err).Msg("unable to open file")
	}
	defer file.Close()

	//key:service name
	services := make(map[string][]*rpcMethod)

	currentService := ""

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()

		if currentService == "" {
			if packageName := extractRPCPackage(s); packageName != "" {
				c.PackageName = packageName
				continue
			}
		}

		if serviceName := extractRPCService(s); serviceName != "" {
			currentService = serviceName
			log.Info().Str("service name", currentService).Msg("service name")
			continue
		}

		if rpcMethod, err := extractRPCMethod(s); err == nil {
			services[currentService] = append(services[currentService], rpcMethod)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Err(err).Msg("unable to read file")
	}

	return services, nil
}

func extractRPCPackage(line string) string {
	s := `package\s*(\w+)\;`
	serviceRegex := regexp.MustCompile(s)
	matches := serviceRegex.FindStringSubmatch(line)
	if len(matches) < 2 {
		return ""
	}

	return matches[1]
}

func extractRPCService(line string) string {
	s := `service\s*(\w+)`

	serviceRegex := regexp.MustCompile(s)
	matches := serviceRegex.FindStringSubmatch(line)
	if len(matches) < 2 {
		return ""
	}

	return matches[1]
}

func extractRPCMethod(line string) (*rpcMethod, error) {
	s := `rpc ([^\)]+)\(([^\)]+)\) returns \(([^\)]+)\)`

	rpcRegex := regexp.MustCompile(s)
	matches := rpcRegex.FindStringSubmatch(line)

	if len(matches) < 4 {
		return &rpcMethod{}, fmt.Errorf("unable to match RPC details")
	}

	return &rpcMethod{
		Name:         matches[1],
		RequestName:  matches[2],
		ResponseName: matches[3],
	}, nil
}
