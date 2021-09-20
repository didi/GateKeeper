package handler

import (
	"github.com/didi/gatekeeper/golang_common/lib"
	"github.com/didi/gatekeeper/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http/httptest"
	"time"
)

var AppManagerHandler *AppManager

func init() {
	AppManagerHandler = NewAppManager()
}

//通知事件
type AppEvent struct {
	DeleteApp []*model.App
	AddApp    []*model.App
	UpdateApp []*model.App
}

//观察者接口
type AppObserver interface {
	Update(*AppEvent)
}

//被观察者接口
type AppSubject interface {
	Regist(ServiceObserver)
	Deregist(ServiceObserver)
	Notify(*AppEvent)
}

func (s *AppManager) Regist(ob AppObserver) {
	s.Observers[ob] = true
}

func (s *AppManager) Deregist(ob AppObserver) {
	delete(s.Observers, ob)
}

func (s *AppManager) Notify(e *AppEvent) {
	for ob, _ := range s.Observers {
		ob.Update(e)
	}
}

type AppManager struct {
	AppMap    map[string]*model.App
	AppSlice  []*model.App
	err       error
	UpdateAt  time.Time
	Observers map[AppObserver]bool
}

func NewAppManager() *AppManager {
	return &AppManager{
		AppMap:   map[string]*model.App{},
		AppSlice: []*model.App{},
	}
}

func (s *AppManager) GetAppList() []*model.App {
	return s.AppSlice
}

func (s *AppManager) LoadApp() *AppManager {
	//log.Printf(" [INFO] AppManager.LoadApp begin\n")
	ns := NewAppManager()
	defer func() {
		if ns.err != nil {
			log.Printf(" [ERROR] AppManager.LoadApp error:%v\n", ns.err)
		}
	}()
	appInfo := &model.App{}
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/", nil)
	tx, err := lib.GetGormPool("default")
	if err != nil {
		ns.err = err
		return ns
	}
	params := &model.APPListInput{PageNo: 1, PageSize: 99999}
	list, _, err := appInfo.APPList(c, tx, params)
	if err != nil {
		ns.err = err
		return ns
	}
	for _, listItem := range list {
		tmpItem := listItem
		ns.AppMap[listItem.AppID] = &tmpItem
		ns.AppSlice = append(ns.AppSlice, &tmpItem)
		if listItem.UpdatedAt.Unix() > ns.UpdateAt.Unix() {
			ns.UpdateAt = listItem.UpdatedAt
		}
	}
	return ns
}

func (s *AppManager) LoadAndWatch() error {
	ns := s.LoadApp()
	if ns.err != nil {
		return ns.err
	}
	s.AppSlice = ns.AppSlice
	s.AppMap = ns.AppMap
	s.UpdateAt = ns.UpdateAt
	go func() {
		for true {
			time.Sleep(10 * time.Second)
			ns := s.LoadApp()
			if ns.err != nil {
				continue
			}
			if ns.UpdateAt != s.UpdateAt || len(ns.AppSlice) != len(s.AppSlice) {
				log.Printf("s.UpdateAt:%v ns.UpdateAt:%v\n", s.UpdateAt.Format(lib.TimeFormat), ns.UpdateAt.Format(lib.TimeFormat))
				e := &AppEvent{}

				//老服务存在，新服务不存在，则为删除
				for _, app := range s.AppSlice {
					matched := false
					for _, newApp := range ns.AppSlice {
						if app.AppID == newApp.AppID {
							matched = true
						}
					}
					if !matched {
						e.DeleteApp = append(e.DeleteApp, app)
					}
				}

				//新服务有，老服务不存在，则为新增
				for _, newApp := range ns.AppSlice {
					matched := false
					for _, app := range s.AppSlice {
						if app.AppID == newApp.AppID {
							matched = true
						}
					}
					if !matched {
						e.AddApp = append(e.AddApp, newApp)
					}
				}

				//服务名相同，更新时间不同，则为更新
				for _, newApp := range ns.AppSlice {
					matched := false
					for _, app := range s.AppSlice {
						if app.AppID == newApp.AppID && app.UpdatedAt != newApp.UpdatedAt {
							matched = true
						}
					}
					if matched {
						e.UpdateApp = append(e.UpdateApp, newApp)
					}
				}
				s.AppSlice = ns.AppSlice
				s.AppMap = ns.AppMap
				s.UpdateAt = ns.UpdateAt

				log.Printf("e:%v\n", e)
				s.Notify(e)
			}
		}
	}()
	return s.err
}
