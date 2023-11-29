package lsp

import (
	"net/url"
	"os"
	"time"

	"github.com/raiguard/luapls/lua/parser"
	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	glspserv "github.com/tliron/glsp/server"

	// To enable logging
	_ "github.com/tliron/commonlog/simple"
)

const LS_NAME = "luapls"

// Type Server contains the state for the LSP session.
type Server struct {
	config   Config
	envs     map[string]*Env
	files    map[string]*parser.File
	handler  protocol.Handler
	log      commonlog.Logger
	rootPath string
	server   *glspserv.Server
}

func Run(logLevel int) {
	commonlog.Configure(logLevel, ptr("/tmp/luapls.log"))

	s := Server{
		envs:  map[string]*Env{},
		files: map[string]*parser.File{},
	}

	s.handler.Initialize = s.initialize
	s.handler.Initialized = s.initialized
	s.handler.Shutdown = s.shutdown
	s.handler.SetTrace = s.setTrace
	s.handler.TextDocumentDidOpen = s.textDocumentDidOpen
	s.handler.TextDocumentDidChange = s.textDocumentDidChange
	s.handler.TextDocumentDocumentHighlight = s.textDocumentHighlight
	s.handler.TextDocumentHover = s.textDocumentHover
	s.handler.TextDocumentCompletion = s.textDocumentCompletion
	s.handler.TextDocumentSelectionRange = s.textDocumentSelectionRange
	s.handler.TextDocumentDefinition = s.textDocumentDefinition

	s.server = glspserv.NewServer(&s.handler, LS_NAME, logLevel > 2)

	s.log = s.server.Log

	s.server.RunStdio()
}

func (s *Server) initialize(ctx *glsp.Context, params *protocol.InitializeParams) (any, error) {
	capabilities := s.handler.CreateServerCapabilities()
	s.rootPath = *params.RootPath

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo:   &protocol.InitializeResultServerInfo{Name: LS_NAME},
	}, nil
}

func (s *Server) initialized(ctx *glsp.Context, params *protocol.InitializedParams) error {
	go func() {
		s.getConfiguration(ctx)

		s.log.Debug("Parsing files")
		before := time.Now()
		for _, env := range s.envs {
			for _, path := range env.getFiles() {
				if s.files[path] != nil {
					continue
				}
				src, err := os.ReadFile(path)
				if err != nil {
					s.log.Errorf("Failed to parse file %s: %s", path, err)
					continue
				}
				file := parser.New(string(src)).ParseFile()
				s.files[path] = &file
				s.log.Debugf("Parsed file '%s'", path)
			}
		}
		s.log.Debugf("Initial parse: %s", time.Since(before).String())
	}()
	return nil
}

func (s *Server) shutdown(ctx *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func (s *Server) setTrace(ctx *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}

func (s *Server) getFile(uri string) *parser.File {
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		s.log.Errorf("Failed to parse file URI: %s", err)
	}
	file := s.files[u.Path]
	if file == nil {
		s.log.Errorf("File '%s' does not belong to any environment", u.Path)
	}
	return file
}
