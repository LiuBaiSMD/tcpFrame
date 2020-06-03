package consul

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/encoder/json"
	"github.com/micro/go-micro/config/source"
	"github.com/micro/go-micro/config/source/file"
	"github.com/micro/go-micro/util/log"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

var (
	m                      sync.RWMutex
	inited                 bool
	err                    error
	consulAddr             consulConfig
	consulConfigCenterAddr string
)

// consulConfig 配置结构
type consulConfig struct {
	Enabled    bool   `json:"enabled"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	KVLocation string `json:"kv_location"`
	DockerHost string `json:"docker_host"`
}

func init() {
	// 上传服务注册配置
	UpdataConfig("127.0.0.1", 8500, ".", "service.json", "serverRegistry")

	UpdataConfig("127.0.0.1", 8500, ".", "plugin.json", "plugin")
}

// Init 初始化配置
func Init() {
	m.Lock()
	//进行配置推送检测，是否已经推送过配置
	defer m.Unlock()
	if inited {
		log.Logf("[Init] 配置已经初始化过")
		return
	}

	// 加载yml默认配置
	// 先加载基础配置
	//appPath, _ := filepath.Abs(filepath.Dir(filepath.Join("../", string(filepath.Separator))))
	var configs []string
	if err := FindFile("service.json", ".", &configs); err != nil {
		log.Log("寻找配置文件失败！")
		return
	}
	//现在先默认使用一个配置
	appPath := configs[0]
	e := json.NewEncoder()
	log.Log(appPath)
	fileSource := file.NewSource(
		file.WithPath(appPath),
		source.WithEncoder(e),
	)
	conf := config.NewConfig()
	// 加载micro.yml文件
	if err = conf.Load(fileSource); err != nil {
		panic(err)
	}
	log.Log("conf.Bytes():	", string(conf.Bytes()))

	//scan将配置读入到放入的变量consulAddr之中
	if err := conf.Get("consul_config").Scan(&consulAddr); err != nil {
		panic(err)
	}
	// 拼接配置的地址和 KVcenter 存储路径,对本地以及docker环境进行判断
	consulConfigCenterAddr = consulAddr.Host + ":" + strconv.Itoa(consulAddr.Port)

	url := fmt.Sprintf("http://%s/v1/kv/%s", consulConfigCenterAddr, consulAddr.KVLocation)
	_, err, _ := PutConfig(url, string(conf.Bytes()))
	if err != nil {
		log.Fatalf("http 发送模块异常，%s", err)
		panic(err)
	}
	// 侦听文件变动
	watcher, err := conf.Watch()
	if err != nil {
		log.Fatalf("[Init] 开始侦听应用配置文件变动 异常，%s", err)
		panic(err)
	}

	oldStrMap := make(map[string]string)
	oldStrMap = conf.Get().StringMap(oldStrMap)
	go func() {
		for {
			v, err := watcher.Next()
			if err != nil {
				log.Fatalf("[loadAndWatchConfigFile] 侦听应用配置文件变动 异常， %s", err)
				return
			}
			if err = conf.Load(fileSource); err != nil {
				panic(err)
			}
			log.Logf("[loadAndWatchConfigFile] 文件变动，%s", string(v.Bytes()))

			////本部分代码还有部分问题 1.对于底层修改、增删的部分只会认为是change
			strMap := make(map[string]string)
			newMapConf := v.StringMap(strMap)
			findConfDif(oldStrMap, newMapConf)

			_, err, _ = PutConfig(url, string(v.Bytes()))
			if err != nil {
				log.Fatalf("http 发送模块异常，%s", err)
				panic(err)
			}
			log.Log("配置重新上传完毕！")
			oldStrMap = deepCopy(newMapConf)
		}
		log.Log("Init over go!")
	}()
	// 标记已经初始化
	inited = true

	return
}

func UpdataConfig(consulServerIp string, port int, DirName, FileName, KVAddr string) error {
	m.Lock()
	//进行配置推送检测，是否已经推送过配置
	defer m.Unlock()
	if inited {
		log.Logf("[Init] 配置已经初始化过")
		return errors.New("[Init] 配置已经初始化过")
	}

	// 加载yml默认配置
	// 先加载基础配置
	//appPath, _ := filepath.Abs(filepath.Dir(filepath.Join("../", string(filepath.Separator))))
	var configs []string
	if err := FindFile(FileName, DirName, &configs); err != nil {
		log.Log("寻找配置文件失败！")
		return errors.New("寻找配置文件失败！")
	}
	//现在先默认使用一个配置
	if len(configs) <= 0 {
		log.Info("配置加载失败:", FileName)
		return nil
	}
	appPath := configs[0]
	e := json.NewEncoder()
	log.Log(appPath)
	fileSource := file.NewSource(
		//file.WithPath(appPath+"/conf/micro.yml"),
		file.WithPath(appPath),
		//file.WithPath("./conf/micro.yml"),
		source.WithEncoder(e),
	)
	conf := config.NewConfig()
	// 加载micro.yml文件
	if err = conf.Load(fileSource); err != nil {
		panic(err)
	}

	// 拼接配置的地址和 KVcenter 存储路径,对本地以及docker环境进行判断
	url := fmt.Sprintf("http://%s:%d/v1/kv/%s", consulServerIp, port, KVAddr)
	_, err, _ := PutConfig(url, string(conf.Bytes()))
	if err != nil {
		log.Fatalf("http 发送模块异常，%s", err)
		return errors.New(fmt.Sprintf("http 发送模块异常，%s", err))
	}
	return nil
}

func PutConfig(url, body string) (ret string, err error, resp *http.Response) {
	buf := bytes.NewBufferString(body)
	req, err := http.NewRequest("PUT", url, buf)
	if err != nil {
		panic(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
	if err != nil {
		log.Log(err.Error())
		return "", err, resp
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err, resp
	}

	return string(data), nil, resp
}

func findConfDif(oldConf map[string]string, newConf map[string]string) (addConf map[string]string, subConf map[string]string, changeConf map[string]string) {
	//遍历旧配置一遍查看减少的配置,和改变的配置
	addConf = make(map[string]string)
	subConf = make(map[string]string)
	changeConf = make(map[string]string)
	for key, value := range oldConf {
		if newData, ok := newConf[key]; ok {
			if newData != value {
				//在旧配置中存在却不相等的配置  changeConf
				changeConf[string(key)] = string(value)
			}
		} else {
			//旧配置中不存在的配置  subConf
			subConf[string(key)] = string(value)
		}
	}
	//遍历新配置  查看增加的配置
	for key, value := range newConf {
		//log.Log(key, ":", value)
		if _, ok := oldConf[key]; !ok {
			addConf[string(key)] = string(value)
		}
	}
	log.Log("add---------->", addConf)
	log.Log("sub---------->", subConf)
	log.Log("change------->", changeConf)
	return addConf, subConf, changeConf
}

func deepCopy(oldMap map[string]string) (newMap map[string]string) {
	//map[string]string只使用一层拷贝即可
	newMap = make(map[string]string)
	for key, value := range oldMap {
		newMap[key] = value
	}
	return newMap
}

func FindFile(filename, pathname string, filesPath *[]string) error {
	rd, err := ioutil.ReadDir(pathname)
	for _, fi := range rd {
		if fi.IsDir() {
			FindFile(filename, pathname+"/"+fi.Name(), filesPath)
		} else {
			if fi.Name() == filename {
				*filesPath = append(*filesPath, pathname+"/"+fi.Name())
			}
		}
	}
	return err
}
