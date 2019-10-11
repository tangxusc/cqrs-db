package handler

import protocol "github.com/tangxusc/mongo-protocol"

type ListDatabases struct {
}

//TODO:更细致的判断
/*
db.runCommand( { listDatabases: 1 } )
{"Header":{"MessageLength":0,"RequestID":0,"ResponseTo":0,"OpCode":0},"Flags":4,"FullCollectionName":"admin.$cmd","NumberToSkip":0,"NumberToReturn":1,"Query":{"$query":{"listDatabases":1,"nameOnly":true},"$readPreference":{"mode":"secondaryPreferred"}},"ReturnFieldsSelector":null}
{
	"databases" : [
		{
			"name" : "admin",
			"sizeOnDisk" : 102400,
			"empty" : false
		},
		{
			"name" : "config",
			"sizeOnDisk" : 98304,
			"empty" : false
		},
		{
			"name" : "local",
			"sizeOnDisk" : 73728,
			"empty" : false
		}
	],
	"totalSize" : 274432,
	"ok" : 1
}
*/
func (l *ListDatabases) Support(query *protocol.Query) bool {
	_, ok := query.Query["$query"]
	return ok
}

func (l *ListDatabases) Process(query *protocol.Query, reply *protocol.Reply) {
	reply.NumberReturned = 1
	reply.Documents = map[string]interface{}{
		"totalSize": 8888,
		"ok":        1,
		"databases": []interface{}{
			map[string]interface{}{
				"name":       "aggregate",
				"sizeOnDisk": 8888,
				"empty":      false,
			},
		},
	}
}

type IsMaster struct {
}

func (i *IsMaster) Support(query *protocol.Query) bool {
	_, ok := query.Query["isMaster"]
	return ok
}

func (i *IsMaster) Process(query *protocol.Query, reply *protocol.Reply) {
	reply.NumberReturned = 1
	reply.Documents = map[string]interface{}{"isMaster": 1, "ok": 1}
}

type Whatsmyuri struct {
}

func (w *Whatsmyuri) Support(query *protocol.Query) bool {
	_, ok := query.Query["whatsmyuri"]
	return ok
}

//TODO:返回客户端更为准确的信息
func (w *Whatsmyuri) Process(query *protocol.Query, reply *protocol.Reply) {
	reply.NumberReturned = 1
	reply.Documents = map[string]interface{}{"you": "118.114.245.36:48780", "ok": 1}
}

type BuildInfo struct {
}

func (w *BuildInfo) Support(query *protocol.Query) bool {
	_, ok := query.Query["buildinfo"]
	if !ok {
		_, ok = query.Query["buildInfo"]
	}
	return ok
}
func (w *BuildInfo) Process(query *protocol.Query, reply *protocol.Reply) {
	reply.NumberReturned = 1
	reply.Documents = map[string]interface{}{
		"version":          "4.2.0",
		"gitVersion":       "a4b751dcf51dd249c5865812b390cfd1c0129c30",
		"modules":          make([]string, 0),
		"allocator":        "tcmalloc",
		"javascriptEngine": "mozjs",
		"sysInfo":          "deprecated",
		"versionArray":     [4]int32{4, 2, 0, 0},
		"openssl": map[string]interface{}{
			"running":  "OpenSSL 1.1.1  11 Sep 2018",
			"compiled": "OpenSSL 1.1.1  11 Sep 2018",
		},
		"buildEnvironment": map[string]interface{}{
			"distmod":     "ubuntu1804",
			"distarch":    "x86_64",
			"cc":          "/opt/mongodbtoolchain/v3/bin/gcc: gcc (GCC) 8.2.0",
			"ccflags":     "-fno-omit-frame-pointer -fno-strict-aliasing -ggdb -pthread -Wall -Wsign-compare -Wno-unknown-pragmas -Winvalid-pch -Werror -O2 -Wno-unused-local-typedefs -Wno-unused-function -Wno-deprecated-declarations -Wno-unused-const-variable -Wno-unused-but-set-variable -Wno-missing-braces -fstack-protector-strong -fno-builtin-memcmp",
			"cxx":         "/opt/mongodbtoolchain/v3/bin/g++: g++ (GCC) 8.2.0",
			"cxxflags":    "-Woverloaded-virtual -Wno-maybe-uninitialized -fsized-deallocation -std=c++17",
			"linkflags":   "-pthread -Wl,-z,now -rdynamic -Wl,--fatal-warnings -fstack-protector-strong -fuse-ld=gold -Wl,--build-id -Wl,--hash-style=gnu -Wl,-z,noexecstack -Wl,--warn-execstack -Wl,-z,relro",
			"target_arch": "x86_64",
			"target_os":   "linux",
		},
		"bits":              64,
		"debug":             false,
		"maxBsonObjectSize": 16777216,
		"storageEngines": [4]string{
			"biggie",
			"devnull",
			"ephemeralForTest",
			"wiredTiger",
		},
		"ok": 1,
	}
}

type ServerStatus struct {
}

func (w *ServerStatus) Support(query *protocol.Query) bool {
	_, ok := query.Query["serverStatus"]
	return ok
}

//TODO:考虑和返回客户详细信息进行同样的处理
func (w *ServerStatus) Process(query *protocol.Query, reply *protocol.Reply) {
	reply.NumberReturned = 1
	reply.Documents = map[string]interface{}{"you": "118.114.245.36:48780", "ok": 1}
}

func GetBaseQueryHandler() []protocol.QueryHandler {
	return []protocol.QueryHandler{&ListDatabases{}, &IsMaster{}, &Whatsmyuri{}, &BuildInfo{}, &ServerStatus{}}
}
