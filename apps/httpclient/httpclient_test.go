package httpclient

import (
	"testing"
	"fmt"
)

func TestGet(t *testing.T){
	//url := "http://pv.mama100.com/pointstats"
	url := "http://10.50.115.19:2059/template/area-activity/getAdminRole?token=fd925122-2cb8-4295-9796-ebff44c1b833"
	fmt.Println("url resposne: \n ", string(Get(url)))

	///
}

//TODO ....导购模板添加管理员权限TEST
func TestPostWithJsonBody(t *testing.T){
	//http://10.50.115.19:2059
	url := "http://10.50.115.19:2059/template/activity-area/addAdminRole?token=0a88c9cd-0826-4037-a8fd-9aefea3981ce"
	//url := "http://192.168.2.65:2386/template/activity-area/addAdminRole?token=89b35d98-6e78-4a06-ba12-e52605eda1e1"
	data := []byte(`{"employeeId":"13123"}`)
	result, _ :=PostWithJsonBody(url, data)
	fmt.Println("url resposne: \n ", string(result))
}
